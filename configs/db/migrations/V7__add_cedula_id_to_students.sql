-- V7: Add cedula_id to students table
ALTER TABLE students ADD COLUMN cedula_id VARCHAR(50);

-- Update existing students if any (though project is in development)
UPDATE students SET cedula_id = id::text WHERE cedula_id IS NULL;

-- Make it NOT NULL and UNIQUE
ALTER TABLE students ALTER COLUMN cedula_id SET NOT NULL;
ALTER TABLE students ADD CONSTRAINT students_cedula_id_unique UNIQUE (cedula_id);
