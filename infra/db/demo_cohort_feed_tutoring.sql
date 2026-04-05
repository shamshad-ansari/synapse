-- Peer users for the MIT school cohort.
-- Each user has their own courses, flashcards, and review history.
-- Mastery scores are the result of the mastery model running on their review events.
-- These rows represent the current state of their accounts.

BEGIN;

-- Fixed IDs
-- Alex:  1ceeff29-f010-484d-9bbd-e51c4991289f  alex@mit.edu
-- School: 00000000-0000-0000-0000-000000000001

-- ---------------------------------------------------------------------------
-- Cohort users (MIT school)
-- ---------------------------------------------------------------------------
INSERT INTO users (id, school_id, name, email, password_hash) VALUES
  ('c1000001-0000-4000-8000-000000000001', '00000000-0000-0000-0000-000000000001', 'Jordan Lee', 'jordan.lee@mit.edu', '$2a$12$h9hY66qscZS8MHnmYRdDqOeUTPjObHEfABKyyn.Gx7q5Sapsl20gu'),
  ('c1000002-0000-4000-8000-000000000002', '00000000-0000-0000-0000-000000000001', 'Priya Sharma', 'priya.sharma@mit.edu', '$2a$12$h9hY66qscZS8MHnmYRdDqOeUTPjObHEfABKyyn.Gx7q5Sapsl20gu'),
  ('c1000003-0000-4000-8000-000000000003', '00000000-0000-0000-0000-000000000001', 'Marcus Chen', 'marcus.chen@mit.edu', '$2a$12$h9hY66qscZS8MHnmYRdDqOeUTPjObHEfABKyyn.Gx7q5Sapsl20gu'),
  ('c1000004-0000-4000-8000-000000000004', '00000000-0000-0000-0000-000000000001', 'Sofia Alvarez', 'sofia.alvarez@mit.edu', '$2a$12$h9hY66qscZS8MHnmYRdDqOeUTPjObHEfABKyyn.Gx7q5Sapsl20gu'),
  ('c1000005-0000-4000-8000-000000000005', '00000000-0000-0000-0000-000000000001', 'Devon Wright', 'devon.wright@mit.edu', '$2a$12$h9hY66qscZS8MHnmYRdDqOeUTPjObHEfABKyyn.Gx7q5Sapsl20gu'),
  ('c1000006-0000-4000-8000-000000000006', '00000000-0000-0000-0000-000000000001', 'Nina Okonkwo', 'nina.okonkwo@mit.edu', '$2a$12$h9hY66qscZS8MHnmYRdDqOeUTPjObHEfABKyyn.Gx7q5Sapsl20gu'),
  ('c1000007-0000-4000-8000-000000000007', '00000000-0000-0000-0000-000000000001', 'Eli Rosen', 'eli.rosen@mit.edu', '$2a$12$h9hY66qscZS8MHnmYRdDqOeUTPjObHEfABKyyn.Gx7q5Sapsl20gu')
ON CONFLICT (id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- One course per peer (lms_course_id NULL) + topics for tutor matching
-- Topic names align with confusion / weak-topic strings (ILIKE in API).
-- ---------------------------------------------------------------------------
INSERT INTO courses (id, school_id, user_id, name, term, lms_course_id) VALUES
  ('d2000001-0000-4000-8000-000000000001', '00000000-0000-0000-0000-000000000001', 'c1000001-0000-4000-8000-000000000001', 'African History', '2026 Spring', NULL),
  ('d2000002-0000-4000-8000-000000000002', '00000000-0000-0000-0000-000000000001', 'c1000002-0000-4000-8000-000000000002', 'African History', '2026 Spring', NULL),
  ('d2000003-0000-4000-8000-000000000003', '00000000-0000-0000-0000-000000000001', 'c1000003-0000-4000-8000-000000000003', 'Machine Learning', '2026 Spring', NULL),
  ('d2000004-0000-4000-8000-000000000004', '00000000-0000-0000-0000-000000000001', 'c1000004-0000-4000-8000-000000000004', 'African History', '2026 Spring', NULL),
  ('d2000005-0000-4000-8000-000000000005', '00000000-0000-0000-0000-000000000001', 'c1000005-0000-4000-8000-000000000005', 'Introduction to Algorithms', '2026 Spring', NULL),
  ('d2000006-0000-4000-8000-000000000006', '00000000-0000-0000-0000-000000000001', 'c1000006-0000-4000-8000-000000000006', 'Machine Learning', '2026 Spring', NULL),
  ('d2000007-0000-4000-8000-000000000007', '00000000-0000-0000-0000-000000000001', 'c1000007-0000-4000-8000-000000000007', 'African History', '2026 Spring', NULL)
ON CONFLICT (id) DO UPDATE SET
  school_id = EXCLUDED.school_id,
  user_id = EXCLUDED.user_id,
  name = EXCLUDED.name,
  term = EXCLUDED.term,
  lms_course_id = EXCLUDED.lms_course_id;

INSERT INTO topics (id, school_id, course_id, user_id, name, exam_weight, source) VALUES
  ('a3010001-0000-4000-8000-000000000001', '00000000-0000-0000-0000-000000000001', 'd2000001-0000-4000-8000-000000000001', 'c1000001-0000-4000-8000-000000000001', 'Kingdom of Kush', 0.10, 'manual'),
  ('a3010002-0000-4000-8000-000000000002', '00000000-0000-0000-0000-000000000001', 'd2000001-0000-4000-8000-000000000001', 'c1000001-0000-4000-8000-000000000001', 'Transformer Architectures', 0.15, 'manual'),
  ('a3020001-0000-4000-8000-000000000001', '00000000-0000-0000-0000-000000000001', 'd2000002-0000-4000-8000-000000000002', 'c1000002-0000-4000-8000-000000000002', 'Kingdom of Kush', 0.10, 'manual'),
  ('a3030001-0000-4000-8000-000000000001', '00000000-0000-0000-0000-000000000001', 'd2000003-0000-4000-8000-000000000003', 'c1000003-0000-4000-8000-000000000003', 'Transformer Architectures', 0.15, 'manual'),
  ('a3030002-0000-4000-8000-000000000002', '00000000-0000-0000-0000-000000000001', 'd2000003-0000-4000-8000-000000000003', 'c1000003-0000-4000-8000-000000000003', 'B-Trees & Indexing', 0.20, 'manual'),
  ('a3040001-0000-4000-8000-000000000001', '00000000-0000-0000-0000-000000000001', 'd2000004-0000-4000-8000-000000000004', 'c1000004-0000-4000-8000-000000000004', 'Kingdom of Kush', 0.10, 'manual'),
  ('a3050001-0000-4000-8000-000000000001', '00000000-0000-0000-0000-000000000001', 'd2000005-0000-4000-8000-000000000005', 'c1000005-0000-4000-8000-000000000005', 'B-Trees & Indexing', 0.20, 'manual'),
  ('a3060001-0000-4000-8000-000000000001', '00000000-0000-0000-0000-000000000001', 'd2000006-0000-4000-8000-000000000006', 'c1000006-0000-4000-8000-000000000006', 'Transformer Architectures', 0.15, 'manual'),
  ('a3070001-0000-4000-8000-000000000001', '00000000-0000-0000-0000-000000000001', 'd2000007-0000-4000-8000-000000000007', 'c1000007-0000-4000-8000-000000000007', 'Kingdom of Kush', 0.10, 'manual'),
  ('a3070002-0000-4000-8000-000000000002', '00000000-0000-0000-0000-000000000001', 'd2000007-0000-4000-8000-000000000007', 'c1000007-0000-4000-8000-000000000007', 'Rhetorical Analysis', 0.10, 'manual')
ON CONFLICT (id) DO NOTHING;

-- Mastery snapshot for tutor matching.
--   • Peers >= 0.75 on "Kingdom of Kush" only → suggested tutors when Alex's weak topic is Kush.
--   • Peers < 0.75 on Transformers / B-Trees / Rhetoric → they do not outrank Alex on his teaching topics; incoming requests to Alex stay plausible.
INSERT INTO topic_mastery (school_id, user_id, topic_id, mastery_score, confusion_rate, review_count) VALUES
  ('00000000-0000-0000-0000-000000000001', 'c1000001-0000-4000-8000-000000000001', 'a3010001-0000-4000-8000-000000000001', 0.90, 0.05, 42),
  ('00000000-0000-0000-0000-000000000001', 'c1000001-0000-4000-8000-000000000001', 'a3010002-0000-4000-8000-000000000002', 0.62, 0.12, 18),
  ('00000000-0000-0000-0000-000000000001', 'c1000002-0000-4000-8000-000000000002', 'a3020001-0000-4000-8000-000000000001', 0.88, 0.07, 35),
  ('00000000-0000-0000-0000-000000000001', 'c1000003-0000-4000-8000-000000000003', 'a3030001-0000-4000-8000-000000000001', 0.64, 0.14, 22),
  ('00000000-0000-0000-0000-000000000001', 'c1000003-0000-4000-8000-000000000003', 'a3030002-0000-4000-8000-000000000002', 0.68, 0.13, 20),
  ('00000000-0000-0000-0000-000000000001', 'c1000004-0000-4000-8000-000000000004', 'a3040001-0000-4000-8000-000000000001', 0.87, 0.06, 38),
  ('00000000-0000-0000-0000-000000000001', 'c1000005-0000-4000-8000-000000000005', 'a3050001-0000-4000-8000-000000000001', 0.66, 0.14, 21),
  ('00000000-0000-0000-0000-000000000001', 'c1000006-0000-4000-8000-000000000006', 'a3060001-0000-4000-8000-000000000001', 0.63, 0.13, 19),
  ('00000000-0000-0000-0000-000000000001', 'c1000007-0000-4000-8000-000000000007', 'a3070001-0000-4000-8000-000000000001', 0.86, 0.07, 33),
  ('00000000-0000-0000-0000-000000000001', 'c1000007-0000-4000-8000-000000000007', 'a3070002-0000-4000-8000-000000000002', 0.61, 0.15, 16)
ON CONFLICT (school_id, user_id, topic_id) DO UPDATE SET
  mastery_score = EXCLUDED.mastery_score,
  confusion_rate = EXCLUDED.confusion_rate,
  review_count = EXCLUDED.review_count,
  updated_at = now();

-- ---------------------------------------------------------------------------
-- Feed posts (authors: Alex + cohort). Sort in API: upvotes DESC, created_at DESC.
-- ---------------------------------------------------------------------------
INSERT INTO feed_posts (id, school_id, user_id, course_id, topic_id, post_type, title, body, upvotes, created_at) VALUES
  ('f1000001-0000-4000-8000-000000000001', '00000000-0000-0000-0000-000000000001', 'c1000003-0000-4000-8000-000000000003', 'e501a47c-d907-4f82-9ffe-61a221212f55', 'a1111111-f010-484d-9bbd-e51c4991289f', 'share', 'Best explainer on attention is still the original paper', 'Linking the Vaswani figures with the 6.790 lecture notes made the Q/K/V click for me. Happy to share my annotated PDF in office hours.', 7, now() - interval '2 hours'),
  ('f1000002-0000-4000-8000-000000000002', '00000000-0000-0000-0000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', 'a4444444-f010-484d-9bbd-e51c4991289f', 'question', 'Kandake vs pharaoh terminology', 'Are we expected to use Kandake for every ruler slide or only when the source uses it? Want to match the rubric on the midterm.', 7, now() - interval '5 hours'),
  ('f1000003-0000-4000-8000-000000000003', '00000000-0000-0000-0000-000000000001', 'c1000002-0000-4000-8000-000000000002', '847e3a18-c7a4-4986-b991-b21439b7a9ff', 'a2222222-f010-484d-9bbd-e51c4991289f', 'update', 'B+ tree range query proof sketch', 'I wrote a one-page outline for why leaf links give O(log n + k). Posting here in case it saves someone an hour before PS4.', 7, now() - interval '1 day'),
  ('f1000004-0000-4000-8000-000000000004', '00000000-0000-0000-0000-000000000001', 'c1000001-0000-4000-8000-000000000001', 'e501a47c-d907-4f82-9ffe-61a221212f55', 'a1111111-f010-484d-9bbd-e51c4991289f', 'question', 'Positional encoding vs learned embeddings', 'For the project are we allowed to swap sinusoidal PE for learned positions if we ablate both?', 7, now() - interval '1 day 3 hours'),
  ('f1000005-0000-4000-8000-000000000005', '00000000-0000-0000-0000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'b45b4a4a-f123-479b-8cfe-ecfbc26546e0', 'a3333333-f010-484d-9bbd-e51c4991289f', 'share', 'Observer pattern in the wild', 'If you are debugging the Canvas notification mock, the observer lecture code maps almost 1:1 to their event bus.', 5, now() - interval '30 minutes'),
  ('f1000006-0000-4000-8000-000000000006', '00000000-0000-0000-0000-000000000001', 'c1000006-0000-4000-8000-000000000006', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', 'a4444444-f010-484d-9bbd-e51c4991289f', 'update', 'Meroë iron production timeline', 'Compiled dates from Garlake and the course reader—useful for the essay thesis.', 5, now() - interval '6 hours'),
  ('f1000007-0000-4000-8000-000000000007', '00000000-0000-0000-0000-000000000001', 'c1000004-0000-4000-8000-000000000004', '61f43c8a-154e-4af3-9e77-81a25cf30c1a', 'a5555555-f010-484d-9bbd-e51c4991289f', 'question', 'Thesis for rhetorical analysis draft', 'Is it OK if my thesis names two appeals or should it stay to one primary appeal for Comp 2?', 5, now() - interval '8 hours'),
  ('f1000008-0000-4000-8000-000000000008', '00000000-0000-0000-0000-000000000001', 'c1000005-0000-4000-8000-000000000005', '847e3a18-c7a4-4986-b991-b21439b7a9ff', 'a2222222-f010-484d-9bbd-e51c4991289f', 'share', 'Visualization for B+ tree splits', 'Found an animation that matches our lecture notation—link in comments on Stellar.', 5, now() - interval '10 hours'),
  ('f1000009-0000-4000-8000-000000000009', '00000000-0000-0000-0000-000000000001', 'c1000007-0000-4000-8000-000000000007', 'e501a47c-d907-4f82-9ffe-61a221212f55', 'a1111111-f010-484d-9bbd-e51c4991289f', 'question', 'Multi-head attention implementation detail', 'Do we fuse the heads before the residual or after the second linear in the reference impl?', 4, now() - interval '2 days'),
  ('f1000010-0000-4000-8000-000000000010', '00000000-0000-0000-0000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', '847e3a18-c7a4-4986-b991-b21439b7a9ff', 'a2222222-f010-484d-9bbd-e51c4991289f', 'update', 'PS3 sanity check', 'Fan-out calculation on Q3—if you got 341 you probably forgot ceiling on the split. Double-check with the practice key.', 4, now() - interval '2 days 4 hours'),
  ('f1000011-0000-4000-8000-000000000011', '00000000-0000-0000-0000-000000000001', 'c1000003-0000-4000-8000-000000000003', 'b45b4a4a-f123-479b-8cfe-ecfbc26546e0', 'a3333333-f010-484d-9bbd-e51c4991289f', 'share', 'Strategy pattern for the capstone', 'We wrapped each grading strategy behind an interface—made the service layer much cleaner.', 3, now() - interval '3 days'),
  ('f1000012-0000-4000-8000-000000000012', '00000000-0000-0000-0000-000000000001', 'c1000002-0000-4000-8000-000000000002', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', 'a4444444-f010-484d-9bbd-e51c4991289f', 'question', 'Primary sources for Kandakes', 'Does anyone have the translated Stele excerpt bookmarked? Library scan is down for me.', 3, now() - interval '3 days 12 hours'),
  ('f1000013-0000-4000-8000-000000000013', '00000000-0000-0000-0000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', '61f43c8a-154e-4af3-9e77-81a25cf30c1a', 'a5555555-f010-484d-9bbd-e51c4991289f', 'question', 'Peer review workshop', 'Anyone want to swap Comp essays tonight? I can do a 30-minute slot after 9pm.', 2, now() - interval '4 days'),
  ('f1000014-0000-4000-8000-000000000014', '00000000-0000-0000-0000-000000000001', 'c1000005-0000-4000-8000-000000000005', 'e501a47c-d907-4f82-9ffe-61a221212f55', 'a1111111-f010-484d-9bbd-e51c4991289f', 'update', 'Transformer training memes aside', 'Finally got my batch to converge after fixing LR warmup. Loss curve looks sane now.', 2, now() - interval '5 days'),
  ('f1000015-0000-4000-8000-000000000015', '00000000-0000-0000-0000-000000000001', 'c1000001-0000-4000-8000-000000000001', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', 'a4444444-f010-484d-9bbd-e51c4991289f', 'question', 'Map quiz: Nile cataracts', 'Which week should we memorize the cataract numbering—midterm scope sheet is ambiguous.', 1, now() - interval '6 days'),
  ('f1000016-0000-4000-8000-000000000016', '00000000-0000-0000-0000-000000000001', 'c1000006-0000-4000-8000-000000000006', 'b45b4a4a-f123-479b-8cfe-ecfbc26546e0', 'a3333333-f010-484d-9bbd-e51c4991289f', 'share', 'GoF cheat sheet', 'One-page PDF for behavioral patterns—might help before the quiz.', 1, now() - interval '7 days'),
  ('f1000017-0000-4000-8000-000000000017', '00000000-0000-0000-0000-000000000001', 'c1000007-0000-4000-8000-000000000007', '847e3a18-c7a4-4986-b991-b21439b7a9ff', 'a2222222-f010-484d-9bbd-e51c4991289f', 'update', 'OH queue for indexing', 'I will be at the 2pm OH—first come first served for B-tree proofs.', 0, now() - interval '8 days'),
  ('f1000018-0000-4000-8000-000000000018', '00000000-0000-0000-0000-000000000001', 'c1000004-0000-4000-8000-000000000004', 'e501a47c-d907-4f82-9ffe-61a221212f55', 'a1111111-f010-484d-9bbd-e51c4991289f', 'question', 'Extra credit reading', 'Did anyone finish the optional ViT paper—worth summarizing for discussion?', 0, now() - interval '10 days')
ON CONFLICT (id) DO NOTHING;

-- Upvotes: each row must be unique (user_id, post_id). Counts match feed_posts.upvotes.
-- Voters for "author = Alex" are all 7 peers. For author = peer, voters = Alex + 6 peers (exclude author).

INSERT INTO feed_upvotes (user_id, post_id) VALUES
  -- f1 Marcus post: Alex + P1,P2,P4,P5,P6,P7 (not P3)
  ('1ceeff29-f010-484d-9bbd-e51c4991289f', 'f1000001-0000-4000-8000-000000000001'),
  ('c1000001-0000-4000-8000-000000000001', 'f1000001-0000-4000-8000-000000000001'),
  ('c1000002-0000-4000-8000-000000000002', 'f1000001-0000-4000-8000-000000000001'),
  ('c1000004-0000-4000-8000-000000000004', 'f1000001-0000-4000-8000-000000000001'),
  ('c1000005-0000-4000-8000-000000000005', 'f1000001-0000-4000-8000-000000000001'),
  ('c1000006-0000-4000-8000-000000000006', 'f1000001-0000-4000-8000-000000000001'),
  ('c1000007-0000-4000-8000-000000000007', 'f1000001-0000-4000-8000-000000000001'),
  -- f2 Alex: all peers
  ('c1000001-0000-4000-8000-000000000001', 'f1000002-0000-4000-8000-000000000002'),
  ('c1000002-0000-4000-8000-000000000002', 'f1000002-0000-4000-8000-000000000002'),
  ('c1000003-0000-4000-8000-000000000003', 'f1000002-0000-4000-8000-000000000002'),
  ('c1000004-0000-4000-8000-000000000004', 'f1000002-0000-4000-8000-000000000002'),
  ('c1000005-0000-4000-8000-000000000005', 'f1000002-0000-4000-8000-000000000002'),
  ('c1000006-0000-4000-8000-000000000006', 'f1000002-0000-4000-8000-000000000002'),
  ('c1000007-0000-4000-8000-000000000007', 'f1000002-0000-4000-8000-000000000002'),
  -- f3 Priya: Alex + not Priya
  ('1ceeff29-f010-484d-9bbd-e51c4991289f', 'f1000003-0000-4000-8000-000000000003'),
  ('c1000001-0000-4000-8000-000000000001', 'f1000003-0000-4000-8000-000000000003'),
  ('c1000003-0000-4000-8000-000000000003', 'f1000003-0000-4000-8000-000000000003'),
  ('c1000004-0000-4000-8000-000000000004', 'f1000003-0000-4000-8000-000000000003'),
  ('c1000005-0000-4000-8000-000000000005', 'f1000003-0000-4000-8000-000000000003'),
  ('c1000006-0000-4000-8000-000000000006', 'f1000003-0000-4000-8000-000000000003'),
  ('c1000007-0000-4000-8000-000000000007', 'f1000003-0000-4000-8000-000000000003'),
  -- f4 Jordan
  ('1ceeff29-f010-484d-9bbd-e51c4991289f', 'f1000004-0000-4000-8000-000000000004'),
  ('c1000002-0000-4000-8000-000000000002', 'f1000004-0000-4000-8000-000000000004'),
  ('c1000003-0000-4000-8000-000000000003', 'f1000004-0000-4000-8000-000000000004'),
  ('c1000004-0000-4000-8000-000000000004', 'f1000004-0000-4000-8000-000000000004'),
  ('c1000005-0000-4000-8000-000000000005', 'f1000004-0000-4000-8000-000000000004'),
  ('c1000006-0000-4000-8000-000000000006', 'f1000004-0000-4000-8000-000000000004'),
  ('c1000007-0000-4000-8000-000000000007', 'f1000004-0000-4000-8000-000000000004'),
  -- f5 Alex: 5 votes
  ('c1000001-0000-4000-8000-000000000001', 'f1000005-0000-4000-8000-000000000005'),
  ('c1000002-0000-4000-8000-000000000002', 'f1000005-0000-4000-8000-000000000005'),
  ('c1000003-0000-4000-8000-000000000003', 'f1000005-0000-4000-8000-000000000005'),
  ('c1000004-0000-4000-8000-000000000004', 'f1000005-0000-4000-8000-000000000005'),
  ('c1000005-0000-4000-8000-000000000005', 'f1000005-0000-4000-8000-000000000005'),
  -- f6 Nina: 5
  ('1ceeff29-f010-484d-9bbd-e51c4991289f', 'f1000006-0000-4000-8000-000000000006'),
  ('c1000001-0000-4000-8000-000000000001', 'f1000006-0000-4000-8000-000000000006'),
  ('c1000002-0000-4000-8000-000000000002', 'f1000006-0000-4000-8000-000000000006'),
  ('c1000003-0000-4000-8000-000000000003', 'f1000006-0000-4000-8000-000000000006'),
  ('c1000004-0000-4000-8000-000000000004', 'f1000006-0000-4000-8000-000000000006'),
  -- f7 Sofia: 5
  ('1ceeff29-f010-484d-9bbd-e51c4991289f', 'f1000007-0000-4000-8000-000000000007'),
  ('c1000001-0000-4000-8000-000000000001', 'f1000007-0000-4000-8000-000000000007'),
  ('c1000002-0000-4000-8000-000000000002', 'f1000007-0000-4000-8000-000000000007'),
  ('c1000003-0000-4000-8000-000000000003', 'f1000007-0000-4000-8000-000000000007'),
  ('c1000005-0000-4000-8000-000000000005', 'f1000007-0000-4000-8000-000000000007'),
  -- f8 Devon: 5
  ('1ceeff29-f010-484d-9bbd-e51c4991289f', 'f1000008-0000-4000-8000-000000000008'),
  ('c1000001-0000-4000-8000-000000000001', 'f1000008-0000-4000-8000-000000000008'),
  ('c1000002-0000-4000-8000-000000000002', 'f1000008-0000-4000-8000-000000000008'),
  ('c1000003-0000-4000-8000-000000000003', 'f1000008-0000-4000-8000-000000000008'),
  ('c1000004-0000-4000-8000-000000000004', 'f1000008-0000-4000-8000-000000000008'),
  -- f9 Eli: 4
  ('1ceeff29-f010-484d-9bbd-e51c4991289f', 'f1000009-0000-4000-8000-000000000009'),
  ('c1000001-0000-4000-8000-000000000001', 'f1000009-0000-4000-8000-000000000009'),
  ('c1000002-0000-4000-8000-000000000002', 'f1000009-0000-4000-8000-000000000009'),
  ('c1000003-0000-4000-8000-000000000003', 'f1000009-0000-4000-8000-000000000009'),
  -- f10 Alex: 4
  ('c1000001-0000-4000-8000-000000000001', 'f1000010-0000-4000-8000-000000000010'),
  ('c1000002-0000-4000-8000-000000000002', 'f1000010-0000-4000-8000-000000000010'),
  ('c1000004-0000-4000-8000-000000000004', 'f1000010-0000-4000-8000-000000000010'),
  ('c1000005-0000-4000-8000-000000000005', 'f1000010-0000-4000-8000-000000000010'),
  -- f11 Marcus: 3
  ('1ceeff29-f010-484d-9bbd-e51c4991289f', 'f1000011-0000-4000-8000-000000000011'),
  ('c1000001-0000-4000-8000-000000000001', 'f1000011-0000-4000-8000-000000000011'),
  ('c1000002-0000-4000-8000-000000000002', 'f1000011-0000-4000-8000-000000000011'),
  -- f12 Priya: 3
  ('1ceeff29-f010-484d-9bbd-e51c4991289f', 'f1000012-0000-4000-8000-000000000012'),
  ('c1000003-0000-4000-8000-000000000003', 'f1000012-0000-4000-8000-000000000012'),
  ('c1000004-0000-4000-8000-000000000004', 'f1000012-0000-4000-8000-000000000012'),
  -- f13 Alex: 2
  ('c1000001-0000-4000-8000-000000000001', 'f1000013-0000-4000-8000-000000000013'),
  ('c1000005-0000-4000-8000-000000000005', 'f1000013-0000-4000-8000-000000000013'),
  -- f14 Devon: 2
  ('1ceeff29-f010-484d-9bbd-e51c4991289f', 'f1000014-0000-4000-8000-000000000014'),
  ('c1000002-0000-4000-8000-000000000002', 'f1000014-0000-4000-8000-000000000014'),
  -- f15 Jordan: 1
  ('c1000003-0000-4000-8000-000000000003', 'f1000015-0000-4000-8000-000000000015'),
  -- f16 Nina: 1
  ('1ceeff29-f010-484d-9bbd-e51c4991289f', 'f1000016-0000-4000-8000-000000000016')
ON CONFLICT (user_id, post_id) DO NOTHING;

-- ---------------------------------------------------------------------------
-- Tutoring requests (incoming: peers -> Alex; outgoing: Alex -> peers)
-- ---------------------------------------------------------------------------
-- Incoming to Alex: topics match Alex's teaching profile (strong mastery). One Kush row removed so Priya asks for Design Patterns instead.
INSERT INTO tutor_requests (id, school_id, requester_id, tutor_id, topic_id, topic_name, status, message, created_at) VALUES
  ('b1000001-0000-4000-8000-000000000001', '00000000-0000-0000-0000-000000000001', 'c1000002-0000-4000-8000-000000000002', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a3333333-f010-484d-9bbd-e51c4991289f', 'Design Patterns', 'pending', 'UML draft is due Friday—could we whiteboard Strategy vs Template for 15m after lecture?', now() - interval '45 minutes'),
  ('b1000002-0000-4000-8000-000000000002', '00000000-0000-0000-0000-000000000001', 'c1000003-0000-4000-8000-000000000003', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a1111111-f010-484d-9bbd-e51c4991289f', 'Transformer Architectures', 'pending', 'Want to walk through multi-head shapes for the final project checkpoint.', now() - interval '3 hours'),
  ('b1000003-0000-4000-8000-000000000003', '00000000-0000-0000-0000-000000000001', 'c1000005-0000-4000-8000-000000000005', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a2222222-f010-484d-9bbd-e51c4991289f', 'B-Trees & Indexing', 'pending', 'Stuck on fan-out vs order in the practice exam—could use a quick intuition pass.', now() - interval '5 hours'),
  ('b1000004-0000-4000-8000-000000000004', '00000000-0000-0000-0000-000000000001', 'c1000004-0000-4000-8000-000000000004', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a5555555-f010-484d-9bbd-e51c4991289f', 'Rhetorical Analysis', 'pending', 'Happy to trade essay feedback if you have time for ethos paragraphs.', now() - interval '1 day'),
  ('b1000005-0000-4000-8000-000000000005', '00000000-0000-0000-0000-000000000001', 'c1000006-0000-4000-8000-000000000006', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a3333333-f010-484d-9bbd-e51c4991289f', 'Design Patterns', 'pending', 'Need help choosing Strategy vs Template method for the UML portion.', now() - interval '1 day 2 hours'),
  ('b1000006-0000-4000-8000-000000000006', '00000000-0000-0000-0000-000000000001', 'c1000007-0000-4000-8000-000000000007', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a5555555-f010-484d-9bbd-e51c4991289f', 'Rhetorical Analysis', 'accepted', 'Thanks for accepting—see you at the writing center at 5.', now() - interval '2 days'),
  ('b1000007-0000-4000-8000-000000000007', '00000000-0000-0000-0000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'c1000001-0000-4000-8000-000000000001', 'a1111111-f010-484d-9bbd-e51c4991289f', 'Transformer Architectures', 'pending', 'You crushed the PE homework—could I buy you coffee and ask about attention masks?', now() - interval '6 hours'),
  ('b1000008-0000-4000-8000-000000000008', '00000000-0000-0000-0000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'c1000003-0000-4000-8000-000000000003', 'a2222222-f010-484d-9bbd-e51c4991289f', 'B-Trees & Indexing', 'pending', 'Trying to match your split walkthrough from the feed—30m this weekend?', now() - interval '1 day'),
  ('b1000009-0000-4000-8000-000000000009', '00000000-0000-0000-0000-000000000001', 'c1000001-0000-4000-8000-000000000001', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a1111111-f010-484d-9bbd-e51c4991289f', 'Transformer Architectures', 'pending', 'Office hours overlap tomorrow? I have a question on layer norm placement.', now() - interval '30 minutes')
ON CONFLICT (id) DO UPDATE SET
  requester_id = EXCLUDED.requester_id,
  tutor_id = EXCLUDED.tutor_id,
  topic_id = EXCLUDED.topic_id,
  topic_name = EXCLUDED.topic_name,
  status = EXCLUDED.status,
  message = EXCLUDED.message,
  created_at = EXCLUDED.created_at;

-- Jordan and Priya have studied Kingdom of Kush heavily.
-- Their high review counts reflect weeks of active study on this topic.
UPDATE topic_mastery SET review_count = 58, updated_at = now()
WHERE user_id = 'c1000001-0000-4000-8000-000000000001'
  AND topic_id = 'a3010001-0000-4000-8000-000000000001'
  AND school_id = '00000000-0000-0000-0000-000000000001';

UPDATE topic_mastery SET review_count = 51, updated_at = now()
WHERE user_id = 'c1000002-0000-4000-8000-000000000002'
  AND topic_id = 'a3020001-0000-4000-8000-000000000001'
  AND school_id = '00000000-0000-0000-0000-000000000001';

UPDATE topic_mastery SET review_count = 49, updated_at = now()
WHERE user_id = 'c1000004-0000-4000-8000-000000000004'
  AND topic_id = 'a3040001-0000-4000-8000-000000000001'
  AND school_id = '00000000-0000-0000-0000-000000000001';

UPDATE topic_mastery SET review_count = 46, updated_at = now()
WHERE user_id = 'c1000007-0000-4000-8000-000000000007'
  AND topic_id = 'a3070001-0000-4000-8000-000000000001'
  AND school_id = '00000000-0000-0000-0000-000000000001';

-- Marcus, Devon, Nina are weaker on Alex's strong topics.
-- This is why they send Alex tutoring requests based on their mastery history.
UPDATE topic_mastery SET mastery_score = 0.55, updated_at = now()
WHERE user_id = 'c1000003-0000-4000-8000-000000000003'
  AND topic_id = 'a3030001-0000-4000-8000-000000000001'
  AND school_id = '00000000-0000-0000-0000-000000000001';

UPDATE topic_mastery SET mastery_score = 0.58, updated_at = now()
WHERE user_id = 'c1000005-0000-4000-8000-000000000005'
  AND topic_id = 'a3050001-0000-4000-8000-000000000001'
  AND school_id = '00000000-0000-0000-0000-000000000001';

UPDATE topic_mastery SET mastery_score = 0.52, updated_at = now()
WHERE user_id = 'c1000006-0000-4000-8000-000000000006'
  AND topic_id = 'a3060001-0000-4000-8000-000000000001'
  AND school_id = '00000000-0000-0000-0000-000000000001';

COMMIT;
