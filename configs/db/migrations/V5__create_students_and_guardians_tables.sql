-- Fleet Service: Tablas de estudiantes y representantes (guardians)
-- Depende de: PostGIS, uuid-ossp, V4__create_schools_and_contacts_tables.sql

-- Activar la extensión PostGIS si no está activa
CREATE EXTENSION IF NOT EXISTS postgis;


CREATE TABLE IF NOT EXISTS students (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name       VARCHAR(100)     NOT NULL,
    last_name        VARCHAR(100)     NOT NULL,
    grade            VARCHAR(20),
    school_id        UUID             NOT NULL REFERENCES schools(id) ON DELETE RESTRICT,
    pickup_location  GEOMETRY(Point, 4326),  -- PostGIS: Punto de recogida del estudiante
    pickup_address   VARCHAR(300),
    photo_url        TEXT,
    is_active        BOOLEAN          NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_students_school_id       ON students (school_id);
CREATE INDEX idx_students_is_active       ON students (is_active);
CREATE INDEX idx_students_pickup_location ON students USING GIST (pickup_location);

CREATE TRIGGER trigger_students_updated_at
    BEFORE UPDATE ON students
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- -------------------------------------------------------
-- Tipo ENUM para la relación del representante
-- -------------------------------------------------------
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'guardian_relation') THEN
        CREATE TYPE guardian_relation AS ENUM ('father', 'mother', 'legal_guardian', 'other');
    END IF;
END;
$$;

CREATE TABLE IF NOT EXISTS guardians (
    id          UUID              PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id     UUID              NOT NULL,     -- FK lógica → users(id) en Auth Service
    student_id  UUID              NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    relation    guardian_relation NOT NULL DEFAULT 'other',
    is_primary  BOOLEAN           NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ       NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_guardians_user_student ON guardians (user_id, student_id);
CREATE INDEX idx_guardians_student_id          ON guardians (student_id);
CREATE INDEX idx_guardians_user_id             ON guardians (user_id);
