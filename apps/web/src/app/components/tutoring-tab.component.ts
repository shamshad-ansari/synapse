import { Component, ChangeDetectionStrategy, OnInit, computed, inject, signal } from '@angular/core';
import { LucideAngularModule } from 'lucide-angular';
import { TutoringService, type TutorMatch } from '../features/tutoring/tutoring.service';
import { AuthService } from '../core/auth/auth.service';
import { LearningService } from '../features/learning/learning.service';

type RequestsTab = 'incoming' | 'sent' | 'accepted';

@Component({
  selector: 'app-tutoring-tab',
  standalone: true,
  imports: [LucideAngularModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `
    <div class="flex flex-col overflow-y-auto overflow-x-hidden" style="padding: 44px 56px 56px; gap: 36px">
      <!-- Topbar -->
      <div class="flex flex-col gap-4">
        <div class="flex items-center justify-between gap-4">
          <div>
            <div style="font-size: 26px; font-weight: 700; letter-spacing: -0.6px; font-family: var(--font-display); color: var(--ink)">Tutoring</div>
            <div style="font-size: 13.5px; color: var(--ink-muted); margin-top: 4px">Graph-matched peer learning powered by mastery signals</div>
          </div>
          <button
            type="button"
            class="post-request-btn flex items-center gap-2 transition-all"
            style="font-size: 13px; padding: 8px 16px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 600; cursor: pointer; box-shadow: var(--shadow-sm); transition: var(--transition-base)"
            (click)="showRequestForm.set(true)"
          >
            <lucide-icon name="plus" [size]="15" [strokeWidth]="2" /> Post Request
          </button>
        </div>

        @if (showRequestForm()) {
          <div
            style="padding: 20px 22px; border-radius: var(--r-xl); border: 1px solid var(--divider); background: var(--card-bg); box-shadow: var(--shadow-sm)"
          >
            <div style="font-size: 14px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 14px">New tutoring request</div>
            <div class="flex flex-col gap-3">
              <div>
                <label style="font-size: 11px; font-weight: 600; color: var(--ink-muted); text-transform: uppercase; letter-spacing: 0.5px">Tutor</label>
                <select
                  [value]="postTutorId()"
                  (change)="onPostTutorChange($any($event.target).value)"
                  style="margin-top: 6px; width: 100%; padding: 10px 12px; border-radius: var(--r-md); border: 1px solid var(--divider); font-size: 13px; font-family: var(--font); background: var(--bg); color: var(--ink)"
                >
                  <option value="">Select a suggested tutor…</option>
                  @for (t of tutoringService.tutorMatches(); track t.user_id) {
                    <option [value]="t.user_id">{{ t.name }} — {{ t.topic_name }}</option>
                  }
                </select>
                @if (tutoringService.tutorMatches().length === 0 && !tutoringService.matchesLoading()) {
                  <div style="font-size: 12px; color: var(--ink-muted); margin-top: 6px">
                    No matches yet for your weak topic—try again after more review data, or pick a topic that has strong peers in your school.
                  </div>
                }
              </div>
              <div>
                <label style="font-size: 11px; font-weight: 600; color: var(--ink-muted); text-transform: uppercase; letter-spacing: 0.5px">Topic</label>
                <input
                  type="text"
                  [value]="requestTopic()"
                  (input)="requestTopic.set($any($event.target).value)"
                  placeholder="e.g. Recursion"
                  style="margin-top: 6px; width: 100%; padding: 10px 12px; border-radius: var(--r-md); border: 1px solid var(--divider); font-size: 13px; font-family: var(--font); background: var(--bg); color: var(--ink)"
                />
              </div>
              <div>
                <label style="font-size: 11px; font-weight: 600; color: var(--ink-muted); text-transform: uppercase; letter-spacing: 0.5px">Message</label>
                <textarea
                  [value]="requestMessage()"
                  (input)="requestMessage.set($any($event.target).value)"
                  rows="3"
                  placeholder="What do you want to work on?"
                  style="margin-top: 6px; width: 100%; padding: 10px 12px; border-radius: var(--r-md); border: 1px solid var(--divider); font-size: 13px; font-family: var(--font); background: var(--bg); color: var(--ink); resize: vertical"
                ></textarea>
              </div>
              @if (tutoringService.error()) {
                <div style="font-size: 12px; color: var(--red)">{{ tutoringService.error() }}</div>
              }
              <div class="flex items-center gap-2 justify-end">
                <button
                  type="button"
                  style="padding: 8px 14px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: transparent; color: var(--ink-muted); font-size: 13px; font-weight: 600; cursor: pointer"
                  (click)="showRequestForm.set(false)"
                >
                  Cancel
                </button>
                <button
                  type="button"
                  style="padding: 8px 16px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-size: 13px; font-weight: 600; cursor: pointer; opacity: tutoringService.loading() ? 0.7 : 1"
                  [disabled]="tutoringService.loading() || !requestTopic().trim() || !postTutorId()"
                  (click)="onSubmitPostRequest()"
                >
                  Submit
                </button>
              </div>
            </div>
          </div>
        }
      </div>

      <!-- Narrator -->
      <div
        class="flex items-center gap-4"
        style="padding: 18px 22px; border-radius: var(--r-xl); background: var(--emerald-light); border: 1px solid var(--emerald-border); box-shadow: var(--shadow-sm)"
      >
        <div class="pulse-dot"></div>
        <div style="font-size: 13px; color: var(--ink-2); line-height: 1.6; font-weight: 500">
          <strong style="color: var(--ink); font-family: var(--font-display)">{{ narratorName() }}</strong
          >, peer tutoring suggestions prioritize your weakest area:
          <strong style="color: var(--red)">{{ topWeakTopic() }}</strong
          >.
        </div>
      </div>

      <!-- Two-Column Layout -->
      <div class="grid grid-cols-2 gap-10">
        <!-- Left: Topics You Can Teach -->
        <div>
          <div class="mb-6">
            <div style="font-size: 18px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 6px">
              Topics You Can Teach
            </div>
            <div style="font-size: 12.5px; color: var(--ink-muted)">
              Based on mastery scores above 75%
            </div>
          </div>

          @if (tutoringService.listsLoading() && tutoringService.teachingTopics().length === 0) {
            <div style="font-size: 13px; color: var(--ink-muted)">Loading your teaching profile…</div>
          }

          <div
            class="teaching-profile"
            style="padding: 20px 22px; border-radius: var(--r-xl); background: var(--emerald-light); border: 1px solid var(--emerald-border); box-shadow: var(--shadow-sm); margin-bottom: 20px"
          >
            <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 0.8px; color: var(--emerald); font-weight: 700; margin-bottom: 8px">
              Your Teaching Profile
            </div>
            <div style="font-size: 17px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 4px">
              @if (teachingCount() === 0) {
                No topics at 75%+ yet
              } @else {
                You qualify to tutor in {{ teachingCount() }} {{ teachingCount() === 1 ? 'topic' : 'topics' }}
              }
            </div>
            <div style="font-size: 12.5px; color: var(--ink-2)">
              Students can request help from you when your mastery for a topic is strong enough
            </div>
          </div>

          <div class="flex flex-col gap-2.5">
            @for (topic of tutoringService.teachingTopics(); track topic.topic_id) {
              <div
                class="teaching-topic-card flex items-center gap-3 cursor-pointer"
                style="padding: 14px 16px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: var(--card-bg); box-shadow: var(--shadow-xs); transition: all var(--transition-base)"
              >
                <div
                  class="flex items-center justify-center"
                  style="width: 30px; height: 30px; border-radius: var(--r-md); background: var(--emerald-light); border: 1px solid var(--emerald-border); flex-shrink: 0"
                >
                  <lucide-icon name="check-circle" [size]="16" [strokeWidth]="2" style="color: var(--emerald)" />
                </div>
                <div class="flex-1" style="font-size: 13.5px; font-weight: 600; color: var(--ink)">{{ topic.topic_name }}</div>
                <div style="width: 80px; height: 5px; background: var(--surface-sub); border-radius: 5px; overflow: hidden; border: 1px solid var(--divider)">
                  <div
                    class="progress-fill"
                    [style.width.%]="masteryPercent(topic.mastery)"
                    style="height: 100%; border-radius: 4px; background: var(--emerald)"
                  ></div>
                </div>
                <div style="font-size: 13.5px; font-family: var(--font-display); width: 38px; text-align: right; font-weight: 700; color: var(--emerald)">
                  {{ masteryPercent(topic.mastery) }}%
                </div>
              </div>
            } @empty {
              @if (!tutoringService.listsLoading()) {
                <div style="font-size: 13px; color: var(--ink-muted); line-height: 1.5">
                  Complete reviews to build mastery—topics at 75% or higher appear here automatically.
                </div>
              }
            }
          </div>

          <!-- Requests: tabs -->
          <div class="flex flex-col gap-4" style="margin-top: 22px">
            <div class="flex flex-wrap items-center gap-2" style="border-bottom: 1px solid var(--divider); padding-bottom: 10px">
              @for (tab of requestTabs; track tab.id) {
                <button
                  type="button"
                  (click)="requestsTab.set(tab.id)"
                  style="padding: 8px 14px; border-radius: var(--r-lg); font-size: 13px; font-weight: 600; cursor: pointer; border: 1px solid var(--divider); font-family: var(--font-display)"
                  [style.background]="requestsTab() === tab.id ? 'var(--navy)' : 'transparent'"
                  [style.color]="requestsTab() === tab.id ? '#fff' : 'var(--ink-muted)'"
                >
                  {{ tab.label }}
                  @if (tab.id === 'incoming' && tutoringService.incomingRequests().length > 0) {
                    <span style="margin-left: 6px; font-size: 10px; font-family: var(--mono); opacity: 0.9"
                      >({{ tutoringService.incomingRequests().length }})</span
                    >
                  }
                  @if (tab.id === 'sent' && outgoingPendingCount() > 0) {
                    <span style="margin-left: 6px; font-size: 10px; font-family: var(--mono); opacity: 0.9">({{ outgoingPendingCount() }})</span>
                  }
                </button>
              }
            </div>

            @if (tutoringService.listsLoading()) {
              <div style="font-size: 13px; color: var(--ink-muted)">Loading requests…</div>
            }

            @switch (requestsTab()) {
              @case ('incoming') {
                <div class="flex flex-col gap-2.5">
                  @for (request of tutoringService.incomingRequests(); track request.id) {
                    <div
                      class="request-card flex items-center gap-3.5"
                      style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 16px 20px; background: var(--card-bg); box-shadow: var(--shadow-sm); transition: all var(--transition-base)"
                    >
                      <div
                        class="flex items-center justify-center"
                        [style.background]="avatarColor($index)"
                        style="width: 34px; height: 34px; border-radius: 50%; font-size: 12px; font-weight: 700; color: #fff; flex-shrink: 0"
                      >
                        {{ initials(request.requester_name) }}
                      </div>
                      <div class="flex-1 min-w-0">
                        <div style="font-size: 13.5px; font-weight: 600; color: var(--ink)">{{ request.requester_name }}</div>
                        <div style="font-size: 12px; color: var(--ink-muted); margin-top: 2px">{{ request.topic_name }}</div>
                      </div>
                      <div style="font-size: 10.5px; color: var(--ink-faint); font-family: var(--mono); flex-shrink: 0">
                        {{ formatRequestTime(request.created_at) }}
                      </div>
                      <div class="flex gap-2 flex-shrink-0">
                        <button
                          type="button"
                          class="accept-btn"
                          style="padding: 7px 14px; border-radius: var(--r-xl); background: var(--navy); color: #fff; border: none; font-size: 12px; font-weight: 700; cursor: pointer; box-shadow: var(--shadow-xs); font-family: var(--font-display)"
                          [disabled]="tutoringService.loading()"
                          (click)="onAccept(request.id)"
                        >
                          Accept
                        </button>
                        <button
                          type="button"
                          class="decline-btn"
                          style="padding: 7px 14px; border-radius: var(--r-xl); background: transparent; border: 1px solid var(--divider); color: var(--ink-muted); font-size: 12px; font-weight: 600; cursor: pointer"
                          [disabled]="tutoringService.loading()"
                          (click)="onDecline(request.id)"
                        >
                          Decline
                        </button>
                      </div>
                    </div>
                  } @empty {
                    @if (!tutoringService.listsLoading()) {
                      <div style="font-size: 13px; color: var(--ink-muted)">No pending requests—students will show up here when they ask you for help.</div>
                    }
                  }
                </div>
              }
              @case ('sent') {
                <div class="flex flex-col gap-2.5">
                  @for (request of outgoingPendingList(); track request.id) {
                    <div
                      class="request-card flex items-center gap-3.5"
                      style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 16px 20px; background: var(--card-bg); box-shadow: var(--shadow-sm)"
                    >
                      <div
                        class="flex items-center justify-center"
                        [style.background]="avatarColor($index)"
                        style="width: 34px; height: 34px; border-radius: 50%; font-size: 12px; font-weight: 700; color: #fff; flex-shrink: 0"
                      >
                        {{ initials(request.tutor_name) }}
                      </div>
                      <div class="flex-1 min-w-0">
                        <div style="font-size: 13.5px; font-weight: 600; color: var(--ink)">To: {{ request.tutor_name }}</div>
                        <div style="font-size: 12px; color: var(--ink-muted); margin-top: 2px">{{ request.topic_name }}</div>
                      </div>
                      <div
                        style="font-size: 10px; font-weight: 700; text-transform: uppercase; padding: 4px 8px; border-radius: 10px; background: var(--amber-light); color: var(--amber); border: 1px solid var(--divider)"
                      >
                        Waiting
                      </div>
                      <div style="font-size: 10.5px; color: var(--ink-faint); font-family: var(--mono)">{{ formatRequestTime(request.created_at) }}</div>
                    </div>
                  } @empty {
                    @if (!tutoringService.listsLoading()) {
                      <div style="font-size: 13px; color: var(--ink-muted)">You have no outgoing requests waiting on a tutor.</div>
                    }
                  }
                </div>
              }
              @case ('accepted') {
                <div class="flex flex-col gap-8">
                  <div>
                    <div style="font-size: 14px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 12px">
                      You're tutoring
                    </div>
                    <div class="flex flex-col gap-2.5">
                      @for (request of tutoringService.incomingAccepted(); track request.id) {
                        <div
                          style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 16px 20px; background: var(--card-bg); box-shadow: var(--shadow-sm)"
                        >
                          <div style="font-size: 13.5px; font-weight: 600; color: var(--ink)">{{ request.requester_name }}</div>
                          <div style="font-size: 12px; color: var(--ink-muted); margin-top: 4px">{{ request.topic_name }}</div>
                          @if (request.message) {
                            <div style="font-size: 12px; color: var(--ink-2); margin-top: 8px; line-height: 1.5">{{ request.message }}</div>
                          }
                        </div>
                      } @empty {
                        @if (!tutoringService.listsLoading()) {
                          <div style="font-size: 13px; color: var(--ink-muted)">No accepted sessions where you're the tutor yet.</div>
                        }
                      }
                    </div>
                  </div>
                  <div>
                    <div style="font-size: 14px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 12px">
                      You're learning (they accepted)
                    </div>
                    <div class="flex flex-col gap-2.5">
                      @for (request of tutoringService.outgoingAccepted(); track request.id) {
                        <div
                          style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 16px 20px; background: var(--card-bg); box-shadow: var(--shadow-sm)"
                        >
                          <div style="font-size: 13.5px; font-weight: 600; color: var(--ink)">With {{ request.tutor_name }}</div>
                          <div style="font-size: 12px; color: var(--ink-muted); margin-top: 4px">{{ request.topic_name }}</div>
                          @if (request.message) {
                            <div style="font-size: 12px; color: var(--ink-2); margin-top: 8px; line-height: 1.5">{{ request.message }}</div>
                          }
                        </div>
                      } @empty {
                        @if (!tutoringService.listsLoading()) {
                          <div style="font-size: 13px; color: var(--ink-muted)">No tutors have accepted your requests yet.</div>
                        }
                      }
                    </div>
                  </div>
                </div>
              }
            }
          </div>
        </div>

        <!-- Right: Suggested Tutors for You -->
        <div>
          <div class="mb-6">
            <div style="font-size: 18px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 6px">
              Suggested Tutors for You
            </div>
            <div style="font-size: 12.5px; color: var(--ink-muted)">
              Matched to your weak topic: <strong style="color: var(--red)">{{ matchTopicLabel() }}</strong>
            </div>
          </div>

          @if (tutoringService.matchesLoading()) {
            <div style="padding: 24px; border: 1px dashed var(--divider); border-radius: var(--r-xl); text-align: center; color: var(--ink-muted); font-size: 13px">
              Finding tutors with strong mastery on this topic…
            </div>
          } @else if (tutoringService.tutorMatches().length === 0) {
            <div
              style="padding: 24px; border: 1px dashed var(--divider); border-radius: var(--r-xl); text-align: center; color: var(--ink-muted); font-size: 13px; line-height: 1.6"
            >
              No peers in your school list this topic with 75%+ mastery yet, or your weak topic could not be resolved. Keep reviewing—matches improve as the class graph fills in.
            </div>
          } @else {
            <div class="flex flex-col gap-4">
              @for (tutor of tutoringService.tutorMatches(); track tutor.user_id; let ti = $index) {
                <div
                  class="tutor-card"
                  style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 24px 26px; background: var(--card-bg); box-shadow: var(--shadow-sm); transition: all var(--transition-base)"
                >
                  <div class="flex items-center gap-4 mb-4">
                    <div
                      class="flex items-center justify-center"
                      [style.background]="avatarColor(ti)"
                      style="width: 46px; height: 46px; border-radius: 50%; font-size: 16px; font-weight: 700; color: #fff; flex-shrink: 0"
                    >
                      {{ initials(tutor.name) }}
                    </div>
                    <div class="flex-1">
                      <div class="flex items-center gap-2 mb-1">
                        <span style="font-size: 15px; font-weight: 700; color: var(--ink); font-family: var(--font-display)">{{ tutor.name }}</span>
                        <div
                          class="reputation-badge flex items-center gap-1"
                          style="font-size: 10.5px; font-weight: 700; font-family: var(--mono); color: var(--emerald); background: var(--emerald-light); border: 1px solid var(--emerald-border); padding: 2px 7px; border-radius: 12px"
                        >
                          <lucide-icon name="arrow-up" [size]="9" [strokeWidth]="3" /> {{ tutor.reputation }}
                        </div>
                      </div>
                      <div style="font-size: 12px; color: var(--ink-muted)">{{ tutor.topic_name }}</div>
                    </div>
                  </div>

                  <div class="mb-4">
                    <div style="font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.7px; color: var(--ink-faint); margin-bottom: 10px">
                      Mastery Pulse
                    </div>
                    <div style="padding: 12px 14px; background: var(--surface-sub); border-radius: var(--r-lg); border: 1px solid var(--divider)">
                      <div style="width: 100%; height: 10px; background: var(--surface-sub); border-radius: 5px; overflow: hidden; border: 1px solid var(--divider)">
                        <div
                          class="progress-fill mastery-hbar"
                          [style.width.%]="masteryPct(tutor)"
                          [style.background]="getBarColor(masteryPct(tutor))"
                          style="height: 100%; border-radius: 4px; transform-origin: left"
                        ></div>
                      </div>
                    </div>
                  </div>

                  <div class="flex gap-0 pt-4 mb-4" style="border-top: 1px solid var(--divider)">
                    <div class="flex-1 text-center">
                      <div style="font-size: 18px; font-weight: 800; font-family: var(--font-display); color: var(--navy)">
                        {{ tutor.sessions }}
                      </div>
                      <div style="font-size: 9.5px; color: var(--ink-faint); text-transform: uppercase; letter-spacing: 0.8px; margin-top: 3px; font-weight: 700">
                        Sessions
                      </div>
                    </div>
                    <div class="flex-1 text-center">
                      <div style="font-size: 18px; font-weight: 800; font-family: var(--font-display); color: var(--amber)">
                        {{ tutor.reputation }}
                      </div>
                      <div style="font-size: 9.5px; color: var(--ink-faint); text-transform: uppercase; letter-spacing: 0.8px; margin-top: 3px; font-weight: 700">
                        Reputation
                      </div>
                    </div>
                  </div>

                  <button
                    type="button"
                    class="book-session-btn w-full"
                    style="padding: 10px 18px; border-radius: var(--r-xl); background: var(--navy); color: #fff; border: none; font-size: 13px; font-weight: 700; cursor: pointer; box-shadow: var(--shadow-sm); font-family: var(--font-display); letter-spacing: -0.2px"
                    [disabled]="tutoringService.loading()"
                    (click)="onBookSession(tutor)"
                  >
                    Book Session
                  </button>
                </div>
              }
            </div>
          }
        </div>
      </div>

    </div>
  `,
  styles: [`
    :host { display: flex; flex-direction: column; overflow: hidden; }

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

    @keyframes fillBar {
      from { transform: scaleX(0); }
      to { transform: scaleX(1); }
    }

    .pulse-dot {
      width: 10px;
      height: 10px;
      border-radius: 50%;
      background: var(--emerald);
      flex-shrink: 0;
      animation: pulse 2s ease-in-out infinite;
    }

    @keyframes pulse {
      0%, 100% { opacity: 1; transform: scale(1); }
      50% { opacity: 0.6; transform: scale(1.1); }
    }

    .teaching-profile {
      animation: fadeInUp 0.3s ease forwards;
    }

    .progress-fill {
      transform-origin: left;
      transform: scaleX(0);
      animation: fillBar 0.6s ease 0.2s forwards;
    }

    .mastery-hbar.progress-fill {
      animation: fillBar 0.6s ease 0.2s forwards;
    }

    .teaching-topic-card:hover {
      transform: translateX(3px);
      background: var(--hover-bg) !important;
    }

    .request-card:hover {
      transform: translateY(-2px);
      box-shadow: var(--shadow-md) !important;
      border-color: var(--emerald) !important;
    }

    .tutor-card:hover {
      transform: translateY(-2px);
      box-shadow: var(--shadow-md) !important;
      border-color: var(--emerald) !important;
    }

    .post-request-btn:hover {
      box-shadow: var(--shadow-md) !important;
    }

    .accept-btn,
    .decline-btn {
      transition: transform 0.15s ease;
    }

    .accept-btn:hover,
    .decline-btn:hover {
      transform: scale(1.05);
    }

    .accept-btn:active,
    .decline-btn:active {
      transform: scale(0.95);
    }

    .reputation-badge {
      transition: transform 0.15s ease;
    }

    .reputation-badge:hover {
      transform: scale(1.05);
    }

    .book-session-btn {
      transition: transform 0.15s ease, box-shadow 0.15s ease;
    }

    .book-session-btn:hover {
      transform: scale(1.03);
      box-shadow: var(--shadow-md) !important;
    }

    .book-session-btn:active {
      transform: scale(0.97);
    }

    .book-session-btn:disabled {
      opacity: 0.55;
      cursor: not-allowed;
      transform: none;
    }
  `],
})
export default class TutoringTabComponent implements OnInit {
  protected readonly tutoringService = inject(TutoringService);
  protected readonly authService = inject(AuthService);
  protected readonly learningService = inject(LearningService);

  protected readonly showRequestForm = signal(false);
  protected readonly requestTopic = signal('');
  protected readonly requestMessage = signal('');
  protected readonly postTutorId = signal('');
  protected readonly requestsTab = signal<RequestsTab>('incoming');

  protected readonly requestTabs: { id: RequestsTab; label: string }[] = [
    { id: 'incoming', label: 'Incoming' },
    { id: 'sent', label: 'Sent' },
    { id: 'accepted', label: 'Accepted' },
  ];

  protected readonly topWeakTopic = computed(
    () => this.learningService.confusionInsight()?.top_topic_name ?? 'your weak topics',
  );

  protected readonly narratorName = computed(() => {
    const n = this.authService.currentUser()?.name?.trim();
    if (!n) return 'You';
    const first = n.split(/\s+/)[0];
    return first || 'You';
  });

  protected readonly teachingCount = computed(() => this.tutoringService.teachingTopics().length);

  protected readonly matchTopicLabel = computed(() => {
    const w = this.learningService.confusionInsight()?.top_topic_name?.trim();
    if (w) return w;
    const first = this.tutoringService.teachingTopics()[0]?.topic_name;
    return first ?? 'your weak topics';
  });

  protected readonly outgoingPendingCount = computed(() =>
    this.tutoringService.outgoingRequests().filter((r) => r.status === 'pending').length,
  );

  protected readonly outgoingPendingList = computed(() =>
    this.tutoringService.outgoingRequests().filter((r) => r.status === 'pending'),
  );

  ngOnInit(): void {
    void (async () => {
      await this.learningService.loadConfusionInsight();
      const weak = this.learningService.confusionInsight()?.top_topic_name?.trim();
      await this.tutoringService.refreshAll(weak || undefined);
    })();
  }

  protected masteryPercent(mastery: number): number {
    const m = mastery <= 1 ? mastery * 100 : mastery;
    return Math.round(Math.min(100, Math.max(0, m)));
  }

  protected masteryPct(tutor: { mastery: number }): number {
    const m = tutor.mastery;
    return m <= 1 ? m * 100 : Math.min(100, m);
  }

  protected initials(name: string): string {
    const parts = name.trim().split(/\s+/).filter(Boolean);
    if (parts.length === 0) return '?';
    if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
  }

  protected avatarColor(index: number): string {
    const colors = ['var(--navy)', 'var(--emerald)', 'var(--purple)'];
    return colors[index % colors.length];
  }

  protected formatRequestTime(iso: string): string {
    if (!iso) return '';
    const t = Date.parse(iso);
    if (Number.isNaN(t)) return '';
    const diffMs = Date.now() - t;
    const mins = Math.floor(diffMs / 60000);
    if (mins < 1) return 'just now';
    if (mins < 60) return `${mins}m ago`;
    const hrs = Math.floor(mins / 60);
    if (hrs < 24) return `${hrs}h ago`;
    const days = Math.floor(hrs / 24);
    return `${days}d ago`;
  }

  getBarColor(height: number): string {
    return height >= 90 ? 'var(--emerald)' : height >= 80 ? '#F59E0B' : 'var(--navy)';
  }

  protected onPostTutorChange(userId: string): void {
    this.postTutorId.set(userId);
    const m = this.tutoringService.tutorMatches().find((x) => x.user_id === userId);
    if (m) {
      this.requestTopic.set(m.topic_name);
    }
  }

  async onAccept(id: string): Promise<void> {
    await this.tutoringService.respondToRequest(id, 'accepted');
  }

  async onDecline(id: string): Promise<void> {
    await this.tutoringService.respondToRequest(id, 'declined');
  }

  async onBookSession(tutor: TutorMatch): Promise<void> {
    const first = tutor.name.trim().split(/\s+/)[0] || 'there';
    try {
      await this.tutoringService.createRequest({
        tutor_id: tutor.user_id,
        topic_name: tutor.topic_name,
        message: `Hi ${first}, I'd like to book a tutoring session on ${tutor.topic_name}.`,
      });
      this.requestsTab.set('sent');
      await this.tutoringService.loadOutgoingRequests('all');
    } catch {
      /* surfaced via tutoringService.error */
    }
  }

  async onSubmitPostRequest(): Promise<void> {
    try {
      await this.tutoringService.createRequest({
        tutor_id: this.postTutorId(),
        topic_name: this.requestTopic().trim(),
        message: this.requestMessage().trim(),
      });
      this.requestTopic.set('');
      this.requestMessage.set('');
      this.postTutorId.set('');
      this.showRequestForm.set(false);
      this.requestsTab.set('sent');
      await this.tutoringService.loadOutgoingRequests('all');
    } catch {
      /* error surfaced via tutoringService.error */
    }
  }
}
