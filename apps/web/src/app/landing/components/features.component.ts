import { Component } from '@angular/core';
import {
  LucideAngularModule,
  Target,
  Users,
  BrainCircuit,
  Activity,
  Network,
  LineChart,
  Link,
  Zap,
} from 'lucide-angular';
import { LandingCardComponent } from './landing-card.component';

@Component({
  selector: 'app-features',
  standalone: true,
  imports: [LucideAngularModule, LandingCardComponent],
  templateUrl: './features.component.html',
})
export class FeaturesComponent {
  readonly Target = Target;
  readonly Users = Users;
  readonly BrainCircuit = BrainCircuit;
  readonly Activity = Activity;
  readonly Network = Network;
  readonly LineChart = LineChart;
  readonly Link = Link;
  readonly Zap = Zap;
}
