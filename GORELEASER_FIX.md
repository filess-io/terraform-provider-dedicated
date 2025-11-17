# GoReleaser Archive Name Fix

## ❌ Problema Original

El error ocurría porque GoReleaser intentaba crear múltiples archivos ARM con el mismo nombre:

```
archive named dist/terraform-provider_1.0.0_linux_arm.zip already exists
```

Esto sucede cuando se compilan múltiples versiones de ARM (v6 y v7) pero el template de nombre no las diferencia.

## ✅ Solución Aplicada

Se actualizó el archivo `.goreleaser.yml` con un `name_template` que diferencia correctamente todas las arquitecturas:

### Antes:
```yaml
archives:
  - format: zip
    name_template: "{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}"
```

**Problema**: `armv6` y `armv7` ambos generaban `linux_arm`

### Después:
```yaml
archives:
  - format: zip
    name_template: >-
      {{ .ProjectName }}_{{ .Version }}_
      {{- if eq .Os "darwin" }}darwin
      {{- else if eq .Os "linux" }}linux
      {{- else if eq .Os "windows" }}windows
      {{- else }}{{ .Os }}{{ end }}_
      {{- if eq .Arch "amd64" }}amd64
      {{- else if eq .Arch "386" }}386
      {{- else if eq .Arch "arm64" }}arm64
      {{- else if eq .Arch "arm" }}armv{{ .Arm }}
      {{- else }}{{ .Arch }}{{ end }}
```

**Solución**: Ahora genera nombres únicos:
- `linux_armv6`
- `linux_armv7`

## 📦 Archivos Generados

Después del fix, GoReleaser creará:

```
terraform-provider-filess_1.0.0_linux_amd64.zip
terraform-provider-filess_1.0.0_linux_arm64.zip
terraform-provider-filess_1.0.0_linux_armv6.zip
terraform-provider-filess_1.0.0_linux_armv7.zip
terraform-provider-filess_1.0.0_windows_amd64.zip
terraform-provider-filess_1.0.0_darwin_amd64.zip
terraform-provider-filess_1.0.0_darwin_arm64.zip
```

Cada archivo tiene un nombre único que identifica correctamente la plataforma y arquitectura.

## 🔧 Configuración Completa

### `.goreleaser.yml` Actualizado

Incluye:
- ✅ Nombres únicos para cada plataforma/arquitectura
- ✅ Formato compatible con Terraform Registry
- ✅ Checksums SHA256
- ✅ Soporte para firma GPG (opcional)
- ✅ Changelog automático agrupado por tipo
- ✅ Release draft/prerelease automático

### Plataformas Soportadas

| OS | Arquitectura | Archivo Generado |
|---|---|---|
| Linux | amd64 | `linux_amd64.zip` |
| Linux | arm64 | `linux_arm64.zip` |
| Linux | armv6 | `linux_armv6.zip` |
| Linux | armv7 | `linux_armv7.zip` |
| macOS | amd64 | `darwin_amd64.zip` |
| macOS | arm64 (Apple Silicon) | `darwin_arm64.zip` |
| Windows | amd64 | `windows_amd64.zip` |

## 🚀 Cómo Crear un Release

1. **Commit y push cambios**:
   ```bash
   git add .
   git commit -m "feat: your changes"
   git push origin main
   ```

2. **Crear y push tag**:
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```

3. **GitHub Actions ejecutará automáticamente**:
   - Compilará binarios para todas las plataformas
   - Creará archivos ZIP con nombres únicos
   - Generará checksums SHA256
   - Creará el release en GitHub
   - Adjuntará todos los archivos

## 🔍 Verificar el Release

Una vez completado el workflow:

1. Ve a: `https://github.com/YOUR_ORG/terraform-provider-filess/releases`
2. Verifica que todos los archivos estén presentes:
   - 7 archivos `.zip` (uno por plataforma)
   - 1 archivo `SHA256SUMS`
   - Opcionalmente 1 archivo `.sig` si usas GPG

## ⚙️ Configuración Opcional: GPG Signing

Para firmar los releases con GPG:

1. **Generar clave GPG**:
   ```bash
   gpg --full-generate-key
   ```

2. **Obtener fingerprint**:
   ```bash
   gpg --list-secret-keys --keyid-format=long
   ```

3. **Exportar clave pública**:
   ```bash
   gpg --armor --export YOUR_KEY_ID
   ```

4. **Agregar a GitHub Secrets**:
   - `GPG_PRIVATE_KEY`: La clave privada
   - `GPG_FINGERPRINT`: El fingerprint
   - `PASSPHRASE`: La contraseña (si la tiene)

5. **Actualizar workflow** (ya configurado en `.goreleaser.yml`)

## 📝 Notas Importantes

1. **Nombres de Archivos**: Los nombres siguen la convención de Terraform Registry
2. **Versión en el Binario**: Se incluye la versión en el nombre del binario interno
3. **Changelog**: Se genera automáticamente desde commits
4. **Compatible**: Funciona con Terraform Registry requirements

## ✅ Estado

- [x] `.goreleaser.yml` actualizado con template correcto
- [x] Workflow de GitHub Actions configurado
- [x] Nombres únicos para todas las arquitecturas
- [x] Formato compatible con Terraform Registry
- [x] Changelog automático configurado

## 🐛 Troubleshooting

### Error: "archive already exists"
**Solución**: Ya está resuelto con el nuevo template de nombres

### Error: "GPG signing failed"
**Solución**: Si no quieres usar GPG, comenta la sección `signs:` en `.goreleaser.yml`

### Error: "permission denied"
**Solución**: Asegúrate de que el workflow tiene `permissions: contents: write`

---

**Estado: ✅ RESUELTO Y LISTO PARA RELEASE**

