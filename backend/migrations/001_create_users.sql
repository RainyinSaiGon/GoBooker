-- ============================================================
-- Migration: 001_create_users.sql
-- Description: Create the users table for GoBooker
-- ============================================================

CREATE TABLE IF NOT EXISTS users (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email       TEXT        NOT NULL UNIQUE,
    name        TEXT        NOT NULL,
    password    TEXT        NOT NULL,
    role        TEXT        NOT NULL DEFAULT 'customer',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Index for fast lookup by email (e.g. login)
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

