import { Component, signal, computed, inject, OnInit, ChangeDetectionStrategy } from '@angular/core';
import { LucideAngularModule } from 'lucide-angular';
import { AuthService } from '../core/auth/auth.service';
import { ProfileService } from '../features/profile/profile.service';

@Component({
  selector: 'app-profile-tab',
  standalone: true,
  imports: [LucideAngularModule],
  changeDetection: ChangeDetectionStrategy.OnPush,
  styles: [`
    :host { display: flex; flex-direction: column; overflow: hidden; }
    .hoverable { background: transparent; }
    .hoverable:hover { background: var(--hover-bg); }
  `],
  template: `
    <div class="flex flex-col overflow-y-auto overflow-x-hidden" style="padding: 40px 52px; gap: 32px">

      <!-- Topbar -->
      <div class="flex items-center justify-between gap-4">
        <div>
          <div style="font-size: 22px; font-weight: 700; letter-spacing: -0.5px; color: var(--ink)">Profile</div>
          <div style="font-size: 13px; color: var(--ink-muted); margin-top: 3px">Your reputation, contributions &amp; mastery</div>
        </div>
        <div class="flex items-center gap-2">
          <button
            type="button"
            class="flex items-center gap-1.5 transition-all duration-110"
            style="font-size: 12.5px; padding: 7px 14px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: transparent; color: var(--ink-2); font-weight: 500; cursor: pointer"
          >
            <lucide-icon name="shield" [size]="14" /> Privacy
          </button>
          <button
            type="button"
            class="flex items-center gap-1.5 transition-all duration-110"
            style="font-size: 12.5px; padding: 7px 14px; border-radius: var(--r-lg); border: none; background: var(--navy); color: #fff; font-weight: 500; cursor: pointer"
          >
            <lucide-icon name="edit-2" [size]="14" /> Edit Profile
          </button>
        </div>
      </div>

      @if (profileService.error()) {
        <div style="font-size: 13px; color: var(--red); padding: 12px 16px; border-radius: var(--r-lg); border: 1px solid var(--divider); background: var(--surface-sub)">
          {{ profileService.error() }}
        </div>
      }

      <!-- Profile Hero -->
      <div class="flex items-start gap-6 pb-8" style="border-bottom: 1px solid var(--divider); margin-bottom: 4px">
        <div
          class="flex items-center justify-center"
          style="width: 72px; height: 72px; border-radius: 50%; background: var(--navy); font-size: 26px; font-weight: 700; color: #fff; flex-shrink: 0"
        >
          {{ userInitials() }}
        </div>
        <div class="flex-1">
          <div style="font-size: 24px; font-weight: 700; letter-spacing: -0.5px; margin-bottom: 4px; color: var(--ink)">{{ displayName() }}</div>
          <div style="font-size: 13.5px; color: var(--ink-muted); margin-bottom: 12px">{{ userSubline() }}</div>
          <div class="flex gap-1.5 flex-wrap">
            @for (pill of pills(); track pill.text) {
              <span
                style="display: inline-flex; align-items: center; gap: 4px; padding: 2px 8px; border-radius: 100px; font-size: 11.5px; font-weight: 500"
                [style.background]="pillColorMap[pill.color].bg"
                [style.color]="pillColorMap[pill.color].color"
                [style.border]="'1px solid ' + pillColorMap[pill.color].border"
              >{{ pill.text }}</span>
            }
          </div>
        </div>
        <div class="flex gap-8 text-right">
          <div>
            <div style="font-size: 26px; font-weight: 700; font-family: var(--mono); letter-spacing: -0.5px; color: var(--navy)">{{ heroReputation() }}</div>
            <div style="font-size: 10px; color: var(--ink-faint); text-transform: uppercase; letter-spacing: 0.6px; margin-top: 2px; font-weight: 600">
              Reputation
            </div>
          </div>
          <div>
            <div style="font-size: 26px; font-weight: 700; font-family: var(--mono); letter-spacing: -0.5px; color: var(--emerald)">{{ heroCardsShared() }}</div>
            <div style="font-size: 10px; color: var(--ink-faint); text-transform: uppercase; letter-spacing: 0.6px; margin-top: 2px; font-weight: 600">
              Cards shared
            </div>
          </div>
          <div>
            <div style="font-size: 26px; font-weight: 700; font-family: var(--mono); letter-spacing: -0.5px; color: var(--ink)">{{ heroSessions() }}</div>
            <div style="font-size: 10px; color: var(--ink-faint); text-transform: uppercase; letter-spacing: 0.6px; margin-top: 2px; font-weight: 600">
              Sessions done
            </div>
          </div>
        </div>
      </div>

      @if (profileService.loading()) {
        <div style="font-size: 13px; color: var(--ink-muted)">Loading your profile stats…</div>
      }

      <!-- Content Grid -->
      <div class="grid grid-cols-2 gap-10 items-start">

        <!-- Left Column -->
        <div class="flex flex-col gap-8">

          <!-- Reputation Breakdown -->
          <div>
            <div class="flex items-baseline gap-2 mb-4">
              <div style="font-size: 14px; font-weight: 600; color: var(--ink)">Reputation Breakdown</div>
              <div style="font-size: 12px; color: var(--ink-faint); font-family: var(--mono)">{{ breakdownTotalLabel() }}</div>
            </div>
            @if (!profileService.loading() && repItems().length === 0) {
              <div style="font-size: 13px; color: var(--ink-muted)">No reputation activity yet.</div>
            }
            <div class="flex flex-col gap-3">
              @for (item of repItems(); track item.label) {
                <div class="flex items-center gap-3">
                  <div style="font-size: 13px; color: var(--ink-muted); width: 170px; flex-shrink: 0">{{ item.label }}</div>
                  <div class="flex-1" style="height: 5px; background: #F3F4F6; border-radius: 5px; overflow: hidden">
                    <div style="height: 100%; border-radius: 5px" [style.width]="item.percent + '%'" [style.background]="item.color"></div>
                  </div>
                  <div style="font-size: 12.5px; font-family: var(--mono); color: var(--ink-muted); width: 36px; text-align: right; flex-shrink: 0; font-weight: 600">
                    {{ item.value }}
                  </div>
                </div>
              }
            </div>
          </div>

          <!-- Mastery Overview -->
          <div>
            <div class="flex items-baseline gap-2 mb-4">
              <div style="font-size: 14px; font-weight: 600; color: var(--ink)">Mastery Overview</div>
            </div>
            @if (!profileService.loading() && masteryRows().length === 0) {
              <div style="font-size: 13px; color: var(--ink-muted)">No topic mastery yet — add courses and review flashcards.</div>
            }
            @for (topic of masteryRows(); track topic.name) {
              <div
                class="hoverable flex items-center gap-3 transition-all duration-100"
                style="padding: 10px 14px; border-radius: var(--r-lg)"
              >
                <div class="flex-1" style="font-size: 13.5px; font-weight: 500; color: var(--ink)">{{ topic.name }}</div>
                <div style="width: 90px; height: 4px; background: #F3F4F6; border-radius: 4px; overflow: hidden">
                  <div style="height: 100%; border-radius: 4px" [style.width]="topic.mastery + '%'" [style.background]="masteryColorMap[topic.band]"></div>
                </div>
                <div
                  style="font-size: 13px; font-family: var(--mono); width: 36px; text-align: right; font-weight: 600"
                  [style.color]="masteryColorMap[topic.band]"
                >{{ topic.mastery }}%</div>
              </div>
            }
          </div>

        </div>

        <!-- Right Column -->
        <div class="flex flex-col gap-8">

          <!-- Contributions -->
          <div>
            <div class="flex items-baseline gap-2 mb-4">
              <div style="font-size: 14px; font-weight: 600; color: var(--ink)">Contributions</div>
            </div>
            @for (item of contributions(); track item.title) {
              <div
                class="hoverable flex items-center gap-3 transition-all duration-100"
                style="padding: 11px 14px; border-radius: var(--r-lg)"
              >
                <div
                  class="flex items-center justify-center"
                  style="width: 32px; height: 32px; border-radius: var(--r-md); flex-shrink: 0"
                  [style.background]="item.icon_bg"
                >
                  <lucide-icon [name]="item.icon" [size]="14" [style.color]="item.icon_color" />
                </div>
                <div class="flex-1">
                  <div style="font-size: 13.5px; font-weight: 500; color: var(--ink)">{{ item.title }}</div>
                  <div style="font-size: 12px; color: var(--ink-muted)">{{ item.subtitle }}</div>
                </div>
                <div style="font-size: 16px; font-weight: 700; font-family: var(--mono)" [style.color]="item.value_color">{{ item.value }}</div>
              </div>
            }
          </div>

          <!-- Privacy Controls -->
          <div>
            <div class="flex items-baseline gap-2 mb-4">
              <div style="font-size: 14px; font-weight: 600; color: var(--ink)">Privacy Controls</div>
            </div>
            <div style="font-size: 12px; color: var(--ink-faint); margin-bottom: 8px">Preferences are stored on this device until account settings sync is available.</div>
            <div class="flex flex-col gap-1">
              @for (row of privacyRows; track row.key) {
                <div
                  class="hoverable flex items-center justify-between transition-all duration-100"
                  style="padding: 13px 16px; border-radius: var(--r-lg)"
                >
                  <div>
                    <div style="font-size: 13.5px; font-weight: 500; color: var(--ink)">{{ row.title }}</div>
                    <div style="font-size: 12px; color: var(--ink-muted)">{{ row.subtitle }}</div>
                  </div>
                  <div
                    class="cursor-pointer transition-all duration-180"
                    style="width: 36px; height: 20px; border-radius: 10px; position: relative; flex-shrink: 0"
                    [style.background]="toggles()[row.key] ? 'var(--navy)' : 'var(--ink-ghost)'"
                    (click)="toggleSwitch(row.key)"
                  >
                    <div
                      class="transition-all duration-180"
                      style="width: 16px; height: 16px; border-radius: 50%; background: #fff; position: absolute; top: 2px; box-shadow: 0 1px 3px rgba(0,0,0,0.18)"
                      [style.left]="toggles()[row.key] ? 'calc(100% - 18px)' : '2px'"
                    ></div>
                  </div>
                </div>
              }
            </div>
          </div>

        </div>
      </div>

    </div>
  `,
})
export default class ProfileTabComponent implements OnInit {
  private readonly authService = inject(AuthService);
  readonly profileService = inject(ProfileService);

  protected readonly displayName = computed(() =>
    this.authService.currentUser()?.name ?? 'Student',
  );
  protected readonly userInitials = computed(() => {
    const name = this.authService.currentUser()?.name;
    if (!name) return '??';
    const parts = name.trim().split(/\s+/);
    return ((parts[0]?.[0] ?? '') + (parts.length > 1 ? parts[parts.length - 1][0] : '')).toUpperCase();
  });
  protected readonly userSubline = computed(() => {
    const u = this.authService.currentUser();
    if (!u) return '';
    const school = u.school_name?.trim() || 'Your school';
    return `${school} · ${u.email}`;
  });

  protected readonly pills = computed(() => {
    const school = this.authService.currentUser()?.school_name?.trim();
    const t = this.toggles();
    const list: { color: 'blue' | 'green'; text: string }[] = [];
    if (school) {
      list.push({ color: 'blue', text: `🏫 School: ${school}` });
    }
    list.push({ color: 'green', text: t.global ? '🌐 Global: Active' : '🌐 Global: Off' });
    list.push({ color: 'blue', text: t.mastery ? 'Mastery: Shared' : 'Mastery: Private' });
    return list;
  });

  protected readonly heroReputation = computed(() => {
    const s = this.profileService.summary();
    if (this.profileService.loading() && !s) return '…';
    return s != null ? String(s.reputation_total) : '0';
  });
  protected readonly heroCardsShared = computed(() => {
    const s = this.profileService.summary();
    if (this.profileService.loading() && !s) return '…';
    return s != null ? String(s.cards_shared) : '0';
  });
  protected readonly heroSessions = computed(() => {
    const s = this.profileService.summary();
    if (this.profileService.loading() && !s) return '…';
    return s != null ? String(s.sessions_done) : '0';
  });

  protected readonly breakdownTotalLabel = computed(() => {
    const s = this.profileService.summary();
    if (!s) return '— pts total';
    return `${s.reputation_total} pts total`;
  });

  protected readonly repItems = computed(() => this.profileService.summary()?.reputation_breakdown ?? []);

  protected readonly masteryRows = computed(() => this.profileService.summary()?.mastery ?? []);

  protected readonly contributions = computed(() => this.profileService.summary()?.contributions ?? []);

  toggles = signal<{ mastery: boolean; global: boolean; feed: boolean }>({
    mastery: true,
    global: true,
    feed: false,
  });

  readonly pillColorMap: Record<string, { bg: string; color: string; border: string }> = {
    blue: { bg: 'var(--navy-light)', color: 'var(--navy)', border: 'var(--navy-border)' },
    green: { bg: 'var(--emerald-light)', color: 'var(--emerald)', border: 'var(--emerald-border)' },
  };

  readonly masteryColorMap: Record<string, string> = {
    red: 'var(--red)',
    amber: 'var(--amber)',
    green: 'var(--emerald)',
  };

  readonly privacyRows: { key: 'mastery' | 'global' | 'feed'; title: string; subtitle: string }[] = [
    { key: 'mastery', title: 'Share mastery scores', subtitle: 'Let others see your topic mastery' },
    { key: 'global', title: 'Global opt-in', subtitle: 'Appear in cross-school tutoring' },
    { key: 'feed', title: 'Show in school feed', subtitle: 'Your posts visible to school peers' },
  ];

  ngOnInit(): void {
    void this.profileService.loadSummary();
  }

  toggleSwitch(key: 'mastery' | 'global' | 'feed') {
    this.toggles.update(prev => ({ ...prev, [key]: !prev[key] }));
  }
}
