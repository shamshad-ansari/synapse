import { Component, computed, effect, inject, input, signal } from '@angular/core';
import { cn } from './utils';

@Component({
  selector: 'app-avatar',
  standalone: true,
  template: `<ng-content />`,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"avatar"',
  },
})
export class Avatar {
  readonly userClass = input<string>('', { alias: 'class' });

  /** Shared image load state for child AvatarImage / AvatarFallback. */
  readonly imageStatus = signal<'loading' | 'loaded' | 'error'>('loading');

  protected readonly computedClass = computed(() =>
    cn('relative flex size-10 shrink-0 overflow-hidden rounded-full', this.userClass()),
  );
}

@Component({
  selector: 'app-avatar-image',
  standalone: true,
  template: `
    @if (isLoaded()) {
      <img
        [src]="src()"
        [alt]="alt()"
        [class]="computedClass()"
        data-slot="avatar-image"
      />
    }
  `,
  host: { style: 'display: contents' },
})
export class AvatarImage {
  private readonly avatar = inject(Avatar);

  readonly src = input('');
  readonly alt = input('');
  readonly userClass = input<string>('', { alias: 'class' });

  protected readonly isLoaded = computed(() => this.avatar.imageStatus() === 'loaded');

  protected readonly computedClass = computed(() =>
    cn('aspect-square size-full', this.userClass()),
  );

  constructor() {
    effect(() => {
      const imgSrc = this.src();
      if (!imgSrc) return;

      const img = new Image();
      img.onload = () => this.avatar.imageStatus.set('loaded');
      img.onerror = () => this.avatar.imageStatus.set('error');
      img.src = imgSrc;
    });
  }
}

@Component({
  selector: 'app-avatar-fallback',
  standalone: true,
  template: `@if (isVisible()) { <ng-content /> }`,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"avatar-fallback"',
  },
})
export class AvatarFallback {
  private readonly avatar = inject(Avatar);

  readonly userClass = input<string>('', { alias: 'class' });

  protected readonly isVisible = computed(() => this.avatar.imageStatus() !== 'loaded');

  protected readonly computedClass = computed(() =>
    cn('bg-muted flex size-full items-center justify-center rounded-full', this.userClass()),
  );
}
