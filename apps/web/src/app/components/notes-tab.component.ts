import { Component, signal, inject, computed, effect, OnInit, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { LucideAngularModule } from 'lucide-angular';
import { trigger, transition, style, animate } from '@angular/animations';
import { LearningService, NoteText } from '../features/learning/learning.service';
import { CanvasService } from '../features/canvas/canvas.service';
import { AuthService } from '../core/auth/auth.service';

interface CourseRow {
  id: string;
  name: string;
  iconColor: string;
}

interface NoteOverviewCard {
  note: NoteText;
  topicTag: string;
  readiness: number;
  connectedCards: number;
  confusion: boolean;
  hasActivity: boolean;
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
            <div style="font-size: 13.5px; color: var(--ink-muted); margin-top: 4px;">
              Your knowledge base · {{ learningService.courses().length }} courses · {{ learningService.notes().length }} notes
            </div>
          </div>
          <div class="flex items-center gap-2">
            <button
              class="btn-import flex items-center gap-1.5"
              (click)="onImportFromCanvas()"
              [disabled]="importingCanvas()"
              [style.opacity]="importingCanvas() ? '0.6' : '1'"
              style="font-size: 13px; padding: 8px 16px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: transparent; color: var(--ink-2); font-weight: 600; cursor: pointer;"
            >
              @if (importingCanvas()) {
                <lucide-icon name="refresh-cw" [size]="15" [strokeWidth]="2" class="spin-icon" /> Importing...
              } @else {
                <lucide-icon name="upload" [size]="15" [strokeWidth]="2" /> Import
              }
            </button>
            <button
              class="btn-new-note flex items-center gap-1.5"
              style="font-size: 13px; padding: 8px 16px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 600; cursor: pointer; box-shadow: var(--shadow-sm);"
              (click)="onNewNote()"
              [disabled]="!selectedCourseId()"
              [style.opacity]="!selectedCourseId() ? '0.5' : '1'"
            >
              <lucide-icon name="plus" [size]="15" [strokeWidth]="2" /> New Note
            </button>
          </div>
        </div>
      </div>

      <!-- Synapse Narrator Banner (only shown when no note is selected) -->
      @if (!learningService.activeNote()) {
        <div
          class="synapse-header flex items-center gap-4"
          style="padding: 18px 32px; background: var(--emerald-light); border: 1px solid var(--emerald-border); border-left: none; border-right: none;"
        >
          <div class="pulse-dot" style="width: 10px; height: 10px; border-radius: 50%; background: var(--emerald); flex-shrink: 0;"></div>
          <div style="font-size: 13px; color: var(--ink-2); line-height: 1.6; font-weight: 500;">
            @if (hotspotCount() === 0) {
              <strong style="color: var(--ink); font-family: var(--font-display);">{{ firstName() }}</strong>, you have no confusion hotspots in the last {{ confusionWindowDays() }} days.
            } @else {
              <strong style="color: var(--ink); font-family: var(--font-display);">{{ firstName() }}</strong>, you have
              <strong style="color: var(--red);">{{ hotspotCount() }} confusion hotspots</strong>
              @if (topTopicName()) {
                <span> in {{ topTopicName() }}</span>
              }
              . Resolving these topics should improve your readiness.
            }
          </div>
        </div>
      }

      <!-- Notes Layout -->
      <div class="flex flex-1 overflow-hidden">
        <!-- Notes Sidebar -->
        <div style="width: 206px; flex-shrink: 0; background: var(--sidebar-bg); border-right: 1px solid var(--divider); padding: 0 8px; overflow-y: auto; display: flex; flex-direction: column;">
          <div style="padding: 14px 6px 8px; font-size: 10px; text-transform: uppercase; letter-spacing: 0.7px; color: var(--ink-faint); font-weight: 600;">
            Courses
          </div>

          @if (learningService.loading()) {
            @for (s of [1, 2, 3]; track s) {
              <div class="mb-2" style="padding: 5px 6px;">
                <div class="skeleton" style="height: 14px; border-radius: var(--r-md); width: 140px; background: var(--surface-sub); animation: pulse 1.5s ease-in-out infinite"></div>
              </div>
            }
          } @else {
            @for (course of displayCourses(); track course.id) {
              <div class="mb-0.5">
                <div
                  class="flex items-center gap-1.5 cursor-pointer transition-all duration-100"
                  style="padding: 5px 6px; border-radius: var(--r-md); font-size: 13px; font-weight: 500; color: var(--ink-2);"
                  [style.background]="selectedCourseId() === course.id ? 'var(--active-bg)' : 'transparent'"
                  [style.color]="selectedCourseId() === course.id ? 'var(--navy)' : 'var(--ink-2)'"
                  (click)="toggleCourse(course.id)"
                >
                  <lucide-icon
                    name="chevron-right"
                    [size]="8"
                    style="color: var(--ink-faint); transition: transform 0.16s;"
                    [style.transform]="expandedCourse() === course.id ? 'rotate(90deg)' : 'rotate(0deg)'"
                  />
                  <lucide-icon name="book-open" [size]="14" [strokeWidth]="2" [style.color]="course.iconColor" />
                  <span>{{ course.name }}</span>
                </div>
                @if (expandedCourse() === course.id) {
                  <div @expandCollapse style="padding-left: 14px;">
                    @if (learningService.notesLoading()) {
                      <div style="padding: 4px 6px; font-size: 12px; color: var(--ink-faint);">Loading...</div>
                    } @else if (learningService.notes().length === 0) {
                      <div style="padding: 4px 6px; font-size: 12px; color: var(--ink-faint); font-style: italic;">No notes yet</div>
                    } @else {
                      @for (note of learningService.notes(); track note.id) {
                        <div
                          class="note-item flex items-center gap-1.5 cursor-pointer"
                          [class.note-item-active]="learningService.activeNote()?.id === note.id"
                          (click)="openNote(note)"
                        >
                          <lucide-icon name="file-text" [size]="14" [strokeWidth]="2" style="position: relative; z-index: 1; flex-shrink: 0;" />
                          <span style="position: relative; z-index: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">{{ note.title || 'Untitled' }}</span>
                        </div>
                      }
                    }
                  </div>
                }
              </div>
            }
          }

        </div>

        <!-- Main Panel -->
        <div class="flex-1 overflow-y-auto" style="background: var(--bg);">

          @if (learningService.activeNote(); as note) {
            <!-- NOTE EDITOR -->
            <div class="note-editor" style="max-width: 780px; margin: 0 auto; padding: 40px 48px 80px;">

              <!-- Title -->
              <input
                class="note-title-input"
                [value]="editorTitle()"
                (input)="onTitleInput($event)"
                placeholder="Untitled"
                style="width: 100%; font-size: 32px; font-weight: 700; letter-spacing: -0.8px;
                       font-family: var(--font-display); color: var(--ink); border: none;
                       background: transparent; outline: none; margin-bottom: 4px; padding: 0;"
              />

              <!-- Toolbar -->
              <div class="flex items-center gap-3 mb-6" style="border-bottom: 1px solid var(--divider); padding-bottom: 12px;">
                <button
                  type="button"
                  class="toolbar-btn flex items-center gap-1.5"
                  (click)="onSave()"
                  [disabled]="learningService.noteSaving()"
                  style="font-size: 12px; font-weight: 600; padding: 5px 12px; border-radius: var(--r-md); border: 1px solid var(--divider); background: transparent; color: var(--ink-2); cursor: pointer;"
                >
                  @if (learningService.noteSaving()) {
                    <lucide-icon name="refresh-cw" [size]="12" class="spin-icon" /> Saving...
                  } @else {
                    <lucide-icon name="save" [size]="12" /> Save
                  }
                </button>
                <span style="font-size: 11.5px; color: var(--ink-faint);">{{ saveStatus() }}</span>
                <select
                  [value]="selectedTopicId() ?? ''"
                  (change)="onTopicChange($event)"
                  style="font-size: 12px; font-weight: 500; padding: 5px 10px; border-radius: var(--r-md); border: 1px solid var(--divider); background: #fff; color: var(--ink-2);"
                >
                  <option value="">No topic</option>
                  @for (topic of learningService.topics(); track topic.id) {
                    <option [value]="topic.id">{{ topic.name }}</option>
                  }
                  <option value="__NEW__">+ Create new topic...</option>
                </select>
                <div style="flex: 1;"></div>
                <button
                  type="button"
                  class="toolbar-btn flex items-center gap-1.5"
                  (click)="toggleAiPanel()"
                  style="font-size: 12px; font-weight: 600; padding: 5px 12px; border-radius: var(--r-md); border: none; cursor: pointer;"
                  [style.background]="showAiPanel() ? 'var(--navy)' : 'var(--emerald-light)'"
                  [style.color]="showAiPanel() ? '#fff' : 'var(--emerald)'"
                  [style.border]="showAiPanel() ? 'none' : '1px solid var(--emerald-border)'"
                >
                  <lucide-icon name="sparkles" [size]="12" /> Ask Synapse AI
                </button>
                <button
                  type="button"
                  class="toolbar-btn flex items-center gap-1.5"
                  (click)="onDeleteNote()"
                  style="font-size: 12px; font-weight: 600; padding: 5px 12px; border-radius: var(--r-md); border: 1px solid var(--red-border); background: transparent; color: var(--red); cursor: pointer;"
                >
                  <lucide-icon name="trash-2" [size]="12" /> Delete
                </button>
              </div>

              <!-- AI Panel -->
              @if (showAiPanel()) {
                <div class="ai-panel" style="background: var(--navy-light); border: 1px solid var(--navy-border); border-radius: var(--r-xl); padding: 20px 24px; margin-bottom: 28px;">
                  <div style="font-size: 13px; font-weight: 700; color: var(--navy); margin-bottom: 12px; font-family: var(--font-display);">
                    ✦ Synapse AI · Ask about this note
                  </div>
                  <div style="font-size: 12px; color: var(--ink-muted); margin-bottom: 12px;">
                    I have context of your note, courses, and mastery levels.
                  </div>
                  <div class="flex gap-2">
                    <input
                      class="ai-input"
                      [value]="aiQuestion()"
                      (input)="onAiQuestionInput($event)"
                      (keydown.enter)="onAskAI()"
                      placeholder="e.g. Explain the base case in simpler terms..."
                      style="flex: 1; font-size: 13px; padding: 9px 14px; border-radius: var(--r-md); border: 1px solid var(--divider); background: #fff; color: var(--ink); outline: none; font-family: var(--font);"
                    />
                    <button
                      type="button"
                      (click)="onAskAI()"
                      [disabled]="!aiQuestion().trim() || learningService.noteAiLoading()"
                      style="font-size: 13px; font-weight: 700; padding: 9px 18px; border-radius: var(--r-md); border: none; background: var(--navy); color: #fff; cursor: pointer; white-space: nowrap;"
                      [style.opacity]="!aiQuestion().trim() || learningService.noteAiLoading() ? '0.6' : '1'"
                    >
                      @if (learningService.noteAiLoading()) {
                        <lucide-icon name="refresh-cw" [size]="13" class="spin-icon" />
                      } @else {
                        Ask
                      }
                    </button>
                  </div>
                  @if (learningService.noteAiAnswer()) {
                    <div style="margin-top: 16px; font-size: 13.5px; color: var(--ink); line-height: 1.75; white-space: pre-wrap; background: #fff; border-radius: var(--r-lg); padding: 16px 18px; border: 1px solid var(--divider);">
                      {{ learningService.noteAiAnswer() }}
                    </div>
                  }
                  @if (learningService.noteAiError()) {
                    <div style="margin-top: 12px; font-size: 12.5px; color: var(--red);">{{ learningService.noteAiError() }}</div>
                  }
                </div>
              }

              <!-- Note Content Textarea -->
              <textarea
                class="note-content-textarea"
                [value]="editorContent()"
                (input)="onContentInput($event)"
                placeholder="Start writing your note... (Markdown is supported)"
                style="width: 100%; min-height: 480px; font-size: 15px; line-height: 1.8;
                       font-family: var(--font); color: var(--ink); border: none;
                       background: transparent; outline: none; resize: none; padding: 0;"
              ></textarea>
            </div>

          } @else {
            <!-- NOTES OVERVIEW (shown when no note is selected) -->
            <div style="padding: 32px 40px; background: var(--bg);">
              <div class="flex flex-col gap-5">
                @if (noteCards().length === 0) {
                  <div class="mastery-box" style="cursor: default;">
                    <div style="font-size: 17px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 10px; letter-spacing: -0.3px;">No notes yet</div>
                    <div style="font-size: 13.5px; color: var(--ink-muted); line-height: 1.7;">Create a new note from the top-right button, then generate cards from that note.</div>
                  </div>
                }
                @for (box of noteCards(); track box.note.id) {
                  <div class="mastery-box" (click)="openNote(box.note)">
                    <div class="flex items-start justify-between mb-4">
                      <div style="font-size: 11px; font-weight: 600; padding: 4px 12px; border-radius: var(--r-md); background: var(--navy-light); color: var(--navy); border: 1px solid var(--navy-border);">
                        {{ box.topicTag }}
                      </div>
                      <div
                        class="flex items-center gap-1.5"
                        style="font-size: 12px; font-weight: 700; font-family: var(--font-display); padding: 5px 12px; border-radius: var(--r-md);"
                        [style.background]="!box.hasActivity ? 'var(--surface-sub)' : (box.confusion ? 'var(--red-light)' : 'var(--emerald-light)')"
                        [style.color]="!box.hasActivity ? 'var(--ink-muted)' : (box.confusion ? 'var(--red)' : 'var(--emerald)')"
                        [style.border]="!box.hasActivity ? '1px solid var(--divider)' : (box.confusion ? '1px solid var(--red-border)' : '1px solid var(--emerald-border)')"
                      >
                        @if (!box.hasActivity) {
                          <lucide-icon name="info" [size]="13" [strokeWidth]="2" />
                          No activity yet
                        } @else {
                          @if (box.confusion) { <lucide-icon name="alert-circle" [size]="13" [strokeWidth]="2" /> }
                          {{ box.readiness }}% Ready
                        }
                      </div>
                    </div>
                    <div style="font-size: 17px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 10px; letter-spacing: -0.3px;">{{ box.note.title || 'Untitled' }}</div>
                    <div class="line-clamp-3" style="font-size: 13.5px; color: var(--ink-muted); line-height: 1.7; margin-bottom: 16px;">{{ box.note.content || 'No content yet.' }}</div>
                    <div class="flex items-center justify-between pt-4" style="border-top: 1px solid var(--divider);">
                      <div class="provenance-pill inline-flex items-center gap-1.5 cursor-pointer" style="font-size: 11.5px; color: var(--ink-faint); background: var(--surface-sub); border: 1px solid var(--divider); padding: 5px 11px; border-radius: var(--r-md);">
                        <lucide-icon name="link" [size]="12" [strokeWidth]="2" /> {{ box.connectedCards }} Connected Cards
                      </div>
                      <div class="flex gap-2">
                        <button class="btn-generate" style="font-size: 12px; font-weight: 700; padding: 8px 14px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: transparent; color: var(--ink-2); cursor: pointer; font-family: var(--font-display);" (click)="onGenerateForNote(box.note, $event)">Generate Cards</button>
                        <button class="btn-review" style="font-size: 12px; font-weight: 700; padding: 8px 14px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; cursor: pointer; font-family: var(--font-display);" (click)="navigateToReview($event)">Start Review</button>
                      </div>
                    </div>
                    @if (generatedForNoteId() === box.note.id && learningService.generatedCards().length > 0) {
                      <div style="margin-top: 14px; border-top: 1px solid var(--divider); padding-top: 14px;">
                        <div style="font-size: 12px; font-weight: 700; color: var(--ink); margin-bottom: 8px;">
                          Review Generated Cards ({{ learningService.generatedCards().length }})
                        </div>
                        @for (card of learningService.generatedCards(); track card.prompt) {
                          <div style="border: 1px solid var(--divider); border-radius: var(--r-md); padding: 10px; margin-bottom: 8px; background: #fff;">
                            <div style="font-size: 12.5px; font-weight: 600; color: var(--ink)">{{ card.prompt }}</div>
                            <div style="font-size: 12px; color: var(--ink-muted); margin-top: 4px">{{ card.answer }}</div>
                          </div>
                        }
                        <div style="display:flex; gap:8px; margin-top:8px; flex-wrap: wrap;">
                          <button type="button" (click)="onAcceptGeneratedFromNote(box.note, $event)" style="font-size: 12px; font-weight: 700; padding: 8px 14px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; cursor: pointer; flex: 1; min-width: 120px;">
                            Add All to Deck ({{ learningService.generatedCards().length }})
                          </button>
                          <button type="button" (click)="discardGenerated($event)" style="font-size: 12px; font-weight: 600; padding: 8px 14px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: transparent; color: var(--ink-2); cursor: pointer;">
                            Discard
                          </button>
                        </div>
                      </div>
                    }
                    @if (generatedForNoteId() === box.note.id && learningService.generationError()) {
                      <div style="margin-top: 10px; font-size: 12.5px; color: var(--red);">{{ learningService.generationError() }}</div>
                    }
                  </div>
                }
              </div>
            </div>
          }
        </div>
      </div>
    </div>
  `,
  styles: [
    `
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

    .note-title-input::placeholder { color: var(--ink-faint); }
    .note-content-textarea::placeholder { color: var(--ink-faint); }
    .note-content-textarea { display: block; }
    .toolbar-btn { transition: transform 0.12s ease, opacity 0.12s ease; }
    .toolbar-btn:hover:not([disabled]) { transform: scale(1.03); }
    .toolbar-btn:active:not([disabled]) { transform: scale(0.97); }
    .ai-input:focus { border-color: var(--navy) !important; }
    .ai-panel { animation: fadeInUp 0.2s ease; }
    `,
  ],
})
export default class NotesTabComponent implements OnInit, OnDestroy {
  private readonly router = inject(Router);
  protected readonly learningService = inject(LearningService);
  private readonly canvasService = inject(CanvasService);
  private readonly authService = inject(AuthService);

  expandedCourse = signal('');
  readonly selectedCourseId = signal<string | null>(null);
  readonly editorTitle = signal('');
  readonly editorContent = signal('');
  readonly selectedTopicId = signal<string | null>(null);
  readonly showAiPanel = signal(false);
  readonly aiQuestion = signal('');
  readonly saveStatus = signal('');
  readonly generatedForNoteId = signal<string | null>(null);
  readonly importingCanvas = signal(false);

  private saveTimer: ReturnType<typeof setTimeout> | null = null;

  readonly displayCourses = computed<CourseRow[]>(() =>
    this.learningService.courses().map((c) => ({
      id: c.id,
      name: c.name,
      iconColor: c.color || '#102E67',
    })),
  );
  readonly noteCards = computed<NoteOverviewCard[]>(() => {
    const topicsByID = new Map(this.learningService.topics().map((t) => [t.id, t]));
    const metricsByNoteID = new Map(this.learningService.noteMetrics().map((m) => [m.note_id, m]));
    return this.learningService.notes().map((note) => {
      const topicName = note.topic_id ? topicsByID.get(note.topic_id)?.name : null;
      const topicTag = topicName ? `${this.courseName()} · ${topicName}` : `${this.courseName()} · Unscoped`;
      const metric = metricsByNoteID.get(note.id);
      const readiness = metric?.readiness_pct ?? 0;
      const connectedCards = metric?.connected_cards ?? 0;
      const hasActivity = (metric?.review_count ?? 0) > 0;
      return {
        note,
        topicTag,
        readiness,
        connectedCards,
        confusion: metric?.confusion_flag ?? false,
        hasActivity,
      };
    });
  });
  readonly courseName = computed(() => {
    const courseID = this.selectedCourseId();
    if (!courseID) return 'Course';
    return this.learningService.courses().find((c) => c.id === courseID)?.name ?? 'Course';
  });
  readonly firstName = computed(() => this.authService.currentUser()?.name?.split(' ')[0] ?? 'Student');
  readonly hotspotCount = computed(() => this.learningService.confusionInsight()?.hotspot_count ?? 0);
  readonly topTopicName = computed(() => this.learningService.confusionInsight()?.top_topic_name ?? '');
  readonly confusionWindowDays = computed(() => this.learningService.confusionInsight()?.window_days ?? 14);

  constructor() {
    effect(() => {
      const note = this.learningService.activeNote();
      if (note) {
        this.editorTitle.set(note.title);
        this.editorContent.set(note.content);
        this.selectedTopicId.set(note.topic_id);
        this.saveStatus.set('');
        this.showAiPanel.set(false);
        this.learningService.noteAiAnswer.set(null);
        this.learningService.noteAiError.set(null);
      }
    });

    effect(() => {
      const rows = this.displayCourses();
      const curSelected = this.selectedCourseId();
      
      if (rows.length) {
        if (!curSelected || !rows.some(r => r.id === curSelected)) {
          const first = rows[0].id;
          this.expandedCourse.set(first);
          this.selectedCourseId.set(first);
          void this.learningService.loadNotes(first);
          void this.learningService.loadTopics(first);
        }
      } else if (!rows.length && curSelected) {
        this.expandedCourse.set('');
        this.selectedCourseId.set(null);
      }
    });
  }

  ngOnInit(): void {
    void this.initializeNotesData();
  }

  ngOnDestroy(): void {
    if (this.saveTimer) {
      clearTimeout(this.saveTimer);
    }
  }

  async onImportFromCanvas(): Promise<void> {
    if (this.importingCanvas()) return;
    this.importingCanvas.set(true);
    try {
      await this.canvasService.triggerSync();
      const synced = await this.canvasService.listSyncedCourses();
      if (synced.length) {
        await this.learningService.importFromLMS(
          synced.map((c) => ({
            lms_course_id: c.lms_course_id,
            name: c.course_name,
            term: c.term,
          }))
        );
        await this.learningService.loadCourses();
      }
    } catch {
      // Best-effort sync
    } finally {
      this.importingCanvas.set(false);
    }
  }

  toggleCourse(id: string): void {
    const isExpanding = this.expandedCourse() !== id;
    this.selectedCourseId.set(id);
    this.expandedCourse.set(isExpanding ? id : '');
    
    // Always clear the active note when clicking the course to return to the overview
    this.learningService.activeNote.set(null);

    if (isExpanding) {
      void this.learningService.loadNotes(id);
      void this.learningService.loadTopics(id);
    }
  }

  openNote(note: NoteText): void {
    this.learningService.activeNote.set(note);
    this.aiQuestion.set('');
  }

  async onNewNote(): Promise<void> {
    const cid = this.selectedCourseId();
    if (!cid) return;
    const note = await this.learningService.createNote(cid, 'Untitled', 'Start writing your note...', this.selectedTopicId() ?? undefined);
    this.learningService.activeNote.set(note);
  }

  onTitleInput(ev: Event): void {
    const val = (ev.target as HTMLInputElement).value;
    this.editorTitle.set(val);
    this.scheduleSave();
  }

  onContentInput(ev: Event): void {
    const val = (ev.target as HTMLTextAreaElement).value;
    this.editorContent.set(val);
    this.scheduleSave();
  }

  private scheduleSave(): void {
    this.saveStatus.set('Unsaved changes');
    if (this.saveTimer) clearTimeout(this.saveTimer);
    this.saveTimer = setTimeout(() => void this.onSave(), 1500);
  }

  async onSave(): Promise<void> {
    const note = this.learningService.activeNote();
    if (!note) return;
    await this.learningService.updateNote(
      note.id,
      this.editorTitle(),
      this.editorContent(),
      this.selectedTopicId(),
    );
    this.saveStatus.set('Saved');
  }

  async onDeleteNote(): Promise<void> {
    const note = this.learningService.activeNote();
    if (!note) return;
    if (!confirm(`Delete "${note.title}"? This cannot be undone.`)) return;
    await this.learningService.deleteNote(note.id);
  }

  toggleAiPanel(): void {
    this.showAiPanel.update((v) => !v);
  }

  onAiQuestionInput(ev: Event): void {
    this.aiQuestion.set((ev.target as HTMLInputElement).value);
  }

  async onTopicChange(ev: Event): Promise<void> {
    const raw = (ev.target as HTMLSelectElement).value;
    if (raw === '__NEW__') {
      const topicName = window.prompt('Enter new topic name:');
      if (topicName && topicName.trim()) {
        const courseId = this.selectedCourseId();
        if (courseId) {
          try {
            const newTopic = await this.learningService.createTopic(courseId, topicName.trim());
            this.selectedTopicId.set(newTopic.id);
            this.scheduleSave();
            return;
          } catch (e) {
            console.error('Failed to create topic', e);
          }
        }
      }
      // Revert if cancelled or failed
      (ev.target as HTMLSelectElement).value = this.selectedTopicId() ?? '';
      return;
    }

    this.selectedTopicId.set(raw || null);
    this.scheduleSave();
  }

  onAskAI(): void {
    const note = this.learningService.activeNote();
    if (!note || !this.aiQuestion().trim()) return;
    void this.learningService.askNoteAI(note.id, this.aiQuestion());
  }

  onGenerateForNote(note: NoteText, event: Event): void {
    event.stopPropagation();
    if (!note.content.trim()) {
      this.learningService.generationError.set('Add some note content before generating cards.');
      this.generatedForNoteId.set(note.id);
      return;
    }
    this.generatedForNoteId.set(note.id);
    void this.learningService.generateFlashcards(note.course_id, note.content, note.topic_id ?? undefined);
  }

  onAcceptGeneratedFromNote(note: NoteText, event: Event): void {
    event.stopPropagation();
    void this.learningService.acceptGeneratedCards(
      note.course_id,
      this.learningService.generatedCards(),
      note.topic_id ?? undefined,
    );
  }

  discardGenerated(event: Event): void {
    event.stopPropagation();
    this.learningService.generatedCards.set([]);
    this.learningService.generationError.set(null);
    this.generatedForNoteId.set(null);
  }

  navigateToReview(event: Event): void {
    event.stopPropagation();
    this.router.navigate(['/review']);
  }

  private async initializeNotesData(): Promise<void> {
    // Load courses from Synapse DB only.
    // Canvas sync + import is triggered explicitly via the Import button.
    await this.learningService.loadCourses();
  }
}
