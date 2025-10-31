-- Initialize database for AI-BPMS
-- This script creates initial database setup

-- Create Keycloak database (separate from main BPMS)
CREATE DATABASE keycloak;

-- Grant permissions
GRANT ALL PRIVILEGES ON DATABASE ai_bpms TO bpms_user;
GRANT ALL PRIVILEGES ON DATABASE keycloak TO bpms_user;

-- Create extensions for performance optimization
\c ai_bpms;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Create initial schema structure
CREATE SCHEMA IF NOT EXISTS process_engine;
CREATE SCHEMA IF NOT EXISTS user_management;
CREATE SCHEMA IF NOT EXISTS analytics;
CREATE SCHEMA IF NOT EXISTS audit;

-- Grant schema permissions
GRANT ALL ON SCHEMA process_engine TO bpms_user;
GRANT ALL ON SCHEMA user_management TO bpms_user;
GRANT ALL ON SCHEMA analytics TO bpms_user;
GRANT ALL ON SCHEMA audit TO bpms_user;