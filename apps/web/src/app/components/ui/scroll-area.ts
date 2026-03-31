import { Component, computed, input } from '@angular/core';
import { cn } from './utils';

@Component({
  selector: 'app-scroll-area',
  standalone: true,
  template: `
    <div
      data-slot="scroll-area-viewport"
      [class]="viewportClass()"
      tabindex="0"
    >
      <ng-content />
    </div>
  `,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"scroll-area"',
  },
  styles: [`
    :host {
      display: block;
      overflow: hidden;
    }

    .scroll-viewport {
      scrollbar-width: thin;
      scrollbar-color: var(--border) transparent;
    }

    .scroll-viewport::-webkit-scrollbar {
      width: 10px;
      height: 10px;
    }

    .scroll-viewport::-webkit-scrollbar-track {
      background: transparent;
    }

    .scroll-viewport::-webkit-scrollbar-thumb {
      background-color: var(--border);
      border-radius: 9999px;
      border: 2px solid transparent;
      background-clip: content-box;
    }

    .scroll-viewport::-webkit-scrollbar-corner {
      background: transparent;
    }
  `],
})
export class ScrollArea {
  readonly userClass = input<string>('', { alias: 'class' });
  readonly orientation = input<'vertical' | 'horizontal' | 'both'>('vertical');

  protected readonly computedClass = computed(() =>
    cn('relative', this.userClass()),
  );

  protected readonly viewportClass = computed(() => {
    const overflow =
      this.orientation() === 'vertical'
        ? 'overflow-y-auto overflow-x-hidden'
        : this.orientation() === 'horizontal'
          ? 'overflow-x-auto overflow-y-hidden'
          : 'overflow-auto';

    return cn(
      'scroll-viewport focus-visible:ring-ring/50 size-full rounded-[inherit] transition-[color,box-shadow] outline-none focus-visible:ring-[3px] focus-visible:outline-1',
      overflow,
    );
  });
}
