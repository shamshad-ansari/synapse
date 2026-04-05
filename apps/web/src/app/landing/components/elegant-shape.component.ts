import { Component, Input, ElementRef, ViewChild, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-elegant-shape',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div
      #outerContainer
      class="absolute transition-transform duration-[2400ms] opacity-0"
      [class]="className"
      [style.transform]="'translateY(-150px) rotate(' + (rotate - 15) + 'deg)'"
    >
      <div
        #innerContainer
        class="relative float-animation"
        [style.width.px]="width"
        [style.height.px]="height"
      >
        <div
          class="absolute inset-0 rounded-2xl flex flex-col justify-center gap-[15%] p-[10%]"
          [class]="colorClass"
        >
          <!-- Faint parallel lines to simulate text -->
          <div class="w-full h-[3.5px] bg-white/45 rounded-full"></div>
          <div class="w-5/6 h-[3.5px] bg-white/45 rounded-full"></div>
          <div class="w-11/12 h-[3.5px] bg-white/45 rounded-full"></div>
          <div class="w-4/5 h-[3.5px] bg-white/45 rounded-full"></div>
        </div>
      </div>
    </div>
  `,
  styles: [`
    @keyframes float {
      0%, 100% { transform: translateY(0); }
      50% { transform: translateY(15px); }
    }
    .float-animation {
      animation: float 12s ease-in-out infinite;
    }
  `]
})
export class ElegantShapeComponent implements AfterViewInit {
  @Input() className = '';
  @Input() delay = 0;
  @Input() width = 400;
  @Input() height = 100;
  @Input() rotate = 0;
  @Input() colorClass = 'from-emerald-500/[0.15]';

  @ViewChild('outerContainer') outerContainer!: ElementRef<HTMLDivElement>;

  ngAfterViewInit() {
    setTimeout(() => {
      if (this.outerContainer) {
        this.outerContainer.nativeElement.style.transform = `translateY(0) rotate(${this.rotate}deg)`;
        this.outerContainer.nativeElement.style.opacity = '1';
        this.outerContainer.nativeElement.style.transitionDelay = `${this.delay}s`;
        this.outerContainer.nativeElement.style.transitionProperty = 'transform, opacity';
        this.outerContainer.nativeElement.style.transitionTimingFunction = 'cubic-bezier(0.23, 0.86, 0.39, 0.96), ease';
        this.outerContainer.nativeElement.style.transitionDuration = '2.4s, 1.2s';
      }
    }, 50);
  }
}
