import { Component, computed, inject, signal } from '@angular/core';
import { Router, RouterOutlet, NavigationEnd } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { filter, map } from 'rxjs';
import { trigger, transition, style, animate, query, group } from '@angular/animations';
import { SidebarComponent } from './components/sidebar.component';

export const routeAnimation = trigger('routeAnimation', [
  transition('* <=> *', [
    query(':enter', [
      style({ opacity: 0, transform: 'translateY(8px)' }),
      animate(
        '200ms cubic-bezier(0.4, 0, 0.2, 1)',
        style({ opacity: 1, transform: 'translateY(0)' }),
      ),
    ], { optional: true }),
  ]),
]);

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, SidebarComponent],
  animations: [routeAnimation],
  template: `
    @if (showDashboardShell()) {
      <div class="flex h-screen">
        <app-sidebar [(isOpen)]="sidebarOpen" />
        <main
          class="flex-1 flex flex-col overflow-hidden"
          [@routeAnimation]="currentUrl()"
        >
          <router-outlet />
        </main>
      </div>
    } @else {
      <main class="min-h-screen overflow-y-auto" [@routeAnimation]="currentUrl()">
        <router-outlet />
      </main>
    }
  `,
})
export class AppComponent {
  sidebarOpen = signal(true);

  private readonly router = inject(Router);

  protected readonly currentUrl = toSignal(
    this.router.events.pipe(
      filter((e): e is NavigationEnd => e instanceof NavigationEnd),
      map(e => e.urlAfterRedirects),
    ),
    { initialValue: this.router.url },
  );

  protected readonly showDashboardShell = computed(() => {
    const url = this.currentUrl();
    const fullscreenRoutes = ['/', '/login', '/register', '/canvas/connect', '/canvas/connected'];
    return !fullscreenRoutes.some(r => url === r || url.startsWith(r + '?'));
  });
}
