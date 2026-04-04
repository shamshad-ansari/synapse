import { Component, inject, model, signal, computed, HostListener, OnInit } from '@angular/core';
import { Router, RouterLink, NavigationEnd } from '@angular/router';
import { NgTemplateOutlet } from '@angular/common';
import { toSignal } from '@angular/core/rxjs-interop';
import { filter, map } from 'rxjs';
import { trigger, transition, style, animate } from '@angular/animations';
import { LucideAngularModule } from 'lucide-angular';
import { AuthService } from '../core/auth/auth.service';
import { LearningService } from '../features/learning/learning.service';

interface NavItem {
  icon: string;
  label: string;
  route: string;
  badge?: string;
}

interface NavSection {
  header: string;
  pushDown: boolean;
  items: NavItem[];
}

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [
    RouterLink,
    NgTemplateOutlet,
    LucideAngularModule,
  ],
  animations: [
    trigger('fadeSlide', [
      transition(':enter', [
        style({ opacity: 0, width: 0, overflow: 'hidden' }),
        animate('200ms ease-out', style({ opacity: 1, width: '*' })),
      ]),
      transition(':leave', [
        style({ overflow: 'hidden' }),
        animate('200ms ease-in', style({ opacity: 0, width: 0 })),
      ]),
    ]),
    trigger('sectionSlide', [
      transition(':enter', [
        style({ opacity: 0, height: 0, overflow: 'hidden' }),
        animate('200ms ease-out', style({ opacity: 1, height: '*' })),
      ]),
      transition(':leave', [
        style({ overflow: 'hidden' }),
        animate('200ms ease-in', style({ opacity: 0, height: 0 })),
      ]),
    ]),
  ],
  template: `
    <nav
      class="flex h-full min-h-0 flex-col overflow-y-auto relative"
      [style.width]="isOpen() ? 'var(--sidebar-width)' : '72px'"
      [style.flex-shrink]="0"
      [style.padding]="isOpen() ? '16px 10px 20px' : '16px 12px 20px'"
      [style.background]="'var(--sidebar-bg)'"
      [style.border-right]="'1px solid var(--divider)'"
      [style.transition]="'width 250ms cubic-bezier(0.4,0,0.2,1), padding 250ms cubic-bezier(0.4,0,0.2,1)'"
    >
      <!-- Toggle Button -->
      <div
        class="toggle-btn absolute cursor-pointer flex items-center justify-center"
        style="top: 20px; right: -12px; width: 24px; height: 24px; border-radius: 50%; z-index: 10"
        [style.background]="'var(--card-bg)'"
        [style.border]="'1px solid var(--divider)'"
        [style.color]="'var(--ink-muted)'"
        [style.box-shadow]="'var(--shadow-sm)'"
        (click)="toggle()"
      >
        <lucide-icon [name]="isOpen() ? 'chevrons-left' : 'chevrons-right'" [size]="14" [strokeWidth]="2" />
      </div>

      <!-- Logo -->
      <div class="flex items-center gap-2 mb-2" style="padding: 8px 10px 18px">
        <div
          class="logo-icon flex items-center justify-center"
          style="width: 28px; height: 28px; border-radius: 8px; background: #eef4e8; border: 1px solid #d8e6cc; flex-shrink: 0; cursor: pointer; overflow: hidden"
        >
          <img
            src="assets/synapse-logo.png"
            alt="Synapse logo"
            style="width: 100%; height: 100%; object-fit: cover; object-position: center; transform: scale(1.32)"
          />
        </div>
        @if (isOpen()) {
          <span
            @fadeSlide
            style="font-size: 15px; font-weight: 700; font-family: var(--font-display); color: var(--ink); letter-spacing: -0.3px; white-space: nowrap"
          >Synapse</span>
          <span
            @fadeSlide
            style="margin-left: auto; font-size: 9px; font-family: var(--mono); color: var(--ink-faint); background: transparent; border: 1px solid var(--divider); padding: 2px 6px; border-radius: 4px; white-space: nowrap"
          >β 0.4</span>
        }
      </div>

      <!-- Nav Item Template -->
      <ng-template #navItemTpl let-item let-pushDown="pushDown" let-isFirst="isFirst">
        <a
          [routerLink]="item.route"
          class="nav-item flex items-center gap-2.5 cursor-pointer select-none relative"
          [class.active]="isActive(item.route)"
          [style.padding]="isOpen() ? '7px 10px' : '7px 0'"
          [style.border-radius]="'var(--r-lg)'"
          [style.color]="isActive(item.route) ? 'var(--navy)' : 'var(--ink-muted)'"
          [style.font-size]="'13.5px'"
          [style.font-weight]="isActive(item.route) ? 500 : 400"
          [style.transition]="'all var(--transition-base)'"
          [style.justify-content]="isOpen() ? 'flex-start' : 'center'"
          [style.margin-top]="pushDown && !isOpen() && isFirst ? 'auto' : null"
        >
          <div class="hover-bg"></div>
          @if (isActive(item.route)) {
            <div class="active-bg"></div>
          }
          <span class="icon-wrap" [style.color]="isActive(item.route) ? 'var(--navy)' : 'var(--ink-muted)'">
            <lucide-icon [name]="item.icon" [size]="17" [strokeWidth]="2" />
          </span>
          @if (isOpen()) {
            <span @fadeSlide class="nav-label">{{ item.label }}</span>
          }
          @if (item.badge && isOpen()) {
            <span class="nav-badge">{{ item.badge }}</span>
          }
          @if (item.badge && !isOpen()) {
            <span class="notification-dot"></span>
          }
        </a>
      </ng-template>

      <!-- Sections -->
      @for (section of sections(); track section.header) {
        @if (isOpen()) {
          <div
            @sectionSlide
            class="section-header"
            [style.margin-top]="section.pushDown ? 'auto' : null"
          >{{ section.header }}</div>
        }
        @for (item of section.items; track item.route; let first = $first) {
          <ng-container
            [ngTemplateOutlet]="navItemTpl"
            [ngTemplateOutletContext]="{ $implicit: item, pushDown: section.pushDown, isFirst: first }"
          />
        }
      }

      <!-- User Info -->
      <div
        class="user-section flex items-center gap-2 cursor-pointer relative"
        [class.profile-active]="isActive('/profile')"
        style="margin-top: auto; padding: 16px 10px 0; border-top: 1px solid var(--divider)"
        [style.justify-content]="isOpen() ? 'flex-start' : 'center'"
        (click)="goToProfile()"
      >
        <div
          class="flex items-center justify-center"
          style="width: 30px; height: 30px; border-radius: 50%; font-size: 11px; font-weight: 600; flex-shrink: 0; color: #fff; background: var(--emerald)"
        >{{ userInitials() }}</div>
        @if (isOpen()) {
          <div @fadeSlide>
            <div style="font-size: 13px; font-weight: 500; color: var(--ink); white-space: nowrap">{{ userName() }}</div>
            <div style="font-size: 11px; color: var(--ink-muted); white-space: nowrap">{{ userSchool() }}</div>
          </div>
          <div
            @fadeSlide
            style="margin-left: auto"
            (click)="toggleUserMenu($event)"
          >
            <lucide-icon name="more-horizontal" [size]="16" color="var(--ink-faint)" />
          </div>
        }

        @if (showUserMenu()) {
          <div class="user-menu" (click)="$event.stopPropagation()">
            <button class="user-menu-item" (click)="onLogout()">
              Log out
            </button>
          </div>
        }
      </div>
    </nav>
  `,
  styles: [`
    :host {
      display: block;
      height: 100%;
      min-height: 0;
    }

    .toggle-btn {
      transition: background-color var(--transition-fast), transform var(--transition-fast);
    }
    .toggle-btn:hover {
      transform: scale(1.1);
      background-color: var(--hover-bg) !important;
    }
    .toggle-btn:active { transform: scale(0.95); }

    .logo-icon { transition: transform var(--transition-fast); }
    .logo-icon:hover { transform: scale(1.05) rotate(5deg); }
    .logo-icon:active { transform: scale(0.95); }

    .section-header {
      font-size: 9.5px;
      font-weight: 600;
      letter-spacing: 0.8px;
      text-transform: uppercase;
      color: var(--ink-faint);
      padding: 16px 10px 6px;
    }

    a.nav-item { text-decoration: none; }
    .nav-item:active { transform: scale(0.98); }

    .hover-bg {
      position: absolute; inset: 0;
      background: var(--hover-bg);
      border-radius: var(--r-lg);
      opacity: 0;
      transition: opacity var(--transition-fast);
      pointer-events: none;
    }
    .nav-item:not(.active):hover .hover-bg { opacity: 1; }
    .nav-item:not(.active):hover .nav-label { color: var(--ink); }

    .active-bg {
      position: absolute; inset: 0;
      background: var(--active-bg);
      border-radius: var(--r-lg);
      pointer-events: none;
    }

    .icon-wrap {
      position: relative; z-index: 1;
      display: flex; align-items: center;
    }
    .nav-label {
      position: relative; z-index: 1;
      white-space: nowrap;
    }
    .nav-badge {
      position: relative; z-index: 1;
      margin-left: auto;
      font-size: 10px;
      font-family: var(--mono);
      background: var(--red-light);
      color: var(--red);
      padding: 1px 7px;
      border-radius: 12px;
      font-weight: 600;
      border: 1px solid var(--red-border);
    }

    .notification-dot {
      position: absolute; top: 4px; right: 4px;
      width: 6px; height: 6px;
      background: var(--red);
      border-radius: 50%;
      border: 1.5px solid var(--sidebar-bg);
    }

    .user-section { transition: transform var(--transition-fast); }
    .user-section:hover { transform: scale(1.02); }
    .user-section:active { transform: scale(0.98); }
    .user-section.profile-active {
      background: var(--active-bg);
      border-radius: var(--r-lg);
      border-top-color: transparent !important;
    }

    .user-menu {
      position: absolute;
      bottom: 100%;
      right: 0;
      margin-bottom: 6px;
      min-width: 140px;
      background: #fff;
      border: 1px solid var(--divider);
      border-radius: var(--r-lg);
      box-shadow: var(--shadow-md);
      padding: 4px;
      z-index: 50;
    }
    .user-menu-item {
      display: block;
      width: 100%;
      text-align: left;
      padding: 8px 12px;
      font-size: 13px;
      font-family: var(--font);
      color: var(--ink);
      border: none;
      background: none;
      border-radius: var(--r-md);
      cursor: pointer;
      transition: background var(--transition-fast);
    }
    .user-menu-item:hover {
      background: var(--hover-bg, #f5f5f5);
    }
  `],
})
export class SidebarComponent implements OnInit {
  readonly isOpen = model(true);

  private readonly router = inject(Router);
  private readonly authService = inject(AuthService);
  private readonly learningService = inject(LearningService);

  protected readonly showUserMenu = signal(false);

  protected readonly dueCount = computed(() => this.learningService.dueCards().length);

  ngOnInit() {
    this.learningService.loadDueCards();
  }

  protected readonly userName = computed(() =>
    this.authService.currentUser()?.name ?? 'Loading...',
  );

  protected readonly userSchool = computed(() => {
    const user = this.authService.currentUser();
    return user ? 'University' : '';
  });

  protected readonly userInitials = computed(() => {
    const name = this.authService.currentUser()?.name;
    if (!name) return '??';
    const parts = name.trim().split(/\s+/);
    const first = parts[0]?.[0] ?? '';
    const last = parts.length > 1 ? parts[parts.length - 1][0] : '';
    return (first + last).toUpperCase();
  });

  protected readonly currentUrl = toSignal(
    this.router.events.pipe(
      filter((e): e is NavigationEnd => e instanceof NavigationEnd),
      map(e => e.urlAfterRedirects),
    ),
    { initialValue: this.router.url },
  );

  protected readonly sections = computed<NavSection[]>(() => [
    {
      header: 'Learning',
      pushDown: false,
      items: [
        { icon: 'home', label: 'Today', route: '/today' },
        { icon: 'file-text', label: 'Notes', route: '/notes' },
        {
          icon: 'zap',
          label: 'Review',
          route: '/review',
          badge: this.dueCount() > 0 ? String(this.dueCount()) : undefined
        },
        { icon: 'calendar', label: 'Planner', route: '/planner' },
      ],
    },
    {
      header: 'Community',
      pushDown: false,
      items: [
        { icon: 'message-square', label: 'Feed', route: '/feed', badge: '3' },
        { icon: 'users', label: 'Tutoring', route: '/tutoring' },
      ],
    },
  ]);

  protected isActive(route: string): boolean {
    return this.currentUrl()?.startsWith(route) ?? false;
  }

  toggle() {
    this.isOpen.update(v => !v);
  }

  goToProfile() {
    this.router.navigate(['/profile']);
  }

  toggleUserMenu(event: Event) {
    event.stopPropagation();
    this.showUserMenu.update(v => !v);
  }

  onLogout() {
    this.showUserMenu.set(false);
    this.authService.logout();
  }

  @HostListener('document:click')
  onDocumentClick() {
    if (this.showUserMenu()) {
      this.showUserMenu.set(false);
    }
  }
}
