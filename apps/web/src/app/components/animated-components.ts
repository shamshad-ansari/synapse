import { Component, Directive, input, computed } from '@angular/core';
import { cn } from './ui/utils';

/**
 * InteractiveCard — a wrapper that adds hover scale + shadow and tap scale via CSS.
 * Usage: <app-interactive-card (click)="doSomething()">...</app-interactive-card>
 */
@Component({
  selector: 'app-interactive-card',
  standalone: true,
  template: `<ng-content />`,
  host: {
    '[class]': 'computedClass()',
    '[class.interactive-card--shadow]': 'shadow()',
    '[class.interactive-card--clickable]': 'clickable()',
  },
  styles: [`
    :host {
      display: block;
      transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    }
    :host(.interactive-card--clickable) { cursor: pointer; }
    :host:hover { transform: var(--hover-scale, scale(1.01)); }
    :host(.interactive-card--shadow):hover { box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08); }
    :host:active { transform: var(--tap-scale, scale(0.99)); }
  `],
})
export class InteractiveCardComponent {
  readonly userClass = input('', { alias: 'class' });
  readonly shadow = input(true);
  readonly clickable = input(false);

  protected readonly computedClass = computed(() => this.userClass());
}

/**
 * AnimatedButton — a styled button with hover/tap animation.
 * Usage: <button appAnimatedButton variant="primary" size="md">Click</button>
 */
@Directive({
  selector: '[appAnimatedButton]',
  standalone: true,
  host: {
    '[class]': 'computedClass()',
    '[style.border-radius]': '"var(--r-lg)"',
    '[style.font-weight]': '"500"',
    '[style.cursor]': '"pointer"',
    '[style.font-family]': '"var(--font)"',
  },
})
export class AnimatedButtonDirective {
  readonly variant = input<'primary' | 'secondary' | 'ghost'>('primary');
  readonly size = input<'sm' | 'md' | 'lg'>('md');
  readonly userClass = input('', { alias: 'class' });

  private static readonly variantMap = {
    primary:   'animated-btn--primary',
    secondary: 'animated-btn--secondary',
    ghost:     'animated-btn--ghost',
  } as const;

  private static readonly sizeMap = {
    sm: 'animated-btn--sm',
    md: 'animated-btn--md',
    lg: 'animated-btn--lg',
  } as const;

  protected readonly computedClass = computed(() =>
    cn(
      'flex items-center gap-1.5 animated-btn',
      AnimatedButtonDirective.variantMap[this.variant()],
      AnimatedButtonDirective.sizeMap[this.size()],
      this.userClass(),
    ),
  );
}

/**
 * PulsingDot — a small dot with a pulsing ring animation.
 * Usage: <app-pulsing-dot color="var(--red)" [size]="5" />
 */
@Component({
  selector: 'app-pulsing-dot',
  standalone: true,
  template: `
    <div [style.width.px]="size()" [style.height.px]="size()" [style.border-radius]="'50%'" [style.background]="color()"></div>
    <div
      class="pulse-ring"
      [style.border-color]="color()"
    ></div>
  `,
  host: {
    '[style.position]': '"relative"',
    '[style.display]': '"inline-block"',
    '[style.width.px]': 'size()',
    '[style.height.px]': 'size()',
  },
  styles: [`
    .pulse-ring {
      position: absolute;
      inset: -2px;
      border-radius: 50%;
      border: 2px solid;
      pointer-events: none;
      animation: pulseRing 2s infinite;
    }
    @keyframes pulseRing {
      0%, 100% { transform: scale(1); opacity: 0.7; }
      50% { transform: scale(1.8); opacity: 0; }
    }
  `],
})
export class PulsingDotComponent {
  readonly color = input('var(--red)');
  readonly size = input(5);
}

/**
 * FadeIn — wraps content with a CSS entrance animation.
 * Usage: <app-fade-in direction="up" [delay]="0.1">...</app-fade-in>
 */
@Component({
  selector: 'app-fade-in',
  standalone: true,
  template: `<ng-content />`,
  host: {
    '[style.display]': '"block"',
    '[style.animation]': 'animationValue()',
  },
  styles: [`
    @keyframes fadeInUp    { from { opacity: 0; transform: translateY(20px); }  to { opacity: 1; transform: translateY(0); } }
    @keyframes fadeInDown  { from { opacity: 0; transform: translateY(-20px); } to { opacity: 1; transform: translateY(0); } }
    @keyframes fadeInLeft  { from { opacity: 0; transform: translateX(20px); }  to { opacity: 1; transform: translateX(0); } }
    @keyframes fadeInRight { from { opacity: 0; transform: translateX(-20px); } to { opacity: 1; transform: translateX(0); } }
    @keyframes fadeInNone  { from { opacity: 0; } to { opacity: 1; } }
  `],
})
export class FadeInComponent {
  readonly delay = input(0);
  readonly duration = input(0.3);
  readonly direction = input<'up' | 'down' | 'left' | 'right' | 'none'>('up');

  private static readonly nameMap = {
    up: 'fadeInUp', down: 'fadeInDown', left: 'fadeInLeft',
    right: 'fadeInRight', none: 'fadeInNone',
  } as const;

  protected readonly animationValue = computed(() => {
    const name = FadeInComponent.nameMap[this.direction()];
    return `${name} ${this.duration()}s cubic-bezier(0.4, 0, 0.2, 1) ${this.delay()}s both`;
  });
}

/**
 * Skeleton — a loading placeholder with a shimmer animation.
 * Usage: <app-skeleton width="100%" [height]="16" />
 */
@Component({
  selector: 'app-skeleton',
  standalone: true,
  template: ``,
  host: {
    '[class]': 'userClass()',
    '[style.width]': 'width()',
    '[style.height.px]': 'height()',
    '[style.border-radius]': 'borderRadius()',
    '[style.background]': '"var(--ink-ghost)"',
  },
  styles: [`
    :host {
      display: block;
      animation: shimmer 1.5s ease-in-out infinite;
    }
    @keyframes shimmer {
      0%, 100% { opacity: 0.5; }
      50% { opacity: 0.8; }
    }
  `],
})
export class SkeletonComponent {
  readonly width = input<string | number>('100%');
  readonly height = input<string | number>(16);
  readonly borderRadius = input('var(--r-md)');
  readonly userClass = input('', { alias: 'class' });
}
