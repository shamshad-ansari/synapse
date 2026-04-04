import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { Router } from '@angular/router';
import { LucideAngularModule } from 'lucide-angular';
import { trigger, transition, style, animate } from '@angular/animations';
import { AuthService } from '../core/auth/auth.service';
import { LearningService } from '../features/learning/learning.service';
import { FeedService, type FeedPost } from '../features/feed/feed.service';

@Component({
  selector: 'app-feed-tab',
  standalone: true,
  imports: [LucideAngularModule],
  animations: [
    trigger('expandCollapse', [
      transition(':enter', [
        style({ height: 0, opacity: 0, overflow: 'hidden' }),
        animate('250ms ease-out', style({ height: '*', opacity: 1 })),
      ]),
      transition(':leave', [
        style({ overflow: 'hidden' }),
        animate('250ms ease-out', style({ height: 0, opacity: 0 })),
      ]),
    ]),
  ],
  template: `
    <div class="flex flex-col overflow-y-auto overflow-x-hidden" style="padding: 44px 56px 56px; gap: 36px">

      <!-- Topbar -->
      <div class="flex items-center justify-between gap-4">
        <div>
          <div style="font-size: 26px; font-weight: 700; letter-spacing: -0.6px; font-family: var(--font-display); color: var(--ink)">Feed</div>
          <div style="font-size: 13.5px; color: var(--ink-muted); margin-top: 4px">Intelligence-ranked learning insights from your cohort</div>
        </div>
        <button
          type="button"
          class="post-btn flex items-center gap-2 transition-all"
          style="font-size: 13px; padding: 8px 16px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 600; cursor: pointer; box-shadow: var(--shadow-sm); transition: var(--transition-base)"
          (click)="showPostForm.update(v => !v)"
        >
          <lucide-icon name="plus" [size]="15" [strokeWidth]="2" /> Post
        </button>
      </div>

      @if (showPostForm()) {
        <div
          style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 20px 22px; background: var(--card-bg); box-shadow: var(--shadow-sm)"
        >
          <div style="font-size: 13px; font-weight: 600; color: var(--ink); margin-bottom: 12px; font-family: var(--font-display)">New post</div>
          <label class="block" style="margin-bottom: 10px">
            <span style="font-size: 11px; font-weight: 600; color: var(--ink-muted); text-transform: uppercase; letter-spacing: 0.04em">Title</span>
            <input
              type="text"
              [value]="newPostTitle()"
              (input)="onNewTitle($event)"
              placeholder="What's this about?"
              style="display: block; width: 100%; margin-top: 6px; padding: 10px 12px; border-radius: var(--r-md); border: 1px solid var(--divider); font-size: 14px; font-family: var(--font); background: var(--bg); color: var(--ink)"
            />
          </label>
          <label class="block" style="margin-bottom: 14px">
            <span style="font-size: 11px; font-weight: 600; color: var(--ink-muted); text-transform: uppercase; letter-spacing: 0.04em">Body</span>
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
            >Cancel</button>
            <button
              type="button"
              style="font-size: 13px; padding: 8px 16px; border-radius: var(--r-md); border: none; background: var(--navy); color: #fff; font-weight: 600; cursor: pointer"
              [disabled]="feedService.posting()"
              (click)="onSubmitPost()"
            >{{ feedService.posting() ? 'Posting…' : 'Publish' }}</button>
          </div>
        </div>
      }

      <!-- Narrated Feed Header - Synapse Intelligence -->
      <div
        class="flex items-center gap-4"
        style="padding: 18px 22px; border-radius: var(--r-xl); background: var(--emerald-light); border: 1px solid var(--emerald-border); box-shadow: var(--shadow-sm)"
      >
        <div class="pulse-dot"></div>
        <div style="font-size: 13px; color: var(--ink-2); line-height: 1.6; font-weight: 500">
          <strong style="color: var(--ink); font-family: var(--font-display)">{{ firstName() }}</strong>, this feed is ranked by your confusion hotspots
          @if (topWeakTopic()) {
            in <strong style="color: var(--emerald)">{{ topWeakTopic() }}</strong>
          }
          and your recent study activity.
        </div>
      </div>

      <!-- Filters -->
      <div class="flex gap-0.5" style="background: var(--surface-sub); border: 1px solid var(--divider); border-radius: var(--r-lg); padding: 4px; align-self: flex-start">
        <div
          class="cursor-pointer transition-all"
          [style.padding]="'6px 16px'"
          [style.border-radius]="'var(--r-md)'"
          [style.font-size]="'13px'"
          [style.font-weight]="600"
          [style.color]="activeSegment() === 'school' ? 'var(--navy)' : 'var(--ink-muted)'"
          [style.background]="activeSegment() === 'school' ? 'var(--card-bg)' : 'transparent'"
          [style.box-shadow]="activeSegment() === 'school' ? 'var(--shadow-xs)' : 'none'"
          style="transition: all var(--transition-base)"
          (click)="activeSegment.set('school')"
        >School Feed</div>
        <div
          class="cursor-pointer transition-all"
          [style.padding]="'6px 16px'"
          [style.border-radius]="'var(--r-md)'"
          [style.font-size]="'13px'"
          [style.font-weight]="600"
          [style.color]="activeSegment() === 'global' ? 'var(--navy)' : 'var(--ink-muted)'"
          [style.background]="activeSegment() === 'global' ? 'var(--card-bg)' : 'transparent'"
          [style.box-shadow]="activeSegment() === 'global' ? 'var(--shadow-xs)' : 'none'"
          style="transition: all var(--transition-base)"
          (click)="activeSegment.set('global')"
        >Global Network</div>
      </div>

      <!-- Loading skeleton -->
      @if (feedService.loading() && posts().length === 0) {
        <div class="flex flex-col gap-5">
          @for (_ of [1, 2, 3]; track _) {
            <div class="skeleton" style="height: 160px; border-radius: var(--r-xl); background: var(--surface-sub)"></div>
          }
        </div>
      }

      <!-- Posts -->
      @if (!(feedService.loading() && posts().length === 0)) {
      <div class="flex flex-col gap-5">
        @for (post of posts(); track post.id; let i = $index) {
          <div
            class="post-card cursor-pointer"
            style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 24px 28px; background: var(--card-bg); box-shadow: var(--shadow-sm); transition: all var(--transition-base)"
          >
            <!-- Header -->
            <div class="flex items-start gap-3 mb-4">
              <div
                class="flex items-center justify-center"
                style="width: 34px; height: 34px; border-radius: 50%; font-size: 12px; font-weight: 700; color: #fff; flex-shrink: 0"
                [style.background]="avatarColorForName(post.author_name)"
              >{{ authorInitials(post.author_name) }}</div>
              <div class="flex-1">
                <div class="flex items-center gap-2">
                  <span style="font-size: 14px; font-weight: 600; color: var(--ink)">{{ post.author_name }}</span>
                </div>
                <div style="font-size: 12px; color: var(--ink-muted); margin-top: 2px">{{ formatFeedTime(post.created_at) }}</div>
              </div>
              <!-- Topic Tag -->
              <div style="font-size: 11px; font-weight: 600; padding: 4px 12px; border-radius: var(--r-md); background: var(--navy-light); color: var(--navy); border: 1px solid var(--navy-border); white-space: nowrap">{{ formatPostType(post.post_type) }}</div>
            </div>

            <!-- Content -->
            <div style="font-size: 16px; font-weight: 600; letter-spacing: -0.3px; margin-bottom: 10px; color: var(--ink); font-family: var(--font-display)">{{ post.title }}</div>
            <div style="font-family: var(--serif); font-size: 14px; color: var(--ink-muted); line-height: 1.75; margin-bottom: 16px; font-weight: 300">{{ post.body }}</div>

            <!-- Actions -->
            <div
              class="post-actions flex items-center gap-2"
              style="padding-top: 12px; border-top: 1px solid var(--divider)"
            >
              <!-- Upvote -->
              <div
                class="action-btn flex items-center gap-1.5 cursor-pointer"
                [style.font-size]="'12.5px'"
                [style.color]="post.upvoted ? 'var(--navy)' : 'var(--ink-muted)'"
                [style.padding]="'6px 11px'"
                [style.border-radius]="'var(--r-md)'"
                [style.font-weight]="600"
                [style.background]="post.upvoted ? 'var(--navy-light)' : 'transparent'"
                [style.border]="post.upvoted ? '1px solid var(--navy-border)' : '1px solid transparent'"
                style="transition: all var(--transition-fast)"
                (click)="toggleUpvote(i); $event.stopPropagation()"
              >
                <lucide-icon name="arrow-up" [size]="14" [strokeWidth]="2" /> {{ post.upvotes }}
              </div>
              <!-- Comments -->
              <div
                class="action-btn flex items-center gap-1.5 cursor-pointer"
                style="font-size: 12.5px; color: var(--ink-muted); padding: 6px 11px; border-radius: var(--r-md); font-weight: 500; transition: all var(--transition-fast)"
              >
                <lucide-icon name="message-circle" [size]="14" [strokeWidth]="2" /> 0
              </div>
              <!-- Bookmark -->
              <div
                class="action-btn flex items-center gap-1.5 cursor-pointer"
                style="font-size: 12.5px; color: var(--ink-muted); padding: 6px 11px; border-radius: var(--r-md); font-weight: 500; transition: all var(--transition-fast)"
              >
                <lucide-icon name="bookmark" [size]="14" [strokeWidth]="2" /> Save
              </div>

              @if (false) {
              <!-- Why this ranked? — hidden until API provides ranking metadata -->
              <div
                class="ranking-btn ml-auto flex items-center gap-1.5 cursor-pointer"
                [style.font-size]="'11.5px'"
                [style.color]="isRankingExpanded(i) ? 'var(--emerald)' : 'var(--ink-faint)'"
                [style.padding]="'6px 12px'"
                [style.border-radius]="'var(--r-md)'"
                [style.font-weight]="600"
                [style.background]="isRankingExpanded(i) ? 'var(--emerald-light)' : 'transparent'"
                [style.border]="isRankingExpanded(i) ? '1px solid var(--emerald-border)' : '1px solid var(--divider)'"
                style="transition: all var(--transition-base)"
                (click)="toggleRanking(i); $event.stopPropagation()"
              >
                <lucide-icon name="info" [size]="13" [strokeWidth]="2" />
                Why ranked?
                <lucide-icon
                  name="chevron-down"
                  [size]="13"
                  [strokeWidth]="2"
                  class="chevron-icon"
                  [style.transform]="isRankingExpanded(i) ? 'rotate(180deg)' : 'rotate(0)'"
                />
              </div>
              }
            </div>

            @if (false) {
            <!-- Ranking Explanation -->
            @if (isRankingExpanded(i)) {
              <div
                @expandCollapse
                style="margin-top: 12px; padding: 14px 16px; border-radius: var(--r-lg); background: var(--surface-sub); border: 1px solid var(--divider); font-size: 12.5px; color: var(--ink-2); line-height: 1.7; overflow: hidden"
              ></div>
            }
            }
          </div>
        }
      </div>
      }
    </div>
  `,
  styles: [`
    :host { display: flex; flex-direction: column; overflow: hidden; }

    @keyframes pulseDot {
      0%, 100% {
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

    .reputation-badge:hover {
      transform: scale(1.05);
    }

    .provenance-link:hover {
      border-color: var(--navy) !important;
      color: var(--navy) !important;
    }

    .action-btn:hover {
      background: var(--hover-bg);
      transform: scale(1.05);
    }

    .action-btn:active {
      transform: scale(0.95);
    }

    .ranking-btn:hover {
      transform: scale(1.02);
    }

    .ranking-btn:active {
      transform: scale(0.98);
    }

    .chevron-icon {
      transition: transform 0.2s ease;
      display: inline-flex;
    }

    @keyframes pulse {
      0%, 100% { opacity: 1; }
      50% { opacity: 0.4; }
    }
    .skeleton {
      background: var(--surface-sub);
      animation: pulse 1.5s ease-in-out infinite;
    }
  `],
})
export default class FeedTabComponent implements OnInit {
  private readonly router = inject(Router);
  private readonly authService = inject(AuthService);
  private readonly learningService = inject(LearningService);
  readonly feedService = inject(FeedService);

  protected readonly firstName = computed(() =>
    this.authService.currentUser()?.name?.split(' ')[0] ?? 'there',
  );
  protected readonly topWeakTopic = computed(() =>
    this.learningService.confusionInsight()?.top_topic_name ?? null,
  );

  posts = signal<FeedPost[]>([]);

  activeSegment = signal<string>('school');
  expandedRankings = signal<Set<number>>(new Set<number>());

  showPostForm = signal(false);
  newPostTitle = signal('');
  newPostBody = signal('');

  ngOnInit(): void {
    void this.feedService.loadPosts().then(() => {
      this.posts.set([...this.feedService.posts()]);
    });
  }

  onNewTitle(ev: Event): void {
    this.newPostTitle.set((ev.target as HTMLInputElement).value);
  }

  onNewBody(ev: Event): void {
    this.newPostBody.set((ev.target as HTMLTextAreaElement).value);
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
    this.posts.set([...this.feedService.posts()]);
    this.newPostTitle.set('');
    this.newPostBody.set('');
    this.showPostForm.set(false);
  }

  toggleUpvote(index: number): void {
    this.posts.update((list) =>
      list.map((p, i) => {
        if (i !== index) {
          return p;
        }
        return {
          ...p,
          upvoted: !p.upvoted,
          upvotes: p.upvoted ? p.upvotes - 1 : p.upvotes + 1,
        };
      }),
    );
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

  toggleRanking(index: number): void {
    const current = this.expandedRankings();
    const next = new Set(current);
    if (next.has(index)) {
      next.delete(index);
    } else {
      next.add(index);
    }
    this.expandedRankings.set(next);
  }

  isRankingExpanded(index: number): boolean {
    return this.expandedRankings().has(index);
  }
}
