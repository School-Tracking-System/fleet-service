-- Fleet Service: Tabla de vehículos
-- Depende de: uuid-ossp extension y vehicle_status ENUM (creados en init-db.sql)

CREATE TABLE IF NOT EXISTS vehicles (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plate            VARCHAR(20)  NOT NULL UNIQUE,
    brand            VARCHAR(50)  NOT NULL,
    model            VARCHAR(50)  NOT NULL,
    year             INT          NOT NULL,
    color            VARCHAR(30),
    vehicle_type     VARCHAR(30),                   -- e.g. 'van', 'bus', 'minibus'
    capacity         INT          NOT NULL,
    chassis_num      VARCHAR(100) UNIQUE,
    status           VARCHAR(20)  NOT NULL DEFAULT 'active', -- 'active' | 'maintenance' | 'inactive'
    insurance_exp    DATE,
    tech_review_exp  DATE,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Índices para búsquedas frecuentes
CREATE INDEX idx_vehicles_plate  ON vehicles (plate);
CREATE INDEX idx_vehicles_status ON vehicles (status);

-- Trigger para actualizar updated_at automáticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_vehicles_updated_at
    BEFORE UPDATE ON vehicles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
