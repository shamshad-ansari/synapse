import { Injectable, inject, signal } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface ProfileRepBreakdownRow {
  label: string;
  value: number;
  percent: number;
  color: string;
}

export interface ProfileMasteryRow {
  name: string;
  mastery: number;
  band: 'green' | 'amber' | 'red';
}

export interface ProfileContributionRow {
  icon: string;
  icon_color: string;
  icon_bg: string;
  title: string;
  subtitle: string;
  value: string;
  value_color: string;
}

export interface ProfileSummary {
  reputation_total: number;
  cards_shared: number;
  sessions_done: number;
  reputation_breakdown: ProfileRepBreakdownRow[];
  mastery: ProfileMasteryRow[];
  contributions: ProfileContributionRow[];
}

@Injectable({ providedIn: 'root' })
export class ProfileService {
  private readonly http = inject(HttpClient);
  private readonly apiUrl = environment.apiUrl;

  readonly summary = signal<ProfileSummary | null>(null);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  async loadSummary(): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: ProfileSummary }>(`${this.apiUrl}/v1/profile/summary`),
      );
      if (res?.data == null) {
        this.error.set('Invalid profile response from server.');
        this.summary.set(null);
        return;
      }
      this.summary.set(res.data);
    } catch (err: unknown) {
      this.error.set(this.parseHttpError(err));
      this.summary.set(null);
    } finally {
      this.loading.set(false);
    }
  }

  private parseHttpError(err: unknown): string {
    if (err instanceof HttpErrorResponse) {
      const body = err.error;
      if (body && typeof body === 'object' && body !== null && 'error' in body) {
        const msg = (body as { error?: string | null }).error;
        if (typeof msg === 'string' && msg.length > 0) {
          return msg;
        }
      }
      if (err.status === 0) {
        return 'Cannot reach the API. Run the gateway (e.g. make up) and use ng serve so /v1 is proxied to port 8080.';
      }
      if (typeof err.message === 'string' && err.message.length > 0) {
        return err.message;
      }
    }
    return 'Failed to load profile';
  }
}
