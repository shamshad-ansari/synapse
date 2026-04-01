import { Injectable, inject, signal, computed } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Router } from '@angular/router';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface User {
  id: string;
  school_id: string;
  name: string;
  email: string;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

const ACCESS_TOKEN_KEY = 'synapse_access_token';
const REFRESH_TOKEN_KEY = 'synapse_refresh_token';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);
  private readonly apiUrl = environment.apiUrl;

  readonly currentUser = signal<User | null>(null);
  readonly isAuthenticated = computed(() => this.currentUser() !== null);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  async login(email: string, password: string, schoolDomain: string): Promise<User> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.post<{ data: AuthResponse }>(`${this.apiUrl}/v1/auth/login`, {
          email,
          password,
          school_domain: schoolDomain,
        }),
      );
      this.storeTokens(res.data.access_token, res.data.refresh_token);
      this.currentUser.set(res.data.user);
      return res.data.user;
    } catch (err: unknown) {
      this.error.set(this.apiErrorMessage(err, 'Login failed'));
      throw err;
    } finally {
      this.loading.set(false);
    }
  }

  async register(name: string, email: string, password: string, schoolDomain: string): Promise<User> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.post<{ data: AuthResponse }>(`${this.apiUrl}/v1/auth/register`, {
          name,
          email,
          password,
          school_domain: schoolDomain,
        }),
      );
      this.storeTokens(res.data.access_token, res.data.refresh_token);
      this.currentUser.set(res.data.user);
      return res.data.user;
    } catch (err: unknown) {
      this.error.set(this.apiErrorMessage(err, 'Registration failed'));
      throw err;
    } finally {
      this.loading.set(false);
    }
  }

  logout(): void {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    this.currentUser.set(null);
    this.router.navigate(['/login']);
  }

  async loadCurrentUser(): Promise<void> {
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: User }>(`${this.apiUrl}/v1/me`),
      );
      this.currentUser.set(res.data);
    } catch (err: any) {
      if (err?.status === 401) {
        this.logout();
      } else {
        this.currentUser.set(null);
      }
    }
  }

  getAccessToken(): string | null {
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  }

  async initAuth(): Promise<void> {
    if (this.getAccessToken()) {
      await this.loadCurrentUser();
    }
  }

  private storeTokens(access: string, refresh: string): void {
    localStorage.setItem(ACCESS_TOKEN_KEY, access);
    localStorage.setItem(REFRESH_TOKEN_KEY, refresh);
  }

  /** Parses api-gateway envelope { data, error } or network errors. */
  private apiErrorMessage(err: unknown, fallback: string): string {
    if (err instanceof HttpErrorResponse) {
      const body = err.error;
      if (body && typeof body === 'object' && 'error' in body) {
        const e = (body as { error?: string | null }).error;
        if (typeof e === 'string' && e.length > 0) {
          return e;
        }
      }
      if (err.status === 0) {
        return 'Cannot reach the server. Is the API running and is apiUrl correct in environment?';
      }
      if (typeof err.message === 'string' && err.message.length > 0) {
        return err.message;
      }
    }
    return fallback;
  }
}
