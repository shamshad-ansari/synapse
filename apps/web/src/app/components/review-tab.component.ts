import { Component, signal, computed, inject, effect, OnInit, HostListener } from '@angular/core';
import { Router } from '@angular/router';
import { LucideAngularModule } from 'lucide-angular';
import { trigger, transition, style, animate, query, stagger } from '@angular/animations';
import { LearningService } from '../features/learning/learning.service';

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
      @if (learningService.loading()) {
        <div class="w-full max-w-[720px] flex flex-col gap-4" style="padding-top: 8px">
          <div class="skeleton" style="height: 28px; width: 220px; border-radius: var(--r-md); background: var(--surface-sub)"></div>
          <div class="skeleton" style="height: 300px; width: 100%; border-radius: var(--r-xl); background: var(--surface-sub)"></div>
        </div>
      } @else if (allCards().length === 0) {
        <div class="flex flex-col items-center justify-center gap-5" style="min-height: 320px; max-width: 420px; text-align: center">
          <div style="font-size: 15px; font-weight: 600; color: var(--ink-muted); line-height: 1.6">
            No cards due right now. Create some flashcards to start reviewing.
          </div>
          <button
            type="button"
            style="font-size: 13px; padding: 10px 20px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 600; cursor: pointer; font-family: var(--font-display); box-shadow: var(--shadow-md)"
            (click)="navigateToNotes()"
          >
            Go to Notes
          </button>
        </div>
      } @else {
      @if (sessionDone()) {
        <!-- Session Complete screen -->
        <div class="flex flex-col items-center gap-8" style="max-width: 480px; text-align: center; padding-top: 40px">
          <div style="font-size: 48px">🎉</div>
          <div>
            <div style="font-size: 28px; font-weight: 700; letter-spacing: -0.6px; font-family: var(--font-display); color: var(--ink); margin-bottom: 8px">
              Session Complete
            </div>
            <div style="font-size: 14px; color: var(--ink-muted); line-height: 1.7">
              {{ totalCards() }} cards reviewed. Your mastery model has been updated.
            </div>
          </div>
          <!-- Stats row -->
          <div class="flex gap-6">
            <div class="text-center">
              <div style="font-size: 32px; font-weight: 800; font-family: var(--font-display); color: var(--emerald)">{{ stats().correct }}</div>
              <div style="font-size: 11px; text-transform: uppercase; letter-spacing: 0.8px; color: var(--ink-faint); font-weight: 700; margin-top: 4px">Correct</div>
            </div>
            <div class="text-center">
              <div style="font-size: 32px; font-weight: 800; font-family: var(--font-display); color: var(--red)">{{ stats().again }}</div>
              <div style="font-size: 11px; text-transform: uppercase; letter-spacing: 0.8px; color: var(--ink-faint); font-weight: 700; margin-top: 4px">Again</div>
            </div>
            <div class="text-center">
              <div style="font-size: 32px; font-weight: 800; font-family: var(--font-display); color: var(--amber)">{{ stats().confused }}</div>
              <div style="font-size: 11px; text-transform: uppercase; letter-spacing: 0.8px; color: var(--ink-faint); font-weight: 700; margin-top: 4px">Confused</div>
            </div>
          </div>
          <!-- Accuracy pill -->
          <div
            style="padding: 8px 20px; border-radius: var(--r-xl); font-size: 14px; font-weight: 700; font-family: var(--font-display)"
            [style.background]="accuracy() >= 70 ? 'var(--emerald-light)' : 'var(--red-light)'"
            [style.color]="accuracy() >= 70 ? 'var(--emerald)' : 'var(--red)'"
            [style.border]="'1px solid ' + (accuracy() >= 70 ? 'var(--emerald-border)' : 'var(--red-border)')"
          >
            {{ accuracy() }}% accuracy
          </div>
          <!-- Buttons -->
          <div class="flex gap-3">
            <button
              style="font-size: 13px; padding: 10px 24px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: transparent; color: var(--ink-2); font-weight: 600; cursor: pointer; transition: var(--transition-base)"
              (click)="navigateToToday()"
            >
              Back to Today
            </button>
            <button
              style="font-size: 13px; padding: 10px 24px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 600; cursor: pointer; font-family: var(--font-display); box-shadow: var(--shadow-sm); transition: var(--transition-base)"
              (click)="reviewAgain()"
            >
              Review Again
            </button>
          </div>
        </div>
      } @else {
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
          <span>Card {{ cardIndex() + 1 }} of {{ totalCards() }}</span>
          <span>Est. {{ remaining() }} min remaining</span>
        </div>
      </div>

      <!-- Flashcard with 3D Flip -->
      <div class="w-full max-w-[720px]" [@cardSwap]="cardIndex()">
        <div class="flip-card" [class.is-flipped]="revealed()">
          <div class="flip-card-inner">
            <!-- Front Face -->
            <div
              class="review-card flip-card-front"
              [class.review-card--unrevealed]="!revealed()"
              (click)="handleReveal()"
              (mouseenter)="isCardHovered.set(true)"
              (mouseleave)="isCardHovered.set(false)"
            >
              <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 1px; color: var(--ink-faint); font-family: var(--mono); font-weight: 700">
                Recall
              </div>
              <div class="flex items-center gap-2" style="font-size: 12.5px; color: var(--navy); font-weight: 600">
                <lucide-icon name="book-open" [size]="15" [strokeWidth]="2" /> {{ currentCard().src }}
              </div>
              <div style="font-family: var(--font); font-size: 24px; font-weight: 500; line-height: 1.5; max-width: 560px; color: var(--ink); margin-top: 20px">
                {{ currentCard().q }}
              </div>
              <div style="font-size: 13px; color: var(--ink-faint); margin-top: 40px; font-weight: 500">
                Click or press Space to reveal
              </div>
            </div>

            <!-- Back Face -->
            <div class="review-card flip-card-back">
              <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 1px; color: var(--ink-faint); font-family: var(--mono); font-weight: 700">
                Answer
              </div>
              <div class="flex items-center gap-2" style="font-size: 12.5px; color: var(--emerald); font-weight: 600">
                <lucide-icon name="check-circle" [size]="15" [strokeWidth]="2" /> Correct Response
              </div>
              <div
                style="font-family: var(--font); font-size: 18px; color: var(--ink); line-height: 1.7; max-width: 560px; font-weight: 400; margin-top: 20px"
              >
                {{ currentCard().a }}
              </div>
              <a
                class="review-source-link flex items-center gap-1.5 cursor-pointer transition-all"
                style="margin-top: 40px"
                (click)="$event.stopPropagation(); navigateToNotes()"
              >
                <lucide-icon name="file-text" [size]="14" [strokeWidth]="2" /> Source: {{ currentCard().src }}
              </a>
            </div>
          </div>
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
                  (click)="handleRating(btn.confidence)"
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
      }
      }
    </div>
  `,
  styles: [`
    :host { display: flex; flex-direction: column; overflow: hidden; }

    @keyframes pulse {
      0%, 100% { opacity: 1; }
      50% { opacity: 0.4; }
    }
    .skeleton {
      animation: pulse 1.5s ease-in-out infinite;
    }

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
      min-height: 420px;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 12px;
      background: var(--bg);
      box-shadow: var(--shadow-md);
      transition: all 400ms cubic-bezier(0.23, 1, 0.32, 1);
      text-align: center;
      position: relative;
    }
    
    .flip-card {
      background-color: transparent;
      width: 100%;
      height: 420px;
      perspective: 1500px;
    }

    .flip-card-inner {
      position: relative;
      width: 100%;
      height: 100%;
      text-align: center;
      transition: transform 0.7s cubic-bezier(0.4, 0, 0.2, 1);
      transform-style: preserve-3d;
    }

    .flip-card.is-flipped .flip-card-inner {
      transform: rotateY(180deg);
    }

    .flip-card-front, .flip-card-back {
      position: absolute;
      width: 100%;
      height: 100%;
      -webkit-backface-visibility: hidden;
      backface-visibility: hidden;
      top: 0;
      left: 0;
    }

    .flip-card-back {
      transform: rotateY(180deg);
      background: var(--surface-sub);
      border-color: var(--emerald-border);
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
export default class ReviewTabComponent implements OnInit {
  private readonly router = inject(Router);
  protected readonly learningService = inject(LearningService);

  readonly sessionId = signal(crypto.randomUUID());
  readonly cardStartedAt = signal(Date.now());

  cardIndex = signal(0);
  revealed = signal(false);
  confused = signal(false);
  stats = signal({ correct: 0, again: 0, confused: 0 });
  isCardHovered = signal(false);
  readonly sessionDone = signal(false);

  readonly allCards = computed(() => this.learningService.dueCards());

  readonly totalCards = computed(() => this.allCards().length);

  currentCard = computed(() => {
    const list = this.allCards();
    const i = this.cardIndex();
    const d = list[i];
    if (!d) {
      return { q: '', a: '', src: '' };
    }
    return {
      q: d.prompt,
      a: d.answer,
      src: d.topic_name || 'Review',
    };
  });

  progress = computed(() => {
    const t = this.totalCards();
    if (!t) {
      return 0;
    }
    return ((this.cardIndex() + 1) / t) * 100;
  });

  remaining = computed(() => {
    const t = this.totalCards();
    const idx = this.cardIndex();
    return Math.max(1, Math.ceil((t - idx) * 0.5));
  });

  readonly accuracy = computed(() => {
    const total = this.stats().correct + this.stats().again;
    if (!total) return 0;
    return Math.round((this.stats().correct / total) * 100);
  });

  ratingButtons = [
    { label: '😌 Easy', sub: '+10 days', className: 'easy', delay: 0, confidence: 4 as const },
    { label: '👍 Good', sub: '+4 days', className: 'good', delay: 60, confidence: 3 as const },
    { label: '😬 Hard', sub: '+1 day', className: 'hard', delay: 120, confidence: 2 as const },
    { label: '🔁 Again', sub: '10 min', className: 'again', delay: 180, confidence: 1 as const },
  ];

  constructor() {
    effect(() => {
      this.cardIndex();
      this.allCards();
      this.cardStartedAt.set(Date.now());
    });
  }

  ngOnInit(): void {
    void this.learningService.loadDueCards(undefined, 20);
  }

  handleReveal(): void {
    this.revealed.set(true);
  }

  async handleRating(confidence: number): Promise<void> {
    const list = this.allCards();
    const idx = this.cardIndex();
    const row = list[idx];
    if (!row) {
      return;
    }
    const responseTimeMs = Math.max(0, Date.now() - this.cardStartedAt());
    const correct = confidence >= 3;
    const apiConfused = this.confused() || confidence === 1;
    try {
      await this.learningService.submitReview(
        this.sessionId(),
        row.flashcard_id,
        correct,
        confidence,
        apiConfused,
        responseTimeMs,
      );
    } catch {
      return;
    }

    this.stats.update((s) => ({
      correct: s.correct + (correct ? 1 : 0),
      again: s.again + (confidence === 1 ? 1 : 0),
      confused: s.confused + ((this.confused() || confidence === 1) ? 1 : 0),
    }));

    if (idx + 1 < list.length) {
      this.cardIndex.set(idx + 1);
      this.revealed.set(false);
      this.confused.set(false);
    } else {
      this.sessionDone.set(true);
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

  reviewAgain(): void {
    this.cardIndex.set(0);
    this.revealed.set(false);
    this.confused.set(false);
    this.sessionId.set(crypto.randomUUID());
    this.stats.set({ correct: 0, again: 0, confused: 0 });
    this.sessionDone.set(false);
    void this.learningService.loadDueCards(undefined, 20);
  }

  @HostListener('window:keydown', ['$event'])
  handleKeyboard(event: KeyboardEvent) {
    if (event.code === 'Space' && !this.revealed()) {
      event.preventDefault();
      this.handleReveal();
    } else if (this.revealed()) {
      if (event.key === '1') this.handleRating(1);
      if (event.key === '2') this.handleRating(2);
      if (event.key === '3') this.handleRating(3);
      if (event.key === '4') this.handleRating(4);
    }
  }

}
