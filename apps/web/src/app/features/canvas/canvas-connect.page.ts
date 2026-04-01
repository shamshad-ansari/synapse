import {
  Component, ChangeDetectionStrategy, inject, signal, computed, OnInit,
} from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { LucideAngularModule } from 'lucide-angular';
import { CanvasService } from './canvas.service';
import { environment } from '../../../environments/environment';

/** http(s) URL with a host; includes localhost (no period in hostname). */
function isValidCanvasInstitutionUrl(raw: string): boolean {
  const s = raw.trim();
  if (s.length < 4) {
    return false;
  }
  let toParse = s;
  if (!/^https?:\/\//i.test(s)) {
    toParse = `https://${s}`;
  }
  try {
    const u = new URL(toParse);
    if (u.protocol !== 'http:' && u.protocol !== 'https:') {
      return false;
    }
    return u.hostname.length > 0;
  } catch {
    return false;
  }
}

@Component({
  selector: 'app-canvas-connect-page',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [LucideAngularModule],
  template: `
    <div class="page">
      <div class="page-inner">
        <!-- Brand header: Step 1 only -->
        @if (currentStep() === 1) {
          <div class="brand">
            <img src="assets/synapse-logo.png" width="32" height="32" alt="Synapse" />
            <span class="brand-name">Synapse</span>
          </div>
        }

        <!-- STEP 1: Connect form -->
        @if (currentStep() === 1) {
          <div class="card">
            <div class="card-icon">
              <lucide-icon name="link" [size]="28" />
            </div>
            <h1 class="title">Connect your Canvas account</h1>
            <p class="subtitle">We'll sync your courses, deadlines, and grades automatically.</p>
            @if (canvasService.error()) {
              <p class="error-banner">{{ canvasService.error() }}</p>
            }

            <div class="form">
              <div class="field">
                <label class="label" for="canvas-url">Canvas URL</label>
                <input
                  id="canvas-url"
                  class="input"
                  type="text"
                  [value]="institutionUrl()"
                  (input)="onUrlInput($event)"
                  placeholder="university.instructure.com"
                />
                <span class="helper">Find this in your browser when you're logged into Canvas</span>
                @if (institutionUrl()) {
                  <div class="url-preview">
                    <span class="mono-text">{{ institutionUrl() }}/login/oauth2/auth</span>
                    @if (urlIsValid()) {
                      <lucide-icon name="check-circle" [size]="14" class="check-icon" />
                    }
                  </div>
                }
              </div>

              <button
                type="button"
                class="submit-btn"
                [disabled]="!urlIsValid()"
                (click)="onConnect()"
              >
                Connect Canvas
              </button>
            </div>

            <p class="skip-link">
              <a (click)="onSkip()" class="link">Skip for now &rarr;</a>
            </p>
          </div>
        }

        <!-- STEP 2: Success -->
        @if (currentStep() === 2 && connectionStatus() === 'success') {
          <div class="card card-center">
            <lucide-icon name="check-circle" [size]="48" class="pop-in success-icon" />
            <h1 class="title center">Canvas connected!</h1>
            <p class="subtitle center">Loading your courses...</p>
            <lucide-icon name="refresh-cw" [size]="20" class="spinner muted-spinner" />
          </div>
        }

        <!-- STEP 2: Error -->
        @if (currentStep() === 2 && connectionStatus() === 'error') {
          <div class="card card-center">
            <lucide-icon name="x-circle" [size]="48" class="error-icon" />
            <h1 class="title center">Connection failed</h1>
            <p class="subtitle center">Something went wrong with the Canvas connection.</p>
            <button type="button" class="submit-btn" (click)="onRetry()">
              Try again
            </button>
          </div>
        }

        <!-- STEP 3: All Set -->
        @if (currentStep() === 3) {
          <div class="card card-center">
            <lucide-icon name="check-circle" [size]="48" class="emerald-icon" />
            <h1 class="title center">You're all set!</h1>
            <p class="subtitle center">Your courses are syncing. This takes about a minute.</p>

            <div class="sync-rows">
              <div class="sync-row">
                <span class="pulse-dot" style="animation-delay: 0s"></span>
                <span class="sync-label">Courses</span>
              </div>
              <div class="sync-row">
                <span class="pulse-dot" style="animation-delay: 0.4s"></span>
                <span class="sync-label">Assignments & deadlines</span>
              </div>
              <div class="sync-row">
                <span class="pulse-dot" style="animation-delay: 0.8s"></span>
                <span class="sync-label">Grades</span>
              </div>
            </div>

            <button type="button" class="submit-btn" (click)="onGoDashboard()">
              Go to dashboard &rarr;
            </button>
          </div>
        }
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
    .card-center {
      display: flex;
      flex-direction: column;
      align-items: center;
      text-align: center;
    }
    .card-icon {
      display: flex;
      justify-content: center;
      margin-bottom: 16px;
      color: var(--navy);
    }
    .title {
      font-family: var(--font-display);
      font-size: 20px;
      font-weight: 700;
      color: var(--ink);
      margin: 0 0 4px;
      line-height: 1.3;
    }
    .center { text-align: center; }
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
    .input:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
    .helper {
      font-size: 11.5px;
      color: var(--ink-faint);
      line-height: 1.4;
    }
    .url-preview {
      display: flex;
      align-items: center;
      gap: 6px;
    }
    .mono-text {
      font-size: 11px;
      font-family: var(--mono);
      color: var(--ink-faint);
    }
    .check-icon { color: var(--emerald); }
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
    .muted-spinner {
      color: var(--ink-faint);
      margin-top: 8px;
    }
    @keyframes spin { to { transform: rotate(360deg); } }
    .skip-link {
      text-align: center;
      margin: 20px 0 0;
      font-size: 13px;
    }
    .link {
      color: var(--ink-muted);
      font-weight: 500;
      text-decoration: none;
      cursor: pointer;
    }
    .link:hover { text-decoration: underline; }

    /* Step 2 success animation */
    .success-icon {
      color: var(--emerald);
      margin-bottom: 16px;
    }
    .pop-in {
      animation: popIn 400ms cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
    }
    @keyframes popIn {
      0% { transform: scale(0.5); opacity: 0; }
      70% { transform: scale(1.1); }
      100% { transform: scale(1); opacity: 1; }
    }

    /* Step 2 error */
    .error-icon {
      color: var(--red);
      margin-bottom: 16px;
    }

    /* Step 3 */
    .emerald-icon {
      color: var(--emerald);
      margin-bottom: 16px;
    }
    .sync-rows {
      display: flex;
      flex-direction: column;
      gap: 12px;
      width: 100%;
      margin-bottom: 24px;
    }
    .sync-row {
      display: flex;
      align-items: center;
      gap: 10px;
    }
    .pulse-dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: var(--emerald);
      animation: pulse 1.2s ease-in-out infinite;
    }
    @keyframes pulse {
      0%, 100% { opacity: 0.4; }
      50% { opacity: 1; }
    }
    .sync-label {
      font-size: 13.5px;
      color: var(--ink);
    }
  `],
})
export class CanvasConnectPage implements OnInit {
  protected readonly canvasService = inject(CanvasService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);

  readonly currentStep = signal<1 | 2 | 3>(1);
  readonly connectionStatus = signal<'success' | 'error' | null>(null);

  readonly institutionUrl = signal(environment.production ? '' : 'http://localhost:8082');

  /** Accepts http(s) URLs including localhost (no dot) and hostnames like school.instructure.com. */
  readonly urlIsValid = computed(() => isValidCanvasInstitutionUrl(this.institutionUrl()));

  ngOnInit(): void {
    const url = this.router.url;
    if (url.startsWith('/canvas/connected')) {
      const status = this.route.snapshot.queryParamMap.get('status');
      if (status === 'success') {
        this.connectionStatus.set('success');
        this.currentStep.set(2);
        this.advanceToAllSet();
      } else if (status === 'error') {
        this.connectionStatus.set('error');
        this.currentStep.set(2);
      }
    }
  }

  onUrlInput(event: Event): void {
    this.institutionUrl.set((event.target as HTMLInputElement).value);
  }

  onConnect(): void {
    this.canvasService.connectCanvas(this.institutionUrl());
  }

  onRetry(): void {
    this.canvasService.error.set(null);
    this.connectionStatus.set(null);
    this.currentStep.set(1);
  }

  onSkip(): void {
    this.router.navigate(['/today']);
  }

  onGoDashboard(): void {
    this.router.navigate(['/today']);
  }

  private advanceToAllSet(): void {
    setTimeout(async () => {
      try {
        await this.canvasService.triggerSync();
      } catch { /* sync is best-effort */ }
      await this.canvasService.loadStatus();
      this.currentStep.set(3);
    }, 2000);
  }
}
