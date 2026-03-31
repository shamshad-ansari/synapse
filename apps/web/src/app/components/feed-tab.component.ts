import { Component, signal } from '@angular/core';
import { Router } from '@angular/router';
import { LucideAngularModule } from 'lucide-angular';
import { trigger, transition, style, animate } from '@angular/animations';

interface Post {
  avatar: string;
  author: string;
  reputation: number;
  school: string;
  time: string;
  topic: string;
  title: string;
  body: string;
  upvotes: number;
  upvoted: boolean;
  comments: number;
  rankingReason: string;
  provenanceLink: string;
  avatarColor: string;
}

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
          class="post-btn flex items-center gap-2 transition-all"
          style="font-size: 13px; padding: 8px 16px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 600; cursor: pointer; box-shadow: var(--shadow-sm); transition: var(--transition-base)"
        >
          <lucide-icon name="plus" [size]="15" [strokeWidth]="2" /> Post
        </button>
      </div>

      <!-- Narrated Feed Header - Synapse Intelligence -->
      <div
        class="flex items-center gap-4"
        style="padding: 18px 22px; border-radius: var(--r-xl); background: var(--emerald-light); border: 1px solid var(--emerald-border); box-shadow: var(--shadow-sm)"
      >
        <div class="pulse-dot"></div>
        <div style="font-size: 13px; color: var(--ink-2); line-height: 1.6; font-weight: 500">
          <strong style="color: var(--ink); font-family: var(--font-display)">Alex</strong>, this feed is ranked by your confusion hotspots in <strong style="color: var(--emerald)">Recursion</strong> and a missed review session from yesterday.
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

      <!-- Posts -->
      <div class="flex flex-col gap-5">
        @for (post of posts; track post.author; let i = $index) {
          <div
            class="post-card cursor-pointer"
            style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 24px 28px; background: var(--card-bg); box-shadow: var(--shadow-sm); transition: all var(--transition-base)"
          >
            <!-- Header -->
            <div class="flex items-start gap-3 mb-4">
              <div
                class="flex items-center justify-center"
                style="width: 34px; height: 34px; border-radius: 50%; font-size: 12px; font-weight: 700; color: #fff; flex-shrink: 0"
                [style.background]="post.avatarColor"
              >{{ post.avatar }}</div>
              <div class="flex-1">
                <div class="flex items-center gap-2">
                  <span style="font-size: 14px; font-weight: 600; color: var(--ink)">{{ post.author }}</span>
                  <div
                    class="reputation-badge flex items-center gap-1"
                    style="font-size: 11px; font-weight: 700; font-family: var(--mono); color: var(--emerald); background: var(--emerald-light); border: 1px solid var(--emerald-border); padding: 2px 7px; border-radius: 12px; transition: transform 0.15s ease"
                  >
                    <lucide-icon name="arrow-up" [size]="10" [strokeWidth]="3" /> {{ post.reputation }}
                  </div>
                </div>
                <div style="font-size: 12px; color: var(--ink-muted); margin-top: 2px">{{ post.school }} · {{ post.time }}</div>
              </div>
              <!-- Topic Tag -->
              <div style="font-size: 11px; font-weight: 600; padding: 4px 12px; border-radius: var(--r-md); background: var(--navy-light); color: var(--navy); border: 1px solid var(--navy-border); white-space: nowrap">{{ post.topic }}</div>
            </div>

            <!-- Content -->
            <div style="font-size: 16px; font-weight: 600; letter-spacing: -0.3px; margin-bottom: 10px; color: var(--ink); font-family: var(--font-display)">{{ post.title }}</div>
            <div style="font-family: var(--serif); font-size: 14px; color: var(--ink-muted); line-height: 1.75; margin-bottom: 16px; font-weight: 300">{{ post.body }}</div>

            <!-- Provenance Link -->
            <div
              class="provenance-link inline-flex items-center gap-1.5 cursor-pointer"
              style="font-size: 11.5px; color: var(--ink-faint); background: var(--surface-sub); border: 1px solid var(--divider); padding: 4px 10px; border-radius: var(--r-md); margin-bottom: 16px; transition: var(--transition-fast)"
              (click)="navigateToNotes($event)"
            >
              <lucide-icon name="link" [size]="12" [strokeWidth]="2" /> Provenance: {{ post.provenanceLink }}
            </div>

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
              >
                <lucide-icon name="arrow-up" [size]="14" [strokeWidth]="2" /> {{ post.upvotes }}
              </div>
              <!-- Comments -->
              <div
                class="action-btn flex items-center gap-1.5 cursor-pointer"
                style="font-size: 12.5px; color: var(--ink-muted); padding: 6px 11px; border-radius: var(--r-md); font-weight: 500; transition: all var(--transition-fast)"
              >
                <lucide-icon name="message-circle" [size]="14" [strokeWidth]="2" /> {{ post.comments }}
              </div>
              <!-- Bookmark -->
              <div
                class="action-btn flex items-center gap-1.5 cursor-pointer"
                style="font-size: 12.5px; color: var(--ink-muted); padding: 6px 11px; border-radius: var(--r-md); font-weight: 500; transition: all var(--transition-fast)"
              >
                <lucide-icon name="bookmark" [size]="14" [strokeWidth]="2" /> Save
              </div>

              <!-- Why this ranked? -->
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
            </div>

            <!-- Ranking Explanation -->
            @if (isRankingExpanded(i)) {
              <div
                @expandCollapse
                style="margin-top: 12px; padding: 14px 16px; border-radius: var(--r-lg); background: var(--surface-sub); border: 1px solid var(--divider); font-size: 12.5px; color: var(--ink-2); line-height: 1.7; overflow: hidden"
              >{{ post.rankingReason }}</div>
            }
          </div>
        }
      </div>
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
  `],
})
export default class FeedTabComponent {
  private readonly router: Router;

  activeSegment = signal<string>('school');
  expandedRankings = signal<Set<number>>(new Set<number>());

  posts: Post[] = [
    {
      avatar: 'JL',
      author: 'Jamie Liu',
      reputation: 127,
      school: 'MIT · CS225',
      time: '2h ago',
      topic: 'Physics · Mechanics',
      title: 'When exactly do you use memoization vs tabulation for dynamic programming?',
      body: 'I keep getting confused on this — my professor says they\'re equivalent in complexity but when do you actually prefer one? I find tabulation cleaner but memoization more intuitive. Is there a general heuristic?',
      upvotes: 14,
      upvoted: true,
      comments: 3,
      rankingReason: 'Matches your weak topic: Recursion and your recent confusion spike in DP optimization strategies.',
      provenanceLink: 'Recursion Notes · Section 4',
      avatarColor: 'var(--navy)',
    },
    {
      avatar: 'MR',
      author: 'Maya Roth',
      reputation: 243,
      school: 'MIT · CS225',
      time: '5h ago',
      topic: 'Mathematics · Induction',
      title: 'Complete Guide: Strong vs Weak Induction — with worked examples',
      body: 'Compiled 8 pages of notes covering when to use each, common pitfalls, and 6 worked exam-style problems. Includes a decision framework for picking the right approach on any problem.',
      upvotes: 31,
      upvoted: false,
      comments: 7,
      rankingReason: 'Trending in your cohort. 68% of CS225 students reviewed this material in the last 48 hours.',
      provenanceLink: 'Induction Overview · Section 2',
      avatarColor: 'var(--purple)',
    },
    {
      avatar: 'SK',
      author: 'Sam Kato',
      reputation: 189,
      school: 'MIT · CS225',
      time: '8h ago',
      topic: 'Computer Science · Algorithms',
      title: 'School confusion alert: Recursion base cases',
      body: 'Based on anonymized review data, most confusion is concentrated on off-by-one base cases and empty list recursion. Pset 3 Q4 specifically targets this. Recommend reviewing before Thursday.',
      upvotes: 47,
      upvoted: false,
      comments: 12,
      rankingReason: 'School-wide confusion spike detected. This aligns with your current study focus and upcoming deadline.',
      provenanceLink: 'Recursion Notes · Section 1',
      avatarColor: 'var(--emerald)',
    },
  ];

  constructor(router: Router) {
    this.router = router;
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

  navigateToNotes(event: Event): void {
    event.stopPropagation();
    this.router.navigate(['/notes']);
  }
}
