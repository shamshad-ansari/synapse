import { Component, ChangeDetectionStrategy, inject, signal, DestroyRef } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import {
  ReactiveFormsModule, FormBuilder, Validators,
  AbstractControl, ValidationErrors, ValidatorFn,
} from '@angular/forms';
import { LucideAngularModule } from 'lucide-angular';
import { AuthService } from '../../core/auth/auth.service';
import { CanvasService } from '../canvas/canvas.service';

const passwordsMatch: ValidatorFn = (group: AbstractControl): ValidationErrors | null => {
  const pw = group.get('password')?.value;
  const cpw = group.get('confirm_password')?.value;
  if (!cpw) return null;
  return pw === cpw ? null : { passwordsMismatch: true };
};

@Component({
  selector: 'app-register-page',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [ReactiveFormsModule, RouterLink, LucideAngularModule],
  template: `
    <div class="page">
      <div class="page-inner">
        <div class="brand">
          <img src="assets/synapse-logo.png" width="32" height="32" alt="Synapse" />
          <span class="brand-name">Synapse</span>
        </div>

        <div class="card">
          <h1 class="title">Create your account</h1>
          <p class="subtitle">Start studying smarter</p>

          <form [formGroup]="form" (ngSubmit)="onSubmit()" class="form">
            <div class="field">
              <label class="label" for="reg-name">Full name</label>
              <input
                id="reg-name"
                class="input"
                type="text"
                formControlName="name"
                placeholder="Alex Kim"
                [class.input-error]="showError('name')"
              />
              @if (showError('name')) {
                <span class="field-error">Name is required</span>
              }
            </div>

            <div class="field">
              <label class="label" for="reg-email">Email address</label>
              <input
                id="reg-email"
                class="input"
                type="email"
                formControlName="email"
                placeholder="you&#64;university.edu"
                [class.input-error]="showError('email')"
              />
              @if (showError('email')) {
                <span class="field-error">Enter a valid email</span>
              }
            </div>

            <div class="field">
              <label class="label" for="reg-school">School domain</label>
              <input
                id="reg-school"
                class="input"
                type="text"
                formControlName="school_domain"
                placeholder="university.edu"
                [class.input-error]="showError('school_domain')"
              />
              <span class="helper">Your university's Canvas URL</span>
              @if (showError('school_domain')) {
                <span class="field-error">School domain is required</span>
              }
            </div>

            <div class="field">
              <label class="label" for="reg-password">Password</label>
              <div class="password-wrapper">
                <input
                  id="reg-password"
                  class="input password-input"
                  [type]="showPassword() ? 'text' : 'password'"
                  formControlName="password"
                  placeholder="••••••••"
                  [class.input-error]="showError('password')"
                />
                <button
                  type="button"
                  class="toggle-pw"
                  (click)="showPassword.set(!showPassword())"
                  tabindex="-1"
                >
                  <lucide-icon
                    [name]="showPassword() ? 'eye-off' : 'eye'"
                    [size]="16"
                  />
                </button>
              </div>
              @if (form.controls.password.value) {
                <div class="strength-track">
                  <div
                    class="strength-bar"
                    [style.width]="strengthWidth()"
                    [style.background]="strengthColor()"
                  ></div>
                </div>
              }
              @if (showError('password')) {
                @if (form.controls.password.hasError('required')) {
                  <span class="field-error">Password is required</span>
                } @else {
                  <span class="field-error">Minimum 8 characters</span>
                }
              }
            </div>

            <div class="field">
              <label class="label" for="reg-confirm">Confirm password</label>
              <input
                id="reg-confirm"
                class="input"
                type="password"
                formControlName="confirm_password"
                placeholder="••••••••"
                [class.input-error]="showConfirmError()"
              />
              @if (showConfirmError()) {
                <span class="field-error">Passwords don't match</span>
              }
            </div>

            @if (authService.error()) {
              <div class="error-banner">{{ authService.error() }}</div>
            }

            <button
              type="submit"
              class="submit-btn"
              [disabled]="authService.loading() || form.invalid"
            >
              @if (authService.loading()) {
                <lucide-icon name="refresh-cw" [size]="16" class="spinner" />
                Creating account…
              } @else {
                Create account
              }
            </button>
          </form>

          <div class="divider"><span>or</span></div>

          <p class="footer">
            Already have an account?
            <a routerLink="/login" class="link">Sign in</a>
          </p>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .page {
      min-height: 100vh;
      background: var(--bg);
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 24px;
    }
    .page-inner {
      display: flex;
      flex-direction: column;
      align-items: center;
      width: 100%;
      max-width: 400px;
    }
    .brand {
      display: flex;
      align-items: center;
      gap: 10px;
      margin-bottom: 24px;
    }
    .brand-name {
      font-family: var(--font-display);
      font-size: 20px;
      font-weight: 700;
      color: var(--ink);
    }
    .card {
      width: 100%;
      padding: 36px 32px 32px;
      background: #ffffff;
      border: 1px solid var(--divider);
      border-radius: var(--r-xl);
      box-shadow: var(--shadow-md);
    }
    .title {
      font-family: var(--font-display);
      font-size: 22px;
      font-weight: 700;
      color: var(--ink);
      margin: 0 0 4px;
      line-height: 1.3;
    }
    .subtitle {
      font-family: var(--font);
      font-size: 13.5px;
      color: var(--ink-muted);
      margin: 0 0 24px;
    }
    .form {
      display: flex;
      flex-direction: column;
      gap: 18px;
    }
    .field {
      display: flex;
      flex-direction: column;
      gap: 5px;
    }
    .label {
      font-size: 13px;
      font-weight: 500;
      color: var(--ink);
      line-height: 1.4;
    }
    .input {
      background: var(--surface-sub);
      border: 1px solid var(--divider);
      border-radius: var(--r-md);
      padding: 9px 12px;
      font-size: 14px;
      font-family: var(--font);
      color: var(--ink);
      outline: none;
      transition: border-color 150ms ease;
      width: 100%;
      box-sizing: border-box;
    }
    .input:focus { border-color: var(--navy); }
    .input::placeholder { color: var(--ink-ghost); }
    .input-error { border-color: var(--red); }
    .input:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
    .password-wrapper {
      position: relative;
      display: flex;
      align-items: center;
    }
    .password-input { padding-right: 38px; }
    .toggle-pw {
      position: absolute;
      right: 10px;
      top: 50%;
      transform: translateY(-50%);
      background: none;
      border: none;
      padding: 0;
      cursor: pointer;
      color: var(--ink-faint);
      display: flex;
      align-items: center;
    }
    .toggle-pw:active { transform: translateY(-50%) scale(1); }
    .strength-track {
      height: 3px;
      border-radius: 2px;
      background: var(--divider);
      overflow: hidden;
    }
    .strength-bar {
      height: 100%;
      border-radius: 2px;
      transition: width 200ms ease;
    }
    .helper {
      font-size: 11.5px;
      color: var(--ink-faint);
      line-height: 1.4;
    }
    .field-error {
      font-size: 12px;
      color: var(--red);
      line-height: 1.4;
    }
    .error-banner {
      color: var(--red);
      font-size: 13px;
      padding: 8px 12px;
      background: var(--red-light);
      border-radius: var(--r-md);
    }
    .submit-btn {
      width: 100%;
      padding: 10px 0;
      font-size: 14px;
      font-weight: 600;
      font-family: var(--font);
      color: #fff;
      background: var(--navy);
      border: none;
      border-radius: var(--r-lg);
      cursor: pointer;
      transition: opacity 150ms ease;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
    }
    .submit-btn:hover:not(:disabled) { opacity: 0.9; }
    .submit-btn:disabled { opacity: 0.6; cursor: not-allowed; }
    .submit-btn:active:not(:disabled) { transform: scale(1); }
    .spinner {
      animation: spin 0.8s linear infinite;
    }
    @keyframes spin { to { transform: rotate(360deg); } }
    .divider {
      display: flex;
      align-items: center;
      gap: 12px;
      margin: 24px 0 20px;
      color: var(--ink-faint);
      font-size: 12px;
    }
    .divider::before,
    .divider::after {
      content: '';
      flex: 1;
      height: 1px;
      background: var(--divider);
    }
    .footer {
      text-align: center;
      font-size: 13px;
      color: var(--ink-muted);
      margin: 0;
    }
    .link {
      color: var(--navy);
      font-weight: 500;
      text-decoration: none;
    }
    .link:hover { text-decoration: underline; }
  `],
})
export class RegisterPage {
  protected readonly authService = inject(AuthService);
  private readonly canvasService = inject(CanvasService);
  private readonly router = inject(Router);
  private readonly fb = inject(FormBuilder);
  private readonly destroyRef = inject(DestroyRef);

  readonly showPassword = signal(false);

  readonly form = this.fb.nonNullable.group(
    {
      name: ['', [Validators.required]],
      email: ['', [Validators.required, Validators.email]],
      school_domain: ['', [Validators.required]],
      password: ['', [Validators.required, Validators.minLength(8)]],
      confirm_password: ['', [Validators.required]],
    },
    { validators: passwordsMatch },
  );

  constructor() {
    const sub = this.form.valueChanges.subscribe(() => {
      if (this.authService.error()) {
        this.authService.error.set(null);
      }
    });
    this.destroyRef.onDestroy(() => sub.unsubscribe());
  }

  showError(field: 'name' | 'email' | 'school_domain' | 'password'): boolean {
    const ctrl = this.form.controls[field];
    return ctrl.invalid && ctrl.touched;
  }

  showConfirmError(): boolean {
    const ctrl = this.form.controls.confirm_password;
    return ctrl.touched && (ctrl.invalid || this.form.hasError('passwordsMismatch'));
  }

  private getPasswordStrength(value: string): 0 | 1 | 2 {
    if (value.length >= 10 && /\d/.test(value) && /[^a-zA-Z0-9]/.test(value)) return 2;
    if (value.length >= 8 && /\d/.test(value)) return 1;
    return 0;
  }

  strengthWidth(): string {
    const s = this.getPasswordStrength(this.form.controls.password.value);
    return s === 0 ? '33%' : s === 1 ? '66%' : '100%';
  }

  strengthColor(): string {
    const s = this.getPasswordStrength(this.form.controls.password.value);
    return s === 0 ? 'var(--red)' : s === 1 ? 'var(--amber)' : 'var(--emerald)';
  }

  async onSubmit() {
    this.form.markAllAsTouched();
    if (this.form.invalid) return;

    const { name, email, school_domain, password } = this.form.getRawValue();
    this.authService.error.set(null);
    try {
      await this.authService.register(name, email, password, school_domain);
      await this.canvasService.loadStatus();
      if (this.canvasService.status() === null) {
        this.router.navigate(['/canvas/connect']);
      } else {
        this.router.navigate(['/today']);
      }
    } catch {
      // error displayed via authService.error()
    }
  }
}
