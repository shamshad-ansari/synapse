import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { AuthService } from '../auth/auth.service';
import { environment } from '../../../environments/environment';

const SKIP_PATHS = ['/v1/auth/login', '/v1/auth/register'];

export const jwtInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const token = authService.getAccessToken();

  const isApiUrl = req.url.startsWith(environment.apiUrl) || req.url.startsWith(environment.lmsApiUrl);
  if (token && isApiUrl) {
    const shouldSkip = SKIP_PATHS.some(path => req.url.includes(path));
    if (!shouldSkip) {
      req = req.clone({
        setHeaders: { Authorization: `Bearer ${token}` },
      });
    }
  }

  return next(req);
};
