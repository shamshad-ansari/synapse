import { Component } from '@angular/core';
import { NgStyle } from '@angular/common';
import { LucideAngularModule, ArrowRight } from 'lucide-angular';

@Component({
  selector: 'app-hero',
  standalone: true,
  imports: [NgStyle, LucideAngularModule],
  templateUrl: './hero.component.html',
})
export class HeroComponent {
  readonly ArrowRight = ArrowRight;
  readonly logoImage = 'assets/synapse-logo.png';
  readonly backgroundSvg = 'assets/background2.svg';

  readonly dottedBgStyle: Record<string, string> = {
    'background-image': 'radial-gradient(#e5e7eb 1.5px, transparent 1.5px)',
    'background-size': '32px 32px',
    'background-position': 'center center',
  };

  readonly glowBgStyle: Record<string, string> = {
    'background-image': `url(${this.backgroundSvg})`,
    'background-size': 'cover',
    'background-position': 'center center',
    'background-repeat': 'no-repeat',
  };

  readonly masteryFontStyle: Record<string, string> = {
    'font-family': 'Georgia, serif',
    'font-style': 'italic',
  };
}
