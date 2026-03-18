-- Fleet Service: Tablas para Rutas y Paradas
-- Depende de: vehicles, drivers, schools, students

-- -------------------------------------------------------
-- Tipo ENUM para la dirección de la ruta
-- -------------------------------------------------------
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'route_direction') THEN
        CREATE TYPE route_direction AS ENUM ('to_school', 'from_school');
    END IF;
END;
$$;

CREATE TABLE IF NOT EXISTS routes (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    vehicle_id      UUID REFERENCES vehicles(id) ON DELETE SET NULL,
    driver_id       UUID REFERENCES drivers(id) ON DELETE SET NULL,
    school_id       UUID NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    direction       route_direction NOT NULL DEFAULT 'to_school',
    schedule_time   TIME NOT NULL,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_routes_school_id ON routes (school_id);
CREATE INDEX idx_routes_vehicle_id ON routes (vehicle_id);
CREATE INDEX idx_routes_driver_id ON routes (driver_id);

CREATE TRIGGER trigger_routes_updated_at
    BEFORE UPDATE ON routes
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- -------------------------------------------------------

CREATE TABLE IF NOT EXISTS route_stops (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    route_id        UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    student_id      UUID NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    stop_order      INT NOT NULL,
    location        GEOMETRY(Point, 4326) NOT NULL,
    address         VARCHAR(300),
    est_time        TIME,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_route_stops_route_id ON route_stops (route_id);
CREATE INDEX idx_route_stops_student_id ON route_stops (student_id);
CREATE INDEX idx_route_stops_location ON route_stops USING GIST (location);
