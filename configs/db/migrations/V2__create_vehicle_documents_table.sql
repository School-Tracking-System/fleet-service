-- Fleet Service: Tabla de documentos del vehículo
-- Depende de: V1__create_vehicles_table.sql

CREATE TABLE IF NOT EXISTS vehicle_documents (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vehicle_id  UUID        NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    doc_type    VARCHAR(50) NOT NULL,   -- e.g. 'SOAT', 'technical_review', 'license'
    doc_url     TEXT        NOT NULL,
    expires_at  DATE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_vehicle_documents_vehicle_id ON vehicle_documents (vehicle_id);
CREATE INDEX idx_vehicle_documents_doc_type   ON vehicle_documents (doc_type);
