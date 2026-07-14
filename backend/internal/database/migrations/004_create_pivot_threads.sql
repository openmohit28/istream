CREATE TABLE IF NOT EXISTS pivot_threads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    steps JSONB NOT NULL DEFAULT '[]',
    forked_from UUID REFERENCES pivot_threads(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_pivot_threads_user_updated
    ON pivot_threads (user_id, updated_at DESC);