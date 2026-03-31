import { Component, computed, inject, input, model } from '@angular/core';
import { cn } from './utils';

@Component({
  selector: 'app-tabs',
  standalone: true,
  template: `<ng-content />`,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"tabs"',
  },
})
export class Tabs {
  readonly userClass = input<string>('', { alias: 'class' });

  /** Active tab value. Bind with [(value)] for two-way or value="x" for initial. */
  readonly value = model<string>('');

  protected readonly computedClass = computed(() =>
    cn('flex flex-col gap-2', this.userClass()),
  );

  selectTab(tabValue: string) {
    this.value.set(tabValue);
  }
}

@Component({
  selector: 'app-tabs-list',
  standalone: true,
  template: `<ng-content />`,
  host: {
    'role': 'tablist',
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"tabs-list"',
  },
})
export class TabsList {
  readonly userClass = input<string>('', { alias: 'class' });

  protected readonly computedClass = computed(() =>
    cn(
      'bg-muted text-muted-foreground inline-flex h-9 w-fit items-center justify-center rounded-xl p-[3px] flex',
      this.userClass(),
    ),
  );
}

@Component({
  selector: 'app-tabs-trigger',
  standalone: true,
  template: `<ng-content />`,
  host: {
    'role': 'tab',
    'tabindex': '0',
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"tabs-trigger"',
    '[attr.data-state]': 'isActive() ? "active" : "inactive"',
    '[attr.aria-selected]': 'isActive()',
    '(click)': 'select()',
    '(keydown.enter)': 'select()',
    '(keydown.space)': 'onSpace($event)',
  },
})
export class TabsTrigger {
  private readonly tabs = inject(Tabs);

  readonly value = input.required<string>();
  readonly userClass = input<string>('', { alias: 'class' });

  protected readonly isActive = computed(() => this.tabs.value() === this.value());

  protected readonly computedClass = computed(() =>
    cn(
      'data-[state=active]:bg-card dark:data-[state=active]:text-foreground focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:outline-ring dark:data-[state=active]:border-input dark:data-[state=active]:bg-input/30 text-foreground dark:text-muted-foreground inline-flex h-[calc(100%-1px)] flex-1 items-center justify-center gap-1.5 rounded-xl border border-transparent px-2 py-1 text-sm font-medium whitespace-nowrap transition-[color,box-shadow] focus-visible:ring-[3px] focus-visible:outline-1 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*=\'size-\'])]:size-4',
      this.userClass(),
    ),
  );

  protected select() {
    this.tabs.selectTab(this.value());
  }

  protected onSpace(event: Event) {
    event.preventDefault();
    this.select();
  }
}

@Component({
  selector: 'app-tabs-content',
  standalone: true,
  template: `@if (isActive()) { <ng-content /> }`,
  host: {
    'role': 'tabpanel',
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"tabs-content"',
    '[attr.data-state]': 'isActive() ? "active" : "inactive"',
  },
  styles: [`:host { display: block; }`],
})
export class TabsContent {
  private readonly tabs = inject(Tabs);

  readonly value = input.required<string>();
  readonly userClass = input<string>('', { alias: 'class' });

  protected readonly isActive = computed(() => this.tabs.value() === this.value());

  protected readonly computedClass = computed(() =>
    cn('flex-1 outline-none', this.userClass()),
  );
}
