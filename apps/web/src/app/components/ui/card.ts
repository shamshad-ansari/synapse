import { Component, computed, input } from '@angular/core';
import { cn } from './utils';

@Component({
  selector: 'app-card',
  standalone: true,
  template: `<ng-content />`,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"card"',
  },
})
export class Card {
  readonly userClass = input<string>('', { alias: 'class' });
  protected readonly computedClass = computed(() =>
    cn('bg-card text-card-foreground flex flex-col gap-6 rounded-xl border', this.userClass()),
  );
}

@Component({
  selector: 'app-card-header',
  standalone: true,
  template: `<ng-content />`,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"card-header"',
  },
})
export class CardHeader {
  readonly userClass = input<string>('', { alias: 'class' });
  protected readonly computedClass = computed(() =>
    cn(
      '@container/card-header grid auto-rows-min grid-rows-[auto_auto] items-start gap-1.5 px-6 pt-6 has-data-[slot=card-action]:grid-cols-[1fr_auto] [.border-b]:pb-6',
      this.userClass(),
    ),
  );
}

@Component({
  selector: 'app-card-title',
  standalone: true,
  template: `<ng-content />`,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"card-title"',
    'role': 'heading',
    'aria-level': '4',
  },
  styles: [`:host { display: block; }`],
})
export class CardTitle {
  readonly userClass = input<string>('', { alias: 'class' });
  protected readonly computedClass = computed(() =>
    cn('leading-none', this.userClass()),
  );
}

@Component({
  selector: 'app-card-description',
  standalone: true,
  template: `<ng-content />`,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"card-description"',
  },
  styles: [`:host { display: block; }`],
})
export class CardDescription {
  readonly userClass = input<string>('', { alias: 'class' });
  protected readonly computedClass = computed(() =>
    cn('text-muted-foreground', this.userClass()),
  );
}

@Component({
  selector: 'app-card-action',
  standalone: true,
  template: `<ng-content />`,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"card-action"',
  },
})
export class CardAction {
  readonly userClass = input<string>('', { alias: 'class' });
  protected readonly computedClass = computed(() =>
    cn('col-start-2 row-span-2 row-start-1 self-start justify-self-end', this.userClass()),
  );
}

@Component({
  selector: 'app-card-content',
  standalone: true,
  template: `<ng-content />`,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"card-content"',
  },
  styles: [`:host { display: block; }`],
})
export class CardContent {
  readonly userClass = input<string>('', { alias: 'class' });
  protected readonly computedClass = computed(() =>
    cn('px-6 [&:last-child]:pb-6', this.userClass()),
  );
}

@Component({
  selector: 'app-card-footer',
  standalone: true,
  template: `<ng-content />`,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"card-footer"',
  },
})
export class CardFooter {
  readonly userClass = input<string>('', { alias: 'class' });
  protected readonly computedClass = computed(() =>
    cn('flex items-center px-6 pb-6 [.border-t]:pt-6', this.userClass()),
  );
}
