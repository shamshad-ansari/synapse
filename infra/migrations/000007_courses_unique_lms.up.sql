-- Prevent duplicate Synapse courses from the same LMS course.
-- First, remove any existing duplicates (keep the oldest row per user+lms_course_id).
DELETE FROM courses
WHERE id IN (
  SELECT id FROM (
    SELECT id,
           ROW_NUMBER() OVER (PARTITION BY user_id, lms_course_id ORDER BY created_at ASC) AS rn
    FROM courses
    WHERE lms_course_id IS NOT NULL
  ) numbered
  WHERE rn > 1
);

-- Now add a partial unique index (only for non-null lms_course_id).
CREATE UNIQUE INDEX IF NOT EXISTS idx_courses_user_lms_unique
  ON courses (user_id, lms_course_id)
  WHERE lms_course_id IS NOT NULL;
