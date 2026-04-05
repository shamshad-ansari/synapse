-- Planner: study sessions (calendar grid blocks) and user-created deadlines

CREATE TABLE IF NOT EXISTS study_sessions (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  school_id         UUID NOT NULL REFERENCES schools(id),
  title             TEXT NOT NULL,
  course_id         UUID REFERENCES courses(id) ON DELETE SET NULL,
  topic_id          UUID REFERENCES topics(id) ON DELETE SET NULL,
  scheduled_date    DATE NOT NULL,
  start_time        TIME NOT NULL DEFAULT '09:00',
  duration_minutes  INT NOT NULL DEFAULT 30,
  status            TEXT NOT NULL DEFAULT 'planned'
                    CHECK (status IN ('planned', 'done', 'missed')),
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_study_sessions_user_date
  ON study_sessions(user_id, school_id, scheduled_date);

CREATE TABLE IF NOT EXISTS study_deadlines (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  school_id         UUID NOT NULL REFERENCES schools(id),
  name              TEXT NOT NULL,
  course_name       TEXT,
  due_date          DATE NOT NULL,
  source            TEXT NOT NULL DEFAULT 'manual'
                    CHECK (source IN ('manual', 'lms')),
  lms_assignment_id TEXT,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_study_deadlines_user_date
  ON study_deadlines(user_id, school_id, due_date);
