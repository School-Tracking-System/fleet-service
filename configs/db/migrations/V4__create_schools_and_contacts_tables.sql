-- Fleet Service: Tablas de escuelas y contactos de escuela
-- Depende de: PostGIS extension y uuid-ossp

-- Activar la extensión PostGIS si no está activa
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS schools (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(200) NOT NULL,
    address     VARCHAR(300) NOT NULL,
    location    GEOMETRY(Point, 4326),    -- PostGIS: Coordenadas GPS de la escuela
    phone       VARCHAR(20),
    email       VARCHAR(255),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_schools_name     ON schools (name);
CREATE INDEX idx_schools_location ON schools USING GIST (location);

CREATE TRIGGER trigger_schools_updated_at
    BEFORE UPDATE ON schools
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- -------------------------------------------------------

CREATE TABLE IF NOT EXISTS school_contacts (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    school_id   UUID        NOT NULL REFERENCES schools(id) ON DELETE CASCADE,
    user_id     UUID        NOT NULL,    -- FK lógica → users(id) en Auth Service (staff de escuela)
    position    VARCHAR(100),            -- e.g. 'Coordinator', 'Principal'
    is_active   BOOLEAN     NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_school_contacts_school_id ON school_contacts (school_id);
CREATE INDEX idx_school_contacts_user_id   ON school_contacts (user_id);
