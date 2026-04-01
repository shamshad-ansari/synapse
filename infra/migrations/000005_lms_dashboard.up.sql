-- LMS dashboard support: feed items, submission states, and sync run telemetry

CREATE TABLE IF NOT EXISTS lms_announcements (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  school_id UUID NOT NULL REFERENCES schools(id),
  lms_course_id TEXT NOT NULL,
  lms_announcement_id TEXT NOT NULL,
  title TEXT NOT NULL,
  message TEXT,
  posted_at TIMESTAMPTZ,
  html_url TEXT,
  last_synced_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, lms_announcement_id)
);

CREATE INDEX IF NOT EXISTS idx_lms_announcements_user_school
  ON lms_announcements(user_id, school_id, posted_at DESC);
CREATE INDEX IF NOT EXISTS idx_lms_announcements_course
  ON lms_announcements(school_id, lms_course_id);

CREATE TABLE IF NOT EXISTS lms_discussion_topics (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  school_id UUID NOT NULL REFERENCES schools(id),
  lms_course_id TEXT NOT NULL,
  lms_topic_id TEXT NOT NULL,
  title TEXT NOT NULL,
  message TEXT,
  posted_at TIMESTAMPTZ,
  html_url TEXT,
  last_synced_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, lms_topic_id)
);

CREATE INDEX IF NOT EXISTS idx_lms_discussion_topics_user_school
  ON lms_discussion_topics(user_id, school_id, posted_at DESC);
CREATE INDEX IF NOT EXISTS idx_lms_discussion_topics_course
  ON lms_discussion_topics(school_id, lms_course_id);

CREATE TABLE IF NOT EXISTS lms_submission_states (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  school_id UUID NOT NULL REFERENCES schools(id),
  lms_assignment_id TEXT NOT NULL,
  lms_course_id TEXT NOT NULL,
  workflow_state TEXT,
  missing BOOLEAN NOT NULL DEFAULT false,
  late BOOLEAN NOT NULL DEFAULT false,
  excused BOOLEAN NOT NULL DEFAULT false,
  submitted_at TIMESTAMPTZ,
  graded_at TIMESTAMPTZ,
  score NUMERIC,
  points_possible NUMERIC,
  synced_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, lms_assignment_id)
);

CREATE INDEX IF NOT EXISTS idx_lms_submission_states_user_school
  ON lms_submission_states(user_id, school_id, synced_at DESC);
CREATE INDEX IF NOT EXISTS idx_lms_submission_states_course
  ON lms_submission_states(school_id, lms_course_id);

CREATE TABLE IF NOT EXISTS lms_sync_runs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  school_id UUID NOT NULL REFERENCES schools(id),
  status TEXT NOT NULL,
  error_message TEXT,
  courses_synced INT NOT NULL DEFAULT 0,
  started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  finished_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_lms_sync_runs_user_school
  ON lms_sync_runs(user_id, school_id, started_at DESC);
