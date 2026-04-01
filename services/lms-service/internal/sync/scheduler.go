package sync

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/shamshad-ansari/synapse/services/lms-service/internal/domain"
)

// StartScheduler runs a periodic background sync loop for all LMS connections.
func StartScheduler(ctx context.Context, repo domain.LMSRepository, syncer *Syncer, logger *zap.Logger, interval, startupDelay time.Duration, opts SyncOptions) {
	if interval <= 0 {
		interval = 15 * time.Minute
	}
	if startupDelay < 0 {
		startupDelay = 0
	}

	go func() {
		if startupDelay > 0 {
			timer := time.NewTimer(startupDelay)
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
			}
		}

		runSyncCycle(ctx, repo, syncer, logger, opts)

		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runSyncCycle(ctx, repo, syncer, logger, opts)
			}
		}
	}()
}

func runSyncCycle(ctx context.Context, repo domain.LMSRepository, syncer *Syncer, logger *zap.Logger, opts SyncOptions) {
	connections, err := repo.ListConnectionsForSync(ctx)
	if err != nil {
		logger.Error("scheduled sync: list connections failed", zap.Error(err))
		return
	}
	if len(connections) == 0 {
		return
	}

	for _, conn := range connections {
		if conn == nil {
			continue
		}

		_ = repo.UpdateConnectionSyncStatus(ctx, conn.UserID, conn.SchoolID, "syncing", conn.LastSyncedAt)

		runCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
		err := syncer.SyncUser(runCtx, conn, opts)
		cancel()
		if err != nil {
			logger.Warn(
				"scheduled sync failed",
				zap.String("user_id", conn.UserID.String()),
				zap.String("school_id", conn.SchoolID.String()),
				zap.Error(err),
			)
			_ = repo.UpdateConnectionSyncStatus(ctx, conn.UserID, conn.SchoolID, "error", conn.LastSyncedAt)
			continue
		}
	}
}
