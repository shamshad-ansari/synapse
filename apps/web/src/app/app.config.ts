import { ApplicationConfig, provideZoneChangeDetection, APP_INITIALIZER } from '@angular/core';
import { provideRouter, withComponentInputBinding } from '@angular/router';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { jwtInterceptor } from './core/interceptors/jwt.interceptor';
import { errorInterceptor } from './core/interceptors/error.interceptor';
import { AuthService } from './core/auth/auth.service';
import { LUCIDE_ICONS, LucideIconProvider } from 'lucide-angular';
import {
  Home, FileText, Zap, Calendar,
  MessageSquare, Users, User,
  MoreHorizontal, ChevronsLeft, ChevronsRight,
  Search, Bell, Clock, Target, BookOpen,
  AlertTriangle, RefreshCw, ChevronLeft, ChevronRight,
  Upload, Plus, Link, AlertCircle,
  ArrowUp, CheckCircle, X, HelpCircle,
  MessageCircle, Bookmark, ChevronDown, Info,
  Shield, Edit2, ShieldCheck, Eye, EyeOff,
  XCircle, Key, Save, Sparkles, Trash2,
} from 'lucide-angular';
import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes, withComponentInputBinding()),
    provideAnimationsAsync(),
    provideHttpClient(withInterceptors([jwtInterceptor, errorInterceptor])),
    {
      provide: APP_INITIALIZER,
      useFactory: (authService: AuthService) => () => authService.initAuth(),
      deps: [AuthService],
      multi: true,
    },
    {
      provide: LUCIDE_ICONS,
      multi: true,
      useValue: new LucideIconProvider({
        Home, FileText, Zap, Calendar,
        MessageSquare, Users, User,
        MoreHorizontal, ChevronsLeft, ChevronsRight,
        Search, Bell, Clock, Target, BookOpen,
        AlertTriangle, RefreshCw, ChevronLeft, ChevronRight,
        Upload, Plus, Link, AlertCircle,
        ArrowUp, CheckCircle, X, HelpCircle,
        MessageCircle, Bookmark, ChevronDown, Info,
        Shield, Edit2, ShieldCheck, Eye, EyeOff,
        XCircle, Key, Save, Sparkles, Trash2,
      }),
    },
  ],
};
