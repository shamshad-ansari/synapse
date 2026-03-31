import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-landing-card',
  standalone: true,
  templateUrl: './landing-card.component.html',
})
export class LandingCardComponent {
  @Input() extraClass = '';

  get classes(): string {
    return `bg-white border border-gray-200 rounded-[16px] shadow-[0_4px_20px_-2px_rgba(0,0,0,0.05)] ${this.extraClass}`.trim();
  }
}
