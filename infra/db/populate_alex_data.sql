-- Seed data for Alex (alex@mit.edu)
-- User ID: 1ceeff29-f010-484d-9bbd-e51c4991289f
-- School ID: 00000000-0000-0000-0000-000000000001
--
-- Idempotent: safe to re-run make seed-demo (clears prior demo rows for fixed topic IDs).

BEGIN;

DELETE FROM review_events re
WHERE re.flashcard_id IN (
  SELECT id FROM flashcards
  WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f'
    AND topic_id IN (
      'a1111111-f010-484d-9bbd-e51c4991289f',
      'a2222222-f010-484d-9bbd-e51c4991289f',
      'a3333333-f010-484d-9bbd-e51c4991289f',
      'a4444444-f010-484d-9bbd-e51c4991289f',
      'a5555555-f010-484d-9bbd-e51c4991289f'
    )
);

DELETE FROM scheduler_states ss
WHERE ss.flashcard_id IN (
  SELECT id FROM flashcards
  WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f'
    AND topic_id IN (
      'a1111111-f010-484d-9bbd-e51c4991289f',
      'a2222222-f010-484d-9bbd-e51c4991289f',
      'a3333333-f010-484d-9bbd-e51c4991289f',
      'a4444444-f010-484d-9bbd-e51c4991289f',
      'a5555555-f010-484d-9bbd-e51c4991289f'
    )
);

DELETE FROM flashcards
WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f'
  AND topic_id IN (
    'a1111111-f010-484d-9bbd-e51c4991289f',
    'a2222222-f010-484d-9bbd-e51c4991289f',
    'a3333333-f010-484d-9bbd-e51c4991289f',
    'a4444444-f010-484d-9bbd-e51c4991289f',
    'a5555555-f010-484d-9bbd-e51c4991289f'
  );

DELETE FROM note_texts
WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f'
  AND topic_id IN (
    'a1111111-f010-484d-9bbd-e51c4991289f',
    'a2222222-f010-484d-9bbd-e51c4991289f',
    'a3333333-f010-484d-9bbd-e51c4991289f',
    'a4444444-f010-484d-9bbd-e51c4991289f',
    'a5555555-f010-484d-9bbd-e51c4991289f'
  );

-- 1. Machine Learning: Transformer Architectures
INSERT INTO topics (id, school_id, course_id, user_id, name, exam_weight, source)
VALUES ('a1111111-f010-484d-9bbd-e51c4991289f', '00000000-0000-0000-0000-000000000001', 'e501a47c-d907-4f82-9ffe-61a221212f55', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'Transformer Architectures', 0.15, 'manual')
ON CONFLICT (id) DO NOTHING;

INSERT INTO note_texts (id, school_id, course_id, user_id, topic_id, title, content)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'e501a47c-d907-4f82-9ffe-61a221212f55', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a1111111-f010-484d-9bbd-e51c4991289f', 'Transformers & Attention', 'The core of modern LLMs is the Transformer architecture. It relies on Self-Attention mechanism to weigh input tokens regardless of their distance. Multi-head attention allows the model to jointly attend to information from different representation subspaces at different positions.')
ON CONFLICT DO NOTHING;

-- Flashcards for ML
INSERT INTO flashcards (id, school_id, course_id, user_id, topic_id, card_type, prompt, answer, created_by)
VALUES 
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'e501a47c-d907-4f82-9ffe-61a221212f55', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a1111111-f010-484d-9bbd-e51c4991289f', 'qa', 'What is the primary mechanism that allows Transformers to bypass sequential processing?', 'Self-Attention (allowing parallelization of input data)', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'e501a47c-d907-4f82-9ffe-61a221212f55', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a1111111-f010-484d-9bbd-e51c4991289f', 'qa', 'Name the 3 vectors used in scaled dot-product attention.', 'Query (Q), Key (K), and Value (V)', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'e501a47c-d907-4f82-9ffe-61a221212f55', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a1111111-f010-484d-9bbd-e51c4991289f', 'qa', 'Why is Positional Encoding necessary in Transformers?', 'Because Transformers have no recurrence (RNN) or convolution, they need a way to know the order of tokens.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'e501a47c-d907-4f82-9ffe-61a221212f55', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a1111111-f010-484d-9bbd-e51c4991289f', 'qa', 'What is Multi-Head Attention?', 'Running multiple attention mechanisms (heads) in parallel for different representation subspaces.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'e501a47c-d907-4f82-9ffe-61a221212f55', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a1111111-f010-484d-9bbd-e51c4991289f', 'qa', 'Which paper introduced the Transformer architecture?', 'Attention Is All You Need (Vaswani et al., 2017)', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'e501a47c-d907-4f82-9ffe-61a221212f55', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a1111111-f010-484d-9bbd-e51c4991289f', 'qa', 'What are residual connections used for in Transformers?', 'To help prevent vanishing gradients in deep networks by letting information bypass layers.', 'manual');

-- 2. Algorithms & Systems: B-Trees and Indexing
INSERT INTO topics (id, school_id, course_id, user_id, name, exam_weight, source)
VALUES ('a2222222-f010-484d-9bbd-e51c4991289f', '00000000-0000-0000-0000-000000000001', '847e3a18-c7a4-4986-b991-b21439b7a9ff', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'B-Trees & Indexing', 0.20, 'manual')
ON CONFLICT (id) DO NOTHING;

INSERT INTO note_texts (id, school_id, course_id, user_id, topic_id, title, content)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '847e3a18-c7a4-4986-b991-b21439b7a9ff', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a2222222-f010-484d-9bbd-e51c4991289f', 'Database Indexing with B-Trees', 'B-Trees are self-balancing search trees that facilitate fast data retrieval. B+ Trees are often used in DB systems because data resides only in leaves, and leaf nodes are linked, allowing efficient range queries.')
ON CONFLICT DO NOTHING;

-- Flashcards for Algorithms
INSERT INTO flashcards (id, school_id, course_id, user_id, topic_id, card_type, prompt, answer, created_by)
VALUES 
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '847e3a18-c7a4-4986-b991-b21439b7a9ff', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a2222222-f010-484d-9bbd-e51c4991289f', 'qa', 'Where does the actual data reside in a B+ Tree vs B-Tree?', 'In B+ Tree, data only in leaves. In B-Tree, data can be in internal nodes too.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '847e3a18-c7a4-4986-b991-b21439b7a9ff', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a2222222-f010-484d-9bbd-e51c4991289f', 'qa', 'Why are B+ Trees better for range queries?', 'Because leaf nodes are linked in a Doubly Linked List.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '847e3a18-c7a4-4986-b991-b21439b7a9ff', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a2222222-f010-484d-9bbd-e51c4991289f', 'qa', 'What is the fan-out in a B-Tree?', 'The number of children a node can have.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '847e3a18-c7a4-4986-b991-b21439b7a9ff', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a2222222-f010-484d-9bbd-e51c4991289f', 'qa', 'Average search time in a B-Tree?', 'O(log n)', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '847e3a18-c7a4-4986-b991-b21439b7a9ff', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a2222222-f010-484d-9bbd-e51c4991289f', 'qa', 'What happens when a B-Tree node exceeds max capacity?', 'It splits into two nodes and promotes the median to the parent.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '847e3a18-c7a4-4986-b991-b21439b7a9ff', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a2222222-f010-484d-9bbd-e51c4991289f', 'qa', 'Is a B-Tree always balanced?', 'Yes, by definition (all leaves are at the same depth).', 'manual');

-- 3. Intro to Software Engineering: Design Patterns
INSERT INTO topics (id, school_id, course_id, user_id, name, exam_weight, source)
VALUES ('a3333333-f010-484d-9bbd-e51c4991289f', '00000000-0000-0000-0000-000000000001', 'b45b4a4a-f123-479b-8cfe-ecfbc26546e0', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'Design Patterns', 0.25, 'manual')
ON CONFLICT (id) DO NOTHING;

INSERT INTO note_texts (id, school_id, course_id, user_id, topic_id, title, content)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'b45b4a4a-f123-479b-8cfe-ecfbc26546e0', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a3333333-f010-484d-9bbd-e51c4991289f', 'Observer and Strategy Patterns', 'Observer pattern establishes a 1-to-many dependency so objects change when state changes. Strategy pattern allows switching algorithms at runtime via composition rather than inheritance.')
ON CONFLICT DO NOTHING;

-- Flashcards for SE
INSERT INTO flashcards (id, school_id, course_id, user_id, topic_id, card_type, prompt, answer, created_by)
VALUES 
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'b45b4a4a-f123-479b-8cfe-ecfbc26546e0', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a3333333-f010-484d-9bbd-e51c4991289f', 'qa', 'Difference between Observer and Pub/Sub?', 'Observer typically synchronous/coupled. Pub/Sub uses a Message Broker (asynch/decoupled).', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'b45b4a4a-f123-479b-8cfe-ecfbc26546e0', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a3333333-f010-484d-9bbd-e51c4991289f', 'qa', 'Core principle of the Strategy Pattern?', 'Favor composition over inheritance.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'b45b4a4a-f123-479b-8cfe-ecfbc26546e0', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a3333333-f010-484d-9bbd-e51c4991289f', 'qa', 'What are the 3 categories of GoF patterns?', 'Creational, Structural, Behavioral.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'b45b4a4a-f123-479b-8cfe-ecfbc26546e0', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a3333333-f010-484d-9bbd-e51c4991289f', 'qa', 'Which pattern is used when an object should appear as another type?', 'Adapter Pattern.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'b45b4a4a-f123-479b-8cfe-ecfbc26546e0', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a3333333-f010-484d-9bbd-e51c4991289f', 'qa', 'What is the Singleton Pattern?', 'Ensures a class has only one instance and provides a global access point.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', 'b45b4a4a-f123-479b-8cfe-ecfbc26546e0', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a3333333-f010-484d-9bbd-e51c4991289f', 'qa', 'Pattern for undo/redo functionality?', 'Command Pattern.', 'manual');

-- 4. African History: The Kingdom of Kush
INSERT INTO topics (id, school_id, course_id, user_id, name, exam_weight, source)
VALUES ('a4444444-f010-484d-9bbd-e51c4991289f', '00000000-0000-0000-0000-000000000001', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'Kingdom of Kush', 0.10, 'manual')
ON CONFLICT (id) DO NOTHING;

INSERT INTO note_texts (id, school_id, course_id, user_id, topic_id, title, content)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a4444444-f010-484d-9bbd-e51c4991289f', 'The Rise of Meroë', 'The Kingdom of Kush, once a powerful rival to Egypt, established its capital at Meroë. Known for iron-smelting and steep pyramids, Meroitic culture combined local and Egyptian influences. Queen Mothers, known as Kandakes, played strong political roles.')
ON CONFLICT DO NOTHING;

-- Flashcards for History
INSERT INTO flashcards (id, school_id, course_id, user_id, topic_id, card_type, prompt, answer, created_by)
VALUES 
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a4444444-f010-484d-9bbd-e51c4991289f', 'qa', 'Where was the Kingdom of Kush located?', 'Ancient Nubia (modern-day Sudan).', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a4444444-f010-484d-9bbd-e51c4991289f', 'qa', 'What was the title of the powerful Queen Mothers in Kush?', 'Kandake (Candace).', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a4444444-f010-484d-9bbd-e51c4991289f', 'qa', 'Kush oversaw Egypt as the ___ Dynasty?', '25th Dynasty.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a4444444-f010-484d-9bbd-e51c4991289f', 'qa', 'Capital city known for iron-working?', 'Meroë.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a4444444-f010-484d-9bbd-e51c4991289f', 'qa', 'River vital to the Kingdom?', 'Nile River.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '6015b5e8-f221-4dfe-bd5f-bbdc6749d82f', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a4444444-f010-484d-9bbd-e51c4991289f', 'qa', 'Language group of the Kushites?', 'Nilo-Saharan.', 'manual');

-- 5. Composition II: Rhetorical Analysis
INSERT INTO topics (id, school_id, course_id, user_id, name, exam_weight, source)
VALUES ('a5555555-f010-484d-9bbd-e51c4991289f', '00000000-0000-0000-0000-000000000001', '61f43c8a-154e-4af3-9e77-81a25cf30c1a', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'Rhetorical Analysis', 0.10, 'manual')
ON CONFLICT (id) DO NOTHING;

INSERT INTO note_texts (id, school_id, course_id, user_id, topic_id, title, content)
VALUES (gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '61f43c8a-154e-4af3-9e77-81a25cf30c1a', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a5555555-f010-484d-9bbd-e51c4991289f', 'Aristotle s Appeals: Pathos, Ethos, Logos', 'Rhetorical triangle consists of Ethos (credibility), Pathos (emotion), and Logos (logic). Analysis focuses on how an author uses these appeals to persuade an audience within a specific context (Kairos).')
ON CONFLICT DO NOTHING;

-- Flashcards for Composition
INSERT INTO flashcards (id, school_id, course_id, user_id, topic_id, card_type, prompt, answer, created_by)
VALUES 
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '61f43c8a-154e-4af3-9e77-81a25cf30c1a', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a5555555-f010-484d-9bbd-e51c4991289f', 'qa', 'What does Ethos appeal to?', 'The credibility or character of the speaker.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '61f43c8a-154e-4af3-9e77-81a25cf30c1a', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a5555555-f010-484d-9bbd-e51c4991289f', 'qa', 'What does Pathos appeal to?', 'The audience s emotions.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '61f43c8a-154e-4af3-9e77-81a25cf30c1a', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a5555555-f010-484d-9bbd-e51c4991289f', 'qa', 'What does Logos appeal to?', 'Reason, facts, and logic.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '61f43c8a-154e-4af3-9e77-81a25cf30c1a', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a5555555-f010-484d-9bbd-e51c4991289f', 'qa', 'Define Kairos.', 'The opportunistic or right moment for a message.', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '61f43c8a-154e-4af3-9e77-81a25cf30c1a', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a5555555-f010-484d-9bbd-e51c4991289f', 'qa', 'What is the Rhetorical Situation?', 'The context of a rhetorical act (Audience, Occasion, Tone).', 'manual'),
(gen_random_uuid(), '00000000-0000-0000-0000-000000000001', '61f43c8a-154e-4af3-9e77-81a25cf30c1a', '1ceeff29-f010-484d-9bbd-e51c4991289f', 'a5555555-f010-484d-9bbd-e51c4991289f', 'qa', 'Who defined the 3 primary appeals?', 'Aristotle.', 'manual');

-- Initialize scheduler states for new cards
INSERT INTO scheduler_states (flashcard_id, user_id, school_id, ease_factor, interval_days, due_at)
SELECT id, user_id, school_id, 2.5, 1, now()
FROM flashcards 
WHERE user_id = '1ceeff29-f010-484d-9bbd-e51c4991289f'
AND id NOT IN (SELECT flashcard_id FROM scheduler_states);

COMMIT;
