import { Component } from '@angular/core';
import { NgStyle } from '@angular/common';
import { LucideAngularModule, ArrowRight } from 'lucide-angular';
import { ElegantShapeComponent } from './elegant-shape.component';

@Component({
  selector: 'app-hero',
  standalone: true,
  imports: [NgStyle, LucideAngularModule, ElegantShapeComponent],
  templateUrl: './hero.component.html',
})
export class HeroComponent {
  readonly ArrowRight = ArrowRight;
  readonly logoImage = 'assets/synapse-logo.png';

  readonly masteryFontStyle: Record<string, string> = {
    'font-family': 'Georgia, serif',
    'font-style': 'italic',
  };
}
