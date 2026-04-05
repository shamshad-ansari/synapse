import { Component } from '@angular/core';
import { LandingCardComponent } from './landing-card.component';

@Component({
  selector: 'app-features',
  standalone: true,
  imports: [LandingCardComponent],
  templateUrl: './features.component.html',
})
export class FeaturesComponent {
}
