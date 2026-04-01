import { Component, signal, inject, computed, effect, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { LucideAngularModule } from 'lucide-angular';
import { trigger, transition, style, animate } from '@angular/animations';
import { LearningService } from '../features/learning/learning.service';
import { CanvasService } from '../features/canvas/canvas.service';

interface CourseRow {
  id: string;
  name: string;
  iconColor: string;
  notes: string[];
}

interface MasteryBox {
  topicTag: string;
  masteryScore: number;
  title: string;
  preview: string;
  connectedCards: number;
  confusion: boolean;
}

@Component({
  selector: 'app-notes-tab',
  standalone: true,
  imports: [LucideAngularModule],
  animations: [
    trigger('expandCollapse', [
      transition(':enter', [
        style({ height: 0, opacity: 0, overflow: 'hidden' }),
        animate('200ms cubic-bezier(0.4, 0, 0.2, 1)', style({ height: '*', opacity: 1 })),
      ]),
      transition(':leave', [
        style({ overflow: 'hidden' }),
        animate('200ms cubic-bezier(0.4, 0, 0.2, 1)', style({ height: 0, opacity: 0 })),
      ]),
    ]),
  ],
  template: `
    <div class="flex flex-col flex-1 min-h-0 overflow-hidden" style="padding: 0; gap: 0;">
      <!-- Topbar -->
      <div style="padding: 18px 32px 14px; border-bottom: 1px solid var(--divider); flex-shrink: 0; background: var(--bg);">
        <div class="flex items-center justify-between gap-4">
          <div>
            <div style="font-size: 26px; font-weight: 700; letter-spacing: -0.6px; font-family: var(--font-display); color: var(--ink);">Notes</div>
            <div style="font-size: 13.5px; color: var(--ink-muted); margin-top: 4px;">Your knowledge base · {{ learningService.courses().length }} courses · 14 notes</div>
          </div>
          <div class="flex items-center gap-2">
            <button
              class="btn-import flex items-center gap-1.5"
              style="font-size: 13px; padding: 8px 16px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: transparent; color: var(--ink-2); font-weight: 600; cursor: pointer;"
            >
              <lucide-icon name="upload" [size]="15" [strokeWidth]="2" /> Import
            </button>
            <button
              class="btn-new-note flex items-center gap-1.5"
              style="font-size: 13px; padding: 8px 16px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 600; cursor: pointer; box-shadow: var(--shadow-sm);"
            >
              <lucide-icon name="plus" [size]="15" [strokeWidth]="2" /> New Note
            </button>
          </div>
        </div>
      </div>

      <!-- Synapse Narrator Header -->
      <div
        class="synapse-header flex items-center gap-4"
        style="padding: 18px 32px; background: var(--emerald-light); border: 1px solid var(--emerald-border); border-left: none; border-right: none;"
      >
        <div
          class="pulse-dot"
          style="width: 10px; height: 10px; border-radius: 50%; background: var(--emerald); flex-shrink: 0;"
        ></div>
        <div style="font-size: 13px; color: var(--ink-2); line-height: 1.6; font-weight: 500;">
          <strong style="color: var(--ink); font-family: var(--font-display);">Alex</strong>, you have
          <strong style="color: var(--red);">3 confusion hotspots</strong> in Physics. Resolving these adds
          <strong style="color: var(--emerald);">12%</strong> to your exam readiness.
        </div>
      </div>

      <!-- Notes Layout -->
      <div class="flex flex-1 overflow-hidden">
        <!-- Notes Tree -->
        <div style="width: 206px; flex-shrink: 0; background: var(--sidebar-bg); border-right: 1px solid var(--divider); padding: 0 8px; overflow-y: auto;">
          <div style="padding: 14px 6px 8px; font-size: 10px; text-transform: uppercase; letter-spacing: 0.7px; color: var(--ink-faint); font-weight: 600;">
            Courses
          </div>

          @if (learningService.loading()) {
            @for (s of [1, 2, 3]; track s) {
              <div class="mb-2" style="padding: 5px 6px;">
                <div
                  class="skeleton"
                  style="height: 14px; border-radius: var(--r-md); width: 140px; background: var(--surface-sub); animation: pulse 1.5s ease-in-out infinite"
                ></div>
              </div>
            }
          } @else {
            @for (course of displayCourses(); track course.id) {
              <div class="mb-0.5">
                <div
                  class="flex items-center gap-1.5 cursor-pointer transition-all duration-100"
                  style="padding: 5px 6px; border-radius: var(--r-md); font-size: 13px; font-weight: 500; color: var(--ink-2);"
                  (click)="toggleCourse(course.id)"
                >
                  <lucide-icon
                    name="chevron-right"
                    [size]="8"
                    style="color: var(--ink-faint); transition: transform 0.16s;"
                    [style.transform]="expandedCourse() === course.id ? 'rotate(90deg)' : 'rotate(0deg)'"
                  />
                  <lucide-icon
                    name="book-open"
                    [size]="14"
                    [strokeWidth]="2"
                    [style.color]="course.iconColor"
                  />
                  <span>{{ course.name }}</span>
                </div>
                @if (expandedCourse() === course.id) {
                  <div @expandCollapse style="padding-left: 14px;">
                    @for (note of course.notes; track note) {
                      <div
                        class="note-item flex items-center gap-1.5 cursor-pointer"
                        [class.note-item-active]="course.id === 'discrete' && note === 'Induction Overview'"
                      >
                        <lucide-icon name="file-text" [size]="14" [strokeWidth]="2" style="position: relative; z-index: 1;" />
                        <span style="position: relative; z-index: 1;">{{ note }}</span>
                      </div>
                    }
                  </div>
                }
              </div>
            }
          }

          @if (selectedCourseId()) {
            <div
              class="card"
              style="margin-top: 14px; padding: 14px; border: 1px solid #EAEAEA; border-radius: var(--r-xl); box-shadow: var(--shadow-md); background: #FFFFFF;"
            >
              <div style="font-size: 13px; font-weight: 600; color: var(--ink); margin-bottom: 8px">
                Generate Flashcards with AI
              </div>
              <textarea
                [value]="noteInput()"
                (input)="onNoteInput($event)"
                placeholder="Paste your lecture notes here..."
                style="width:100%; min-height:120px; background:var(--surface-sub);
                       border:1px solid var(--divider); border-radius:var(--r-md);
                       padding:10px 12px; font-size:13.5px; font-family:var(--font);
                       color:var(--ink); resize:vertical; outline:none"
              ></textarea>
              <button
                type="button"
                class="submit-btn"
                [disabled]="!noteInput().trim() || learningService.generating()"
                (click)="onGenerate()"
                style="margin-top: 10px; font-size: 13px; font-weight: 600; padding: 8px 16px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; cursor: pointer; width: 100%; display: flex; align-items: center; justify-content: center; gap: 8px;"
                [style.opacity]="!noteInput().trim() || learningService.generating() ? '0.55' : '1'"
              >
                @if (learningService.generating()) {
                  <lucide-icon name="refresh-cw" [size]="14" class="spin-icon" />
                  Generating...
                } @else {
                  Generate Flashcards
                }
              </button>
              @if (learningService.generationError()) {
                <div style="margin-top: 8px; font-size: 12.5px; color: var(--red);">
                  {{ learningService.generationError() }}
                </div>
              }
            </div>
          }

          @if (learningService.generatedCards().length > 0) {
            <div
              class="card"
              style="margin-top: 12px; padding: 14px; border: 1px solid #EAEAEA; border-radius: var(--r-xl); box-shadow: var(--shadow-md); background: #FFFFFF;"
            >
              <div style="font-size: 13px; font-weight: 600; color: var(--ink); margin-bottom: 12px">
                Review Generated Cards ({{ learningService.generatedCards().length }})
              </div>
              @for (card of learningService.generatedCards(); track card.prompt) {
                <div
                  class="generated-card"
                  style="border: 1px solid var(--divider);
                    border-radius: var(--r-md); padding: 12px; margin-bottom: 8px"
                >
                  <div style="font-size: 13px; font-weight: 500; color: var(--ink)">{{ card.prompt }}</div>
                  <div style="font-size: 12.5px; color: var(--ink-muted); margin-top: 4px">{{ card.answer }}</div>
                </div>
              }
              <div style="display:flex; gap:8px; margin-top:8px; flex-wrap: wrap;">
                <button
                  type="button"
                  (click)="onAcceptAll()"
                  style="font-size: 12px; font-weight: 700; padding: 8px 14px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; cursor: pointer; font-family: var(--font-display); flex: 1; min-width: 120px;"
                >
                  Add All to Deck ({{ learningService.generatedCards().length }})
                </button>
                <button
                  type="button"
                  (click)="learningService.generatedCards.set([])"
                  style="font-size: 12px; font-weight: 600; padding: 8px 14px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: transparent; color: var(--ink-2); cursor: pointer;"
                >
                  Discard
                </button>
              </div>
            </div>
          }
        </div>

        <!-- Note Mastery Boxes View -->
        <div class="flex-1 overflow-y-auto" style="padding: 32px 40px; background: var(--bg);">
          <div class="flex flex-col gap-5">
            @for (box of masteryBoxes; track box.title) {
              <div class="mastery-box">
                <!-- Header -->
                <div class="flex items-start justify-between mb-4">
                  <div
                    style="font-size: 11px; font-weight: 600; padding: 4px 12px; border-radius: var(--r-md); background: var(--navy-light); color: var(--navy); border: 1px solid var(--navy-border);"
                  >
                    {{ box.topicTag }}
                  </div>
                  <div
                    class="flex items-center gap-1.5"
                    style="font-size: 12px; font-weight: 700; font-family: var(--font-display); padding: 5px 12px; border-radius: var(--r-md);"
                    [style.background]="box.confusion ? 'var(--red-light)' : 'var(--emerald-light)'"
                    [style.color]="box.confusion ? 'var(--red)' : 'var(--emerald)'"
                    [style.border]="box.confusion ? '1px solid var(--red-border)' : '1px solid var(--emerald-border)'"
                  >
                    @if (box.confusion) {
                      <lucide-icon name="alert-circle" [size]="13" [strokeWidth]="2" />
                    }
                    {{ box.masteryScore }}% Ready
                  </div>
                </div>

                <!-- Title -->
                <div style="font-size: 17px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 10px; letter-spacing: -0.3px;">
                  {{ box.title }}
                </div>

                <!-- Preview -->
                <div
                  class="line-clamp-3"
                  style="font-family: var(--font); font-size: 13.5px; color: var(--ink-muted); line-height: 1.7; margin-bottom: 16px; font-weight: 400;"
                >
                  {{ box.preview }}
                </div>

                <!-- Footer -->
                <div class="flex items-center justify-between pt-4" style="border-top: 1px solid var(--divider);">
                  <div
                    class="provenance-pill inline-flex items-center gap-1.5 cursor-pointer"
                    style="font-size: 11.5px; color: var(--ink-faint); background: var(--surface-sub); border: 1px solid var(--divider); padding: 5px 11px; border-radius: var(--r-md); transition: var(--transition-fast);"
                  >
                    <lucide-icon name="link" [size]="12" [strokeWidth]="2" /> {{ box.connectedCards }} Connected Cards
                  </div>
                  <div class="flex gap-2">
                    <button
                      class="btn-generate"
                      style="font-size: 12px; font-weight: 700; padding: 8px 14px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: transparent; color: var(--ink-2); cursor: pointer; font-family: var(--font-display);"
                      (click)="navigateToReview($event)"
                    >
                      Generate Cards
                    </button>
                    @if (box.confusion) {
                      <button
                        class="btn-confusing"
                        style="font-size: 12px; font-weight: 700; padding: 8px 14px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; cursor: pointer; font-family: var(--font-display);"
                        (click)="$event.stopPropagation()"
                      >
                        Mark Confusing
                      </button>
                    } @else {
                      <button
                        class="btn-review"
                        style="font-size: 12px; font-weight: 700; padding: 8px 14px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; cursor: pointer; font-family: var(--font-display);"
                        (click)="navigateToReview($event)"
                      >
                        Start Review
                      </button>
                    }
                  </div>
                </div>
              </div>
            }
          </div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    :host { display: flex; flex-direction: column; overflow: hidden; }

    /* Synapse narrator fade-in */
    .synapse-header {
      animation: fadeInUp 0.4s ease-out;
    }

    @keyframes fadeInUp {
      from {
        opacity: 0;
        transform: translateY(10px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }

    /* Pulsing dot */
    .pulse-dot {
      animation: pulseDot 2s infinite;
    }

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

    /* Note item base */
    .note-item {
      padding: 4px 6px;
      border-radius: var(--r-md);
      font-size: 12.5px;
      color: var(--ink-muted);
      background: transparent;
      font-weight: 400;
      position: relative;
      transition: all 0.1s ease;
    }

    .note-item:hover {
      color: var(--ink);
      background: var(--hover-bg);
      transform: translateX(2px) scale(1.01);
    }

    .note-item:active {
      transform: scale(0.98);
    }

    /* Note item active state */
    .note-item-active {
      color: var(--navy) !important;
      background: var(--active-bg) !important;
      font-weight: 600 !important;
    }

    .note-item-active:hover {
      color: var(--navy) !important;
      background: var(--active-bg) !important;
    }

    /* Mastery box */
    .mastery-box {
      border: 1px solid #EAEAEA;
      border-radius: var(--r-xl);
      padding: 28px 32px;
      background: #FFFFFF;
      box-shadow: 0 4px 20px -2px rgba(0, 0, 0, 0.05);
      transition: all var(--transition-base);
      cursor: pointer;
    }

    .mastery-box:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 24px -2px rgba(0, 0, 0, 0.08);
    }

    /* Preview text clamp */
    .line-clamp-3 {
      display: -webkit-box;
      -webkit-line-clamp: 3;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }

    /* Provenance pill hover */
    .provenance-pill {
      transition: var(--transition-fast);
    }

    .provenance-pill:hover {
      transform: scale(1.02);
      border-color: var(--navy) !important;
      color: var(--navy) !important;
    }

    /* Button hover/active effects */
    .btn-import,
    .btn-new-note,
    .btn-generate,
    .btn-review,
    .btn-confusing {
      transition: transform 0.15s cubic-bezier(0.4, 0, 0.2, 1),
                  box-shadow 0.15s cubic-bezier(0.4, 0, 0.2, 1);
    }

    .btn-import:hover {
      transform: scale(1.02);
    }

    .btn-import:active {
      transform: scale(0.98);
    }

    .btn-new-note:hover {
      transform: scale(1.02);
      box-shadow: var(--shadow-md);
    }

    .btn-new-note:active {
      transform: scale(0.98);
    }

    .btn-generate:hover {
      transform: scale(1.05);
    }

    .btn-generate:active {
      transform: scale(0.95);
    }

    .btn-review:hover {
      transform: scale(1.05);
    }

    .btn-review:active {
      transform: scale(0.95);
    }

    .btn-confusing:hover {
      transform: scale(1.05);
    }

    .btn-confusing:active {
      transform: scale(0.95);
    }

    @keyframes pulse {
      0%, 100% { opacity: 1; }
      50% { opacity: 0.4; }
    }

    @keyframes spin {
      to { transform: rotate(360deg); }
    }

    .spin-icon {
      animation: spin 0.8s linear infinite;
      display: inline-block;
    }
  `],
})
export default class NotesTabComponent implements OnInit {
  private readonly router = inject(Router);
  protected readonly learningService = inject(LearningService);
  private readonly canvasService = inject(CanvasService);

  private readonly placeholderNotes = [
    'Induction Overview',
    'Recursion Notes',
    'Set Theory Basics',
    'Logic & Proof',
  ];

  readonly displayCourses = computed<CourseRow[]>(() =>
    this.learningService.courses().map((c) => ({
      id: c.id,
      name: c.name,
      iconColor: c.color || '#102E67',
      notes: this.placeholderNotes,
    })),
  );

  expandedCourse = signal('');
  readonly noteInput = signal('');
  readonly selectedCourseId = signal<string | null>(null);

  constructor() {
    effect(() => {
      const rows = this.displayCourses();
      const cur = this.expandedCourse();
      if (rows.length && !rows.some((r) => r.id === cur)) {
        const first = rows[0].id;
        this.expandedCourse.set(first);
        this.selectedCourseId.set(first);
      }
    });
  }

  ngOnInit(): void {
    void this.learningService.loadCourses();
    void this.canvasService.loadStatus();
  }

  masteryBoxes: MasteryBox[] = [
    {
      topicTag: 'Mathematics · Induction',
      masteryScore: 82,
      title: 'Induction Overview — Complete Guide',
      preview:
        'Mathematical induction is a proof technique used to prove that a statement P(n) is true for all natural numbers n. It consists of two steps: the base case and the inductive step...',
      connectedCards: 5,
      confusion: false,
    },
    {
      topicTag: 'Computer Science · Recursion',
      masteryScore: 29,
      title: 'Recursion Notes — Base Cases & Stack Traces',
      preview:
        'Recursion is when a function calls itself. Every recursive function must have a base case to prevent infinite loops. Understanding the call stack is critical for debugging recursive algorithms...',
      connectedCards: 8,
      confusion: true,
    },
    {
      topicTag: 'Mathematics · Set Theory',
      masteryScore: 78,
      title: 'Set Theory Basics — Operations & Notation',
      preview:
        'A set is a collection of distinct objects. Sets can be defined by enumeration or by property. Common operations include union, intersection, difference, and complement...',
      connectedCards: 4,
      confusion: false,
    },
    {
      topicTag: 'Mathematics · Logic',
      masteryScore: 51,
      title: 'Logic & Proof — Fundamentals',
      preview:
        'Propositional logic deals with statements that are either true or false. Logical connectives include AND, OR, NOT, implies, and if-and-only-if. Proofs can be direct, by contradiction, or by contrapositive...',
      connectedCards: 6,
      confusion: false,
    },
  ];

  toggleCourse(id: string): void {
    this.selectedCourseId.set(id);
    this.expandedCourse.set(this.expandedCourse() === id ? '' : id);
  }

  onNoteInput(ev: Event): void {
    const t = ev.target as HTMLTextAreaElement | null;
    this.noteInput.set(t?.value ?? '');
  }

  onGenerate(): void {
    const cid = this.selectedCourseId();
    if (!cid) {
      return;
    }
    void this.learningService.generateFlashcards(cid, this.noteInput(), undefined);
  }

  onAcceptAll(): void {
    const cid = this.selectedCourseId();
    if (!cid) {
      return;
    }
    void this.learningService
      .acceptGeneratedCards(cid, this.learningService.generatedCards())
      .then(() => this.noteInput.set(''));
  }

  navigateToReview(event: Event): void {
    event.stopPropagation();
    this.router.navigate(['/review']);
  }
}
