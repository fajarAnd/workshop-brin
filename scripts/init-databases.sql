-- Initialize databases for Workshop BRIN
-- This script runs once when PostgreSQL container is first created

-- Create wa_service user and database
CREATE USER wa_service WITH PASSWORD 'workshop2025';
CREATE DATABASE wa_service OWNER wa_service;
GRANT ALL PRIVILEGES ON DATABASE wa_service TO wa_service;

-- Create n8n user and database
CREATE USER n8n WITH PASSWORD 'workshop2025';
CREATE DATABASE n8n OWNER n8n;
GRANT ALL PRIVILEGES ON DATABASE n8n TO n8n;

-- Connect to wa_service database to install extensions
\c wa_service;
-- Install pgvector extension for vector similarity search
CREATE EXTENSION IF NOT EXISTS vector;
-- Enable UUID extension for generating UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create a comment on the database
COMMENT ON DATABASE wa_service IS 'Database for WhatsApp service with vector search capabilities';