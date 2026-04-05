import {
  Component,
  OnInit,
  OnDestroy,
  inject,
  NgZone,
  signal,
  computed,
  viewChild,
  ElementRef,
  afterNextRender,
} from '@angular/core';
import { LucideAngularModule } from 'lucide-angular';
import { AuthService } from '../core/auth/auth.service';
import { LearningService } from '../features/learning/learning.service';
import { FeedService, type FeedComment, type FeedPost } from '../features/feed/feed.service';

@Component({
  selector: 'app-feed-tab',
  standalone: true,
  imports: [LucideAngularModule],
  template: `
    <div
      #feedScroll
      class="flex flex-col flex-1 min-h-0 overflow-y-auto overflow-x-hidden"
      style="padding: 44px 56px 56px; gap: 36px"
    >
      <div class="flex items-center justify-between gap-4">
        <div>
          <div
            style="font-size: 26px; font-weight: 700; letter-spacing: -0.6px; font-family: var(--font-display); color: var(--ink)"
          >
            Feed
          </div>
          <div style="font-size: 13.5px; color: var(--ink-muted); margin-top: 4px">
            Cohort posts with study-aware context
          </div>
        </div>
        <button
          type="button"
          class="post-btn flex items-center gap-2 transition-all"
          style="font-size: 13px; padding: 8px 16px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 600; cursor: pointer; box-shadow: var(--shadow-sm); transition: var(--transition-base)"
          (click)="showPostForm.update((v) => !v)"
        >
          <lucide-icon name="plus" [size]="15" [strokeWidth]="2" /> Post
        </button>
      </div>

      @if (showPostForm()) {
        <div
          style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 20px 22px; background: var(--card-bg); box-shadow: var(--shadow-sm)"
        >
          <div style="font-size: 13px; font-weight: 600; color: var(--ink); margin-bottom: 12px; font-family: var(--font-display)">
            New post
          </div>
          <label class="block" style="margin-bottom: 10px">
            <span style="font-size: 11px; font-weight: 600; color: var(--ink-muted); text-transform: uppercase; letter-spacing: 0.04em"
              >Title</span
            >
            <input
              type="text"
              [value]="newPostTitle()"
              (input)="onNewTitle($event)"
              placeholder="What's this about?"
              style="display: block; width: 100%; margin-top: 6px; padding: 10px 12px; border-radius: var(--r-md); border: 1px solid var(--divider); font-size: 14px; font-family: var(--font); background: var(--bg); color: var(--ink)"
            />
          </label>
          <label class="block" style="margin-bottom: 14px">
            <span style="font-size: 11px; font-weight: 600; color: var(--ink-muted); text-transform: uppercase; letter-spacing: 0.04em"
              >Body</span
            >
            <textarea
              [value]="newPostBody()"
              (input)="onNewBody($event)"
              rows="4"
              placeholder="Share details, context, or questions…"
              style="display: block; width: 100%; margin-top: 6px; padding: 10px 12px; border-radius: var(--r-md); border: 1px solid var(--divider); font-size: 14px; font-family: var(--font); background: var(--bg); color: var(--ink); resize: vertical; min-height: 96px"
            ></textarea>
          </label>
          <div class="flex items-center justify-end gap-2">
            <button
              type="button"
              style="font-size: 13px; padding: 8px 14px; border-radius: var(--r-md); border: 1px solid var(--divider); background: transparent; color: var(--ink-muted); font-weight: 600; cursor: pointer"
              (click)="showPostForm.set(false)"
            >
              Cancel
            </button>
            <button
              type="button"
              style="font-size: 13px; padding: 8px 16px; border-radius: var(--r-md); border: none; background: var(--navy); color: #fff; font-weight: 600; cursor: pointer"
              [disabled]="feedService.posting() || !canSubmitPost()"
              [style.opacity]="feedService.posting() || !canSubmitPost() ? '0.65' : '1'"
              [style.cursor]="feedService.posting() || !canSubmitPost() ? 'not-allowed' : 'pointer'"
              (click)="onSubmitPost()"
            >
              {{ feedService.posting() ? 'Posting…' : 'Publish' }}
            </button>
          </div>
        </div>
      }

      <div
        class="flex items-center gap-4"
        style="padding: 18px 22px; border-radius: var(--r-xl); background: var(--emerald-light); border: 1px solid var(--emerald-border); box-shadow: var(--shadow-sm)"
      >
        <div class="pulse-dot"></div>
        <div style="font-size: 13px; color: var(--ink-2); line-height: 1.6; font-weight: 500">
          <strong style="color: var(--ink); font-family: var(--font-display)">{{ firstName() }}</strong
          >, here’s what your cohort is discussing alongside signals from your confusion hotspots
          @if (topWeakTopic()) {
            around <strong style="color: var(--emerald)">{{ topWeakTopic() }}</strong>
          }
          and recent study activity.
        </div>
      </div>

      <div
        class="flex gap-0.5"
        style="background: var(--surface-sub); border: 1px solid var(--divider); border-radius: var(--r-lg); padding: 4px; align-self: flex-start"
      >
        <div
          [style.padding]="'6px 16px'"
          [style.border-radius]="'var(--r-md)'"
          [style.font-size]="'13px'"
          [style.font-weight]="600"
          [style.color]="'var(--navy)'"
          [style.background]="'var(--card-bg)'"
          [style.box-shadow]="'var(--shadow-xs)'"
          style="transition: all var(--transition-base)"
        >
          School Feed
        </div>
        <div
          class="transition-all"
          [style.padding]="'6px 16px'"
          [style.border-radius]="'var(--r-md)'"
          [style.font-size]="'13px'"
          [style.font-weight]="500"
          [style.color]="'var(--ink-faint)'"
          [style.opacity]="'0.7'"
          style="transition: all var(--transition-base)"
        >
          Global Network Soon
        </div>
      </div>

      @if (feedService.error()) {
        <div
          style="padding: 14px 16px; border-radius: var(--r-lg); border: 1px solid var(--red-border); background: var(--red-light); color: var(--red); font-size: 13px; font-weight: 600"
        >
          {{ feedService.error() }}
        </div>
      }

      @if (feedService.loading() && feedService.posts().length === 0) {
        <div class="flex flex-col gap-5">
          @for (_ of [1, 2, 3]; track _) {
            <div class="skeleton" style="height: 160px; border-radius: var(--r-xl); background: var(--surface-sub)"></div>
          }
        </div>
      }

      @if (!feedService.loading() && feedService.posts().length === 0 && !feedService.error()) {
        <div
          style="border: 1px dashed var(--divider); border-radius: var(--r-xl); padding: 28px; background: var(--surface-sub); color: var(--ink-muted); font-size: 14px; line-height: 1.6"
        >
          No posts yet. Share a question or insight to get the conversation started.
        </div>
      }

      @if (!(feedService.loading() && feedService.posts().length === 0)) {
        <div class="flex flex-col gap-5">
          @for (post of feedService.posts(); track trackPost(post, $index)) {
            <div
              class="post-card"
              style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 24px 28px; background: var(--card-bg); box-shadow: var(--shadow-sm); transition: all var(--transition-base)"
            >
              <div class="flex items-start gap-3 mb-4">
                <div
                  class="flex items-center justify-center"
                  style="width: 34px; height: 34px; border-radius: 50%; font-size: 12px; font-weight: 700; color: #fff; flex-shrink: 0"
                  [style.background]="avatarColorForName(post.author_name)"
                >
                  {{ authorInitials(post.author_name) }}
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2">
                    <span style="font-size: 14px; font-weight: 600; color: var(--ink)">{{ post.author_name }}</span>
                  </div>
                  <div style="font-size: 12px; color: var(--ink-muted); margin-top: 2px">
                    {{ formatFeedTime(post.created_at) }}
                  </div>
                </div>
                <div
                  style="font-size: 11px; font-weight: 600; padding: 4px 12px; border-radius: var(--r-md); background: var(--navy-light); color: var(--navy); border: 1px solid var(--navy-border); white-space: nowrap"
                >
                  {{ formatPostType(post.post_type) }}
                </div>
              </div>

              <div
                style="font-size: 16px; font-weight: 600; letter-spacing: -0.3px; margin-bottom: 10px; color: var(--ink); font-family: var(--font-display)"
              >
                {{ post.title }}
              </div>
              <div
                style="font-family: var(--font); font-size: 14px; color: var(--ink-muted); line-height: 1.75; margin-bottom: 16px; font-weight: 400"
              >
                {{ post.body }}
              </div>

              <div
                class="post-actions flex items-center gap-2 flex-wrap"
                style="padding-top: 12px; border-top: 1px solid var(--divider)"
              >
                <button
                  type="button"
                  class="action-btn flex items-center gap-1.5"
                  [style.font-size]="'12.5px'"
                  [style.color]="post.upvoted ? 'var(--navy)' : 'var(--ink-muted)'"
                  [style.padding]="'6px 11px'"
                  [style.border-radius]="'var(--r-md)'"
                  [style.font-weight]="600"
                  [style.background]="post.upvoted ? 'var(--navy-light)' : 'transparent'"
                  [style.border]="post.upvoted ? '1px solid var(--navy-border)' : '1px solid transparent'"
                  style="cursor: pointer; transition: all var(--transition-fast)"
                  (click)="onUpvote(post.id, $event)"
                >
                  <lucide-icon name="arrow-up" [size]="14" [strokeWidth]="2" /> {{ post.upvotes }}
                </button>
                <button
                  type="button"
                  class="action-btn flex items-center gap-1.5"
                  style="font-size: 12.5px; color: var(--ink-muted); padding: 6px 11px; border-radius: var(--r-md); font-weight: 500; transition: all var(--transition-fast); cursor: pointer; border: 1px solid transparent; background: transparent"
                  (click)="toggleComments(post.id, $event)"
                >
                  <lucide-icon name="message-circle" [size]="14" [strokeWidth]="2" /> {{ post.comment_count }}
                </button>
                @if (isOwnPost(post)) {
                  <button
                    type="button"
                    class="action-btn flex items-center gap-1.5"
                    style="padding: 6px 10px; border-radius: var(--r-md); border: none;
                           background: transparent; cursor: pointer; font-size: 12px;
                           font-weight: 600; color: var(--red); transition: var(--transition-base)"
                    (click)="onDeletePost(post.id, $event)"
                    title="Delete post"
                  >
                    <lucide-icon name="trash-2" [size]="13" [strokeWidth]="2" />
                  </button>
                }
              </div>

              @if (isCommentsExpanded(post.id)) {
                <div style="margin-top: 16px; border-top: 1px solid var(--divider); padding-top: 16px">
                  @if (feedService.commentLoading()[post.id]) {
                    <div style="font-size: 12px; color: var(--ink-muted); margin-bottom: 12px">Loading replies…</div>
                  }

                  <div class="flex flex-col gap-3 mb-4">
                    @for (comment of commentsForPost(post.id); track comment.id) {
                      <div
                        [style.margin-left]="comment.parent_id ? '24px' : '0'"
                        style="display: flex; gap: 10px; align-items: flex-start"
                      >
                        <div
                          class="flex items-center justify-center"
                          style="width: 26px; height: 26px; border-radius: 50%;
                                 font-size: 10px; font-weight: 700; color: #fff;
                                 flex-shrink: 0; background: var(--navy)"
                        >
                          {{ comment.author_name.slice(0, 1).toUpperCase() }}
                        </div>
                        <div class="flex-1 min-w-0">
                          <div style="display: flex; gap: 8px; align-items: center; margin-bottom: 3px">
                            <span style="font-size: 12px; font-weight: 700; color: var(--ink)">{{ comment.author_name }}</span>
                            <span style="font-size: 11px; color: var(--ink-faint)">{{ formatFeedTime(comment.created_at) }}</span>
                          </div>
                          <div style="font-size: 13px; color: var(--ink-2); line-height: 1.55">{{ comment.body }}</div>
                          <button
                            type="button"
                            style="margin-top: 4px; font-size: 11px; color: var(--ink-muted); background: none;
                                   border: none; cursor: pointer; padding: 0; font-weight: 600"
                            (click)="setReplyingTo(post.id, comment.id); $event.stopPropagation()"
                          >
                            Reply
                          </button>
                        </div>
                      </div>
                    } @empty {
                      @if (!feedService.commentLoading()[post.id]) {
                        <div style="font-size: 12px; color: var(--ink-faint)">No replies yet. Be the first.</div>
                      }
                    }
                  </div>

                  @if (replyingTo()[post.id]) {
                    <div style="font-size: 11px; color: var(--ink-muted); margin-bottom: 6px; display: flex; align-items: center; gap: 6px">
                      <lucide-icon name="corner-down-right" [size]="11" [strokeWidth]="2" />
                      Replying to a comment
                      <button
                        type="button"
                        style="font-size: 11px; color: var(--red); background: none; border: none; cursor: pointer; padding: 0"
                        (click)="setReplyingTo(post.id, null)"
                      >
                        Cancel
                      </button>
                    </div>
                  }

                  <div class="flex gap-2 items-start">
                    <textarea
                      [value]="replyDraft()[post.id] ?? ''"
                      (input)="setReplyDraft(post.id, $any($event.target).value)"
                      rows="2"
                      placeholder="Write a reply…"
                      style="flex: 1; padding: 8px 10px; border-radius: var(--r-md);
                             border: 1px solid var(--divider); font-size: 13px;
                             font-family: var(--font); background: var(--bg);
                             color: var(--ink); resize: none"
                      (click)="$event.stopPropagation()"
                    ></textarea>
                    <button
                      type="button"
                      style="padding: 8px 14px; border-radius: var(--r-md); border: none;
                             background: var(--navy); color: #fff; font-size: 12px;
                             font-weight: 700; cursor: pointer; flex-shrink: 0; align-self: flex-end"
                      [disabled]="!(replyDraft()[post.id]?.trim())"
                      (click)="submitComment(post.id, $event)"
                    >
                      Send
                    </button>
                  </div>
                </div>
              }
            </div>
          }
        </div>
      }

      @if (feedService.loadingMore()) {
        <div class="flex justify-center py-2">
          <span style="font-size: 13px; color: var(--ink-muted)">Loading more…</span>
        </div>
      }

      <div #loadMoreSentinel class="h-2 w-full shrink-0" aria-hidden="true"></div>
    </div>
  `,
  styles: [
    `
      :host {
        display: flex;
        flex-direction: column;
        flex: 1;
        min-height: 0;
        overflow: hidden;
      }

      @keyframes pulseDot {
        0%,
        100% {
          transform: scale(1);
          opacity: 0.7;
        }
        50% {
          transform: scale(1.15);
          opacity: 1;
        }
      }

      .pulse-dot {
        width: 10px;
        height: 10px;
        border-radius: 50%;
        background: var(--emerald);
        flex-shrink: 0;
        animation: pulseDot 2s ease-in-out infinite;
      }

      .post-btn:hover {
        box-shadow: var(--shadow-md) !important;
      }

      .post-card:hover {
        box-shadow: var(--shadow-md) !important;
        border-color: var(--emerald) !important;
        transform: translateY(-2px);
      }

      .action-btn:hover {
        background: var(--hover-bg);
        transform: scale(1.05);
      }

      .action-btn:active {
        transform: scale(0.95);
      }

      @keyframes pulse {
        0%,
        100% {
          opacity: 1;
        }
        50% {
          opacity: 0.4;
        }
      }
      .skeleton {
        background: var(--surface-sub);
        animation: pulse 1.5s ease-in-out infinite;
      }
    `,
  ],
})
export default class FeedTabComponent implements OnInit, OnDestroy {
  private readonly authService = inject(AuthService);
  private readonly learningService = inject(LearningService);
  private readonly ngZone = inject(NgZone);
  readonly feedService = inject(FeedService);

  private readonly feedScroll = viewChild<ElementRef<HTMLElement>>('feedScroll');
  private readonly loadMoreSentinel = viewChild<ElementRef<HTMLElement>>('loadMoreSentinel');
  private io: IntersectionObserver | null = null;

  protected readonly firstName = computed(() => this.authService.currentUser()?.name?.split(' ')[0] ?? 'there');
  protected readonly topWeakTopic = computed(() => this.learningService.confusionInsight()?.top_topic_name ?? null);

  showPostForm = signal(false);
  newPostTitle = signal('');
  newPostBody = signal('');
  protected readonly expandedComments = signal<Set<string>>(new Set());
  protected readonly replyDraft = signal<Record<string, string>>({});
  protected readonly replyingTo = signal<Record<string, string | null>>({});

 constructor() {
  afterNextRender(() => {
    this.ngZone.runOutsideAngular(() => this.attachInfiniteScroll());
  });
}

  ngOnInit(): void {
    void this.feedService.loadInitial();
  }

  ngOnDestroy(): void {
    this.io?.disconnect();
    this.io = null;
  }

  private attachInfiniteScroll(): void {
    this.io?.disconnect();
    const root = this.feedScroll()?.nativeElement;
    const sentinel = this.loadMoreSentinel()?.nativeElement;
    if (!root || !sentinel) {
      return;
    }
    this.io = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting) {
          this.ngZone.run(() => {
            void this.feedService.loadMore();
          });
        }
      },
      { root, rootMargin: '240px', threshold: 0 },
    );
    this.io.observe(sentinel);
  }

  onNewTitle(ev: Event): void {
    this.newPostTitle.set((ev.target as HTMLInputElement).value);
  }

  onNewBody(ev: Event): void {
    this.newPostBody.set((ev.target as HTMLTextAreaElement).value);
  }

  /**
   * Prefer the server id so list items keep stable DOM identity across inserts/removals.
   */
  trackPost(post: { id?: string }, index: number): string {
    const id = post?.id?.trim();
    return id || `feed-post:${index}`;
  }

  canSubmitPost(): boolean {
    return !!this.newPostTitle().trim() && !!this.newPostBody().trim();
  }

  async onUpvote(postId: string, ev: Event): Promise<void> {
    ev.stopPropagation();
    await this.feedService.toggleUpvote(postId);
  }

  toggleComments(postId: string, ev?: Event): void {
    ev?.stopPropagation();
    this.expandedComments.update((s) => {
      const next = new Set(s);
      if (next.has(postId)) {
        next.delete(postId);
      } else {
        next.add(postId);
        void this.feedService.loadComments(postId);
      }
      return next;
    });
  }

  isCommentsExpanded(postId: string): boolean {
    return this.expandedComments().has(postId);
  }

  setReplyDraft(postId: string, value: string): void {
    this.replyDraft.update((m) => ({ ...m, [postId]: value }));
  }

  setReplyingTo(postId: string, commentId: string | null): void {
    this.replyingTo.update((m) => ({ ...m, [postId]: commentId }));
  }

  async submitComment(postId: string, ev?: Event): Promise<void> {
    ev?.stopPropagation();
    const body = (this.replyDraft()[postId] ?? '').trim();
    if (!body) {
      return;
    }
    const parentId = this.replyingTo()[postId] ?? undefined;
    await this.feedService.createComment(postId, body, parentId);
    this.setReplyDraft(postId, '');
    this.setReplyingTo(postId, null);
  }

  async onDeletePost(postId: string, ev?: Event): Promise<void> {
    ev?.stopPropagation();
    await this.feedService.deletePost(postId);
    this.expandedComments.update((s) => {
      if (!s.has(postId)) {
        return s;
      }
      const next = new Set(s);
      next.delete(postId);
      return next;
    });
    this.replyDraft.update((m) => {
      if (!(postId in m)) {
        return m;
      }
      const { [postId]: _, ...rest } = m;
      return rest;
    });
    this.replyingTo.update((m) => {
      if (!(postId in m)) {
        return m;
      }
      const { [postId]: _, ...rest } = m;
      return rest;
    });
  }

  isOwnPost(post: FeedPost): boolean {
    return post.user_id === this.authService.currentUser()?.id;
  }

  commentsForPost(postId: string): FeedComment[] {
    return this.feedService.commentsByPost()[postId] ?? [];
  }

  async onSubmitPost(): Promise<void> {
    if (!this.newPostTitle().trim() || !this.newPostBody().trim()) {
      return;
    }
    await this.feedService.createPost({
      title: this.newPostTitle().trim(),
      body: this.newPostBody().trim(),
      post_type: 'question',
    });
    this.newPostTitle.set('');
    this.newPostBody.set('');
    this.showPostForm.set(false);
  }

  authorInitials(name: string): string {
    const parts = name.trim().split(/\s+/).filter(Boolean);
    if (parts.length === 0) {
      return '?';
    }
    if (parts.length === 1) {
      return parts[0].substring(0, 2).toUpperCase();
    }
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
  }

  avatarColorForName(name: string): string {
    const colors = ['var(--navy)', 'var(--emerald)', 'var(--purple)', 'var(--amber)', 'var(--red)'];
    let h = 0;
    for (let i = 0; i < name.length; i++) {
      h = (Math.imul(31, h) + name.charCodeAt(i)) | 0;
    }
    return colors[Math.abs(h) % colors.length];
  }

  formatPostType(t: string): string {
    if (!t) {
      return 'Post';
    }
    return t.charAt(0).toUpperCase() + t.slice(1).toLowerCase();
  }

  formatFeedTime(iso: string): string {
    const d = new Date(iso);
    if (Number.isNaN(d.getTime())) {
      return '';
    }
    const diffMs = Date.now() - d.getTime();
    const sec = Math.floor(diffMs / 1000);
    if (sec < 60) {
      return 'just now';
    }
    const min = Math.floor(sec / 60);
    if (min < 60) {
      return `${min}m ago`;
    }
    const h = Math.floor(min / 60);
    if (h < 24) {
      return `${h}h ago`;
    }
    const days = Math.floor(h / 24);
    if (days < 7) {
      return `${days}d ago`;
    }
    return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
  }
}
