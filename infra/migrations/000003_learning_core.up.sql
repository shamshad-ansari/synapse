-- Learning core: courses, topics, flashcards, SM-2 scheduler, reviews, mastery

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS courses (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id   UUID NOT NULL REFERENCES schools(id),
  user_id     UUID NOT NULL REFERENCES users(id),
  name        TEXT NOT NULL,
  term        TEXT,
  color       TEXT,
  lms_course_id TEXT,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_courses_user ON courses(user_id, school_id);

CREATE TABLE IF NOT EXISTS topics (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id       UUID NOT NULL,
  course_id       UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  user_id         UUID NOT NULL,
  name            TEXT NOT NULL,
  parent_topic_id UUID REFERENCES topics(id),
  exam_weight     NUMERIC,
  source          TEXT NOT NULL DEFAULT 'manual',
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_topics_course ON topics(course_id, school_id);

CREATE TABLE IF NOT EXISTS flashcards (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id        UUID NOT NULL,
  course_id        UUID NOT NULL REFERENCES courses(id),
  user_id          UUID NOT NULL REFERENCES users(id),
  topic_id         UUID REFERENCES topics(id),
  card_type        TEXT NOT NULL DEFAULT 'qa',
  prompt           TEXT NOT NULL,
  answer           TEXT NOT NULL,
  created_by       TEXT NOT NULL DEFAULT 'user',
  visibility       TEXT NOT NULL DEFAULT 'private',
  embedding        VECTOR(1536),
  created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_flashcards_course ON flashcards(course_id, user_id, school_id);
CREATE INDEX IF NOT EXISTS idx_flashcards_topic ON flashcards(topic_id, school_id);

CREATE TABLE IF NOT EXISTS scheduler_states (
  flashcard_id   UUID NOT NULL REFERENCES flashcards(id) ON DELETE CASCADE,
  user_id        UUID NOT NULL,
  school_id      UUID NOT NULL,
  ease_factor    NUMERIC NOT NULL DEFAULT 2.5,
  interval_days  INT NOT NULL DEFAULT 1,
  due_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  lapse_count    INT NOT NULL DEFAULT 0,
  last_review_at TIMESTAMPTZ,
  PRIMARY KEY (flashcard_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_scheduler_due ON scheduler_states(user_id, due_at);

CREATE TABLE IF NOT EXISTS review_events (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id       UUID NOT NULL,
  user_id         UUID NOT NULL REFERENCES users(id),
  flashcard_id    UUID NOT NULL REFERENCES flashcards(id),
  session_id      UUID NOT NULL,
  ts              TIMESTAMPTZ NOT NULL DEFAULT now(),
  correct         BOOLEAN NOT NULL,
  confidence      SMALLINT NOT NULL CHECK (confidence BETWEEN 1 AND 4),
  confused        BOOLEAN NOT NULL DEFAULT false,
  response_time_ms INT NOT NULL,
  ease_before     NUMERIC,
  interval_before INT
);
CREATE INDEX IF NOT EXISTS idx_review_events_user ON review_events(user_id, ts DESC);

CREATE TABLE IF NOT EXISTS topic_mastery (
  school_id      UUID NOT NULL,
  user_id        UUID NOT NULL REFERENCES users(id),
  topic_id       UUID NOT NULL REFERENCES topics(id),
  mastery_score  NUMERIC NOT NULL DEFAULT 0 CHECK (mastery_score BETWEEN 0 AND 1),
  confusion_rate NUMERIC NOT NULL DEFAULT 0,
  review_count   INT NOT NULL DEFAULT 0,
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (school_id, user_id, topic_id)
);
