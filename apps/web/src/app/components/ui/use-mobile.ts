import { signal, DestroyRef, inject } from '@angular/core';

const MOBILE_BREAKPOINT = 768;

export function injectIsMobile() {
  const isMobile = signal(
    typeof window !== 'undefined' ? window.innerWidth < MOBILE_BREAKPOINT : false
  );

  if (typeof window !== 'undefined') {
    const mql = window.matchMedia(`(max-width: ${MOBILE_BREAKPOINT - 1}px)`);
    const onChange = () => isMobile.set(window.innerWidth < MOBILE_BREAKPOINT);
    mql.addEventListener('change', onChange);

    inject(DestroyRef).onDestroy(() => mql.removeEventListener('change', onChange));
  }

  return isMobile.asReadonly();
}
