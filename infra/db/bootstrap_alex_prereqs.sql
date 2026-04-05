-- Prereqs for populate_alex_data.sql: fixed MIT demo school, Alex user, and course rows.
-- Idempotent (ON CONFLICT DO NOTHING). Run after migrations, before populate_alex_data.sql.
--
-- Alex login: alex@mit.edu / demo123 (same bcrypt as demo cohort)

BEGIN;

INSERT INTO schools (id, name, domain)
VALUES ('00000000-0000-0000-0000-000000000001', 'MIT Demo', 'mit.edu')
ON CONFLICT (id) DO NOTHING;

INSERT INTO users (id, school_id, name, email, password_hash)
VALUES (
  '1ceeff29-f010-484d-9bbd-e51c4991289f',
  '00000000-0000-0000-0000-000000000001',
  'Alex Kim',
  'alex@mit.edu',
  '$2a$12$h9hY66qscZS8MHnmYRdDqOeUTPjObHEfABKyyn.Gx7q5Sapsl20gu'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO courses (id, school_id, user_id, name, term, color, lms_course_id) VALUES
  ('e501a47c-d907-4f82-9ffe-61a221212f55', '00000000-0000-0000-0000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', '6.790 Machine Learning', '2026 Spring', '#102E67', NULL),
  ('847e3a18-c7a4-4986-b991-b21439b7a9ff', '00000000-0000-0000-0000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', '6.006 Introduction to Algorithms', '2026 Spring', '#00A344', NULL),
  ('b45b4a4a-f123-479b-8cfe-ecfbc26546e0', '00000000-0000-0000-0000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'Software Construction', '2026 Spring', '#F59E0B', NULL),
  ('6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', '00000000-0000-0000-0000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'African History', '2026 Spring', '#FF5252', NULL),
  ('61f43c8a-154e-4af3-9e77-81a25cf30c1a', '00000000-0000-0000-0000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'Composition II', '2026 Spring', '#6B6B6B', NULL)
ON CONFLICT (id) DO NOTHING;

COMMIT;
