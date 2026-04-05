import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { AuthService } from '../auth/auth.service';
import { environment } from '../../../environments/environment';

const SKIP_AUTH_PATHS = ['/v1/auth/login', '/v1/auth/register'];

/**
 * Attaches Bearer token to Synapse API requests.
 * Supports full `apiUrl` (e.g. http://localhost:8080) or same-origin `/v1/*` when `apiUrl` is empty (dev proxy).
 */
export const jwtInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const token = authService.getAccessToken();

  if (!token) {
    return next(req);
  }

  if (SKIP_AUTH_PATHS.some((path) => req.url.includes(path))) {
    return next(req);
  }

  const url = req.url;
  const apiBase = (environment.apiUrl ?? '').replace(/\/$/, '');
  const useRelativeGateway = apiBase.length === 0;

  const hitsGateway =
    (useRelativeGateway && url.includes('/v1/')) ||
    (apiBase.length > 0 && url.startsWith(apiBase)) ||
    (!!environment.lmsApiUrl && url.startsWith(environment.lmsApiUrl));

  if (hitsGateway) {
    req = req.clone({
      setHeaders: { Authorization: `Bearer ${token}` },
    });
  }

  return next(req);
};
