-- Create workflow_config table for global workflow routing configuration
CREATE TABLE workflow_config (
    id SERIAL PRIMARY KEY,
    workflow_type VARCHAR(20) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default configuration (N8N for backward compatibility)
INSERT INTO workflow_config (id, workflow_type) VALUES (1, 'n8n');