#!/usr/bin/env bash
# Git commit: today Due stat, sidebar school line, profile auth wiring.
# Run from repo root: bash scripts/4.sh

set -e

cd "$(dirname "$0")/.."

echo "=== Commit: Today stat, sidebar school, profile user data ==="
git add \
  apps/web/src/app/components/today-tab.component.ts \
  apps/web/src/app/components/sidebar.component.ts \
  apps/web/src/app/components/profile-tab.component.ts
git commit -m "$(cat <<'EOF'
fix(web): today Due stat, sidebar school line, profile from auth

Show numeric Due count only in This Week stats (label already says Due).
Sidebar subline uses University when user is loaded, empty while loading.
Profile hero uses AuthService for name, initials, and University subline.
EOF
)"
