package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	authorizationCodeTTL = 5 * time.Minute
	accessTokenTTL       = 10 * time.Hour
	staticDevPAT         = "mock-canvas-pat-dev"
)

type issuedCode struct {
	ClientID    string
	RedirectURI string
	ExpiresAt   time.Time
}

var (
	oauthStoreMu      sync.Mutex
	issuedCodes       = make(map[string]issuedCode)
	issuedAccessToken = make(map[string]time.Time)
	issuedRefreshToken = make(map[string]time.Time)
)

// AuthorizePage renders a Canvas-style OAuth consent screen.
// GET /login/oauth2/auth?client_id=&redirect_uri=&response_type=code&state=
func AuthorizePage(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	responseType := r.URL.Query().Get("response_type")
	state := r.URL.Query().Get("state")

	if clientID == "" {
		http.Error(w, "client_id is required", http.StatusBadRequest)
		return
	}
	if redirectURI == "" {
		http.Error(w, "redirect_uri is required", http.StatusBadRequest)
		return
	}
	if responseType != "code" {
		http.Error(w, "response_type must be code", http.StatusBadRequest)
		return
	}

	approveURL := fmt.Sprintf("/login/oauth2/approve?client_id=%s&redirect_uri=%s&state=%s",
		url.QueryEscape(clientID), url.QueryEscape(redirectURI), url.QueryEscape(state))
	cancelURL := buildRedirectWithQuery(redirectURI, map[string]string{
		"error": "access_denied",
		"state": state,
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, consentPageHTML, approveURL, cancelURL)
}

// ApproveAuthorization handles the authorize button click, generating a code
// and redirecting back to the lms-service callback.
// GET /login/oauth2/approve?redirect_uri=&state=
func ApproveAuthorization(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")

	if clientID == "" {
		http.Error(w, "client_id is required", http.StatusBadRequest)
		return
	}
	if redirectURI == "" {
		http.Error(w, "redirect_uri is required", http.StatusBadRequest)
		return
	}

	code := generateMockCode()
	oauthStoreMu.Lock()
	issuedCodes[code] = issuedCode{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		ExpiresAt:   time.Now().Add(authorizationCodeTTL),
	}
	oauthStoreMu.Unlock()

	location := buildRedirectWithQuery(redirectURI, map[string]string{
		"code":  code,
		"state": state,
	})
	http.Redirect(w, r, location, http.StatusFound)
}

// TokenExchange exchanges an authorization code for mock tokens.
// POST /login/oauth2/token (application/x-www-form-urlencoded)
func TokenExchange(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form body", http.StatusBadRequest)
		return
	}

	grantType := r.Form.Get("grant_type")
	clientID := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")
	redirectURI := r.Form.Get("redirect_uri")
	code := r.Form.Get("code")

	if grantType == "refresh_token" {
		handleRefreshGrant(w, r)
		return
	}
	if grantType != "authorization_code" {
		writeOAuthError(w, http.StatusBadRequest, "unsupported_grant_type", "grant_type must be authorization_code or refresh_token")
		return
	}
	if clientID == "" || clientSecret == "" || redirectURI == "" || code == "" {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "client_id, client_secret, redirect_uri, and code are required")
		return
	}

	oauthStoreMu.Lock()
	issued, ok := issuedCodes[code]
	if !ok {
		oauthStoreMu.Unlock()
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "authorization code is invalid")
		return
	}
	delete(issuedCodes, code) // single use
	oauthStoreMu.Unlock()

	if time.Now().After(issued.ExpiresAt) {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "authorization code expired")
		return
	}
	if issued.ClientID != clientID || issued.RedirectURI != redirectURI {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "authorization code does not match client or redirect_uri")
		return
	}

	accessToken := "mock-canvas-token-" + generateHex(16)
	refreshToken := "mock-canvas-refresh-" + generateHex(16)
	registerIssuedAccessToken(accessToken, time.Now().Add(accessTokenTTL))
	registerIssuedRefreshToken(refreshToken, time.Now().Add(7*24*time.Hour))

	writeJSON(w, http.StatusOK, map[string]any{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"refresh_token": refreshToken,
		"expires_in":    36000,
	})
}

func handleRefreshGrant(w http.ResponseWriter, r *http.Request) {
	clientID := r.Form.Get("client_id")
	clientSecret := r.Form.Get("client_secret")
	refreshToken := r.Form.Get("refresh_token")
	if clientID == "" || clientSecret == "" || refreshToken == "" {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "client_id, client_secret, and refresh_token are required")
		return
	}
	if !isIssuedRefreshToken(refreshToken) {
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "refresh_token is invalid or expired")
		return
	}

	newAccessToken := "mock-canvas-token-" + generateHex(16)
	newRefreshToken := "mock-canvas-refresh-" + generateHex(16)
	registerIssuedAccessToken(newAccessToken, time.Now().Add(accessTokenTTL))
	registerIssuedRefreshToken(newRefreshToken, time.Now().Add(7*24*time.Hour))

	writeJSON(w, http.StatusOK, map[string]any{
		"access_token":  newAccessToken,
		"token_type":    "Bearer",
		"refresh_token": newRefreshToken,
		"expires_in":    36000,
	})
}

func generateMockCode() string {
	return "mock-auth-code-" + generateHex(16)
}

func generateHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func registerIssuedAccessToken(token string, expiresAt time.Time) {
	oauthStoreMu.Lock()
	defer oauthStoreMu.Unlock()
	issuedAccessToken[token] = expiresAt
}

func registerIssuedRefreshToken(token string, expiresAt time.Time) {
	oauthStoreMu.Lock()
	defer oauthStoreMu.Unlock()
	issuedRefreshToken[token] = expiresAt
}

func isIssuedAccessToken(token string) bool {
	if token == staticDevPAT {
		return true
	}

	oauthStoreMu.Lock()
	defer oauthStoreMu.Unlock()

	expiresAt, ok := issuedAccessToken[token]
	if !ok {
		return false
	}
	if time.Now().After(expiresAt) {
		delete(issuedAccessToken, token)
		return false
	}
	return true
}

func isIssuedRefreshToken(token string) bool {
	oauthStoreMu.Lock()
	defer oauthStoreMu.Unlock()
	expiresAt, ok := issuedRefreshToken[token]
	if !ok {
		return false
	}
	if time.Now().After(expiresAt) {
		delete(issuedRefreshToken, token)
		return false
	}
	return true
}

func buildRedirectWithQuery(base string, params map[string]string) string {
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	q := u.Query()
	for k, v := range params {
		if v == "" {
			continue
		}
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func writeOAuthError(w http.ResponseWriter, status int, code, description string) {
	writeJSON(w, status, map[string]string{
		"error":             code,
		"error_description": description,
	})
}

const consentPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Authorize — Canvas</title>
<style>
  * { margin: 0; padding: 0; box-sizing: border-box; }
  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: #f5f5f5;
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  .container {
    background: #fff;
    border-radius: 8px;
    box-shadow: 0 2px 12px rgba(0,0,0,0.1);
    max-width: 420px;
    width: 100%%;
    overflow: hidden;
  }
  .header {
    background: #394B58;
    color: #fff;
    padding: 20px 28px;
    font-size: 18px;
    font-weight: 600;
  }
  .body {
    padding: 28px;
  }
  .app-name {
    font-size: 16px;
    font-weight: 600;
    color: #2D3B45;
    margin-bottom: 8px;
  }
  .desc {
    font-size: 14px;
    color: #556572;
    margin-bottom: 20px;
    line-height: 1.5;
  }
  .permissions {
    background: #f9fafb;
    border: 1px solid #e8e8e8;
    border-radius: 6px;
    padding: 16px;
    margin-bottom: 24px;
  }
  .permissions h3 {
    font-size: 13px;
    font-weight: 600;
    color: #2D3B45;
    margin-bottom: 10px;
  }
  .permissions ul {
    list-style: none;
    padding: 0;
  }
  .permissions li {
    font-size: 13px;
    color: #556572;
    padding: 4px 0;
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .permissions li::before {
    content: '\2713';
    color: #0B874B;
    font-weight: 700;
  }
  .actions {
    display: flex;
    gap: 12px;
  }
  .btn {
    flex: 1;
    padding: 10px 16px;
    font-size: 14px;
    font-weight: 600;
    border-radius: 6px;
    cursor: pointer;
    text-align: center;
    text-decoration: none;
    display: inline-block;
    border: none;
  }
  .btn-authorize {
    background: #0B874B;
    color: #fff;
  }
  .btn-authorize:hover { background: #0a7a43; }
  .btn-cancel {
    background: #fff;
    color: #556572;
    border: 1px solid #ccc;
  }
  .btn-cancel:hover { background: #f5f5f5; }
</style>
</head>
<body>
<div class="container">
  <div class="header">Canvas</div>
  <div class="body">
    <div class="app-name">Synapse is requesting access to your account</div>
    <div class="desc">
      Synapse would like to access your Canvas data to sync your courses,
      assignments, and grades.
    </div>
    <div class="permissions">
      <h3>This application will be able to:</h3>
      <ul>
        <li>Access your course enrollments</li>
        <li>View your assignments and deadlines</li>
        <li>Read your grades and submissions</li>
      </ul>
    </div>
    <div class="actions">
      <a href="%s" class="btn btn-authorize">Authorize</a>
      <a href="%s" class="btn btn-cancel">Cancel</a>
    </div>
  </div>
</div>
</body>
</html>`
