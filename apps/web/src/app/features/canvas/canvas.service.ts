import { Injectable, inject, signal } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../../environments/environment';
import { AuthService } from '../../core/auth/auth.service';

export interface LmsConnection {
  id: string;
  user_id: string;
  school_id: string;
  lms_type: string;
  institution_url: string;
  token_expires_at: string;
  last_synced_at: string | null;
  sync_status: string;
  created_at: string;
}

export interface LmsSyncedCourse {
  lms_course_id: string;
  course_name: string;
  term: string;
  last_synced_at: string | null;
}

@Injectable({ providedIn: 'root' })
export class CanvasService {
  private readonly http = inject(HttpClient);
  private readonly authService = inject(AuthService);
  private readonly lmsApiUrl = environment.lmsApiUrl;

  readonly status = signal<LmsConnection | null>(null);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  async loadStatus(): Promise<void> {
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: LmsConnection }>(`${this.lmsApiUrl}/v1/lms/status`),
      );
      this.status.set(res.data);
    } catch (err: any) {
      if (err?.status === 404) {
        this.status.set(null);
      } else {
        const msg = err?.error?.error ?? err?.message ?? 'Failed to load LMS status';
        this.error.set(msg);
      }
    }
  }

  connectCanvas(institutionUrl: string): void {
    const token = this.authService.getAccessToken();
    if (!token) {
      this.error.set('Session expired. Log in again, then connect Canvas.');
      return;
    }
    window.location.href =
      `${this.lmsApiUrl}/v1/lms/connect/canvas` +
      `?institution_url=${encodeURIComponent(institutionUrl)}` +
      `&token=${encodeURIComponent(token)}`;
  }

  async connectToken(institutionUrl: string, accessToken: string): Promise<LmsConnection> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.post<{ data: LmsConnection }>(`${this.lmsApiUrl}/v1/lms/connect/token`, {
          institution_url: institutionUrl,
          access_token: accessToken,
        }),
      );
      this.status.set(res.data);
      return res.data;
    } catch (err: any) {
      const msg = this.extractError(err);
      this.error.set(msg);
      throw err;
    } finally {
      this.loading.set(false);
    }
  }

  async disconnect(): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      await firstValueFrom(
        this.http.delete(`${this.lmsApiUrl}/v1/lms/disconnect`),
      );
      this.status.set(null);
    } catch (err: any) {
      const msg = this.extractError(err);
      this.error.set(msg);
      throw err;
    } finally {
      this.loading.set(false);
    }
  }

  async triggerSync(): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      await firstValueFrom(
        this.http.post(`${this.lmsApiUrl}/v1/lms/sync`, {}),
      );
    } catch (err: any) {
      const msg = this.extractError(err);
      this.error.set(msg);
      throw err;
    } finally {
      this.loading.set(false);
    }
  }

  async listSyncedCourses(): Promise<LmsSyncedCourse[]> {
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: LmsSyncedCourse[] }>(`${this.lmsApiUrl}/v1/lms/courses`),
      );
      return Array.isArray(res.data) ? res.data : [];
    } catch (err: any) {
      const msg = this.extractError(err);
      this.error.set(msg);
      throw err;
    }
  }

  private extractError(err: HttpErrorResponse | any): string {
    return err?.error?.error ?? err?.message ?? 'Something went wrong';
  }
}
