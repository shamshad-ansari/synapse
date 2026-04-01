-- RAG: HNSW indexes on embeddings + note_texts for user note content and embeddings

-- flashcards.embedding exists from 000003; ANN index for similarity search
CREATE INDEX IF NOT EXISTS idx_flashcards_embedding
  ON flashcards USING hnsw (embedding vector_cosine_ops)
  WITH (m = 16, ef_construction = 64);

CREATE TABLE IF NOT EXISTS note_texts (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  school_id   UUID NOT NULL,
  course_id   UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
  user_id     UUID NOT NULL REFERENCES users(id),
  topic_id    UUID REFERENCES topics(id),
  title       TEXT NOT NULL DEFAULT '',
  content     TEXT NOT NULL,
  embedding   VECTOR(1536),
  embedded_at TIMESTAMPTZ,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_note_texts_course ON note_texts(course_id, user_id, school_id);
CREATE INDEX IF NOT EXISTS idx_note_texts_embedding
  ON note_texts USING hnsw (embedding vector_cosine_ops)
  WITH (m = 16, ef_construction = 64);
