import { HttpInterceptorFn, HttpErrorResponse } from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, throwError } from 'rxjs';
import { AuthService } from '../auth/auth.service';

const AUTH_ATTEMPT_PATHS = ['/v1/auth/login', '/v1/auth/register'];

export const errorInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);

  return next(req).pipe(
    catchError((err: HttpErrorResponse) => {
      if (err.status === 401) {
        const isAuthAttempt = AUTH_ATTEMPT_PATHS.some((p) => req.url.includes(p));
        if (!isAuthAttempt) {
          authService.logout();
        }
      }
      return throwError(() => err);
    }),
  );
};
