DROP INDEX IF EXISTS idx_lms_sync_runs_user_school;
DROP TABLE IF EXISTS lms_sync_runs;

DROP INDEX IF EXISTS idx_lms_submission_states_course;
DROP INDEX IF EXISTS idx_lms_submission_states_user_school;
DROP TABLE IF EXISTS lms_submission_states;

DROP INDEX IF EXISTS idx_lms_discussion_topics_course;
DROP INDEX IF EXISTS idx_lms_discussion_topics_user_school;
DROP TABLE IF EXISTS lms_discussion_topics;

DROP INDEX IF EXISTS idx_lms_announcements_course;
DROP INDEX IF EXISTS idx_lms_announcements_user_school;
DROP TABLE IF EXISTS lms_announcements;
