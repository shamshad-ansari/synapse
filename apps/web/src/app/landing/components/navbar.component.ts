import { Component, signal, computed, inject, afterNextRender } from '@angular/core';
import { RouterLink } from '@angular/router';
import { LandingButtonComponent } from './landing-button.component';
import { AuthService } from '../../core/auth/auth.service';

@Component({
  selector: 'app-navbar',
  standalone: true,
  imports: [RouterLink, LandingButtonComponent],
  templateUrl: './navbar.component.html',
  styles: [`
    :host { display: block; }
    nav { transition: box-shadow 200ms ease, border-color 200ms ease; }
    nav.scrolled {
      backdrop-filter: blur(12px);
      box-shadow: 0 1px 12px rgba(0,0,0,0.06);
      border-bottom-color: rgba(0,0,0,0.08);
    }
  `],
})
export class NavbarComponent {
  private readonly authService = inject(AuthService);

  readonly logoImage = '/assets/synapse-logo.png';
  readonly isAuthenticated = computed(() => this.authService.isAuthenticated());
  readonly isScrolled = signal(false);

  constructor() {
    afterNextRender(() => {
      window.addEventListener('scroll', () => {
        this.isScrolled.set(window.scrollY > 80);
      }, { passive: true });
    });
  }
}
