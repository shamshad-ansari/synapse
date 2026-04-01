import { Routes } from '@angular/router';
import { authGuard } from './core/auth/auth.guard';
import { publicGuard } from './core/auth/public.guard';

export type TabType = 'today' | 'notes' | 'review' | 'planner' | 'feed' | 'tutoring' | 'profile';

export const routes: Routes = [
  { path: '', title: 'Synapse', loadComponent: () => import('./landing/landing.component').then(m => m.LandingComponent) },
  { path: 'login', title: 'Log In — Synapse', loadComponent: () => import('./features/auth/login.page').then(m => m.LoginPage), canActivate: [publicGuard] },
  { path: 'register', title: 'Register — Synapse', loadComponent: () => import('./features/auth/register.page').then(m => m.RegisterPage), canActivate: [publicGuard] },
  { path: 'today', title: 'Today', loadComponent: () => import('./components/today-tab.component'), canActivate: [authGuard] },
  { path: 'notes', title: 'Notes', loadComponent: () => import('./components/notes-tab.component'), canActivate: [authGuard] },
  { path: 'planner', title: 'Planner', loadComponent: () => import('./components/planner-tab.component'), canActivate: [authGuard] },
  { path: 'review', title: 'Review', loadComponent: () => import('./components/review-tab.component'), canActivate: [authGuard] },
  { path: 'feed', title: 'Feed', loadComponent: () => import('./components/feed-tab.component'), canActivate: [authGuard] },
  { path: 'tutoring', title: 'Tutoring', loadComponent: () => import('./components/tutoring-tab.component'), canActivate: [authGuard] },
  { path: 'profile', title: 'Profile', loadComponent: () => import('./components/profile-tab.component'), canActivate: [authGuard] },
  { path: 'canvas/connect', title: 'Connect Canvas — Synapse', loadComponent: () => import('./features/canvas/canvas-connect.page').then(m => m.CanvasConnectPage), canActivate: [authGuard] },
  { path: 'canvas/connected', title: 'Connect Canvas — Synapse', loadComponent: () => import('./features/canvas/canvas-connect.page').then(m => m.CanvasConnectPage), canActivate: [authGuard] },
  { path: '**', redirectTo: '' },
];
