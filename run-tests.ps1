# Script para buildear y ejecutar tests de go-infrastructure usando Docker
# Uso: .\run-tests.ps1

Write-Host "`n==============================================================================" -ForegroundColor Cyan
Write-Host "  🧪 BUILDING & TESTING go-infrastructure" -ForegroundColor Green
Write-Host "==============================================================================" -ForegroundColor Cyan
Write-Host ""

# Paso 1: Build de la imagen Docker
Write-Host "📦 Paso 1/3: Construyendo imagen Docker de test..." -ForegroundColor Yellow
docker build -f Dockerfile.test -t go-infrastructure-test:latest .

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Error al construir la imagen Docker" -ForegroundColor Red
    exit 1
}

Write-Host "✅ Imagen construida exitosamente" -ForegroundColor Green
Write-Host ""

# Paso 2: Ejecutar tests
Write-Host "🧪 Paso 2/3: Ejecutando tests..." -ForegroundColor Yellow
docker run --rm -v "${PWD}:/app" go-infrastructure-test:latest

$testExitCode = $LASTEXITCODE

Write-Host ""

# Paso 3: Mostrar resultados
if ($testExitCode -eq 0) {
    Write-Host "==============================================================================" -ForegroundColor Cyan
    Write-Host "  ✅ TODOS LOS TESTS PASARON EXITOSAMENTE" -ForegroundColor Green
    Write-Host "==============================================================================" -ForegroundColor Cyan
    
    # Verificar si se generó coverage
    if (Test-Path "coverage.out") {
        Write-Host ""
        Write-Host "📊 Generando reporte de cobertura HTML..." -ForegroundColor Yellow
        docker run --rm -v "${PWD}:/app" go-infrastructure-test:latest go tool cover -html=coverage.out -o coverage.html
        Write-Host "✅ Reporte generado: coverage.html" -ForegroundColor Green
    }
} else {
    Write-Host "==============================================================================" -ForegroundColor Cyan
    Write-Host "  ❌ ALGUNOS TESTS FALLARON" -ForegroundColor Red
    Write-Host "==============================================================================" -ForegroundColor Cyan
}

Write-Host ""
exit $testExitCode
