CREATE TABLE IF NOT EXISTS tutor_requests (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id    UUID NOT NULL,
  requester_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  tutor_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  topic_id     UUID REFERENCES topics(id) ON DELETE SET NULL,
  topic_name   TEXT NOT NULL DEFAULT '',
  status       TEXT NOT NULL DEFAULT 'pending', -- pending|accepted|declined|completed|cancelled
  message      TEXT NOT NULL DEFAULT '',
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tutor_requests_tutor ON tutor_requests(tutor_id, school_id);
CREATE INDEX IF NOT EXISTS idx_tutor_requests_requester ON tutor_requests(requester_id, school_id);
