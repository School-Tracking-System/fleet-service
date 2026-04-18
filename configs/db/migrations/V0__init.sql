-- Fleet DB — Extensiones y tipos base
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE vehicle_status  AS ENUM ('active', 'inactive', 'maintenance');
CREATE TYPE driver_status   AS ENUM ('active', 'inactive', 'suspended');
CREATE TYPE guardian_relation AS ENUM ('father', 'mother', 'other');
CREATE TYPE trip_direction  AS ENUM ('to_school', 'from_school');
