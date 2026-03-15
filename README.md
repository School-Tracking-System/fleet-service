# Fleet Service

Microservicio encargado de la gestión y ciclo de vida de los **datos maestros** del sistema de tracking escolar: Vehículos, Conductores, Escuelas y Estudiantes.

---

## ¿Por qué gRPC y no HTTP/REST?

Este servicio expone su API a través de **gRPC** en lugar de HTTP/REST. Esta es una decisión arquitectónica deliberada basada en las necesidades del sistema de tracking:

| Característica | HTTP/REST (JSON) | gRPC (Protobuf) |
| --- | --- | --- |
| **Serialización** | JSON (texto, verbose) | Protocol Buffers (binario, ~5x más compacto) |
| **Contrato de API** | Swagger/OpenAPI (opcional) | Archivo `.proto` (obligatorio y versionado) |
| **Streaming** | Polling o SSE | Streaming nativo bidireccional (fundamental para GPS) |
| **Rendimiento** | Más lento bajo alta carga | ~8x más rápido en comunicación interna |
| **Type Safety** | Validación manual o por librerías | Generado por el compilador, 100% type-safe |

> **En nuestro contexto:** Cuando cientos de vehículos envíen coordenadas GPS cada segundo, el overhead del JSON en la comunicación interna sería inaceptable. gRPC asegura que la comunicación entre el **API Gateway → Fleet Service** sea eficiente, de bajo latencia y con contratos bien definidos.

> **El API Gateway** es el único que habla HTTP con el mundo exterior. Por detrás, traduce las peticiones REST del cliente a llamadas gRPC a los microservicios internos.

```
Cliente (App/Web)
      │
      │  HTTP/REST (JSON)
      ▼
┌─────────────┐
│ API Gateway │  Puerto 8000
│  Chi Router │
└──────┬──────┘
       │  gRPC (Protobuf)
       ▼
┌─────────────┐
│Fleet Service│  Puerto 9090
│ gRPC Server │
└──────┬──────┘
       │  SQL
       ▼
┌─────────────┐
│  PostgreSQL │
└─────────────┘
```

---

## Arquitectura Interna (Hexagonal)

```
cmd/
└── api/
    ├── main.go          # Punto de entrada. Arranca fx.App
    └── module.go        # Inyección de dependencias (Uber fx)

internal/
├── core/
│   ├── domain/
│   │   └── models.go   # Entidades puras: Vehicle, Driver, Student, School
│   ├── ports/
│   │   ├── repositories/ # Interfaces de persistencia (VehicleRepository)
│   │   └── services/     # Interfaces de lógica de negocio (VehicleService)
│   └── fleet/
│       └── service.go   # Implementación de la lógica de negocio
└── infrastructure/
    ├── grpc/
    │   ├── server.go    # Arranque del servidor gRPC
    │   └── handlers/
    │       └── vehicle_handler.go # Mapeo: Proto Request → Domain → Proto Response
    └── persistence/
        └── postgres/
            └── vehicle_repo.go # Implementación SQL del VehicleRepository

pkg/
├── api/
│   └── v1/             # Código Go GENERADO por buf desde los archivos .proto
│       ├── vehicle.pb.go
│       └── vehicle_grpc.pb.go
└── env/
    └── env.go          # Gestión de variables de entorno

proto/                  # (raíz del monorepo)
└── fleet/
    └── v1/
        └── vehicle.proto # Contrato gRPC — Fuente de verdad
```

---

## Requisitos Previos

- **Go** `>= 1.23`
- **buf** CLI (para compilar archivos `.proto`)
- **PostgreSQL** `>= 14`
- **grpcurl** *(opcional, para pruebas manuales)*

```bash
# Instalar buf en macOS
brew install bufbuild/buf/buf

# Instalar grpcurl (herramienta para llamar gRPC desde terminal)
brew install grpcurl
```

---

## Variables de Entorno

Copiar la plantilla y ajustar los valores:

```bash
cp .env.template .env
```

| Variable | Descripción | Default |
| --- | --- | --- |
| `SERVICE_NAME` | Nombre del servicio (usado en logs) | `fleet` |
| `HTTP_PORT` | Puerto HTTP (healthcheck futuro) | `8081` |
| `GRPC_PORT` | Puerto principal del servidor gRPC | `9090` |
| `DATABASE_URL` | Cadena de conexión a PostgreSQL | `postgres://postgres:postgres@localhost:5432/school_tracking?sslmode=disable` |
| `ENVIRONMENT` | Entorno de ejecución (`development` / `production`) | `development` |
| `LOG_LEVEL` | Nivel de logging (`debug`, `info`, `warn`, `error`) | `debug` |

---

## Iniciar en Local

```bash
# 1. Clonar e instalar dependencias
go mod download

# 2. Copiar variables de entorno
cp .env.template .env

# 3. Asegurarse de tener PostgreSQL corriendo (el servicio auto-crea el schema)

# 4. Ejecutar el servicio
go run ./cmd/api/...
```

Al iniciar, verás en los logs:
```
{"level":"info","msg":"Starting gRPC server","port":"9090"}
{"level":"info","msg":"Successfully connected to PostgreSQL database"}
```

---

## Compilar los Protobuf (Desde la raíz del monorepo)

El código en `pkg/api/v1/` es **generado automáticamente** y no debe editarse manualmente. Para regenerarlo tras modificar `vehicle.proto`:

```bash
# Desde la raíz del monorepo
cd ../../
buf generate
```

Los archivos generados (`*.pb.go`, `*_grpc.pb.go`) están ubicados en `pkg/api/v1/`.

---

## Probar el Servicio con grpcurl

Con el servidor corriendo en `localhost:9090`:

### Listar servicios disponibles (requiere reflection habilitada)
```bash
grpcurl -plaintext localhost:9090 list
```

### Crear un Vehículo
```bash
grpcurl -plaintext -d '{
  "plate": "ABC-1234",
  "brand": "Mercedes-Benz",
  "model": "Sprinter",
  "year": 2023,
  "capacity": 20
}' localhost:9090 fleet.v1.VehicleService/CreateVehicle
```

### Obtener un Vehículo por ID
```bash
grpcurl -plaintext -d '{"id": "<uuid-del-vehiculo>"}' \
  localhost:9090 fleet.v1.VehicleService/GetVehicle
```

### Listar Vehículos (con paginación)
```bash
grpcurl -plaintext -d '{"limit": 10, "offset": 0}' \
  localhost:9090 fleet.v1.VehicleService/ListVehicles
```

---

## Modelo de Dominio: Vehicle

```go
type Vehicle struct {
    ID        uuid.UUID
    Plate     string        // Placa del vehículo (en mayúsculas, campo requerido)
    Brand     string        // Marca (ej. "Mercedes-Benz")
    Model     string        // Modelo (ej. "Sprinter")
    Year      int           // Año de fabricación
    Capacity  int           // Número de asientos para estudiantes
    Status    VehicleStatus // "active" | "maintenance" | "inactive"
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

## Tecnologías Utilizadas

| Librería | Propósito |
| --- | --- |
| `google.golang.org/grpc` | Servidor y cliente gRPC |
| `google.golang.org/protobuf` | Serialización Protocol Buffers |
| `go.uber.org/fx` | Inyección de dependencias y ciclo de vida |
| `go.uber.org/zap` | Logging estructurado |
| `github.com/lib/pq` | Driver PostgreSQL para `database/sql` |
| `github.com/google/uuid` | Generación de identificadores únicos |
| `github.com/joho/godotenv` | Carga de variables de entorno desde `.env` |
| `github.com/caarlos0/env/v10` | Mapeo de env vars a structs de configuración |

---

## Contribuir / Extender el Servicio

Para agregar una nueva entidad al Fleet Service (ej. `Driver`), seguir el mismo patrón:

1. **Agregar el mensaje y servicio** en `proto/fleet/v1/driver.proto`
2. **Regenerar código** con `buf generate` desde la raíz
3. **Agregar el modelo de dominio** en `internal/core/domain/models.go`
4. **Agregar la interfaz** de repositorio en `internal/core/ports/repositories/`
5. **Implementar el repositorio** en `internal/infrastructure/persistence/postgres/`
6. **Crear el handler gRPC** en `internal/infrastructure/grpc/handlers/`
7. **Registrar** el nuevo handler en `internal/infrastructure/grpc/server.go`
