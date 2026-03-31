import { Component, signal, computed } from '@angular/core';
import { Router } from '@angular/router';
import { LucideAngularModule } from 'lucide-angular';
import { trigger, transition, style, animate, query, stagger } from '@angular/animations';

const cards = [
  { q: 'What is the difference between weak and strong induction?', a: 'Weak induction assumes P(k) to prove P(k+1). Strong induction assumes P(j) for all j ≤ k — giving you all previous cases. Strong induction is used when proving P(k+1) requires more than just the immediately prior case.', src: 'Induction Overview · S3' },
  { q: 'Define the base case in mathematical induction.', a: 'The base case proves P(n₀) directly — usually P(0) or P(1). It is the foundation: without it, the inductive chain has nowhere to anchor.', src: 'Induction Overview · S1' },
  { q: 'When can a recursion have multiple base cases?', a: 'When the inductive step requires P(k) AND P(k-1) to prove P(k+1), you need base cases for both P(0) and P(1). E.g., Fibonacci.', src: 'Recursion Notes · S2' },
  { q: 'What is the induction hypothesis?', a: 'The assumption that P(k) holds for some arbitrary but fixed k ≥ base. You use this fact to then prove P(k+1).', src: 'Induction Overview · S2' },
];

@Component({
  selector: 'app-review-tab',
  standalone: true,
  imports: [LucideAngularModule],
  animations: [
    trigger('cardSwap', [
      transition('* => *', [
        style({ opacity: 0, transform: 'scale(0.96) rotateY(-8deg)' }),
        animate('350ms cubic-bezier(0.4, 0, 0.2, 1)', style({ opacity: 1, transform: 'scale(1) rotateY(0)' })),
      ]),
    ]),
    trigger('fadeInUp', [
      transition(':enter', [
        style({ opacity: 0, transform: 'translateY(10px)' }),
        animate('300ms ease-out', style({ opacity: 1, transform: 'translateY(0)' })),
      ]),
    ]),
    trigger('staggerIn', [
      transition(':enter', [
        style({ opacity: 0, transform: 'translateY(10px)' }),
        animate('250ms {{delay}}ms ease-out', style({ opacity: 1, transform: 'translateY(0)' })),
      ]),
    ]),
  ],
  template: `
    <div class="flex flex-col items-center justify-start overflow-y-auto overflow-x-hidden" style="padding: 44px 56px 56px; gap: 32px">
      <!-- Topbar -->
      <div class="w-full flex items-center justify-between gap-4 max-w-[720px]">
        <div>
          <div style="font-size: 26px; font-weight: 700; letter-spacing: -0.6px; font-family: var(--font-display); color: var(--ink)">Review Session</div>
          <div style="font-size: 13.5px; color: var(--ink-muted); margin-top: 4px">CS225 · Mixed queue · SM-2 adaptive</div>
        </div>
        <button
          class="review-end-btn flex items-center gap-2 transition-all"
          (click)="navigateToToday()"
        >
          <lucide-icon name="x" [size]="15" [strokeWidth]="2" /> End Session
        </button>
      </div>

      <!-- Progress -->
      <div class="w-full max-w-[720px]">
        <div style="height: 5px; background: var(--surface-sub); border-radius: 5px; overflow: hidden; border: 1px solid var(--divider)">
          <div
            class="review-progress-fill"
            [style.width.%]="progress()"
            style="height: 100%; background: var(--emerald); border-radius: 4px"
          ></div>
        </div>
        <div class="flex justify-between mt-2.5" style="font-size: 12px; color: var(--ink-muted)">
          <span>Card {{ cardIndex() + 1 }} of {{ totalCards }}</span>
          <span>Est. {{ remaining() }} min remaining</span>
        </div>
      </div>

      <!-- Flashcard -->
      <div class="w-full max-w-[720px]" [@cardSwap]="cardIndex()">
        <div
          class="review-card w-full text-center"
          [class.review-card--unrevealed]="!revealed()"
          [style.border-color]="isCardHovered() && !revealed() ? 'var(--emerald)' : 'var(--divider)'"
          [style.box-shadow]="isCardHovered() ? 'var(--shadow-lg)' : 'var(--shadow-md)'"
          (click)="!revealed() ? handleReveal() : null"
          (mouseenter)="isCardHovered.set(true)"
          (mouseleave)="isCardHovered.set(false)"
        >
          <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 1px; color: var(--ink-faint); font-family: var(--mono); font-weight: 700">
            Recall
          </div>
          <div class="flex items-center gap-2" style="font-size: 12.5px; color: var(--navy); font-weight: 600">
            <lucide-icon name="book-open" [size]="15" [strokeWidth]="2" /> Induction · CS225
          </div>
          <div style="font-family: var(--font); font-size: 21px; font-weight: 500; line-height: 1.5; max-width: 560px; color: var(--ink)">
            {{ currentCard().q }}
          </div>
          @if (revealed()) {
            <div
              @fadeInUp
              style="font-family: var(--font); font-size: 15.5px; color: var(--ink-2); line-height: 1.7; max-width: 560px; font-weight: 400; border-top: 1px solid var(--divider); padding-top: 18px"
            >
              {{ currentCard().a }}
            </div>
          }
          @if (!revealed()) {
            <div style="font-size: 12.5px; color: var(--ink-faint)">Click anywhere to reveal answer</div>
          }
          <a
            class="review-source-link flex items-center gap-1.5 cursor-pointer transition-all"
            (click)="$event.stopPropagation(); navigateToNotes()"
          >
            <lucide-icon name="file-text" [size]="14" [strokeWidth]="2" /> Source: {{ currentCard().src }}
          </a>
        </div>
      </div>

      <!-- Controls -->
      <div class="w-full max-w-[720px]">
        @if (!revealed()) {
          <div class="flex justify-center w-full">
            <button
              class="review-reveal-btn"
              (click)="handleReveal()"
            >
              Reveal Answer
            </button>
          </div>
        }
        @if (revealed()) {
          <div class="flex flex-col gap-3.5 w-full">
            <div style="font-size: 12.5px; color: var(--ink-muted); text-align: center; font-weight: 600">
              How well did you recall this?
            </div>
            <div class="flex gap-2.5">
              @for (btn of ratingButtons; track btn.className) {
                <div
                  class="review-rating-btn flex-1 text-center cursor-pointer"
                  [class]="'review-rating-btn review-rating-btn--' + btn.className + ' flex-1 text-center cursor-pointer'"
                  [@staggerIn]="{ value: ':enter', params: { delay: btn.delay } }"
                  (click)="handleNextCard()"
                >
                  {{ btn.label }}<br />
                  <span style="font-size: 10.5px; color: var(--ink-faint); font-weight: 500">{{ btn.sub }}</span>
                </div>
              }
            </div>
            <div
              class="review-confused-btn flex items-center justify-center gap-2 cursor-pointer transition-all"
              [style.font-size.px]="12.5"
              [style.color]="confused() ? 'var(--red)' : 'var(--ink-muted)'"
              [style.padding]="'9px 12px'"
              [style.border-radius]="'var(--r-lg)'"
              [style.background]="confused() ? 'var(--red-light)' : 'transparent'"
              [style.border]="confused() ? '1px solid var(--red-border)' : '1px solid transparent'"
              [style.font-weight]="500"
              [class.review-confused-btn--active]="confused()"
              (click)="toggleConfused()"
            >
              <lucide-icon name="help-circle" [size]="16" [strokeWidth]="2" />
              <span>{{ confused() ? 'Marked as confusing ✓' : 'Mark as confusing — updates your mastery model' }}</span>
            </div>
          </div>
        }
      </div>

      <!-- Stats Sidebar -->
      <div class="fixed right-10 top-1/2 -translate-y-1/2 flex flex-col gap-2.5">
        <div class="review-stat-card text-center">
          <div style="font-size: 22px; font-weight: 800; font-family: var(--font-display); color: var(--emerald)">{{ stats().correct }}</div>
          <div style="font-size: 9px; color: var(--ink-faint); text-transform: uppercase; letter-spacing: 0.8px; margin-top: 3px; font-weight: 700">Correct</div>
        </div>
        <div class="review-stat-card text-center">
          <div style="font-size: 22px; font-weight: 800; font-family: var(--font-display); color: var(--red)">{{ stats().again }}</div>
          <div style="font-size: 9px; color: var(--ink-faint); text-transform: uppercase; letter-spacing: 0.8px; margin-top: 3px; font-weight: 700">Again</div>
        </div>
        <div class="review-stat-card text-center">
          <div style="font-size: 22px; font-weight: 800; font-family: var(--font-display); color: var(--amber)">{{ stats().confused }}</div>
          <div style="font-size: 9px; color: var(--ink-faint); text-transform: uppercase; letter-spacing: 0.8px; margin-top: 3px; font-weight: 700">Confused</div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    :host { display: flex; flex-direction: column; overflow: hidden; }

    .review-progress-fill {
      transition: width 0.5s ease-out;
    }

    .review-end-btn {
      font-size: 13px;
      padding: 8px 16px;
      border-radius: var(--r-lg);
      border: 1px solid var(--divider);
      background: transparent;
      color: var(--ink-2);
      font-weight: 600;
      cursor: pointer;
      transition: var(--transition-base);
    }
    .review-end-btn:hover {
      border-color: var(--navy);
      color: var(--navy);
    }

    .review-card {
      border: 1px solid var(--divider);
      border-radius: var(--r-xl);
      padding: 56px 52px;
      min-height: 300px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 18px;
      background: var(--card-bg);
      box-shadow: var(--shadow-md);
      transition: all var(--transition-base);
    }
    .review-card--unrevealed {
      cursor: pointer;
    }
    .review-card--unrevealed:hover {
      transform: translateY(-4px);
    }

    .review-source-link {
      font-size: 11.5px;
      color: var(--ink-faint);
      background: var(--surface-sub);
      border: 1px solid var(--divider);
      padding: 4px 12px;
      border-radius: var(--r-md);
      text-decoration: none;
      transition: var(--transition-base);
    }
    .review-source-link:hover {
      border-color: var(--navy);
      color: var(--navy);
    }

    .review-reveal-btn {
      padding: 14px 56px;
      background: var(--navy);
      color: #fff;
      border: none;
      border-radius: var(--r-lg);
      font-size: 14.5px;
      font-weight: 700;
      cursor: pointer;
      letter-spacing: -0.2px;
      box-shadow: var(--shadow-md);
      font-family: var(--font-display);
      transition: transform 0.15s ease, box-shadow 0.15s ease;
    }
    .review-reveal-btn:hover {
      transform: scale(1.04);
      box-shadow: var(--shadow-lg);
    }
    .review-reveal-btn:active {
      transform: scale(0.96);
    }

    .review-rating-btn {
      padding: 11px 10px;
      border-radius: var(--r-lg);
      border: 1px solid var(--divider);
      background: var(--card-bg);
      color: var(--ink-muted);
      font-size: 13.5px;
      font-weight: 600;
      transition: all var(--transition-base);
    }
    .review-rating-btn:hover {
      transform: scale(1.04) translateY(-2px);
      box-shadow: var(--shadow-sm);
    }
    .review-rating-btn:active {
      transform: scale(0.97);
    }
    .review-rating-btn--easy:hover {
      border-color: var(--emerald-border);
      background: var(--emerald-light);
      color: var(--emerald);
    }
    .review-rating-btn--good:hover {
      border-color: var(--navy-border);
      background: var(--navy-light);
      color: var(--navy);
    }
    .review-rating-btn--hard:hover {
      border-color: var(--amber-border);
      background: var(--amber-light);
      color: var(--amber);
    }
    .review-rating-btn--again:hover {
      border-color: var(--red-border);
      background: var(--red-light);
      color: var(--red);
    }

    .review-confused-btn {
      transition: var(--transition-base);
    }
    .review-confused-btn:not(.review-confused-btn--active):hover {
      border-color: var(--divider) !important;
      background: var(--hover-bg) !important;
    }

    .review-stat-card {
      background: var(--card-bg);
      border: 1px solid var(--divider);
      border-radius: var(--r-xl);
      padding: 16px 18px;
      min-width: 72px;
      box-shadow: var(--shadow-sm);
      cursor: pointer;
      transition: box-shadow var(--transition-base), transform 0.15s ease;
    }
    .review-stat-card:hover {
      box-shadow: var(--shadow-md);
      transform: scale(1.05) translateY(-2px);
    }
  `],
})
export default class ReviewTabComponent {
  private readonly router: Router;

  cardIndex = signal(0);
  revealed = signal(false);
  confused = signal(false);
  stats = signal({ correct: 2, again: 1, confused: 3 });
  isCardHovered = signal(false);

  totalCards = 12;

  currentCard = computed(() => cards[this.cardIndex() % cards.length]);
  progress = computed(() => ((this.cardIndex() + 1) / this.totalCards) * 100);
  remaining = computed(() => Math.max(1, 12 - this.cardIndex()));

  ratingButtons = [
    { label: '😌 Easy', sub: '+10 days', className: 'easy', delay: 0 },
    { label: '👍 Good', sub: '+4 days', className: 'good', delay: 60 },
    { label: '😬 Hard', sub: '+1 day', className: 'hard', delay: 120 },
    { label: '🔁 Again', sub: '10 min', className: 'again', delay: 180 },
  ];

  constructor(router: Router) {
    this.router = router;
  }

  handleReveal(): void {
    this.revealed.set(true);
  }

  handleNextCard(): void {
    if (this.cardIndex() + 1 < this.totalCards) {
      this.cardIndex.set(this.cardIndex() + 1);
      this.revealed.set(false);
      this.confused.set(false);
    }
  }

  toggleConfused(): void {
    this.confused.set(!this.confused());
  }

  navigateToToday(): void {
    this.router.navigate(['/today']);
  }

  navigateToNotes(): void {
    this.router.navigate(['/notes']);
  }
}
