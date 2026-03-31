import { computed, Directive, input } from '@angular/core';
import { cn } from './utils';

@Directive({
  selector: 'label[appLabel]',
  standalone: true,
  host: {
    '[class]': 'computedClass()',
    '[attr.data-slot]': '"label"',
  },
})
export class Label {
  readonly userClass = input<string>('', { alias: 'class' });

  protected readonly computedClass = computed(() =>
    cn(
      'flex items-center gap-2 text-sm leading-none font-medium select-none group-data-[disabled=true]:pointer-events-none group-data-[disabled=true]:opacity-50 peer-disabled:cursor-not-allowed peer-disabled:opacity-50',
      this.userClass(),
    ),
  );
}
