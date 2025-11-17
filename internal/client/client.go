package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	APIToken   string
	HTTPClient *http.Client
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

type APIResponse struct {
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewClient(baseURL, apiToken string) *Client {
	return &Client{
		BaseURL:  baseURL,
		APIToken: apiToken,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) doRequest(method, path string, body interface{}) (*APIResponse, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Asegurar que el token se envía correctamente
	if c.APIToken == "" {
		return nil, fmt.Errorf("API token is empty")
	}
	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Retry logic para errores 401 intermitentes (hasta 3 intentos)
	maxRetries := 3
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Esperar un poco antes de reintentar (exponential backoff)
			time.Sleep(time.Duration(attempt*100) * time.Millisecond)
			// Recrear el request para cada intento (el body solo se puede leer una vez)
			req, err = http.NewRequest(method, c.BaseURL+path, reqBody)
			if err != nil {
				return nil, fmt.Errorf("error recreating request: %w", err)
			}
			req.Header.Set("Authorization", "Bearer "+c.APIToken)
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")
		}

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error making request: %w", err)
			continue
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("error reading response body: %w", err)
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			apiErr := &APIError{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}

			// Intentar parsear el error como JSON para mejor mensaje
			var errorResp map[string]interface{}
			if err := json.Unmarshal(respBody, &errorResp); err == nil {
				if errorMsg, ok := errorResp["error"].(string); ok && errorMsg != "" {
					apiErr.Message = errorMsg
				}
			}

			lastErr = apiErr
			// Si es 401 y no es el último intento, reintentar
			if resp.StatusCode == 401 && attempt < maxRetries-1 {
				continue
			}
			return nil, lastErr
		}

		var apiResp APIResponse
		if err := json.Unmarshal(respBody, &apiResp); err != nil {
			return nil, fmt.Errorf("error unmarshaling response: %w", err)
		}

		return &apiResp, nil
	}

	return nil, lastErr
}

func (c *Client) Get(path string) (*APIResponse, error) {
	return c.doRequest("GET", path, nil)
}

func (c *Client) Post(path string, body interface{}) (*APIResponse, error) {
	return c.doRequest("POST", path, body)
}

func (c *Client) Delete(path string) (*APIResponse, error) {
	return c.doRequest("DELETE", path, nil)
}
