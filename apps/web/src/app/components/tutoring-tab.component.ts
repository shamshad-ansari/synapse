import { Component } from '@angular/core';
import { LucideAngularModule } from 'lucide-angular';

@Component({
  selector: 'app-tutoring-tab',
  standalone: true,
  imports: [LucideAngularModule],
  template: `
    <div class="flex flex-col overflow-y-auto overflow-x-hidden" style="padding: 44px 56px 56px; gap: 36px">
      <!-- Topbar -->
      <div class="flex items-center justify-between gap-4">
        <div>
          <div style="font-size: 26px; font-weight: 700; letter-spacing: -0.6px; font-family: var(--font-display); color: var(--ink)">Tutoring</div>
          <div style="font-size: 13.5px; color: var(--ink-muted); margin-top: 4px">Graph-matched peer learning powered by mastery signals</div>
        </div>
        <button
          class="post-request-btn flex items-center gap-2 transition-all"
          style="font-size: 13px; padding: 8px 16px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 600; cursor: pointer; box-shadow: var(--shadow-sm); transition: var(--transition-base)"
        >
          <lucide-icon name="plus" [size]="15" [strokeWidth]="2" /> Post Request
        </button>
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

          <div
            class="teaching-profile"
            style="padding: 20px 22px; border-radius: var(--r-xl); background: var(--emerald-light); border: 1px solid var(--emerald-border); box-shadow: var(--shadow-sm); margin-bottom: 20px"
          >
            <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 0.8px; color: var(--emerald); font-weight: 700; margin-bottom: 8px">
              Your Teaching Profile
            </div>
            <div style="font-size: 17px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 4px">
              You qualify to tutor in 4 topics
            </div>
            <div style="font-size: 12.5px; color: var(--ink-2)">
              Students can request help from you automatically
            </div>
          </div>

          <div class="flex flex-col gap-2.5">
            @for (topic of teachingTopics; track topic.name) {
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
                <div class="flex-1" style="font-size: 13.5px; font-weight: 600; color: var(--ink)">{{ topic.name }}</div>
                <div style="width: 80px; height: 5px; background: var(--surface-sub); border-radius: 5px; overflow: hidden; border: 1px solid var(--divider)">
                  <div
                    class="progress-fill"
                    [style.width.%]="topic.mastery"
                    style="height: 100%; border-radius: 4px; background: var(--emerald)"
                  ></div>
                </div>
                <div style="font-size: 13.5px; font-family: var(--font-display); width: 38px; text-align: right; font-weight: 700; color: var(--emerald)">
                  {{ topic.mastery }}%
                </div>
              </div>
            }
          </div>

          <!-- Incoming Requests -->
          <div class="mt-8">
            <div class="flex items-baseline gap-2.5 mb-5">
              <div style="font-size: 15px; font-weight: 700; font-family: var(--font-display); color: var(--ink)">Incoming Requests</div>
              <div
                style="font-size: 10px; font-family: var(--mono); background: var(--red-light); color: var(--red); padding: 2px 8px; border-radius: 12px; font-weight: 700; border: 1px solid var(--red-border)"
              >
                2 NEW
              </div>
            </div>
            <div class="flex flex-col gap-2.5">
              @for (request of requests; track request.name) {
                <div
                  class="request-card flex items-center gap-3.5"
                  style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 16px 20px; background: var(--card-bg); box-shadow: var(--shadow-sm); transition: all var(--transition-base)"
                >
                  <div
                    class="flex items-center justify-center"
                    [style.background]="request.bgColor"
                    style="width: 34px; height: 34px; border-radius: 50%; font-size: 12px; font-weight: 700; color: #fff; flex-shrink: 0"
                  >
                    {{ request.avatar }}
                  </div>
                  <div class="flex-1">
                    <div style="font-size: 13.5px; font-weight: 600; color: var(--ink)">{{ request.name }}</div>
                    <div style="font-size: 12px; color: var(--ink-muted); margin-top: 2px">{{ request.topic }}</div>
                  </div>
                  <div style="font-size: 10.5px; color: var(--ink-faint); font-family: var(--mono)">{{ request.time }}</div>
                  <div class="flex gap-2">
                    <button
                      class="accept-btn"
                      style="padding: 7px 14px; border-radius: var(--r-xl); background: var(--navy); color: #fff; border: none; font-size: 12px; font-weight: 700; cursor: pointer; box-shadow: var(--shadow-xs); font-family: var(--font-display)"
                    >
                      Accept
                    </button>
                    <button
                      class="decline-btn"
                      style="padding: 7px 14px; border-radius: var(--r-xl); background: transparent; border: 1px solid var(--divider); color: var(--ink-muted); font-size: 12px; font-weight: 600; cursor: pointer"
                    >
                      Decline
                    </button>
                  </div>
                </div>
              }
            </div>
          </div>
        </div>

        <!-- Right: Suggested Tutors for You -->
        <div>
          <div class="mb-6">
            <div style="font-size: 18px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 6px">
              Suggested Tutors for You
            </div>
            <div style="font-size: 12.5px; color: var(--ink-muted)">
              Matched to your weak topic: <strong style="color: var(--red)">Recursion</strong>
            </div>
          </div>

          <div class="flex flex-col gap-4">
            @for (tutor of tutors; track tutor.name) {
              <div
                class="tutor-card"
                style="border: 1px solid var(--divider); border-radius: var(--r-xl); padding: 24px 26px; background: var(--card-bg); box-shadow: var(--shadow-sm); transition: all var(--transition-base)"
              >
                <!-- Header -->
                <div class="flex items-center gap-4 mb-4">
                  <div
                    class="flex items-center justify-center"
                    [style.background]="tutor.bgColor"
                    style="width: 46px; height: 46px; border-radius: 50%; font-size: 16px; font-weight: 700; color: #fff; flex-shrink: 0"
                  >
                    {{ tutor.avatar }}
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
                    <div style="font-size: 12px; color: var(--ink-muted)">{{ tutor.school }}</div>
                  </div>
                </div>

                <!-- Mastery Pulse Bars -->
                <div class="mb-4">
                  <div style="font-size: 11px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.7px; color: var(--ink-faint); margin-bottom: 10px">
                    Mastery Pulse
                  </div>
                  <div class="flex items-end gap-2.5 h-20" style="padding: 12px 14px; background: var(--surface-sub); border-radius: var(--r-lg); border: 1px solid var(--divider)">
                    @for (bar of tutor.masteryBars; track $index) {
                      <div class="flex-1 flex flex-col items-center gap-2">
                        <div
                          class="mastery-bar"
                          [style.height.%]="bar"
                          [style.background]="getBarColor(bar)"
                          [style.animation-delay]="(0.2 + $index * 0.1) + 's'"
                          style="width: 100%; border-radius: 4px 4px 0 0; min-height: 10px"
                        ></div>
                        <div style="font-size: 10px; color: var(--ink-faint); text-align: center; line-height: 1.2; font-weight: 600">
                          {{ tutor.topicLabels[$index] }}
                        </div>
                      </div>
                    }
                  </div>
                </div>

                <!-- Stats -->
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
                      {{ tutor.rating }}
                    </div>
                    <div style="font-size: 9.5px; color: var(--ink-faint); text-transform: uppercase; letter-spacing: 0.8px; margin-top: 3px; font-weight: 700">
                      Rating
                    </div>
                  </div>
                </div>

                <!-- Book Session Button -->
                <button
                  class="book-session-btn w-full"
                  style="padding: 10px 18px; border-radius: var(--r-xl); background: var(--navy); color: #fff; border: none; font-size: 13px; font-weight: 700; cursor: pointer; box-shadow: var(--shadow-sm); font-family: var(--font-display); letter-spacing: -0.2px"
                >
                  Book Session
                </button>
              </div>
            }
          </div>
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

    @keyframes growBar {
      from { transform: scaleY(0); }
      to { transform: scaleY(1); }
    }

    .teaching-profile {
      animation: fadeInUp 0.3s ease forwards;
    }

    .progress-fill {
      transform-origin: left;
      transform: scaleX(0);
      animation: fillBar 0.6s ease 0.2s forwards;
    }

    .mastery-bar {
      transform-origin: bottom;
      transform: scaleY(0);
      animation: growBar 0.6s ease-out forwards;
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
  `],
})
export default class TutoringTabComponent {
  teachingTopics = [
    { name: 'Set Theory', mastery: 78 },
    { name: 'Logic & Proofs', mastery: 82 },
    { name: 'Graph Theory', mastery: 76 },
    { name: 'Combinatorics', mastery: 79 },
  ];

  requests = [
    { avatar: 'PL', name: 'Priya Lal', topic: 'Set Theory · Basic operations', time: '10 min ago', bgColor: 'var(--navy)' },
    { avatar: 'TC', name: 'Tom Chen', topic: 'Logic proofs · Contradiction', time: '2h ago', bgColor: 'var(--emerald)' },
  ];

  tutors = [
    {
      avatar: 'JL',
      name: 'Jamie Liu',
      school: 'MIT · Year 3 · CS',
      reputation: 156,
      sessions: 23,
      rating: 4.9,
      masteryBars: [94, 88, 91],
      topicLabels: ['Recursion', 'Algorithms', 'DP'],
      bgColor: 'var(--navy)',
    },
    {
      avatar: 'MR',
      name: 'Maya Roth',
      school: 'MIT · Year 2 · Math+CS',
      reputation: 243,
      sessions: 18,
      rating: 4.8,
      masteryBars: [91, 95, 89],
      topicLabels: ['Logic', 'Induction', 'Proofs'],
      bgColor: 'var(--purple)',
    },
    {
      avatar: 'SK',
      name: 'Sam Kato',
      school: 'MIT · Year 4 · CS',
      reputation: 189,
      sessions: 31,
      rating: 4.7,
      masteryBars: [88, 92, 85],
      topicLabels: ['Algorithms', 'Graphs', 'DP'],
      bgColor: 'var(--emerald)',
    },
  ];

  getBarColor(height: number): string {
    return height >= 90 ? 'var(--emerald)' : height >= 80 ? '#F59E0B' : 'var(--navy)';
  }
}
