-- Cleanup old events to ensure a clean baseline for metrics
DELETE FROM review_events WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f';
DELETE FROM lms_assignments WHERE school_id = '00000000-0000-0000-0000-000000000001';

-- Alex review history over the last 5 weeks (evening sessions, fixed session IDs).

-- Transformer Architectures (~48 reviews, ~88% correct, ~4% confused)
WITH transformer_cards AS (
  SELECT id, row_number() OVER (ORDER BY id) AS rn
  FROM flashcards
  WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f'
    AND topic_id = 'a1111111-f010-484d-9bbd-e51c4991289f'
), transformer_sessions AS (
  SELECT
    sp.week_offset,
    sp.session_no,
    sp.cards_per_session,
    (
      substr(md5('alex-transformers-' || sp.week_offset || '-' || sp.session_no), 1, 8) || '-' ||
      substr(md5('alex-transformers-' || sp.week_offset || '-' || sp.session_no), 9, 4) || '-' ||
      substr(md5('alex-transformers-' || sp.week_offset || '-' || sp.session_no), 13, 4) || '-' ||
      substr(md5('alex-transformers-' || sp.week_offset || '-' || sp.session_no), 17, 4) || '-' ||
      substr(md5('alex-transformers-' || sp.week_offset || '-' || sp.session_no), 21, 12)
    )::uuid AS session_id,
    CURRENT_DATE
      - (sp.week_offset || ' weeks')::interval
      + (sp.session_hour || ' hours')::interval
      + (((sp.week_offset * 13 + sp.session_no * 17) % 40) || ' minutes')::interval AS session_ts
  FROM (
    VALUES
      (4, 1, 6, 19), (4, 2, 6, 21),
      (3, 1, 6, 19), (3, 2, 6, 22),
      (2, 1, 6, 20), (2, 2, 6, 21),
      (1, 1, 6, 19),
      (0, 1, 6, 20)
  ) AS sp(week_offset, session_no, cards_per_session, session_hour)
  JOIN generate_series(0, 4) AS week_idx(week_offset) ON week_idx.week_offset = sp.week_offset
), transformer_events AS (
  SELECT
    ts.session_id,
    ts.session_ts,
    tc.id AS flashcard_id,
    row_number() OVER (ORDER BY ts.session_ts, ts.session_id, tc.id) AS seq
  FROM transformer_sessions ts
  JOIN transformer_cards tc ON tc.rn <= ts.cards_per_session
)
INSERT INTO review_events
  (user_id, school_id, flashcard_id, correct, confidence, confused, response_time_ms, ts, session_id)
SELECT
  '1ceeff29-f010-484d-9bbd-e51c4991289f',
  '00000000-0000-0000-0000-000000000001',
  te.flashcard_id,
  te.seq % 8 <> 0,
  CASE
    WHEN te.seq % 8 = 0 THEN 2
    WHEN te.seq % 3 = 0 THEN 4
    ELSE 3
  END,
  te.seq IN (11, 37),
  1450 + ((te.seq * 173) % 1850),
  te.session_ts,
  te.session_id
FROM transformer_events te;

-- B-Trees & Indexing (~40 reviews, ~82% correct, ~5% confused)
WITH btree_cards AS (
  SELECT id, row_number() OVER (ORDER BY id) AS rn
  FROM flashcards
  WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f'
    AND topic_id = 'a2222222-f010-484d-9bbd-e51c4991289f'
), btree_sessions AS (
  SELECT
    sp.week_offset,
    sp.session_no,
    sp.cards_per_session,
    (
      substr(md5('alex-btrees-' || sp.week_offset || '-' || sp.session_no), 1, 8) || '-' ||
      substr(md5('alex-btrees-' || sp.week_offset || '-' || sp.session_no), 9, 4) || '-' ||
      substr(md5('alex-btrees-' || sp.week_offset || '-' || sp.session_no), 13, 4) || '-' ||
      substr(md5('alex-btrees-' || sp.week_offset || '-' || sp.session_no), 17, 4) || '-' ||
      substr(md5('alex-btrees-' || sp.week_offset || '-' || sp.session_no), 21, 12)
    )::uuid AS session_id,
    CURRENT_DATE
      - (sp.week_offset || ' weeks')::interval
      + (sp.session_hour || ' hours')::interval
      + (((sp.week_offset * 11 + sp.session_no * 19) % 42) || ' minutes')::interval AS session_ts
  FROM (
    VALUES
      (4, 1, 5, 19), (4, 2, 5, 21),
      (3, 1, 5, 20), (3, 2, 5, 22),
      (2, 1, 5, 19), (2, 2, 5, 21),
      (1, 1, 5, 20),
      (0, 1, 5, 19)
  ) AS sp(week_offset, session_no, cards_per_session, session_hour)
  JOIN generate_series(0, 4) AS week_idx(week_offset) ON week_idx.week_offset = sp.week_offset
), btree_events AS (
  SELECT
    bs.session_id,
    bs.session_ts,
    bc.id AS flashcard_id,
    row_number() OVER (ORDER BY bs.session_ts, bs.session_id, bc.id) AS seq
  FROM btree_sessions bs
  JOIN btree_cards bc ON bc.rn <= bs.cards_per_session
)
INSERT INTO review_events
  (user_id, school_id, flashcard_id, correct, confidence, confused, response_time_ms, ts, session_id)
SELECT
  '1ceeff29-f010-484d-9bbd-e51c4991289f',
  '00000000-0000-0000-0000-000000000001',
  be.flashcard_id,
  be.seq NOT IN (6, 12, 18, 24, 30, 35, 40),
  CASE
    WHEN be.seq IN (6, 12, 18, 24, 30, 35, 40) THEN 2
    WHEN be.seq % 4 = 0 THEN 4
    ELSE 3
  END,
  be.seq IN (9, 28),
  1650 + ((be.seq * 149) % 2100),
  be.session_ts,
  be.session_id
FROM btree_events be;

-- Design Patterns (~32 reviews, ~87% correct, ~3% confused)
WITH design_cards AS (
  SELECT id, row_number() OVER (ORDER BY id) AS rn
  FROM flashcards
  WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f'
    AND topic_id = 'a3333333-f010-484d-9bbd-e51c4991289f'
), design_sessions AS (
  SELECT
    sp.week_offset,
    sp.session_no,
    sp.cards_per_session,
    (
      substr(md5('alex-design-patterns-' || sp.week_offset || '-' || sp.session_no), 1, 8) || '-' ||
      substr(md5('alex-design-patterns-' || sp.week_offset || '-' || sp.session_no), 9, 4) || '-' ||
      substr(md5('alex-design-patterns-' || sp.week_offset || '-' || sp.session_no), 13, 4) || '-' ||
      substr(md5('alex-design-patterns-' || sp.week_offset || '-' || sp.session_no), 17, 4) || '-' ||
      substr(md5('alex-design-patterns-' || sp.week_offset || '-' || sp.session_no), 21, 12)
    )::uuid AS session_id,
    CURRENT_DATE
      - (sp.week_offset || ' weeks')::interval
      + (sp.session_hour || ' hours')::interval
      + (((sp.week_offset * 9 + sp.session_no * 13) % 36) || ' minutes')::interval AS session_ts
  FROM (
    VALUES
      (4, 1, 5, 19),
      (3, 1, 5, 20),
      (2, 1, 5, 21),
      (1, 1, 5, 19), (1, 2, 6, 21),
      (0, 1, 6, 20)
  ) AS sp(week_offset, session_no, cards_per_session, session_hour)
  JOIN generate_series(0, 4) AS week_idx(week_offset) ON week_idx.week_offset = sp.week_offset
), design_events AS (
  SELECT
    ds.session_id,
    ds.session_ts,
    dc.id AS flashcard_id,
    row_number() OVER (ORDER BY ds.session_ts, ds.session_id, dc.id) AS seq
  FROM design_sessions ds
  JOIN design_cards dc ON dc.rn <= ds.cards_per_session
)
INSERT INTO review_events
  (user_id, school_id, flashcard_id, correct, confidence, confused, response_time_ms, ts, session_id)
SELECT
  '1ceeff29-f010-484d-9bbd-e51c4991289f',
  '00000000-0000-0000-0000-000000000001',
  de.flashcard_id,
  de.seq NOT IN (8, 16, 24, 31),
  CASE
    WHEN de.seq IN (8, 16, 24, 31) THEN 2
    WHEN de.seq % 3 = 0 THEN 4
    ELSE 3
  END,
  de.seq = 14,
  1500 + ((de.seq * 161) % 1750),
  de.session_ts,
  de.session_id
FROM design_events de;

-- Rhetorical Analysis (~24 reviews, ~80% correct, ~8% confused)
WITH rhetoric_cards AS (
  SELECT id, row_number() OVER (ORDER BY id) AS rn
  FROM flashcards
  WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f'
    AND topic_id = 'a5555555-f010-484d-9bbd-e51c4991289f'
), rhetoric_sessions AS (
  SELECT
    sp.week_offset,
    sp.session_no,
    sp.cards_per_session,
    (
      substr(md5('alex-rhetoric-' || sp.week_offset || '-' || sp.session_no), 1, 8) || '-' ||
      substr(md5('alex-rhetoric-' || sp.week_offset || '-' || sp.session_no), 9, 4) || '-' ||
      substr(md5('alex-rhetoric-' || sp.week_offset || '-' || sp.session_no), 13, 4) || '-' ||
      substr(md5('alex-rhetoric-' || sp.week_offset || '-' || sp.session_no), 17, 4) || '-' ||
      substr(md5('alex-rhetoric-' || sp.week_offset || '-' || sp.session_no), 21, 12)
    )::uuid AS session_id,
    CURRENT_DATE
      - (sp.week_offset || ' weeks')::interval
      + (sp.session_hour || ' hours')::interval
      + (((sp.week_offset * 7 + sp.session_no * 23) % 32) || ' minutes')::interval AS session_ts
  FROM (
    VALUES
      (4, 1, 5, 20),
      (3, 1, 5, 19),
      (2, 1, 5, 21),
      (1, 1, 4, 20),
      (0, 1, 5, 19)
  ) AS sp(week_offset, session_no, cards_per_session, session_hour)
  JOIN generate_series(0, 4) AS week_idx(week_offset) ON week_idx.week_offset = sp.week_offset
), rhetoric_events AS (
  SELECT
    rs.session_id,
    rs.session_ts,
    rc.id AS flashcard_id,
    row_number() OVER (ORDER BY rs.session_ts, rs.session_id, rc.id) AS seq
  FROM rhetoric_sessions rs
  JOIN rhetoric_cards rc ON rc.rn <= rs.cards_per_session
)
INSERT INTO review_events
  (user_id, school_id, flashcard_id, correct, confidence, confused, response_time_ms, ts, session_id)
SELECT
  '1ceeff29-f010-484d-9bbd-e51c4991289f',
  '00000000-0000-0000-0000-000000000001',
  re.flashcard_id,
  re.seq NOT IN (5, 10, 15, 20, 24),
  CASE
    WHEN re.seq IN (5, 10, 15, 20, 24) THEN 2
    WHEN re.seq % 4 = 0 THEN 4
    ELSE 3
  END,
  re.seq IN (7, 18),
  1750 + ((re.seq * 137) % 1950),
  re.session_ts,
  re.session_id
FROM rhetoric_events re;

-- Kingdom of Kush (~12 reviews, started late, ~35% correct, ~33% confused)
WITH kush_cards AS (
  SELECT id, row_number() OVER (ORDER BY id) AS rn
  FROM flashcards
  WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f'
    AND topic_id = 'a4444444-f010-484d-9bbd-e51c4991289f'
), kush_sessions AS (
  SELECT
    sp.week_offset,
    sp.session_no,
    sp.cards_per_session,
    (
      substr(md5('alex-kush-' || sp.week_offset || '-' || sp.session_no), 1, 8) || '-' ||
      substr(md5('alex-kush-' || sp.week_offset || '-' || sp.session_no), 9, 4) || '-' ||
      substr(md5('alex-kush-' || sp.week_offset || '-' || sp.session_no), 13, 4) || '-' ||
      substr(md5('alex-kush-' || sp.week_offset || '-' || sp.session_no), 17, 4) || '-' ||
      substr(md5('alex-kush-' || sp.week_offset || '-' || sp.session_no), 21, 12)
    )::uuid AS session_id,
    CURRENT_DATE
      - (sp.week_offset || ' weeks')::interval
      + (sp.session_hour || ' hours')::interval
      + (((sp.week_offset * 5 + sp.session_no * 29) % 35) || ' minutes')::interval AS session_ts
  FROM (
    VALUES
      (1, 1, 4, 21),
      (0, 1, 4, 19),
      (0, 2, 4, 21)
  ) AS sp(week_offset, session_no, cards_per_session, session_hour)
  JOIN generate_series(0, 4) AS week_idx(week_offset) ON week_idx.week_offset = sp.week_offset
), kush_events AS (
  SELECT
    ks.session_id,
    ks.session_ts,
    kc.id AS flashcard_id,
    row_number() OVER (ORDER BY ks.session_ts, ks.session_id, kc.id) AS seq
  FROM kush_sessions ks
  JOIN kush_cards kc ON kc.rn <= ks.cards_per_session
)
INSERT INTO review_events
  (user_id, school_id, flashcard_id, correct, confidence, confused, response_time_ms, ts, session_id)
SELECT
  '1ceeff29-f010-484d-9bbd-e51c4991289f',
  '00000000-0000-0000-0000-000000000001',
  ke.flashcard_id,
  ke.seq IN (1, 4, 7, 10),
  CASE
    WHEN ke.seq IN (1, 4, 7, 10) THEN 2
    WHEN ke.seq IN (2, 5, 8, 11) THEN 1
    ELSE 2
  END,
  ke.seq IN (2, 5, 8, 11),
  3200 + ((ke.seq * 241) % 2600),
  ke.session_ts,
  ke.session_id
FROM kush_events ke;

-- Add some study deadlines (idempotent: unique on lms_assignment_id + school_id)
INSERT INTO lms_assignments (school_id, lms_assignment_id, lms_course_id, title, due_at, points_possible, assignment_group)
VALUES
('00000000-0000-0000-0000-000000000001', 'A-001', '1001', 'Meroitic Script Analysis', now() + interval '2 days', 50, 'Homework'),
('00000000-0000-0000-0000-000000000001', 'A-002', '1007', 'CNN Architecture Project', now() + interval '8 days', 100, 'Project')
ON CONFLICT (lms_assignment_id, school_id) DO UPDATE SET
  lms_course_id = EXCLUDED.lms_course_id,
  title = EXCLUDED.title,
  due_at = EXCLUDED.due_at,
  points_possible = EXCLUDED.points_possible,
  assignment_group = EXCLUDED.assignment_group,
  last_synced_at = now();

DELETE FROM topic_mastery
WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f'
  AND school_id = '00000000-0000-0000-0000-000000000001';

-- Alex's mastery scores as computed from his review history.
INSERT INTO topic_mastery
  (school_id, user_id, topic_id, mastery_score, confusion_rate, review_count)
VALUES
  -- Transformer Architectures: 5 weeks of consistent, accurate review
  ('00000000-0000-0000-0000-000000000001',
   '1ceeff29-f010-484d-9bbd-e51c4991289f',
   'a1111111-f010-484d-9bbd-e51c4991289f',
   0.82, 0.04, 48),

  -- B-Trees & Indexing: regular review, solid understanding
  ('00000000-0000-0000-0000-000000000001',
   '1ceeff29-f010-484d-9bbd-e51c4991289f',
   'a2222222-f010-484d-9bbd-e51c4991289f',
   0.79, 0.06, 40),

  -- Design Patterns: strong, reviewed ahead of SE project deadline
  ('00000000-0000-0000-0000-000000000001',
   '1ceeff29-f010-484d-9bbd-e51c4991289f',
   'a3333333-f010-484d-9bbd-e51c4991289f',
   0.84, 0.03, 32),

  -- Rhetorical Analysis: reviewed less but still above 75%
  ('00000000-0000-0000-0000-000000000001',
   '1ceeff29-f010-484d-9bbd-e51c4991289f',
   'a5555555-f010-484d-9bbd-e51c4991289f',
   0.77, 0.07, 24),

  -- Kingdom of Kush: struggling, started late with high confusion
  ('00000000-0000-0000-0000-000000000001',
   '1ceeff29-f010-484d-9bbd-e51c4991289f',
   'a4444444-f010-484d-9bbd-e51c4991289f',
   0.41, 0.28, 12)

ON CONFLICT (school_id, user_id, topic_id) DO UPDATE SET
  mastery_score  = EXCLUDED.mastery_score,
  confusion_rate = EXCLUDED.confusion_rate,
  review_count   = EXCLUDED.review_count,
  updated_at     = now();
