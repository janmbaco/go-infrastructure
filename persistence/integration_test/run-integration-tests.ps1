# Script para ejecutar tests de integración con Docker databases
# Uso: .\run-integration-tests.ps1

param(
    [switch]$SkipCleanup,
    [switch]$KeepContainers
)

Write-Host "`n==============================================================================" -ForegroundColor Cyan
Write-Host "  🧪 INTEGRATION TESTING go-infrastructure with Docker DBs" -ForegroundColor Green
Write-Host "==============================================================================" -ForegroundColor Cyan
Write-Host ""

# Función para esperar a que un servicio esté listo
function Wait-ForService {
    param(
        [string]$ServiceName,
        [string]$HostName,
        [int]$Port,
        [int]$TimeoutSeconds = 60
    )

    Write-Host "⏳ Waiting for $ServiceName to be ready..." -ForegroundColor Yellow

    $startTime = Get-Date
    $timeout = New-TimeSpan -Seconds $TimeoutSeconds

    while ((Get-Date) - $startTime -lt $timeout) {
        try {
            $tcpClient = New-Object System.Net.Sockets.TcpClient
            $tcpClient.Connect($HostName, $Port)
            $tcpClient.Close()
            Write-Host "✅ $ServiceName is ready!" -ForegroundColor Green
            return $true
        }
        catch {
            Start-Sleep -Seconds 2
        }
    }

    Write-Host "❌ Timeout waiting for $ServiceName" -ForegroundColor Red
    return $false
}

# Función para verificar conectividad de base de datos
function Test-DatabaseConnection {
    param(
        [string]$ServiceName,
        [string]$Engine,
        [string]$HostName,
        [string]$Port,
        [string]$User,
        [string]$Password,
        [string]$Database,
        [int]$TimeoutSeconds = 30
    )

    Write-Host "🔍 Testing $ServiceName database connection..." -ForegroundColor Yellow

    $startTime = Get-Date
    $timeout = New-TimeSpan -Seconds $TimeoutSeconds

    while ((Get-Date) - $startTime -lt $timeout) {
        try {
            # Intentar ejecutar el programa de test de conexión
            $env:CGO_ENABLED = 0
            $testResult = & go run ./cmd/testdbconnection $Engine $HostName $Port $User $Password $Database 2>$null

            if ($LASTEXITCODE -eq 0) {
                Write-Host "✅ $ServiceName database connection successful!" -ForegroundColor Green
                return $true
            }
        }
        catch {
            # Ignorar errores y continuar intentando
        }

        Start-Sleep -Seconds 2
    }

    Write-Host "❌ Failed to connect to $ServiceName database" -ForegroundColor Red
    return $false
}

try {
    # Paso 1: Verificar que Docker esté corriendo
    Write-Host "🐳 Paso 1/6: Verificando Docker..." -ForegroundColor Yellow
    $dockerVersion = docker --version 2>$null
    if ($LASTEXITCODE -ne 0) {
        throw "Docker no está instalado o no está corriendo. Por favor instala Docker Desktop."
    }
    Write-Host "✅ Docker está disponible" -ForegroundColor Green
    Write-Host ""

    # Paso 2: Limpiar contenedores previos si existen
    if (-not $SkipCleanup) {
        Write-Host "🧹 Paso 2/6: Limpiando contenedores previos..." -ForegroundColor Yellow
        docker-compose -f docker-compose.test.yml down -v 2>$null | Out-Null
        Write-Host "✅ Limpieza completada" -ForegroundColor Green
    } else {
        Write-Host "⏭️  Paso 2/6: Saltando limpieza de contenedores previos" -ForegroundColor Yellow
    }
    Write-Host ""

    # Paso 3: Levantar servicios de base de datos
    Write-Host "🚀 Paso 3/6: Levantando servicios de base de datos..." -ForegroundColor Yellow
    docker-compose -f docker-compose.test.yml up -d

    if ($LASTEXITCODE -ne 0) {
        throw "Error al levantar los servicios de Docker"
    }
    Write-Host "✅ Servicios levantados" -ForegroundColor Green
    Write-Host ""

    # Paso 4: Esperar a que los servicios estén listos
    Write-Host "⏳ Paso 4/6: Esperando a que las bases de datos estén listas..." -ForegroundColor Yellow

    # PostgreSQL
    if (-not (Wait-ForService -ServiceName "PostgreSQL" -HostName "localhost" -Port 5432)) {
        throw "PostgreSQL no está listo"
    }

    if (-not (Test-DatabaseConnection -ServiceName "PostgreSQL" -Engine "postgres" -HostName "localhost" -Port "5432" -User "testuser" -Password "testpass" -Database "testdb")) {
        throw "PostgreSQL database connection failed"
    }

    # MySQL
    if (-not (Wait-ForService -ServiceName "MySQL" -HostName "localhost" -Port 3306)) {
        throw "MySQL no está listo"
    }

    if (-not (Test-DatabaseConnection -ServiceName "MySQL" -Engine "mysql" -HostName "localhost" -Port "3306" -User "testuser" -Password "testpass" -Database "testdb")) {
        throw "MySQL database connection failed"
    }

    # SQL Server (este puede tomar más tiempo)
    if (-not (Wait-ForService -ServiceName "SQL Server" -HostName "localhost" -Port 1433 -TimeoutSeconds 120)) {
        throw "SQL Server no está listo"
    }

    if (-not (Test-DatabaseConnection -ServiceName "SQL Server" -Engine "sqlserver" -HostName "localhost" -Port "1433" -User "sa" -Password "StrongPass123!" -Database "master")) {
        throw "SQL Server database connection failed"
    }

    Write-Host "✅ Todas las bases de datos están listas y conectables" -ForegroundColor Green
    Write-Host ""

    # Paso 5: Ejecutar tests de integración
    Write-Host "🧪 Paso 5/6: Ejecutando tests de integración..." -ForegroundColor Yellow

    # Configurar variables de entorno para las bases de datos
    $env:POSTGRES_HOST = "localhost"
    $env:POSTGRES_PORT = "5432"
    $env:POSTGRES_USER = "testuser"
    $env:POSTGRES_PASSWORD = "testpass"
    $env:POSTGRES_DB = "testdb"

    $env:MYSQL_HOST = "localhost"
    $env:MYSQL_PORT = "3306"
    $env:MYSQL_USER = "testuser"
    $env:MYSQL_PASSWORD = "testpass"
    $env:MYSQL_DB = "testdb"

    $env:SQLSERVER_HOST = "localhost"
    $env:SQLSERVER_PORT = "1433"
    $env:SQLSERVER_USER = "sa"
    $env:SQLSERVER_PASSWORD = "StrongPass123!"
    $env:SQLSERVER_DB = "master"

    # Ejecutar tests con tag integration
    $env:CGO_ENABLED = 0
    Push-Location ../..
    try {
        go test -tags=integration -v ./persistence/integration_test
    }
    finally {
        Pop-Location
    }

    $testExitCode = $LASTEXITCODE

    Write-Host ""

    # Paso 6: Resultados
    if ($testExitCode -eq 0) {
        Write-Host "==============================================================================" -ForegroundColor Cyan
        Write-Host "  ✅ TODOS LOS TESTS DE INTEGRACIÓN PASARON EXITOSAMENTE" -ForegroundColor Green
        Write-Host "==============================================================================" -ForegroundColor Cyan
    } else {
        Write-Host "==============================================================================" -ForegroundColor Cyan
        Write-Host "  ❌ ALGUNOS TESTS DE INTEGRACIÓN FALLARON" -ForegroundColor Red
        Write-Host "==============================================================================" -ForegroundColor Cyan
        $scriptExitCode = 1
    }

} catch {
    Write-Host "❌ Error durante la ejecución: $($_.Exception.Message)" -ForegroundColor Red
    $scriptExitCode = 1
} finally {
    Write-Host ""

    # Paso final: Limpiar contenedores
    if (-not $KeepContainers) {
        Write-Host "🧹 Limpiando contenedores..." -ForegroundColor Yellow
        docker-compose -f docker-compose.test.yml down -v 2>$null | Out-Null
        Write-Host "✅ Contenedores limpiados" -ForegroundColor Green
    } else {
        Write-Host "📦 Manteniendo contenedores activos (usar -KeepContainers:$false para limpiar)" -ForegroundColor Yellow
    }

    Write-Host ""
    exit $scriptExitCode
}
