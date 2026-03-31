-- Phase 2: Canvas LMS integration tables

CREATE TABLE IF NOT EXISTS lms_connections (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  school_id UUID NOT NULL REFERENCES schools(id),
  lms_type TEXT NOT NULL DEFAULT 'canvas',
  institution_url TEXT NOT NULL,
  access_token TEXT NOT NULL,
  refresh_token TEXT NOT NULL,
  token_expires_at TIMESTAMPTZ NOT NULL,
  last_synced_at TIMESTAMPTZ,
  sync_status TEXT NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, lms_type)
);

CREATE INDEX IF NOT EXISTS idx_lms_connections_school_id ON lms_connections(school_id);
CREATE INDEX IF NOT EXISTS idx_lms_connections_user_id ON lms_connections(user_id);

CREATE TABLE IF NOT EXISTS lms_courses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  synapse_course_id UUID,
  user_id UUID NOT NULL REFERENCES users(id),
  school_id UUID NOT NULL,
  lms_course_id TEXT NOT NULL,
  lms_course_name TEXT NOT NULL,
  lms_term TEXT,
  enrollment_type TEXT NOT NULL DEFAULT 'student',
  last_synced_at TIMESTAMPTZ,
  UNIQUE (user_id, lms_course_id)
);

CREATE INDEX IF NOT EXISTS idx_lms_courses_school_id ON lms_courses(school_id);
CREATE INDEX IF NOT EXISTS idx_lms_courses_user_id ON lms_courses(user_id);

CREATE TABLE IF NOT EXISTS lms_assignments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id UUID NOT NULL,
  lms_assignment_id TEXT NOT NULL,
  lms_course_id TEXT NOT NULL,
  title TEXT NOT NULL,
  due_at TIMESTAMPTZ,
  points_possible NUMERIC,
  assignment_group TEXT,
  last_synced_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (lms_assignment_id, school_id)
);

CREATE INDEX IF NOT EXISTS idx_lms_assignments_school_id ON lms_assignments(school_id);
CREATE INDEX IF NOT EXISTS idx_lms_assignments_course_id ON lms_assignments(lms_course_id);

CREATE TABLE IF NOT EXISTS lms_grade_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id),
  school_id UUID NOT NULL,
  lms_assignment_id TEXT NOT NULL,
  lms_course_id TEXT NOT NULL,
  score NUMERIC,
  points_possible NUMERIC,
  submitted_at TIMESTAMPTZ,
  graded_at TIMESTAMPTZ,
  grade_type TEXT,
  synced_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, lms_assignment_id)
);

CREATE INDEX IF NOT EXISTS idx_lms_grade_events_school_id ON lms_grade_events(school_id);
CREATE INDEX IF NOT EXISTS idx_lms_grade_events_user_id ON lms_grade_events(user_id);

CREATE TABLE IF NOT EXISTS syllabus_parse_jobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id UUID NOT NULL,
  raw_file_url TEXT,
  raw_text TEXT,
  status TEXT NOT NULL DEFAULT 'pending',
  parsed_topics JSONB,
  parsed_exam_dates JSONB,
  confidence_score NUMERIC,
  human_reviewed BOOLEAN NOT NULL DEFAULT false,
  error_message TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_syllabus_parse_jobs_school_id ON syllabus_parse_jobs(school_id);
