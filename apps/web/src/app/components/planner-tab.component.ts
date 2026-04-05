import { Component, ChangeDetectionStrategy, inject, signal, computed, OnInit } from '@angular/core';
import { LucideAngularModule } from 'lucide-angular';
import { PlannerService, StudySession, UpcomingDeadline } from '../features/planner/planner.service';

interface SessionBlock {
  id: string;
  status: 'done' | 'planned' | 'missed';
  title: string;
  duration: string;
}

interface CalendarRow {
  time: string;
  cells: (SessionBlock | null)[];
}

@Component({
  selector: 'app-planner-tab',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [LucideAngularModule],
  template: `
    <div class="flex flex-col overflow-y-auto overflow-x-hidden" style="padding: 40px 52px; gap: 32px">
      <!-- Topbar -->
      <div class="flex items-center justify-between gap-4">
        <div>
          <div style="font-size: 22px; font-weight: 700; letter-spacing: -0.5px; color: var(--ink)">Study Planner</div>
          <div style="font-size: 13px; color: var(--ink-muted); margin-top: 3px">{{ subtitle() }}</div>
        </div>
        <div class="flex items-center gap-2">
          <button
            class="flex items-center gap-1.5 transition-all duration-110 missed-btn"
            style="font-size: 12.5px; padding: 7px 14px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: transparent; color: var(--ink-2); font-weight: 500; cursor: pointer"
            (click)="onMissedYesterday()"
          >
            <lucide-icon name="alert-triangle" [size]="14" /> I missed yesterday
          </button>
          <button
            class="flex items-center gap-1.5 transition-all duration-110 regen-btn"
            style="font-size: 12.5px; padding: 7px 14px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 500; cursor: pointer"
            (click)="onRegeneratePlan()"
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
              (click)="prevWeek()"
            >
              <lucide-icon name="chevron-left" [size]="14" />
            </div>
            <div style="font-size: 13.5px; font-weight: 600; color: var(--ink)">{{ weekLabel() }}</div>
            <div
              class="flex items-center justify-center cursor-pointer transition-all duration-110 nav-arrow"
              style="width: 28px; height: 28px; background: transparent; border: 1px solid var(--divider); border-radius: var(--r-md); color: var(--ink-muted)"
              (click)="nextWeek()"
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
          @for (day of weekDays(); track day.name; let i = $index) {
            <div class="text-center" style="background: var(--surface-sub); padding: 10px 8px">
              <div style="font-size: 10px; font-weight: 700; text-transform: uppercase; letter-spacing: 0.6px; color: var(--ink-faint)">{{ day.name }}</div>
              <div
                style="font-size: 16px; font-weight: 700; font-family: var(--mono); margin-top: 2px"
                [style.color]="day.isToday ? 'var(--navy)' : 'var(--ink)'"
              >
                {{ day.date }}
              </div>
            </div>
          }

          <!-- Time-slot rows -->
          @for (row of calendarRows(); track row.time) {
            <div style="background: var(--surface-sub); padding: 8px">
              <div style="font-size: 9.5px; color: var(--ink-faint); font-family: var(--mono); text-align: right; padding-top: 4px">{{ row.time }}</div>
            </div>
            @for (cell of row.cells; track $index) {
              <div style="background: var(--bg); padding: 8px; min-height: 56px">
                @if (cell) {
                  <div class="session-block" style="display: flex; flex-direction: column; gap: 2px; cursor: pointer" (click)="toggleSession(cell)">
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
        </div>
        @for (d of deadlineRows(); track d.id) {
          <div
            class="flex items-center gap-3 cursor-pointer transition-all duration-110 deadline-row"
            style="padding: 12px 16px; border-radius: var(--r-lg)"
          >
            <div
              style="width: 7px; height: 7px; border-radius: 50%; flex-shrink: 0"
              [style.background]="getUrgencyColor(d.urgency)"
            ></div>
            <div class="flex-1" style="font-size: 13.5px; font-weight: 500; color: var(--ink)">{{ d.name }}</div>
            <div style="font-size: 12.5px; color: var(--ink-muted); font-family: var(--mono)">{{ d.dateStr }}</div>
            <div
              style="font-size: 12.5px; font-weight: 600"
              [style.color]="getUrgencyColor(d.urgency)"
            >{{ d.daysStr }}</div>
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
export default class PlannerTabComponent implements OnInit {
  private readonly plannerService = inject(PlannerService);

  // ── Week offset from the current week (0 = this week, -1 = last week, etc.)
  readonly weekOffset = signal(0);

  // ── Computed: week days for the header row
  readonly weekDays = computed(() => {
    const offset = this.weekOffset();
    const today = new Date();
    const dayOfWeek = today.getDay(); // 0=Sun
    const mondayDiff = dayOfWeek === 0 ? -6 : 1 - dayOfWeek;
    const monday = new Date(today);
    monday.setDate(today.getDate() + mondayDiff + offset * 7);

    const dayNames = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
    return dayNames.map((name, i) => {
      const d = new Date(monday);
      d.setDate(monday.getDate() + i);
      const todayStr = today.toDateString();
      return { name, date: d.getDate(), isToday: d.toDateString() === todayStr, fullDate: d };
    });
  });

  // ── Computed: week label (e.g. "Mar 24 – Mar 30, 2025")
  readonly weekLabel = computed(() => {
    const days = this.weekDays();
    const mon = days[0].fullDate;
    const sun = days[6].fullDate;
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return `${months[mon.getMonth()]} ${mon.getDate()} – ${months[sun.getMonth()]} ${sun.getDate()}, ${sun.getFullYear()}`;
  });

  // ── Computed: date range for API calls
  readonly weekStart = computed(() => this.formatDate(this.weekDays()[0].fullDate));
  readonly weekEnd = computed(() => this.formatDate(this.weekDays()[6].fullDate));

  // ── Computed: calendar grid rows from sessions signal
  readonly calendarRows = computed<CalendarRow[]>(() => {
    const sessions = this.plannerService.sessions();
    const days = this.weekDays();

    // Group sessions by date+time bucket
    const dateMap = new Map<string, StudySession[]>();
    for (const s of sessions) {
      const dateKey = s.scheduled_date.substring(0, 10); // YYYY-MM-DD
      if (!dateMap.has(dateKey)) dateMap.set(dateKey, []);
      dateMap.get(dateKey)!.push(s);
    }

    const timeSlots = [
      { time: '9am', hourMatch: '09' },
      { time: '2pm', hourMatch: '14' },
    ];

    return timeSlots.map(slot => {
      const cells: (SessionBlock | null)[] = days.map(day => {
        const dateKey = this.formatDate(day.fullDate);
        const daySessions = dateMap.get(dateKey) ?? [];
        const match = daySessions.find(s => s.start_time.startsWith(slot.hourMatch));
        if (!match) return null;
        return {
          id: match.id,
          status: match.status as 'done' | 'planned' | 'missed',
          title: match.title,
          duration: match.status === 'done' ? `${match.duration_minutes} min · done`
                  : match.status === 'missed' ? 'missed'
                  : `${match.duration_minutes} min`,
        };
      });
      return { time: slot.time, cells };
    });
  });

  // ── Computed: deadline display rows from deadlines signal
  readonly deadlineRows = computed(() => {
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return this.plannerService.deadlines().map(d => {
      const dueDate = new Date(d.due_date);
      return {
        id: d.id,
        name: d.course_name ? `${d.name} · ${d.course_name}` : d.name,
        dateStr: `${months[dueDate.getMonth()]} ${dueDate.getDate()}`,
        daysStr: d.days_until <= 0 ? 'Today' : `${d.days_until} day${d.days_until === 1 ? '' : 's'}`,
        urgency: d.urgency,
      };
    });
  });

  // ── Computed: subtitle showing nearest exam/deadline
  readonly subtitle = computed(() => {
    const dls = this.plannerService.deadlines();
    if (dls.length === 0) return 'No upcoming deadlines';
    const nearest = dls[0];
    const d = new Date(nearest.due_date);
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return `Your path to ${months[d.getMonth()]} ${d.getDate()} ${nearest.name} · ${nearest.days_until} days remaining`;
  });

  ngOnInit() {
    this.loadData();
  }

  async loadData() {
    const start = this.weekStart();
    const end = this.weekEnd();
    await Promise.all([
      this.plannerService.loadSessions(start, end),
      this.plannerService.loadDeadlines(),
    ]);
  }

  prevWeek() {
    this.weekOffset.update(v => v - 1);
    this.loadData();
  }

  nextWeek() {
    this.weekOffset.update(v => v + 1);
    this.loadData();
  }

  async toggleSession(cell: SessionBlock) {
    const newStatus = cell.status === 'planned' ? 'done' : cell.status === 'done' ? 'planned' : 'planned';
    await this.plannerService.updateSessionStatus(cell.id, newStatus);
    await this.plannerService.loadSessions(this.weekStart(), this.weekEnd());
  }

  async onMissedYesterday() {
    await this.plannerService.markMissedYesterday();
    await this.plannerService.loadSessions(this.weekStart(), this.weekEnd());
  }

  async onRegeneratePlan() {
    await this.plannerService.regeneratePlan();
    await this.loadData();
  }

  getUrgencyColor(urgency: string): string {
    return urgency === 'urgent' ? 'var(--red)' : urgency === 'soon' ? 'var(--amber)' : 'var(--emerald)';
  }

  private formatDate(d: Date): string {
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
  }
}
