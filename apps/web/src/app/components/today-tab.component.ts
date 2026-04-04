import {
  Component,
  ElementRef,
  ChangeDetectionStrategy,
  OnInit,
  afterNextRender,
  computed,
  effect,
  inject,
  viewChild,
} from '@angular/core';
import { Router } from '@angular/router';
import { LucideAngularModule } from 'lucide-angular';
import { AuthService } from '../core/auth/auth.service';
import { TodayService, type Contract } from '../features/today/today.service';

interface PeerInfo {
  avatar: string;
  name: string;
  subject: string;
  rating: string;
  bgColor: string;
}


@Component({
  selector: 'app-today-tab',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [LucideAngularModule],
  template: `
    <div
      class="flex flex-col overflow-y-auto overflow-x-hidden"
      [style.padding]="todayService.loading() ? '0' : '44px 56px 56px'"
      style="gap: 36px"
    >
      @if (todayService.loading()) {
        <div style="padding: 44px 56px">
          <div class="skeleton" style="height: 32px; width: 280px; border-radius: var(--r-md); margin-bottom: 12px"></div>
          <div class="skeleton" style="height: 16px; width: 200px; border-radius: var(--r-md)"></div>
        </div>
      } @else {

      <!-- SVG Gradient Definitions - Synapse Navy to Emerald -->
      <svg style="position: absolute; width: 0; height: 0">
        <defs>
          <linearGradient id="rg-flourish" x1="0%" y1="0%" x2="100%" y2="0%">
            <stop offset="0%" stop-color="#102E67" />
            <stop offset="100%" stop-color="#00A344" />
          </linearGradient>
        </defs>
      </svg>

      <!-- Topbar -->
      <div class="flex items-center justify-between gap-4">
        <div>
          <div style="font-size: 26px; font-weight: 700; letter-spacing: -0.6px; font-family: var(--font-display); color: var(--ink)">
            Good {{ timeOfDay }}, {{ userFirstName() }}
          </div>
          <div style="font-size: 13.5px; color: var(--ink-muted); margin-top: 4px">
            {{ todayDateStr }} · Here's what matters right now
          </div>
        </div>
        <div class="flex items-center gap-2.5">
          <div
            class="search-box flex items-center gap-2 cursor-text"
            style="background: var(--surface-sub); border: 1px solid var(--divider); border-radius: var(--r-lg); padding: 8px 14px; color: var(--ink-faint); font-size: 13px; width: 210px"
          >
            <lucide-icon name="search" [size]="14" [strokeWidth]="2" style="color: var(--ink-faint)" />
            Search anything…
          </div>
          <div
            class="bell-btn relative flex items-center justify-center cursor-pointer"
            style="width: 34px; height: 34px; background: transparent; border: 1px solid var(--divider); border-radius: var(--r-lg); color: var(--ink-muted)"
          >
            <lucide-icon name="bell" [size]="16" [strokeWidth]="2" />
            <div
              class="notification-dot"
              style="position: absolute; top: 7px; right: 7px; width: 6px; height: 6px; background: var(--red); border-radius: 50%; border: 2px solid var(--bg)"
            ></div>
          </div>
        </div>
      </div>

      <!-- Contract Banner -->
      @if (todayService.contract(); as c) {
      <div
        class="contract-banner relative flex items-center gap-8 overflow-hidden"
        style="border: 1px solid #EAEAEA; border-radius: var(--r-xl); padding: 32px 36px; background: #FFFFFF; box-shadow: 0 4px 20px -2px rgba(0,0,0,0.05); transition: border-color var(--transition-base), box-shadow var(--transition-base)"
      >
        <div class="flex-1">
          <div style="font-size: 10px; font-weight: 700; letter-spacing: 1px; text-transform: uppercase; color: var(--navy); margin-bottom: 6px">
            Active Contract
          </div>
          <div style="font-size: 22px; font-weight: 700; letter-spacing: -0.5px; font-family: var(--font-display); color: var(--ink)">
            {{ c.course_name }}
          </div>
          <div class="flex items-center gap-4 mt-3" style="font-size: 13px; color: var(--ink-muted)">
            <div class="flex items-center gap-1.5">
              <lucide-icon name="calendar" [size]="14" [strokeWidth]="2" style="color: var(--navy)" /> Exam · {{ c.exam_date }} <span style="color: var(--ink-faint); margin-left: 4px">{{ c.days_until }} days</span>
            </div>
            <div class="flex items-center gap-1.5">
              <lucide-icon name="clock" [size]="14" [strokeWidth]="2" style="color: var(--navy)" /> {{ dailyBudgetMinutes(c) }} min budget today
            </div>
            <span
              class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full"
              [style.background]="contractStatusPillBg(c)"
              [style.color]="contractStatusPillFg(c)"
              [style.border]="'1px solid ' + contractStatusPillBorder(c)"
              style="font-size: 11.5px; font-weight: 600"
            >
              <div style="width: 5px; height: 5px; border-radius: 50%; background: currentColor"></div> {{ contractStatusLabel(c) }}
            </span>
          </div>
        </div>
        <div style="flex-shrink: 0">
          <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 0.7px; color: var(--ink-faint); margin-bottom: 8px; font-weight: 700">
            Week Budget
          </div>
          <div class="flex justify-between" style="font-size: 12.5px; color: var(--ink-muted); margin-bottom: 8px">
            <strong style="color: var(--ink); font-family: var(--font-display)">{{ c.hours_done }}h</strong> done <span style="color: var(--ink-faint); margin-left: 4px">/ {{ c.weekly_hours_budget }}h target</span>
          </div>
          <div style="height: 6px; background: var(--surface-sub); border-radius: 6px; width: 160px; overflow: hidden; border: 1px solid var(--divider)">
            <div
              style="height: 100%; background: var(--emerald); border-radius: 5px"
              [style.width.%]="weekProgressPercent(c)"
            ></div>
          </div>
        </div>
        <div class="relative" style="width: 76px; height: 76px; flex-shrink: 0">
          <svg viewBox="0 0 62 62" width="76" height="76" style="transform: rotate(-90deg)">
            <circle cx="31" cy="31" r="27" fill="none" stroke="var(--ink-ghost)" stroke-width="5" />
            <circle #ring cx="31" cy="31" r="27" fill="none" stroke="url(#rg-flourish)" stroke-width="5" stroke-linecap="round" style="transition: stroke-dashoffset 1.2s cubic-bezier(0.4,0,0.2,1) 0.4s" />
          </svg>
          <div
            class="absolute inset-0 flex items-center justify-center"
            style="font-size: 16px; font-weight: 800; font-family: var(--font-display); color: var(--ink)"
          >
            {{ c.readiness }}%
          </div>
        </div>
        <div style="flex-shrink: 0">
          <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 0.8px; color: var(--ink-faint); font-weight: 700; margin-bottom: 6px">
            Readiness
          </div>
          <div style="font-size: 12.5px; color: var(--ink-muted); line-height: 1.7">
            Focus area:<br />
            <strong style="color: var(--ink); font-family: var(--font-display)">{{ topWeakTopic() }}</strong>
          </div>
        </div>
      </div>
      }

      <!-- Main Content Layout -->
      <div class="flex gap-12 flex-1 min-h-0">

        <!-- Left Column -->
        <div class="flex-1 flex flex-col gap-8">

          <!-- Do This Now - Next Best Actions -->
          <div>
            <div class="flex items-baseline gap-2.5 mb-5">
              <div style="font-size: 15px; font-weight: 700; font-family: var(--font-display); color: var(--ink)">Do This Now</div>
              <div style="font-size: 11.5px; color: var(--ink-faint); font-family: var(--mono)">{{ todayService.actions().length }} actions</div>
              <div
                class="see-all-link ml-auto cursor-pointer transition-all"
                style="font-size: 13px; color: var(--navy); font-weight: 600; transition: var(--transition-base)"
              >
                See all →
              </div>
            </div>
            <div class="flex flex-col gap-3">
              @for (action of todayService.actions(); track $index) {
                <div
                  class="action-card relative flex items-center gap-4 cursor-pointer overflow-hidden"
                  style="padding: 18px 22px; border-radius: var(--r-xl); border: 1px solid #EAEAEA; background: #FFFFFF; box-shadow: 0 4px 20px -2px rgba(0,0,0,0.05); transition: all var(--transition-base)"
                  [style.animation-delay]="(0.2 + $index * 0.1) + 's'"
                >
                  <div
                    class="flex items-center justify-center"
                    style="width: 32px; height: 32px; border-radius: var(--r-md); flex-shrink: 0; color: var(--navy)"
                  >
                    <lucide-icon [name]="action.icon" [size]="16" [strokeWidth]="2" />
                  </div>
                  <div class="flex-1">
                    <div style="font-size: 14px; font-weight: 600; color: var(--ink); margin-bottom: 3px">{{ action.title }}</div>
                    <div style="font-size: 12.5px; color: var(--ink-muted)">{{ action.reason }}</div>
                  </div>
                  <div class="flex flex-col items-end gap-2">
                    <div style="font-size: 11px; font-family: var(--mono); color: var(--ink-faint)">{{ action.duration }}</div>
                    <button
                      class="action-btn"
                      style="font-size: 12px; font-weight: 700; padding: 8px 16px; border-radius: var(--r-lg); border: none; cursor: pointer; font-family: var(--font-display); background: var(--navy); color: #fff; transition: all var(--transition-fast)"
                      (click)="$event.stopPropagation(); navigateTo(action.route)"
                    >
                      {{ action.buttonText }}
                    </button>
                  </div>
                </div>
              }
            </div>
          </div>

          <!-- Weak Topics + Mastery Pulse Grid -->
          <div class="grid grid-cols-2 gap-8">

            <!-- Weak Topics with Mastery Pulse -->
            <div>
              <div class="flex items-baseline gap-2.5 mb-5">
                <div style="font-size: 15px; font-weight: 700; font-family: var(--font-display); color: var(--ink)">Weak Topics</div>
                <div style="font-size: 11.5px; color: var(--ink-faint); font-family: var(--mono)">{{ todayService.weakTopics().length }} flagged</div>
                <div class="ml-auto cursor-pointer" style="font-size: 13px; color: var(--navy); font-weight: 600">View all →</div>
              </div>
              @if (todayService.weakTopics().length === 0) {
                <div style="padding: 16px; border-radius: var(--r-lg); border: 1px dashed var(--divider); background: var(--surface-sub); color: var(--ink-muted); font-size: 12.5px; line-height: 1.6">
                  No weak topics flagged yet. Complete a review session to generate mastery signals.
                </div>
              } @else {
                @for (topic of todayService.weakTopics(); track topic.name) {
                  <div
                    class="weak-topic-card flex items-center gap-4 cursor-pointer"
                    style="padding: 14px 16px; border-radius: var(--r-lg); border: 1px solid #EAEAEA; background: #FFFFFF; box-shadow: 0 4px 20px -2px rgba(0,0,0,0.05); transition: all var(--transition-base); margin-bottom: 8px"
                    [style.animation-delay]="(0.5 + $index * 0.1) + 's'"
                  >
                    <div class="flex items-end gap-1 h-10" style="width: 60px; flex-shrink: 0">
                      @for (bar of topic.bars; track $index) {
                        <div
                          class="weak-bar"
                          style="flex: 1; border-radius: 2px 2px 0 0; min-height: 4px"
                          [style.height]="bar + '%'"
                          [style.background]="getBarColor(bar)"
                          [style.animation-delay]="(0.2 + $index * 0.05) + 's'"
                        ></div>
                      }
                    </div>
                    <div class="flex-1">
                      <div style="font-size: 13.5px; font-weight: 600; color: var(--ink)">{{ topic.name }}</div>
                      <div style="font-size: 11px; color: var(--ink-faint); margin-top: 1px; text-transform: uppercase; letter-spacing: 0.3px">
                        {{ getMasteryStatus(topic.mastery) }}
                      </div>
                    </div>
                    <div class="text-right">
                      <div
                        style="font-size: 14px; font-weight: 700; font-family: var(--font-display)"
                        [style.color]="getMasteryColor(topic.mastery)"
                      >{{ topic.mastery }}%</div>
                    </div>
                  </div>
                }
              }
            </div>

            <!-- Mastery Pulse - Full View -->
            <div>
              <div class="flex items-baseline gap-2.5 mb-5">
                <div style="font-size: 15px; font-weight: 700; font-family: var(--font-display); color: var(--ink)">Mastery Pulse</div>
                <div style="font-size: 11.5px; color: var(--ink-faint); font-family: var(--mono)">Last 7 topics</div>
              </div>
              @if (masteryBars().length === 0) {
                <div class="flex items-center justify-center" style="height: 128px; padding: 16px; background: var(--surface-sub); border-radius: var(--r-lg); border: 1px dashed var(--divider); color: var(--ink-muted); font-size: 12.5px">
                  Mastery pulse appears after topic-linked review activity.
                </div>
              } @else {
                <div class="flex items-end gap-2 h-32" style="padding: 16px; background: var(--surface-sub); border-radius: var(--r-lg); border: 1px solid var(--divider)">
                  @for (height of masteryBars(); track $index) {
                    <div
                      class="mastery-bar"
                      style="flex: 1; border-radius: 4px 4px 0 0; min-height: 8px"
                      [style.height]="height + '%'"
                      [style.background]="getBarColor(height)"
                      [style.animation-delay]="(0.6 + $index * 0.08) + 's'"
                    ></div>
                  }
                </div>
                <div class="flex justify-between mt-2" style="font-size: 10px; color: var(--ink-faint); font-family: var(--mono)">
                  @for (label of masteryLabels(); track $index) {
                    <span>{{ label }}</span>
                  }
                </div>
              }
            </div>
          </div>

          <!-- Tutors + Interventions Grid -->
          <div class="grid grid-cols-2 gap-8">

            <!-- Suggested Tutors -->
            <div>
              <div class="flex items-baseline gap-2.5 mb-5">
                <div style="font-size: 15px; font-weight: 700; font-family: var(--font-display); color: var(--ink)">Suggested Tutors</div>
                <div style="font-size: 11.5px; color: var(--ink-faint); font-family: var(--mono)">for {{ topWeakTopic() }}</div>
                <div class="ml-auto cursor-pointer" style="font-size: 13px; color: var(--navy); font-weight: 600" (click)="navigateTo('/tutoring')">Browse →</div>
              </div>
              @for (peer of peers; track peer.name) {
                <div
                  class="peer-item flex items-center gap-3 cursor-pointer"
                  style="padding: 11px 14px; border-radius: var(--r-lg); background: transparent; transition: background var(--transition-base)"
                >
                  <div
                    class="flex items-center justify-center"
                    style="width: 32px; height: 32px; border-radius: 50%; font-size: 11.5px; font-weight: 700; flex-shrink: 0; color: #fff"
                    [style.background]="peer.bgColor"
                  >
                    {{ peer.avatar }}
                  </div>
                  <div class="flex-1">
                    <div style="font-size: 13px; font-weight: 600; color: var(--ink)">{{ peer.name }}</div>
                    <div style="font-size: 11.5px; color: var(--ink-muted)">{{ peer.subject }}</div>
                  </div>
                  <div style="font-size: 11.5px; font-weight: 700; font-family: var(--mono); color: #F59E0B">★ {{ peer.rating }}</div>
                </div>
              }
            </div>

            <!-- Interventions -->
            <div>
              <div class="flex items-baseline gap-2.5 mb-5">
                <div style="font-size: 15px; font-weight: 700; font-family: var(--font-display); color: var(--ink)">Interventions</div>
              </div>
              <div style="padding: 16px 18px; border-radius: var(--r-lg); background: var(--navy-light); border: 1px solid var(--navy-border); margin-bottom: 12px">
                <div class="flex items-center gap-1.5" style="font-size: 13.5px; font-weight: 600; color: var(--ink); margin-bottom: 4px">
                  <lucide-icon name="book-open" [size]="15" [strokeWidth]="2" style="color: var(--navy)" /> Illusion of competence
                </div>
                <div style="font-size: 12.5px; color: var(--ink-muted); line-height: 1.6">
                  High accuracy on Induction but slow response — harder practice cards recommended to confirm real mastery.
                </div>
              </div>
              <div style="padding: 16px 18px; border-radius: var(--r-lg); background: var(--emerald-light); border: 1px solid var(--emerald-border)">
                <div class="flex items-center gap-1.5" style="font-size: 13.5px; font-weight: 600; color: var(--ink); margin-bottom: 4px">
                  <lucide-icon name="target" [size]="15" [strokeWidth]="2" style="color: var(--emerald)" /> 3-day streak!
                </div>
                <div style="font-size: 12.5px; color: var(--ink-muted); line-height: 1.6">
                  Ahead of 72% of your school cohort. Keep today's session alive.
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Sidebar -->
        <div class="flex flex-col gap-7" style="width: 264px; flex-shrink: 0">

          <!-- This Week Stat Cards -->
          <div>
            <div class="flex items-baseline gap-2.5 mb-5">
              <div style="font-size: 15px; font-weight: 700; font-family: var(--font-display); color: var(--ink)">This Week</div>
            </div>
            <div class="flex gap-3 mb-3">
              @for (stat of statsRow1(); track $index) {
                <div class="flex-1" style="padding: 16px; border-radius: var(--r-lg); border: 1px solid #EAEAEA; background: #FFFFFF; box-shadow: 0 4px 20px -2px rgba(0,0,0,0.05)">
                  <div [style.color]="stat.color" style="font-size: 20px; font-weight: 800; font-family: var(--font-display); margin-bottom: 4px">{{ stat.value }}</div>
                  <div style="font-size: 10.5px; color: var(--ink-faint); text-transform: uppercase; letter-spacing: 0.6px; font-weight: 600">{{ stat.label }}</div>
                </div>
              }
            </div>
            <div class="flex gap-3">
              @for (stat of statsRow2(); track $index) {
                <div class="flex-1" style="padding: 16px; border-radius: var(--r-lg); border: 1px solid #EAEAEA; background: #FFFFFF; box-shadow: 0 4px 20px -2px rgba(0,0,0,0.05)">
                  <div [style.color]="stat.color" style="font-size: 20px; font-weight: 800; font-family: var(--font-display); margin-bottom: 4px">{{ stat.value }}</div>
                  <div style="font-size: 10.5px; color: var(--ink-faint); text-transform: uppercase; letter-spacing: 0.6px; font-weight: 600">{{ stat.label }}</div>
                </div>
              }
            </div>
          </div>

          <!-- Streak -->
          <div>
            <div class="flex items-baseline gap-2.5 mb-5">
              <div style="font-size: 15px; font-weight: 700; font-family: var(--font-display); color: var(--ink)">Streak</div>
              <div style="font-size: 11.5px; color: var(--ink-faint); font-family: var(--mono)">Mon → Sun</div>
            </div>
            <div class="flex gap-2 justify-center">
              @for (status of streakDays(); track $index) {
                <div
                  style="width: 9px; height: 9px; border-radius: 50%"
                  [style.background]="getStreakBackground(status)"
                  [style.box-shadow]="getStreakBoxShadow(status)"
                ></div>
              }
            </div>
          </div>

          <!-- Deadline Alert -->
          @if (todayService.deadlineAlert(); as d) {
          <div style="padding: 20px 22px; border-radius: var(--r-xl); background: var(--red-light); border: 1px solid var(--red-border); box-shadow: var(--shadow-sm)">
            <div style="font-size: 10px; text-transform: uppercase; letter-spacing: 0.8px; color: var(--red); font-weight: 700; margin-bottom: 8px">
              Deadline Alert
            </div>
            <div style="font-size: 16px; font-weight: 700; font-family: var(--font-display); color: var(--ink); margin-bottom: 5px">
              {{ d.title }}
            </div>
            <div style="font-size: 12.5px; color: var(--ink-muted)">
              Due in <strong style="color: var(--red); font-family: var(--font-display)">{{ d.days }} days</strong> · {{ d.course }}
            </div>
            <button
              class="deadline-btn w-full flex items-center justify-center transition-all"
              style="margin-top: 14px; font-size: 13px; padding: 8px 16px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: #FFFFFF; color: var(--ink-2); font-weight: 600; cursor: pointer; transition: var(--transition-base)"
              (click)="navigateTo('/planner')"
            >
              View in Planner
            </button>
          </div>
          }
        </div>
      </div>
      }
    </div>
  `,
  styles: [`
    :host { display: flex; flex-direction: column; overflow: hidden; }

    /* Entrance animations */
    @keyframes fadeInUp {
      from { opacity: 0; transform: translateY(20px); }
      to { opacity: 1; transform: translateY(0); }
    }

    @keyframes fadeInLeft {
      from { opacity: 0; transform: translateX(-20px); }
      to { opacity: 1; transform: translateX(0); }
    }

    @keyframes fadeInLeftSmall {
      from { opacity: 0; transform: translateX(-10px); }
      to { opacity: 1; transform: translateX(0); }
    }

    /* Notification dot pulse */
    @keyframes notificationPulse {
      0%, 100% { transform: scale(1); }
      50% { transform: scale(1.2); }
    }

    /* Progress bar fill */
    @keyframes progressFill {
      from { width: 0; }
      to { width: 56%; }
    }

    /* Bar grow (mastery bars & weak topic bars) */
    @keyframes barGrow {
      from { transform: scaleY(0); }
      to { transform: scaleY(1); }
    }

    /* Contract banner */
    .contract-banner {
      animation: fadeInUp 0.4s ease-out 0.1s both;
    }
    .contract-banner:hover {
      border-color: var(--emerald) !important;
      box-shadow: 0 6px 24px -2px rgba(0,0,0,0.08) !important;
    }

    /* Search box */
    .search-box {
      transition: border-color 0.22s;
    }
    .search-box:hover {
      border-color: var(--navy) !important;
    }

    /* Bell / notification button */
    .bell-btn {
      transition: transform 0.15s, border-color 0.15s;
    }
    .bell-btn:hover {
      transform: scale(1.05);
      border-color: var(--navy) !important;
    }
    .bell-btn:active {
      transform: scale(0.95);
    }

    /* Notification dot */
    .notification-dot {
      animation: notificationPulse 2s infinite;
    }

    /* Progress bar fill */
    .progress-fill {
      animation: progressFill 1s ease-out 0.3s both;
    }

    /* Action cards */
    .action-card {
      animation: fadeInLeft 0.3s ease-out both;
    }
    .action-card:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 24px -2px rgba(0,0,0,0.08) !important;
    }
    .action-card:active {
      transform: scale(0.99);
    }

    /* Action card button */
    .action-btn {
      transition: transform 0.15s;
    }
    .action-btn:hover {
      transform: scale(1.05);
    }
    .action-btn:active {
      transform: scale(0.95);
    }

    /* "See all" link */
    .see-all-link:hover {
      text-decoration: underline !important;
      color: var(--emerald) !important;
    }

    /* Weak topic cards */
    .weak-topic-card {
      animation: fadeInLeftSmall 0.3s ease-out both;
    }
    .weak-topic-card:hover {
      transform: translateX(3px);
      background: var(--hover-bg) !important;
    }

    /* Weak topic mini bars */
    .weak-bar {
      animation: barGrow 0.6s ease-out both;
      transform-origin: bottom;
    }

    /* Mastery pulse bars */
    .mastery-bar {
      animation: barGrow 0.6s ease-out both;
      transform-origin: bottom;
    }

    /* Peer items */
    .peer-item {
      transition: transform 0.15s, background 0.15s;
    }
    .peer-item:hover {
      transform: translateX(3px);
      background: var(--hover-bg) !important;
    }

    /* Deadline button */
    .deadline-btn:hover {
      border-color: var(--navy) !important;
      color: var(--navy) !important;
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
export default class TodayTabComponent implements OnInit {
  private readonly router = inject(Router);
  protected readonly todayService = inject(TodayService);
  protected readonly authService = inject(AuthService);
  readonly ringRef = viewChild<ElementRef<SVGCircleElement>>('ring');

  protected readonly todayDateStr = new Date().toLocaleDateString('en-US', {
    weekday: 'long',
    month: 'short',
    day: 'numeric',
  });

  private readonly statsWithDue = computed(() =>
    this.todayService.stats().map((s) =>
      s.label === 'Due' ? { ...s, value: String(this.todayService.dueCardCount()) } : s,
    ),
  );

  readonly statsRow1 = computed(() => this.statsWithDue().slice(0, 2));
  readonly statsRow2 = computed(() => this.statsWithDue().slice(2, 4));

  readonly peers: PeerInfo[] = [
    { avatar: 'JL', name: 'Jamie Liu', subject: 'Recursion · 94% mastery', rating: '4.9', bgColor: 'var(--navy)' },
    { avatar: 'SK', name: 'Sam Kato', subject: 'Algorithms · 88% mastery', rating: '4.7', bgColor: 'var(--emerald)' },
    { avatar: 'MR', name: 'Maya Roth', subject: 'Logic · 91% mastery', rating: '4.8', bgColor: 'var(--purple)' },
  ];

  readonly streakDays = computed<(boolean | string)[]>(() => {
    const count = Math.max(0, Math.min(14, this.todayService.streak()));
    const out: (boolean | string)[] = [];
    for (let i = 6; i >= 1; i -= 1) {
      out.push(count >= i);
    }
    out.push('today');
    return out;
  });

  readonly masteryBars = computed(() => {
    const topics = this.todayService.weakTopics();
    return topics.map((t) => t.mastery).slice(0, 7);
  });

  readonly masteryLabels = computed(() => {
    return this.todayService.weakTopics().map((t) => t.name).slice(0, 7);
  });

  readonly topWeakTopic = computed(() => this.todayService.weakTopics()[0]?.name ?? 'No weak topics yet');

  constructor() {
    afterNextRender(() => {
      this.applyReadinessRing(this.todayService.contract()?.readiness ?? 0);
    });
    effect(() => {
      const pct = this.todayService.contract()?.readiness ?? 0;
      this.applyReadinessRing(pct);
    });
  }

  ngOnInit(): void {
    void this.todayService.loadToday();
  }

  protected get timeOfDay(): string {
    const h = new Date().getHours();
    return h < 12 ? 'morning' : h < 17 ? 'afternoon' : 'evening';
  }

  protected userFirstName(): string {
    const todayName = this.todayService.greetingName();
    if (todayName?.trim()) {
      return todayName.split(' ')[0];
    }
    const parts = this.authService.currentUser()?.name?.split(' ');
    return parts?.[0] ?? '';
  }

  dailyBudgetMinutes(c: Contract): number {
    const daily = (c.weekly_hours_budget || 0) * 60 / 7;
    return Math.max(0, Math.round(daily));
  }

  weekProgressPercent(c: Contract): number {
    if (!c.weekly_hours_budget) {
      return 0;
    }
    return (c.hours_done / c.weekly_hours_budget) * 100;
  }

  contractStatusLabel(c: Contract): string {
    switch (c.status) {
      case 'on_track':
        return 'On Track';
      case 'at_risk':
        return 'At Risk';
      case 'off_track':
        return 'Off Track';
      default:
        return 'On Track';
    }
  }

  contractStatusPillBg(c: Contract): string {
    switch (c.status) {
      case 'on_track':
        return 'var(--emerald-light)';
      case 'at_risk':
        return 'rgba(245, 158, 11, 0.15)';
      case 'off_track':
        return 'var(--red-light)';
      default:
        return 'var(--emerald-light)';
    }
  }

  contractStatusPillFg(c: Contract): string {
    switch (c.status) {
      case 'on_track':
        return 'var(--emerald)';
      case 'at_risk':
        return '#F59E0B';
      case 'off_track':
        return 'var(--red)';
      default:
        return 'var(--emerald)';
    }
  }

  contractStatusPillBorder(c: Contract): string {
    switch (c.status) {
      case 'on_track':
        return 'var(--emerald-border)';
      case 'at_risk':
        return 'rgba(245, 158, 11, 0.35)';
      case 'off_track':
        return 'var(--red-border)';
      default:
        return 'var(--emerald-border)';
    }
  }

  private applyReadinessRing(readinessPercent: number): void {
    const ring = this.ringRef()?.nativeElement;
    if (!ring) {
      return;
    }
    const circ = 2 * Math.PI * 27;
    ring.style.strokeDasharray = String(circ);
    ring.style.strokeDashoffset = String(circ);
    requestAnimationFrame(() => {
      setTimeout(() => {
        ring.style.strokeDashoffset = String(circ * (1 - readinessPercent / 100));
      }, 60);
    });
  }

  navigateTo(route: string): void {
    this.router.navigate([route]);
  }

  getBarColor(height: number): string {
    return height >= 70 ? 'var(--emerald)' : height >= 50 ? '#F59E0B' : 'var(--red)';
  }

  getMasteryColor(mastery: number): string {
    return mastery >= 70 ? 'var(--emerald)' : mastery >= 50 ? '#F59E0B' : 'var(--red)';
  }

  getMasteryStatus(mastery: number): string {
    return mastery < 40 ? 'UNSTABLE · deteriorating' : mastery < 60 ? 'FAMILIAR · needs review' : 'SOLID · stable';
  }

  getStreakBackground(status: boolean | string): string {
    return status === 'today' ? 'var(--navy)' : status ? 'var(--emerald)' : 'var(--ink-ghost)';
  }

  getStreakBoxShadow(status: boolean | string): string {
    return status === 'today' ? '0 0 0 3px rgba(16,46,103,0.15)' : 'none';
  }
}
