import { Injectable, inject, signal } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface TutorRequest {
  id: string;
  requester_id: string;
  tutor_id: string;
  topic_name: string;
  status: string;
  message: string;
  requester_name: string;
  tutor_name: string;
  created_at: string;
}

export interface TutorMatch {
  user_id: string;
  name: string;
  topic_name: string;
  mastery: number;
  sessions: number;
  reputation: number;
}

export interface TeachingTopic {
  topic_id: string;
  topic_name: string;
  mastery: number;
}

export interface CreateRequestBody {
  tutor_id: string;
  topic_name: string;
  message: string;
  topic_id?: string;
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

function asRecord(v: unknown): Record<string, unknown> {
  return v !== null && typeof v === 'object' ? (v as Record<string, unknown>) : {};
}

function normalizeTutorRequest(raw: Record<string, unknown>): TutorRequest {
  return {
    id: pickStr(raw, 'id', 'ID'),
    requester_id: pickStr(raw, 'requester_id', 'RequesterID'),
    tutor_id: pickStr(raw, 'tutor_id', 'TutorID'),
    topic_name: pickStr(raw, 'topic_name', 'TopicName'),
    status: pickStr(raw, 'status', 'Status'),
    message: pickStr(raw, 'message', 'Message'),
    requester_name: pickStr(raw, 'requester_name', 'RequesterName'),
    tutor_name: pickStr(raw, 'tutor_name', 'TutorName'),
    created_at: pickStr(raw, 'created_at', 'CreatedAt'),
  };
}

function normalizeTutorMatch(raw: Record<string, unknown>): TutorMatch {
  return {
    user_id: pickStr(raw, 'user_id', 'UserID'),
    name: pickStr(raw, 'name', 'Name'),
    topic_name: pickStr(raw, 'topic_name', 'TopicName'),
    mastery: pickNum(raw, 'mastery', 'Mastery'),
    sessions: Math.round(pickNum(raw, 'sessions', 'Sessions')),
    reputation: Math.round(pickNum(raw, 'reputation', 'Reputation')),
  };
}

function normalizeTeachingTopic(raw: Record<string, unknown>): TeachingTopic {
  return {
    topic_id: pickStr(raw, 'topic_id', 'TopicID'),
    topic_name: pickStr(raw, 'topic_name', 'TopicName'),
    mastery: pickNum(raw, 'mastery', 'Mastery'),
  };
}

export type RequestStatusFilter = 'pending' | 'accepted' | 'declined' | 'completed' | 'cancelled' | 'all';

@Injectable({ providedIn: 'root' })
export class TutoringService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = environment.apiUrl;

  readonly incomingRequests = signal<TutorRequest[]>([]);
  readonly incomingAccepted = signal<TutorRequest[]>([]);
  readonly outgoingRequests = signal<TutorRequest[]>([]);
  readonly tutorMatches = signal<TutorMatch[]>([]);
  readonly teachingTopics = signal<TeachingTopic[]>([]);

  /** True while initial / tab refresh loads are running */
  readonly listsLoading = signal(false);
  /** True while tutor match query runs */
  readonly matchesLoading = signal(false);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  async loadTeachingTopics(): Promise<void> {
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: unknown }>(`${this.apiUrl}/v1/tutoring/teaching-topics`),
      );
      const arr = Array.isArray(res.data) ? res.data : [];
      this.teachingTopics.set(arr.map((x) => normalizeTeachingTopic(asRecord(x))));
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load teaching topics');
      this.teachingTopics.set([]);
    }
  }

  async loadIncomingRequests(status: RequestStatusFilter = 'pending'): Promise<void> {
    this.error.set(null);
    try {
      const params = new HttpParams().set('status', status);
      const res = await firstValueFrom(
        this.http.get<{ data: unknown }>(`${this.apiUrl}/v1/tutoring/requests/incoming`, { params }),
      );
      const arr = Array.isArray(res.data) ? res.data : [];
      const mapped = arr.map((x) => normalizeTutorRequest(asRecord(x)));
      if (status === 'accepted') {
        this.incomingAccepted.set(mapped);
      } else {
        this.incomingRequests.set(mapped);
      }
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load incoming requests');
      if (status === 'accepted') {
        this.incomingAccepted.set([]);
      } else {
        this.incomingRequests.set([]);
      }
    }
  }

  async loadOutgoingRequests(status: RequestStatusFilter = 'all'): Promise<void> {
    this.error.set(null);
    try {
      const params = new HttpParams().set('status', status);
      const res = await firstValueFrom(
        this.http.get<{ data: unknown }>(`${this.apiUrl}/v1/tutoring/requests/outgoing`, { params }),
      );
      const arr = Array.isArray(res.data) ? res.data : [];
      this.outgoingRequests.set(arr.map((x) => normalizeTutorRequest(asRecord(x))));
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load outgoing requests');
      this.outgoingRequests.set([]);
    }
  }

  async loadTutorMatches(topicName: string, limit = 8): Promise<void> {
    const topic = topicName.trim();
    if (!topic) {
      this.tutorMatches.set([]);
      return;
    }
    this.matchesLoading.set(true);
    this.error.set(null);
    try {
      const params = new HttpParams().set('topic', topic).set('limit', String(limit));
      const res = await firstValueFrom(
        this.http.get<{ data: unknown }>(`${this.apiUrl}/v1/tutoring/match`, { params }),
      );
      const arr = Array.isArray(res.data) ? res.data : [];
      this.tutorMatches.set(arr.map((x) => normalizeTutorMatch(asRecord(x))));
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load tutor matches');
      this.tutorMatches.set([]);
    } finally {
      this.matchesLoading.set(false);
    }
  }

  /** Loads teaching topics, request queues, and tutor matches for the weak-topic string (or fallback). */
  async refreshAll(weakTopicHint: string | null | undefined): Promise<void> {
    this.listsLoading.set(true);
    this.error.set(null);
    try {
      await Promise.all([
        this.loadTeachingTopics(),
        this.loadIncomingRequests('pending'),
        this.loadIncomingRequests('accepted'),
        this.loadOutgoingRequests('all'),
      ]);
      const hint = (weakTopicHint ?? '').trim();
      await this.loadTutorMatches(hint);
    } finally {
      this.listsLoading.set(false);
    }
  }

  async createRequest(body: CreateRequestBody): Promise<TutorRequest> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const payload: Record<string, unknown> = {
        tutor_id: body.tutor_id ?? '',
        topic_name: body.topic_name,
        message: body.message,
      };
      if (body.topic_id) {
        payload['topic_id'] = body.topic_id;
      }
      const res = await firstValueFrom(
        this.http.post<{ data: unknown }>(`${this.apiUrl}/v1/tutoring/requests`, payload),
      );
      const created = normalizeTutorRequest(asRecord(res.data));
      this.outgoingRequests.update((list) => [created, ...list]);
      return created;
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to create request');
      throw err;
    } finally {
      this.loading.set(false);
    }
  }

  async respondToRequest(requestId: string, status: 'accepted' | 'declined'): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      await firstValueFrom(
        this.http.patch<{ data: unknown }>(`${this.apiUrl}/v1/tutoring/requests/${requestId}`, {
          status,
        }),
      );
      this.incomingRequests.update((list) => list.filter((r) => r.id !== requestId));
      if (status === 'accepted') {
        await this.loadIncomingRequests('accepted');
      }
    } finally {
      this.loading.set(false);
    }
  }

  outgoingPending(): TutorRequest[] {
    return this.outgoingRequests().filter((r) => r.status === 'pending');
  }

  outgoingAccepted(): TutorRequest[] {
    return this.outgoingRequests().filter((r) => r.status === 'accepted');
  }
}
