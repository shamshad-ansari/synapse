import { Component, Input } from '@angular/core';

type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost';
type ButtonSize = 'sm' | 'md' | 'lg';

@Component({
  selector: 'app-landing-button',
  standalone: true,
  templateUrl: './landing-button.component.html',
})
export class LandingButtonComponent {
  @Input() variant: ButtonVariant = 'primary';
  @Input() size: ButtonSize = 'md';
  @Input() extraClass = '';
  @Input() disabled = false;
  @Input() type: 'button' | 'submit' | 'reset' = 'button';

  private readonly baseStyles =
    "inline-flex items-center justify-center font-['Inter'] font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-slate-900 focus-visible:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none rounded-full";

  private readonly variants: Record<ButtonVariant, string> = {
    primary:
      'bg-[#0F172A] text-white hover:bg-[#1E293B] shadow-[0_4px_20px_-2px_rgba(0,0,0,0.05)]',
    secondary:
      'bg-[#F9FAFB] text-[#0F172A] hover:bg-gray-100 border border-gray-200',
    outline:
      'border border-gray-200 bg-transparent hover:bg-gray-50 text-[#0F172A]',
    ghost:
      'bg-transparent hover:bg-gray-50 text-gray-600 hover:text-[#0F172A]',
  };

  private readonly sizes: Record<ButtonSize, string> = {
    sm: 'h-9 px-4 text-sm',
    md: 'h-11 px-6 text-base',
    lg: 'h-14 px-8 text-lg',
  };

  get classes(): string {
    return `${this.baseStyles} ${this.variants[this.variant]} ${this.sizes[this.size]} ${this.extraClass}`.trim();
  }
}
