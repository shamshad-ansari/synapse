-- Enable extensions needed for Synapse.
-- Run once on DB init.
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "vector";