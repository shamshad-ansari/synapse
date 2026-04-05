-- Planner dev seed data for the logged-in user.
-- Uses subqueries to dynamically find the user's UUID + school_id,
-- so this works no matter what the actual UUIDs are.
--
-- Idempotent: remove prior demo planner rows for Alex before re-inserting (make seed-demo safe).

DELETE FROM study_sessions s
USING users u
WHERE s.user_id = u.id
  AND u.email = 'alex@mit.edu'
  AND s.title IN (
    'Recursion Review',
    'Logic Review',
    'Induction Cards',
    'Induction Deep',
    'Set Theory',
    'Mixed Review',
    'Graph Theory Intro',
    'PS4 Prep',
    'Eigenvalue Drills',
    'Proof Practice',
    'BFS/DFS Review',
    'PS3 Solutions'
  );

DELETE FROM study_deadlines d
USING users u
WHERE d.user_id = u.id
  AND u.email = 'alex@mit.edu'
  AND d.name IN ('Midterm Study Goal', 'Complete practice exams');

-- Step 1: Insert study sessions for the current week (and last week)
-- Uses CURRENT_DATE-relative dates.

-- Yesterday & before: done or missed
INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'Recursion Review', CURRENT_DATE - INTERVAL '3 days', '09:00', 45, 'done'
FROM users u WHERE u.email = 'alex@mit.edu';

INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'Logic Review', CURRENT_DATE - INTERVAL '3 days', '14:00', 30, 'done'
FROM users u WHERE u.email = 'alex@mit.edu';

INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'Induction Cards', CURRENT_DATE - INTERVAL '2 days', '09:00', 30, 'done'
FROM users u WHERE u.email = 'alex@mit.edu';

INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'Induction Deep', CURRENT_DATE - INTERVAL '1 day', '14:00', 45, 'missed'
FROM users u WHERE u.email = 'alex@mit.edu';

INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'Set Theory', CURRENT_DATE - INTERVAL '1 day', '09:00', 45, 'done'
FROM users u WHERE u.email = 'alex@mit.edu';

-- Today and upcoming: planned
INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'Mixed Review', CURRENT_DATE, '09:00', 60, 'planned'
FROM users u WHERE u.email = 'alex@mit.edu';

INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'Graph Theory Intro', CURRENT_DATE, '14:00', 45, 'planned'
FROM users u WHERE u.email = 'alex@mit.edu';

INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'PS4 Prep', CURRENT_DATE + INTERVAL '1 day', '09:00', 45, 'planned'
FROM users u WHERE u.email = 'alex@mit.edu';

INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'Eigenvalue Drills', CURRENT_DATE + INTERVAL '1 day', '14:00', 30, 'planned'
FROM users u WHERE u.email = 'alex@mit.edu';

INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'Proof Practice', CURRENT_DATE + INTERVAL '2 days', '09:00', 45, 'planned'
FROM users u WHERE u.email = 'alex@mit.edu';

INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'BFS/DFS Review', CURRENT_DATE + INTERVAL '3 days', '09:00', 60, 'planned'
FROM users u WHERE u.email = 'alex@mit.edu';

INSERT INTO study_sessions (user_id, school_id, title, scheduled_date, start_time, duration_minutes, status)
SELECT u.id, u.school_id, 'PS3 Solutions', CURRENT_DATE + INTERVAL '3 days', '14:00', 90, 'planned'
FROM users u WHERE u.email = 'alex@mit.edu';

-- Step 2: Custom study deadlines (manual ones, not from LMS)
INSERT INTO study_deadlines (user_id, school_id, name, course_name, due_date, source)
SELECT u.id, u.school_id, 'Midterm Study Goal', 'CS225', CURRENT_DATE + INTERVAL '14 days', 'manual'
FROM users u WHERE u.email = 'alex@mit.edu';

INSERT INTO study_deadlines (user_id, school_id, name, course_name, due_date, source)
SELECT u.id, u.school_id, 'Complete practice exams', '18.06', CURRENT_DATE + INTERVAL '7 days', 'manual'
FROM users u WHERE u.email = 'alex@mit.edu';