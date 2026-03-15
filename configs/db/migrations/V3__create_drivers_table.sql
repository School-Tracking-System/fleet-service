-- Fleet Service: Tabla de conductores
-- Depende de: users (auth service) y uuid-ossp extension
-- Nota: user_id referencia la tabla users del servicio Auth (cross-service FK lógica, no física).

CREATE TABLE IF NOT EXISTS drivers (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id         UUID        NOT NULL UNIQUE,   -- FK lógica → users(id) en Auth Service
    license_number  VARCHAR(50) NOT NULL UNIQUE,
    license_type    VARCHAR(20) NOT NULL,           -- e.g. 'B', 'C', 'D'
    license_expiry  DATE        NOT NULL,
    cedula_id       VARCHAR(20) NOT NULL UNIQUE,    -- Documento de identidad
    emergency_phone VARCHAR(20),
    status          VARCHAR(20) NOT NULL DEFAULT 'active', -- 'active' | 'suspended' | 'inactive'
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_drivers_user_id        ON drivers (user_id);
CREATE INDEX idx_drivers_license_number ON drivers (license_number);
CREATE INDEX idx_drivers_status         ON drivers (status);

CREATE TRIGGER trigger_drivers_updated_at
    BEFORE UPDATE ON drivers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
