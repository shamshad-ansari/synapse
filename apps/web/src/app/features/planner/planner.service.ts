import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../../environments/environment';

// ── Backend response types (snake_case from JSON) ──────────────

export interface StudySession {
  id: string;
  title: string;
  scheduled_date: string; // YYYY-MM-DDT00:00:00Z
  start_time: string;     // "09:00"
  duration_minutes: number;
  status: 'planned' | 'done' | 'missed';
  created_at: string;
  updated_at: string;
}

export interface UpcomingDeadline {
  id: string;
  name: string;
  course_name: string;
  due_date: string;
  days_until: number;
  source: 'lms' | 'manual';
  urgency: 'urgent' | 'soon' | 'safe';
}

// ── Request body types ──────────────────────────────────────────

export interface CreateSessionBody {
  title: string;
  scheduled_date: string;
  start_time: string;
  duration_minutes: number;
}

export interface CreateDeadlineBody {
  name: string;
  course_name?: string;
  due_date: string;
}

@Injectable({ providedIn: 'root' })
export class PlannerService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = environment.apiUrl;

  // ── Signals ─────────────────────────────────────────────────
  readonly sessions = signal<StudySession[]>([]);
  readonly deadlines = signal<UpcomingDeadline[]>([]);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  // ── Sessions ────────────────────────────────────────────────

  async loadSessions(start: string, end: string): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: StudySession[] }>(
          `${this.apiUrl}/v1/planner/sessions?start=${start}&end=${end}`
        ),
      );
      this.sessions.set(res.data ?? []);
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load sessions');
    } finally {
      this.loading.set(false);
    }
  }

  async createSession(body: CreateSessionBody): Promise<void> {
    try {
      await firstValueFrom(
        this.http.post<{ data: StudySession }>(
          `${this.apiUrl}/v1/planner/sessions`, body
        ),
      );
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to create session');
    }
  }

  async updateSessionStatus(id: string, status: string): Promise<void> {
    try {
      await firstValueFrom(
        this.http.patch(
          `${this.apiUrl}/v1/planner/sessions/${id}/status`,
          { status }
        ),
      );
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to update session');
    }
  }

  async deleteSession(id: string): Promise<void> {
    try {
      await firstValueFrom(
        this.http.delete(`${this.apiUrl}/v1/planner/sessions/${id}`),
      );
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to delete session');
    }
  }

  // ── Missed yesterday ───────────────────────────────────────

  async markMissedYesterday(): Promise<number> {
    try {
      const res = await firstValueFrom(
        this.http.post<{ data: { marked: number } }>(
          `${this.apiUrl}/v1/planner/missed-yesterday`, {}
        ),
      );
      return res.data?.marked ?? 0;
    } catch {
      return 0;
    }
  }

  // ── Deadlines ──────────────────────────────────────────────

  async loadDeadlines(limit = 10): Promise<void> {
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: UpcomingDeadline[] }>(
          `${this.apiUrl}/v1/planner/deadlines?limit=${limit}`
        ),
      );
      this.deadlines.set(res.data ?? []);
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load deadlines');
    }
  }

  async createDeadline(body: CreateDeadlineBody): Promise<void> {
    try {
      await firstValueFrom(
        this.http.post(`${this.apiUrl}/v1/planner/deadlines`, body),
      );
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to create deadline');
    }
  }

  async deleteDeadline(id: string): Promise<void> {
    try {
      await firstValueFrom(
        this.http.delete(`${this.apiUrl}/v1/planner/deadlines/${id}`),
      );
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to delete deadline');
    }
  }

  // ── Regenerate plan ────────────────────────────────────────

  async regeneratePlan(): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      await firstValueFrom(
        this.http.post<{ data: StudySession[] }>(
          `${this.apiUrl}/v1/planner/regenerate`, {}
        ),
      );
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to regenerate plan');
    } finally {
      this.loading.set(false);
    }
  }
}
