import { Injectable, inject, signal } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpParams } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface Course {
  id: string;
  name: string;
  term: string;
  color: string;
  lms_course_id: string | null;
  created_at: string;
}

export interface Topic {
  id: string;
  course_id: string;
  name: string;
  exam_weight: number | null;
}

export interface Flashcard {
  id: string;
  course_id: string;
  topic_id: string | null;
  card_type: string;
  prompt: string;
  answer: string;
  created_by: string;
  created_at: string;
}

export interface DueCard {
  flashcard_id: string;
  prompt: string;
  answer: string;
  card_type: string;
  topic_name: string;
  ease_factor: number;
  interval_days: number;
  due_at: string;
}

export interface ReviewResult {
  ease_factor: number;
  interval_days: number;
  due_at: string;
}

export interface GeneratedCardPayload {
  prompt: string;
  answer: string;
  card_type: string;
}

export interface NoteText {
  id: string;
  course_id: string;
  topic_id: string | null;
  title: string;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface ConfusionInsight {
  course_id: string | null;
  course_name: string;
  hotspot_count: number;
  top_topic_name: string;
  confused_events: number;
  window_days: number;
}

export interface NoteMetric {
  note_id: string;
  readiness_pct: number;
  connected_cards: number;
  review_count: number;
  confused_count: number;
  confusion_flag: boolean;
  last_review_at: string | null;
}

function pickStr(raw: Record<string, unknown>, ...keys: string[]): string {
  for (const k of keys) {
    const v = raw[k];
    if (v !== undefined && v !== null) {
      return String(v);
    }
  }
  return '';
}

function pickNum(raw: Record<string, unknown>, ...keys: string[]): number {
  for (const k of keys) {
    const v = raw[k];
    if (typeof v === 'number' && !Number.isNaN(v)) {
      return v;
    }
    if (typeof v === 'string' && v !== '') {
      const n = Number(v);
      if (!Number.isNaN(n)) {
        return n;
      }
    }
  }
  return 0;
}

function pickStrOrNull(raw: Record<string, unknown>, ...keys: string[]): string | null {
  for (const k of keys) {
    const v = raw[k];
    if (v === undefined || v === null) {
      continue;
    }
    return String(v);
  }
  return null;
}

function asRecord(v: unknown): Record<string, unknown> {
  return v !== null && typeof v === 'object' ? (v as Record<string, unknown>) : {};
}

function normalizeCourse(raw: Record<string, unknown>): Course {
  return {
    id: pickStr(raw, 'id', 'ID'),
    name: pickStr(raw, 'name', 'Name'),
    term: pickStr(raw, 'term', 'Term'),
    color: pickStr(raw, 'color', 'Color') || '#102E67',
    lms_course_id: pickStrOrNull(raw, 'lms_course_id', 'LMSCourseID'),
    created_at: pickStr(raw, 'created_at', 'CreatedAt'),
  };
}

function normalizeTopic(raw: Record<string, unknown>): Topic {
  const examWeight = raw['exam_weight'] ?? raw['ExamWeight'];
  let exam_weight: number | null = null;
  if (examWeight !== undefined && examWeight !== null && examWeight !== '') {
    const n = typeof examWeight === 'number' ? examWeight : Number(examWeight);
    exam_weight = Number.isNaN(n) ? null : n;
  }
  return {
    id: pickStr(raw, 'id', 'ID'),
    course_id: pickStr(raw, 'course_id', 'CourseID'),
    name: pickStr(raw, 'name', 'Name'),
    exam_weight,
  };
}

function normalizeFlashcard(raw: Record<string, unknown>): Flashcard {
  const topic = raw['topic_id'] ?? raw['TopicID'];
  return {
    id: pickStr(raw, 'id', 'ID'),
    course_id: pickStr(raw, 'course_id', 'CourseID'),
    topic_id: topic === undefined || topic === null ? null : String(topic),
    card_type: pickStr(raw, 'card_type', 'CardType') || 'qa',
    prompt: pickStr(raw, 'prompt', 'Prompt'),
    answer: pickStr(raw, 'answer', 'Answer'),
    created_by: pickStr(raw, 'created_by', 'CreatedBy'),
    created_at: pickStr(raw, 'created_at', 'CreatedAt'),
  };
}

function normalizeDueCard(raw: Record<string, unknown>): DueCard {
  return {
    flashcard_id: pickStr(raw, 'flashcard_id', 'FlashcardID'),
    prompt: pickStr(raw, 'prompt', 'Prompt'),
    answer: pickStr(raw, 'answer', 'Answer'),
    card_type: pickStr(raw, 'card_type', 'CardType'),
    topic_name: pickStr(raw, 'topic_name', 'TopicName'),
    ease_factor: pickNum(raw, 'ease_factor', 'EaseFactor'),
    interval_days: Math.round(pickNum(raw, 'interval_days', 'IntervalDays')),
    due_at: pickStr(raw, 'due_at', 'DueAt'),
  };
}

function normalizeReviewResult(raw: Record<string, unknown>): ReviewResult {
  return {
    ease_factor: pickNum(raw, 'ease_factor', 'EaseFactor'),
    interval_days: Math.round(pickNum(raw, 'interval_days', 'IntervalDays')),
    due_at: pickStr(raw, 'due_at', 'DueAt'),
  };
}

function normalizeNoteText(raw: Record<string, unknown>): NoteText {
  const topic = raw['topic_id'] ?? raw['TopicID'];
  return {
    id: pickStr(raw, 'id', 'ID'),
    course_id: pickStr(raw, 'course_id', 'CourseID'),
    topic_id: topic === undefined || topic === null ? null : String(topic),
    title: pickStr(raw, 'title', 'Title'),
    content: pickStr(raw, 'content', 'Content'),
    created_at: pickStr(raw, 'created_at', 'CreatedAt'),
    updated_at: pickStr(raw, 'updated_at', 'UpdatedAt'),
  };
}

function normalizeConfusionInsight(raw: Record<string, unknown>): ConfusionInsight {
  return {
    course_id: pickStrOrNull(raw, 'course_id', 'CourseID'),
    course_name: pickStr(raw, 'course_name', 'CourseName'),
    hotspot_count: pickNum(raw, 'hotspot_count', 'HotspotCount'),
    top_topic_name: pickStr(raw, 'top_topic_name', 'TopTopicName'),
    confused_events: pickNum(raw, 'confused_events', 'ConfusedEvents'),
    window_days: pickNum(raw, 'window_days', 'WindowDays'),
  };
}

function normalizeNoteMetric(raw: Record<string, unknown>): NoteMetric {
  const lastReviewRaw = raw['last_review_at'] ?? raw['LastReviewAt'];
  return {
    note_id: pickStr(raw, 'note_id', 'NoteID'),
    readiness_pct: Math.max(0, Math.min(100, Math.round(pickNum(raw, 'readiness_pct', 'ReadinessPct')))),
    connected_cards: Math.max(0, Math.round(pickNum(raw, 'connected_cards', 'ConnectedCards'))),
    review_count: Math.max(0, Math.round(pickNum(raw, 'review_count', 'ReviewCount'))),
    confused_count: Math.max(0, Math.round(pickNum(raw, 'confused_count', 'ConfusedCount'))),
    confusion_flag: Boolean(raw['confusion_flag'] ?? raw['ConfusionFlag']),
    last_review_at: lastReviewRaw === undefined || lastReviewRaw === null ? null : String(lastReviewRaw),
  };
}

@Injectable({ providedIn: 'root' })
export class LearningService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = environment.apiUrl;
  private readonly localNotesByCourse = new Map<string, NoteText[]>();

  readonly courses = signal<Course[]>([]);
  readonly flashcards = signal<Flashcard[]>([]);
  readonly dueCards = signal<DueCard[]>([]);
  readonly topics = signal<Topic[]>([]);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  readonly generatedCards = signal<GeneratedCardPayload[]>([]);
  readonly generating = signal(false);
  readonly generationError = signal<string | null>(null);

  readonly notes = signal<NoteText[]>([]);
  readonly notesLoading = signal(false);
  readonly activeNote = signal<NoteText | null>(null);
  readonly noteSaving = signal(false);
  readonly noteAiLoading = signal(false);
  readonly noteAiAnswer = signal<string | null>(null);
  readonly noteAiError = signal<string | null>(null);
  readonly confusionInsight = signal<ConfusionInsight | null>(null);
  readonly noteMetrics = signal<NoteMetric[]>([]);

  private makeLocalNote(courseId: string, title: string, content: string, topicId?: string): NoteText {
    const now = new Date().toISOString();
    return {
      id: typeof crypto !== 'undefined' && 'randomUUID' in crypto ? crypto.randomUUID() : `${Date.now()}-${Math.random()}`,
      course_id: courseId,
      topic_id: topicId ?? null,
      title,
      content,
      created_at: now,
      updated_at: now,
    };
  }

  private readLocalNotes(courseId: string): NoteText[] {
    return this.localNotesByCourse.get(courseId) ?? [];
  }

  async loadCourses(): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: unknown }>(`${this.apiUrl}/v1/courses`),
      );
      const rawList = res.data;
      const arr = Array.isArray(rawList) ? rawList : [];
      this.courses.set(arr.map((x) => normalizeCourse(asRecord(x))));
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load courses');
      this.courses.set([]);
    } finally {
      this.loading.set(false);
    }
  }

  async loadFlashcards(courseId: string): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: unknown }>(`${this.apiUrl}/v1/courses/${courseId}/flashcards`),
      );
      const rawList = res.data;
      const arr = Array.isArray(rawList) ? rawList : [];
      this.flashcards.set(arr.map((x) => normalizeFlashcard(asRecord(x))));
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load flashcards');
      this.flashcards.set([]);
    } finally {
      this.loading.set(false);
    }
  }

  async loadTopics(courseId: string): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: unknown }>(`${this.apiUrl}/v1/courses/${courseId}/topics`),
      );
      const rawList = res.data;
      const arr = Array.isArray(rawList) ? rawList : [];
      this.topics.set(arr.map((x) => normalizeTopic(asRecord(x))));
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load topics');
      this.topics.set([]);
    } finally {
      this.loading.set(false);
    }
  }

  async loadDueCards(courseId?: string, limit = 20): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      let params = new HttpParams().set('limit', String(limit));
      if (courseId) {
        params = params.set('courseId', courseId);
      }
      const res = await firstValueFrom(
        this.http.get<{ data: unknown }>(`${this.apiUrl}/v1/review/due`, { params }),
      );
      const rawList = res.data;
      const arr = Array.isArray(rawList) ? rawList : [];
      this.dueCards.set(arr.map((x) => normalizeDueCard(asRecord(x))));
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load due cards');
      this.dueCards.set([]);
    } finally {
      this.loading.set(false);
    }
  }

  async createFlashcard(
    courseId: string,
    prompt: string,
    answer: string,
    topicId?: string,
  ): Promise<Flashcard> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const body: Record<string, unknown> = {
        course_id: courseId,
        card_type: 'qa',
        prompt,
        answer,
      };
      if (topicId) {
        body['topic_id'] = topicId;
      }
      const res = await firstValueFrom(
        this.http.post<{ data: unknown }>(`${this.apiUrl}/v1/flashcards`, body),
      );
      return normalizeFlashcard(asRecord(res.data));
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to create flashcard');
      throw err;
    } finally {
      this.loading.set(false);
    }
  }

  async submitReview(
    sessionId: string,
    flashcardId: string,
    correct: boolean,
    confidence: number,
    confused: boolean,
    responseTimeMs: number,
  ): Promise<ReviewResult> {
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.post<{ data: unknown }>(`${this.apiUrl}/v1/review/events`, {
          session_id: sessionId,
          flashcard_id: flashcardId,
          correct,
          confidence,
          confused,
          response_time_ms: responseTimeMs,
        }),
      );
      return normalizeReviewResult(asRecord(res.data));
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to submit review');
      throw err;
    }
  }

  async importFromLMS(
    courses: { lms_course_id: string; name: string; term: string }[],
  ): Promise<Course[]> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.post<{ data: unknown }>(`${this.apiUrl}/v1/courses/import-from-lms`, {
          courses,
        }),
      );
      const rawList = res.data;
      const arr = Array.isArray(rawList) ? rawList : [];
      return arr.map((x) => normalizeCourse(asRecord(x)));
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to import courses');
      throw err;
    } finally {
      this.loading.set(false);
    }
  }

  async generateFlashcards(
    courseId: string,
    noteContent: string,
    topicId?: string,
  ): Promise<void> {
    this.generating.set(true);
    this.generationError.set(null);
    try {
      const body: Record<string, unknown> = {
        course_id: courseId,
        note_content: noteContent,
      };
      if (topicId) {
        body['topic_id'] = topicId;
      }
      const res = await firstValueFrom(
        this.http.post<{ data: { candidates?: GeneratedCardPayload[] } }>(
          `${this.apiUrl}/v1/flashcards/generate`,
          body,
        ),
      );
      const raw = res.data as { candidates?: GeneratedCardPayload[] };
      this.generatedCards.set(raw?.candidates ?? []);
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string; status?: number };
      this.generationError.set(
        e?.error?.error ?? e?.message ?? 'Failed to generate flashcards',
      );
      this.generatedCards.set([]);
    } finally {
      this.generating.set(false);
    }
  }

  async loadNotes(courseId: string): Promise<void> {
    this.notesLoading.set(true);
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: unknown }>(`${this.apiUrl}/v1/courses/${courseId}/notes`),
      );
      const arr = Array.isArray(res.data) ? res.data : [];
      const apiNotes = arr.map((x) => normalizeNoteText(asRecord(x)));
      this.localNotesByCourse.set(courseId, apiNotes);
      this.notes.set(apiNotes);
    } catch {
      this.notes.set(this.readLocalNotes(courseId));
    } finally {
      this.notesLoading.set(false);
    }
    await Promise.all([
      this.loadConfusionInsight(courseId),
      this.loadNoteMetrics(courseId),
    ]);
  }

  async loadConfusionInsight(courseId?: string, windowDays = 14): Promise<void> {
    try {
      let params = new HttpParams().set('windowDays', String(windowDays));
      if (courseId) {
        params = params.set('courseId', courseId);
      }
      const res = await firstValueFrom(
        this.http.get<{ data: unknown }>(`${this.apiUrl}/v1/insights/confusion`, { params }),
      );
      this.confusionInsight.set(normalizeConfusionInsight(asRecord(res.data)));
    } catch {
      this.confusionInsight.set(null);
    }
  }

  async loadNoteMetrics(courseId: string, windowDays = 30): Promise<void> {
    try {
      const params = new HttpParams().set('windowDays', String(windowDays));
      const res = await firstValueFrom(
        this.http.get<{ data: unknown }>(`${this.apiUrl}/v1/courses/${courseId}/notes/metrics`, { params }),
      );
      const arr = Array.isArray(res.data) ? res.data : [];
      this.noteMetrics.set(arr.map((x) => normalizeNoteMetric(asRecord(x))));
    } catch {
      this.noteMetrics.set([]);
    }
  }

  async createNote(courseId: string, title: string, content: string, topicId?: string): Promise<NoteText> {
    try {
      const body: Record<string, unknown> = {
        course_id: courseId,
        title,
        content,
      };
      if (topicId) {
        body['topic_id'] = topicId;
      }
      const res = await firstValueFrom(
        this.http.post<{ data: unknown }>(`${this.apiUrl}/v1/notes`, body),
      );
      const note = normalizeNoteText(asRecord(res.data));
      this.notes.update((list) => [note, ...list]);
      this.localNotesByCourse.set(courseId, this.notes());
      await this.loadNoteMetrics(courseId);
      return note;
    } catch {
      const note = this.makeLocalNote(courseId, title, content, topicId);
      this.notes.update((list) => [note, ...list]);
      this.localNotesByCourse.set(courseId, this.notes());
      return note;
    }
  }

  async updateNote(noteId: string, title: string, content: string, topicId?: string | null): Promise<void> {
    this.noteSaving.set(true);
    try {
      const body: Record<string, unknown> = { title, content };
      if (topicId) {
        body['topic_id'] = topicId;
      } else {
        body['topic_id'] = null;
      }
      const res = await firstValueFrom(
        this.http.put<{ data: unknown }>(`${this.apiUrl}/v1/notes/${noteId}`, body),
      );
      const updated = normalizeNoteText(asRecord(res.data));
      this.notes.update((list) => list.map((n) => (n.id === noteId ? updated : n)));
      if (this.activeNote()?.id === noteId) {
        this.activeNote.set(updated);
      }
      this.localNotesByCourse.set(updated.course_id, this.notes());
      await this.loadNoteMetrics(updated.course_id);
    } catch {
      let courseId = '';
      this.notes.update((list) =>
        list.map((n) => {
          if (n.id !== noteId) {
            return n;
          }
          courseId = n.course_id;
          return {
            ...n,
            title,
            content,
            topic_id: topicId ?? null,
            updated_at: new Date().toISOString(),
          };
        }),
      );
      const updatedLocal = this.notes().find((n) => n.id === noteId) ?? null;
      if (updatedLocal && this.activeNote()?.id === noteId) {
        this.activeNote.set(updatedLocal);
      }
      if (courseId) {
        this.localNotesByCourse.set(courseId, this.notes());
      }
      if (courseId) {
        await this.loadNoteMetrics(courseId);
      }
    } finally {
      this.noteSaving.set(false);
    }
  }

  async deleteNote(noteId: string): Promise<void> {
    let courseId = '';
    try {
      await firstValueFrom(this.http.delete(`${this.apiUrl}/v1/notes/${noteId}`));
    } catch {
      // Continue and delete locally so the editor still works if notes API is unavailable.
    } finally {
      this.notes.update((list) =>
        list.filter((n) => {
          const keep = n.id !== noteId;
          if (!keep) {
            courseId = n.course_id;
          }
          return keep;
        }),
      );
      if (this.activeNote()?.id === noteId) {
        this.activeNote.set(null);
      }
      if (courseId) {
        this.localNotesByCourse.set(courseId, this.notes());
        await this.loadNoteMetrics(courseId);
      }
    }
  }

  async askNoteAI(noteId: string, question: string): Promise<void> {
    this.noteAiLoading.set(true);
    this.noteAiAnswer.set(null);
    this.noteAiError.set(null);
    try {
      const res = await firstValueFrom(
        this.http.post<{ data: { answer: string } }>(`${this.apiUrl}/v1/notes/${noteId}/ask`, {
          question,
        }),
      );
      this.noteAiAnswer.set(res.data?.answer ?? '');
    } catch (err: unknown) {
      const e = err as HttpErrorResponse & {
        error?: { error?: string };
      };
      this.noteAiError.set(
        e?.error?.error ?? e?.message ?? 'Failed to get AI response',
      );
      this.noteAiAnswer.set(null);
    } finally {
      this.noteAiLoading.set(false);
    }
  }

  async acceptGeneratedCards(
    courseId: string,
    cards: GeneratedCardPayload[],
    topicId?: string,
  ): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const body: Record<string, unknown> = {
        course_id: courseId,
        cards: cards.map((c) => ({
          prompt: c.prompt,
          answer: c.answer,
          card_type: c.card_type || 'qa',
        })),
      };
      if (topicId) {
        body['topic_id'] = topicId;
      }
      await firstValueFrom(
        this.http.post<{ data: unknown }>(`${this.apiUrl}/v1/flashcards/generate/accept`, body),
      );
      await this.loadFlashcards(courseId);
      this.generatedCards.set([]);
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to save flashcards');
      throw err;
    } finally {
      this.loading.set(false);
    }
  }
}
