import { Component } from '@angular/core';

@Component({
  selector: 'app-trusted-by',
  standalone: true,
  template: `
    <section class="bg-white py-16 md:py-24 overflow-hidden relative">
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-10">
          <p class="text-sm uppercase tracking-[0.2em] font-semibold" style="color: var(--ink-faint);">
            Trusted by students from top universities
          </p>
        </div>

        <div class="relative w-full overflow-hidden flex items-center h-16 [mask-image:_linear-gradient(to_right,transparent_0,_black_128px,_black_calc(100%-128px),transparent_100%)]">
          <div class="flex items-center gap-16 md:gap-24 scroller pr-16 md:pr-24">
            @for (school of schools; track school) {
              <div class="flex items-center justify-center grayscale opacity-60 hover:opacity-100 hover:grayscale-0 transition-all duration-300 whitespace-nowrap">
                <span class="text-2xl md:text-3xl font-display font-medium" style="color: var(--ink-muted);">
                  {{ school }}
                </span>
              </div>
            }
            <!-- Duplicate for seamless scroll -->
            @for (school of schools; track school + '-duplicate') {
              <div class="flex items-center justify-center grayscale opacity-60 hover:opacity-100 hover:grayscale-0 transition-all duration-300 whitespace-nowrap">
                <span class="text-2xl md:text-3xl font-display font-medium" style="color: var(--ink-muted);">
                  {{ school }}
                </span>
              </div>
            }
          </div>
        </div>
      </div>
    </section>
  `,
  styles: `
    .scroller {
      animation: scroll 30s linear infinite;
      width: max-content;
    }
    @keyframes scroll {
      0% {
        transform: translateX(0);
      }
      100% {
        transform: translateX(-50%);
      }
    }
  `
})
export class TrustedByComponent {
  schools = [
    'Stanford University',
    'MIT',
    'Vanderbilt',
    'Yale University',
    'Harvard',
    'UC Berkeley',
    'Columbia University',
    'Princeton'
  ];
}
