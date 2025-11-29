# Integration Tests

Esta carpeta contiene tests de integración que verifican el funcionamiento del `dataaccess` con diferentes bases de datos usando Docker.

## Requisitos

- Docker y Docker Compose instalados
- PowerShell (para Windows) o Bash (para Linux/Mac)

## Bases de Datos Soportadas

Los tests de integración cubren las siguientes bases de datos:

- **PostgreSQL 15**
- **MySQL 8.0**
- **SQL Server 2022**

## Ejecutar Tests de Integración

### Opción 1: Usando Make (recomendado)

```bash
make test-integration
```

### Opción 2: Usando PowerShell directamente

```powershell
# Desde la carpeta persistence/integration_test
.\run-integration-tests.ps1
```

### Opción 3: Usando Go directamente

Primero, levantar las bases de datos:

```bash
# Desde la carpeta persistence/integration_test
docker-compose -f docker-compose.test.yml up -d
```

Esperar a que las bases de datos estén listas, luego ejecutar:

```bash
# Desde la raíz del proyecto
go test -tags=integration -v ./persistence/integration_test -run TestDataAccessIntegration
```

## Configuración de Variables de Entorno

Los tests usan las siguientes variables de entorno (con valores por defecto):

### PostgreSQL
- `POSTGRES_HOST=localhost`
- `POSTGRES_PORT=5432`
- `POSTGRES_USER=testuser`
- `POSTGRES_PASSWORD=testpass`
- `POSTGRES_DB=testdb`

### MySQL
- `MYSQL_HOST=localhost`
- `MYSQL_PORT=3306`
- `MYSQL_USER=testuser`
- `MYSQL_PASSWORD=testpass`
- `MYSQL_DB=testdb`

### SQL Server
- `SQLSERVER_HOST=localhost`
- `SQLSERVER_PORT=1433`
- `SQLSERVER_USER=sa`
- `SQLSERVER_PASSWORD=StrongPass123!`
- `SQLSERVER_DB=master`

## Opciones del Script

### PowerShell Script Options

- `-SkipCleanup`: No limpia contenedores previos
- `-KeepContainers`: Mantiene los contenedores activos después de los tests

```powershell
# Saltar limpieza y mantener contenedores
.\run-integration-tests.ps1 -SkipCleanup -KeepContainers
```

## Lo que se prueba

Los tests de integración verifican:

1. **Conexiones a base de datos**: Verifica que se pueda conectar a cada tipo de base de datos
2. **Operaciones CRUD**:
   - Crear registros
   - Leer registros (por ID y con filtros)
   - Actualizar registros
   - Eliminar registros
3. **Relaciones**: Tests con asociaciones entre tablas usando `preload`
4. **Manejo de errores**:
   - Constraints únicos
   - Tipos de datos inválidos
   - Conexiones fallidas
5. **Funciones genéricas**: Verifica que las funciones con generics (`InsertRow`, `SelectRows`, etc.) funcionen correctamente

## Funciones Genéricas Probadas

```go
// Crear DataAccess con type safety
dataAccess := NewTypedDataAccess[TestUser](db)

// Operaciones CRUD con generics
err := InsertRow(dataAccess, user)
users, err := SelectRows[TestUser](dataAccess, filter)
err = UpdateRow(dataAccess, filter, updatedUser)
err = DeleteRows(dataAccess, filter)
```

## Limpieza

Los contenedores se limpian automáticamente después de ejecutar los tests. Si quieres mantenerlos activos para debugging, usa la opción `-KeepContainers`.

Para limpiar manualmente:

```bash
docker-compose -f docker-compose.test.yml down -v
```

## Troubleshooting

### Error de conexión
- Verifica que Docker esté corriendo
- Espera a que los contenedores estén completamente inicializados (especialmente SQL Server que toma más tiempo)
- Revisa los logs de Docker: `docker-compose -f docker-compose.test.yml logs`

### Tests fallan
- Asegúrate de que no haya otros procesos usando los puertos 5432, 3306, 1433
- Verifica que las credenciales en `docker-compose.test.yml` coincidan con las variables de entorno
- Revisa los logs detallados: `go test -tags=integration -v ./persistence/integration_test`

### Performance
- Los tests pueden tomar varios minutos en la primera ejecución debido a la descarga de imágenes Docker
- SQL Server especialmente toma tiempo en inicializarse completamente