import { Injectable, inject, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { firstValueFrom } from 'rxjs';
import { environment } from '../../../environments/environment';
import { LearningService } from '../learning/learning.service';
import { TutoringService } from '../tutoring/tutoring.service';

/** API row shape (snake_case from JSON). */
interface NbaActionRow {
  icon: string;
  title: string;
  reason: string;
  duration: string;
  button_text: string;
  route: string;
}

export interface NbaAction {
  icon: string;
  title: string;
  reason: string;
  duration: string;
  buttonText: string;
  route: string;
}

export interface Contract {
  course_name: string;
  exam_date: string;
  days_until: number;
  status: 'on_track' | 'at_risk' | 'off_track' | 'no_data';
  weekly_hours_budget: number;
  hours_done: number;
  readiness: number;
}

export interface WeakTopic {
  name: string;
  mastery: number;
  bars: number[];
}

export interface DeadlineAlert {
  course: string;
  title: string;
  days: number;
  type: string;
}

export interface StatInfo {
  value: string;
  label: string;
  color: string;
}

export interface TodayData {
  greeting_name: string;
  actions: NbaActionRow[];
  contract: Contract | null;
  weak_topics: WeakTopic[];
  stats: StatInfo[];
  streak: number;
  deadline_alert: DeadlineAlert | null;
}

@Injectable({ providedIn: 'root' })
export class TodayService {
  private readonly http = inject(HttpClient);
  private readonly learningService = inject(LearningService);
  private readonly tutoringService = inject(TutoringService);
  private readonly apiUrl = environment.apiUrl;

  readonly actions = signal<NbaAction[]>([]);
  readonly contract = signal<Contract | null>(null);
  readonly weakTopics = signal<WeakTopic[]>([]);
  readonly stats = signal<StatInfo[]>([]);
  readonly streak = signal(0);
  readonly deadlineAlert = signal<DeadlineAlert | null>(null);
  readonly greetingName = signal('');
  readonly loading = signal(true);
  readonly error = signal<string | null>(null);

  readonly dueCardCount = computed(() => this.learningService.dueCards().length);

  async loadToday(): Promise<void> {
    this.loading.set(true);
    this.error.set(null);
    try {
      const res = await firstValueFrom(
        this.http.get<{ data: TodayData }>(`${this.apiUrl}/v1/autopilot/today`),
      );
      const d = res.data;
      this.greetingName.set(d.greeting_name);
      this.actions.set(
        d.actions.map((a) => ({
          icon: a.icon,
          title: a.title,
          reason: a.reason,
          duration: a.duration,
          buttonText: a.button_text,
          route: a.route,
        })),
      );
      this.contract.set(d.contract);
      this.weakTopics.set(d.weak_topics);
      this.stats.set(d.stats);
      this.streak.set(d.streak);
      const alert = d.deadline_alert;
      this.deadlineAlert.set(alert?.title?.trim() ? alert : null);
      const weakName = d.weak_topics?.[0]?.name?.trim() ?? '';
      await this.tutoringService.loadTutorMatches(weakName);
      try {
        await this.learningService.loadDueCards(undefined, 20);
      } catch {
        /* optional: today still usable without due-card count */
      }
    } catch (err: unknown) {
      const e = err as { error?: { error?: string }; message?: string };
      this.error.set(e?.error?.error ?? e?.message ?? 'Failed to load today');
    } finally {
      this.loading.set(false);
    }
  }
}
