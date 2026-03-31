import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { LucideAngularModule } from 'lucide-angular';

interface SessionBlock {
  status: 'done' | 'planned' | 'missed';
  title: string;
  duration: string;
}

interface CalendarRow {
  time: string;
  cells: (SessionBlock | null)[];
}

interface Deadline {
  name: string;
  date: string;
  days: string;
  urgency: 'urgent' | 'soon' | 'safe';
}

@Component({
  selector: 'app-planner-tab',
  standalone: true,
  imports: [LucideAngularModule],
  template: `
    <div class="flex flex-col overflow-y-auto overflow-x-hidden" style="padding: 40px 52px; gap: 32px">
      <!-- Topbar -->
      <div class="flex items-center justify-between gap-4">
        <div>
          <div style="font-size: 22px; font-weight: 700; letter-spacing: -0.5px; color: var(--ink)">Study Planner</div>
          <div style="font-size: 13px; color: var(--ink-muted); margin-top: 3px">Your path to Mar 12 exam · 16 days remaining</div>
        </div>
        <div class="flex items-center gap-2">
          <button
            class="flex items-center gap-1.5 transition-all duration-110 missed-btn"
            style="font-size: 12.5px; padding: 7px 14px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: transparent; color: var(--ink-2); font-weight: 500; cursor: pointer"
          >
            <lucide-icon name="alert-triangle" [size]="14" /> I missed yesterday
          </button>
          <button
            class="flex items-center gap-1.5 transition-all duration-110 regen-btn"
            style="font-size: 12.5px; padding: 7px 14px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 500; cursor: pointer"
          >
            <lucide-icon name="refresh-cw" [size]="14" /> Regenerate Plan
          </button>
        </div>
      </div>

      <div>
        <!-- Week Navigation -->
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2.5">
            <div
              class="flex items-center justify-center cursor-pointer transition-all duration-110 nav-arrow"
              style="width: 28px; height: 28px; background: transparent; border: 1px solid var(--divider); border-radius: var(--r-md); color: var(--ink-muted)"
            >
              <lucide-icon name="chevron-left" [size]="14" />
            </div>
            <div style="font-size: 13.5px; font-weight: 600; color: var(--ink)">Feb 24 – Mar 2, 2025</div>
            <div
              class="flex items-center justify-center cursor-pointer transition-all duration-110 nav-arrow"
              style="width: 28px; height: 28px; background: transparent; border: 1px solid var(--divider); border-radius: var(--r-md); color: var(--ink-muted)"
            >
              <lucide-icon name="chevron-right" [size]="14" />
            </div>
          </div>
          <div class="flex gap-3.5" style="font-size: 12px; color: var(--ink-muted)">
            <div class="flex items-center gap-1.5">
              <div style="width: 8px; height: 8px; border-radius: 2px; background: var(--emerald); opacity: 0.7"></div>
              Done
            </div>
            <div class="flex items-center gap-1.5">
              <div style="width: 8px; height: 8px; border-radius: 2px; background: var(--navy); opacity: 0.7"></div>
              Planned
            </div>
            <div class="flex items-center gap-1.5">
              <div style="width: 8px; height: 8px; border-radius: 2px; background: var(--red); opacity: 0.7"></div>
              Missed
            </div>
          </div>
        </div>

        <!-- Calendar Grid -->
        <div
          class="overflow-hidden"
          style="display: grid; grid-template-columns: 52px repeat(7, 1fr); gap: 1px; background: var(--divider); border-radius: var(--r-xl); border: 1px solid var(--divider)"
        >
          <!-- Header row -->
          <div style="background: var(--surface-sub); padding: 10px 8px"></div>
          @for (day of weekDays; track day.name; let i = $index) {
            <div class="text-center" style="background: var(--surface-sub); padding: 10px 8px">
              <div style="font-size: 10px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.6px; color: var(--ink-faint)">{{ day.name }}</div>
              <div
                style="font-size: 16px; font-weight: 700; font-family: var(--mono); margin-top: 2px"
                [style.color]="i === 1 ? 'var(--navy)' : 'var(--ink)'"
              >
                {{ day.date }}
              </div>
            </div>
          }

          <!-- Time-slot rows -->
          @for (row of calendarRows; track row.time) {
            <div style="background: var(--surface-sub); padding: 8px">
              <div style="font-size: 9.5px; color: var(--ink-faint); font-family: var(--mono); text-align: right; padding-top: 4px">{{ row.time }}</div>
            </div>
            @for (cell of row.cells; track $index) {
              <div style="background: var(--bg); padding: 8px; min-height: 56px">
                @if (cell) {
                  <div class="session-block" style="display: flex; flex-direction: column; gap: 2px;">
                    <div
                      style="font-size: 11px; font-weight: 600"
                      [style.color]="cell.status === 'missed' ? 'var(--red)' : 'var(--ink)'"
                    >{{ cell.title }}</div>
                    <div style="font-size: 10px; color: var(--ink-muted)">{{ cell.duration }}</div>
                  </div>
                }
              </div>
            }
          }
        </div>
      </div>

      <!-- Upcoming Deadlines -->
      <div>
        <div class="flex items-baseline gap-2 mb-4">
          <div style="font-size: 14px; font-weight: 600; color: var(--ink)">Upcoming Deadlines</div>
          <div class="ml-auto cursor-pointer" style="font-size: 12.5px; color: var(--navy); font-weight: 500">+ Add deadline</div>
        </div>
        @for (d of deadlines; track d.name) {
          <div
            class="flex items-center gap-3 cursor-pointer transition-all duration-110 deadline-row"
            style="padding: 12px 16px; border-radius: var(--r-lg)"
          >
            <div
              style="width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0"
              [style.background]="getUrgencyColor(d.urgency)"
            ></div>
            <div class="flex-1" style="font-size: 13.5px; font-weight: 500; color: var(--ink)">{{ d.name }}</div>
            <div style="font-size: 12.5px; color: var(--ink-muted); font-family: var(--mono)">{{ d.date }}</div>
            <div
              style="font-size: 12.5px; font-weight: 600"
              [style.color]="getUrgencyColor(d.urgency)"
            >{{ d.days }}</div>
          </div>
        }
      </div>
    </div>
  `,
  styles: [`
    :host { display: flex; flex-direction: column; overflow: hidden; }

    .deadline-row:hover {
      background: var(--hover-bg);
    }
  `],
})
export default class PlannerTabComponent {
  constructor(private router: Router) {}

  weekDays = [
    { name: 'Mon', date: 24 },
    { name: 'Tue', date: 25 },
    { name: 'Wed', date: 26 },
    { name: 'Thu', date: 27 },
    { name: 'Fri', date: 28 },
    { name: 'Sat', date: 29 },
    { name: 'Sun', date: 30 },
  ];

  calendarRows: CalendarRow[] = [
    {
      time: '9am',
      cells: [
        { status: 'done', title: 'Recursion Review', duration: '45 min · done' },
        { status: 'planned', title: 'Induction Cards', duration: '30 min' },
        { status: 'planned', title: 'Set Theory', duration: '45 min' },
        null,
        { status: 'planned', title: 'Mixed Review', duration: '60 min' },
        null,
        { status: 'planned', title: 'Proof Practice', duration: '45 min' },
      ],
    },
    {
      time: '2pm',
      cells: [
        { status: 'done', title: 'Logic Review', duration: '30 min · done' },
        null,
        null,
        { status: 'missed', title: 'Induction Deep', duration: 'missed' },
        { status: 'planned', title: 'Graph Theory', duration: '45 min' },
        { status: 'planned', title: 'PS3 Review', duration: '90 min' },
        null,
      ],
    },
  ];

  deadlines: Deadline[] = [
    { name: 'Problem Set 3', date: 'Feb 27', days: '2 days', urgency: 'urgent' },
    { name: 'Midterm Exam · CS225', date: 'Mar 12', days: '16 days', urgency: 'soon' },
    { name: '18.06 Problem Set 2', date: 'Mar 5', days: '9 days', urgency: 'safe' },
  ];

  getUrgencyColor(urgency: string): string {
    return urgency === 'urgent' ? 'var(--red)' : urgency === 'soon' ? 'var(--amber)' : 'var(--emerald)';
  }
}
