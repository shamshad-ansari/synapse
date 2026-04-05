import { Injectable, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../../environments/environment';
import { AuthService } from '../../core/auth/auth.service';

export interface FeedPost {
  id: string;
  school_id: string;
  user_id: string;
  course_id: string | null;
  topic_id: string | null;
  post_type: string;
  title: string;
  body: string;
  upvotes: number;
  upvoted: boolean;
  author_name: string;
  created_at: string;
  comment_count: number;
}

export interface CreatePostBody {
  title: string;
  body: string;
  post_type: string;
  course_id?: string;
  topic_id?: string;
}

export interface FeedComment {
  id: string;
  post_id: string;
  parent_id: string | null;
  user_id: string;
  author_name: string;
  body: string;
  created_at: string;
}

const PAGE_SIZE = 20;

function unwrapData<T>(raw: unknown): T | null {
  if (raw === null || raw === undefined) {
    return null;
  }
  if (typeof raw === 'object' && raw !== null && 'data' in raw) {
    return (raw as { data: T }).data;
  }
  return raw as T;
}

@Injectable({ providedIn: 'root' })
export class FeedService {
  private readonly http = inject(HttpClient);
  private readonly authService = inject(AuthService);
  private readonly apiUrl = environment.apiUrl;

  /** Offset of the next page to fetch from the API (does not move when a post is prepended locally). */
  private loadedThroughOffset = 0;

  readonly posts = signal<FeedPost[]>([]);
  readonly loading = signal(false);
  readonly loadingMore = signal(false);
  readonly hasMore = signal(true);
  readonly error = signal<string | null>(null);
  readonly posting = signal(false);
  readonly commentsByPost = signal<Record<string, FeedComment[]>>({});
  readonly commentLoading = signal<Record<string, boolean>>({});

  private applyPage(batch: FeedPost[]): void {
    this.posts.set(batch);
    this.loadedThroughOffset = batch.length;
    this.hasMore.set(batch.length === PAGE_SIZE);
  }

  private async refreshFromServer(): Promise<void> {
    const batch = await this.fetchPage(0);
    this.applyPage(batch);
  }

  async loadInitial(): Promise<void> {
    this.loading.set(true);
    this.loadingMore.set(false);
    this.error.set(null);
    this.loadedThroughOffset = 0;
    this.hasMore.set(true);
    try {
      const batch = await this.fetchPage(0);
      this.applyPage(batch);
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load feed');
    } finally {
      this.loading.set(false);
    }
  }

  async loadMore(): Promise<void> {
    if (!this.hasMore() || this.loadingMore() || this.loading()) {
      return;
    }
    this.loadingMore.set(true);
    this.error.set(null);
    try {
      const batch = await this.fetchPage(this.loadedThroughOffset);
      if (batch.length === 0) {
        this.hasMore.set(false);
        return;
      }
      this.posts.update((list) => {
        const seen = new Set(list.map((p) => p.id));
        const merged = [...list];
        for (const p of batch) {
          if (!seen.has(p.id)) {
            seen.add(p.id);
            merged.push(p);
          }
        }
        return merged;
      });
      this.loadedThroughOffset += batch.length;
      this.hasMore.set(batch.length === PAGE_SIZE);
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load feed');
    } finally {
      this.loadingMore.set(false);
    }
  }

  private async fetchPage(offset: number): Promise<FeedPost[]> {
    const res = await firstValueFrom(
      this.http.get<{ data?: unknown } | unknown[]>(
        `${this.apiUrl}/v1/feed?limit=${PAGE_SIZE}&offset=${offset}`,
      ),
    );
    const raw = res as { data?: unknown } | unknown[];
    const arr = Array.isArray((raw as { data?: unknown }).data)
      ? (raw as { data: unknown[] }).data
      : Array.isArray(raw)
        ? raw
        : [];
    return arr.map((x) => this.normalize(x as Record<string, unknown>));
  }

  async createPost(body: CreatePostBody): Promise<void> {
    this.posting.set(true);
    try {
      await firstValueFrom(this.http.post(`${this.apiUrl}/v1/feed`, body));
      await this.refreshFromServer();
    } finally {
      this.posting.set(false);
    }
  }

  async toggleUpvote(postId: string): Promise<void> {
    try {
      await firstValueFrom(this.http.post(`${this.apiUrl}/v1/feed/${postId}/upvote`, {}));
      await this.refreshFromServer();
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to update post');
    }
  }

  async loadComments(postId: string): Promise<void> {
    this.commentLoading.update((m) => ({ ...m, [postId]: true }));
    try {
      const res = await firstValueFrom(
        this.http.get<{ data?: FeedComment[] } | FeedComment[]>(`${this.apiUrl}/v1/feed/${postId}/comments`),
      );
      const unwrapped = unwrapData<FeedComment[]>(res);
      const arr = Array.isArray(unwrapped)
        ? unwrapped
        : Array.isArray(res)
          ? res
          : [];
      this.commentsByPost.update((m) => ({ ...m, [postId]: arr }));
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load comments');
    } finally {
      this.commentLoading.update((m) => ({ ...m, [postId]: false }));
    }
  }

  async createComment(postId: string, body: string, parentId?: string): Promise<FeedComment> {
    const payload: Record<string, unknown> = { body };
    if (parentId) payload['parent_id'] = parentId;
    const res = await firstValueFrom(
      this.http.post<{ data?: FeedComment } | FeedComment>(`${this.apiUrl}/v1/feed/${postId}/comments`, payload),
    );
    const rawCreated = unwrapData<FeedComment>(res) ?? (res as FeedComment);
    const created = {
      ...rawCreated,
      author_name: rawCreated.author_name || this.authService.currentUser()?.name || 'You',
    };
    this.commentsByPost.update((m) => ({
      ...m,
      [postId]: [...(m[postId] ?? []), created],
    }));
    this.posts.update((list) =>
      list.map((p) =>
        p.id === postId ? { ...p, comment_count: (p.comment_count ?? 0) + 1 } : p,
      ),
    );
    return created;
  }

  async deletePost(postId: string): Promise<void> {
    await firstValueFrom(this.http.delete(`${this.apiUrl}/v1/feed/${postId}`));
    await this.refreshFromServer();
  }

  private normalize(raw: Record<string, unknown>, defaults: Record<string, unknown> = {}): FeedPost {
    const pick = (...keys: string[]): unknown => {
      for (const k of keys) {
        const v = raw[k] ?? defaults[k];
        if (v !== undefined && v !== null && v !== '') {
          return v;
        }
      }
      return undefined;
    };
    const str = (...keys: string[]) => String(pick(...keys) ?? '');
    const num = (...keys: string[]) => Number(pick(...keys) ?? 0);
    const course = pick('course_id', 'courseId');
    const topic = pick('topic_id', 'topicId');
    return {
      id: str('id'),
      school_id: str('school_id', 'schoolId'),
      user_id: str('user_id', 'userId'),
      course_id: course !== undefined && course !== null && String(course) !== '' ? String(course) : null,
      topic_id: topic !== undefined && topic !== null && String(topic) !== '' ? String(topic) : null,
      post_type: str('post_type', 'postType') || 'question',
      title: str('title'),
      body: str('body'),
      upvotes: num('upvotes'),
      upvoted: Boolean(pick('upvoted')),
      author_name: str('author_name', 'authorName'),
      created_at: str('created_at', 'createdAt'),
      comment_count: Math.max(0, Math.floor(num('comment_count', 'commentCount'))),
    };
  }
}
