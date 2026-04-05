package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

type PostgresProfileRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresProfileRepo(pool *pgxpool.Pool) *PostgresProfileRepo {
	return &PostgresProfileRepo{pool: pool}
}

// GetSummary aggregates profile metrics for the authenticated user.
// ReputationTotal = feed_upvotes_sum + 3*cards_shared + 5*tutoring_completed + 2*feed_comments + 2*sessions_done
// (documented formula — adjust weights in one place if product tuning is needed).
func (r *PostgresProfileRepo) GetSummary(ctx context.Context, userID, schoolID uuid.UUID) (*domain.ProfileSummary, error) {
	var cardsShared, sessionsDone, feedUpvotesSum, tutoringDone, feedComments int

	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*)::int FROM flashcards
		 WHERE user_id = $1 AND school_id = $2 AND visibility <> 'private'`,
		userID, schoolID,
	).Scan(&cardsShared)
	if err != nil {
		return nil, fmt.Errorf("GetSummary cards_shared: %w", err)
	}

	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(DISTINCT session_id)::int FROM review_events
		 WHERE user_id = $1 AND school_id = $2`,
		userID, schoolID,
	).Scan(&sessionsDone)
	if err != nil {
		return nil, fmt.Errorf("GetSummary sessions_done: %w", err)
	}

	err = r.pool.QueryRow(ctx,
		`SELECT COALESCE(SUM(upvotes), 0)::int FROM feed_posts
		 WHERE user_id = $1 AND school_id = $2`,
		userID, schoolID,
	).Scan(&feedUpvotesSum)
	if err != nil {
		return nil, fmt.Errorf("GetSummary feed_upvotes: %w", err)
	}

	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*)::int FROM tutor_requests
		 WHERE school_id = $1 AND status = 'completed' AND (requester_id = $2 OR tutor_id = $2)`,
		schoolID, userID,
	).Scan(&tutoringDone)
	if err != nil {
		return nil, fmt.Errorf("GetSummary tutoring: %w", err)
	}

	err = r.pool.QueryRow(ctx,
		`SELECT COUNT(*)::int FROM feed_comments
		 WHERE user_id = $1 AND school_id = $2`,
		userID, schoolID,
	).Scan(&feedComments)
	if err != nil {
		return nil, fmt.Errorf("GetSummary feed_comments: %w", err)
	}

	reputationTotal := feedUpvotesSum + 3*cardsShared + 5*tutoringDone + 2*feedComments + 2*sessionsDone

	breakdownVals := []int{cardsShared, feedUpvotesSum, tutoringDone, feedComments}
	breakdownLabels := []string{"Flashcards shared", "Upvotes on posts", "Tutoring completed", "Feed comments"}
	breakdownColors := []string{"var(--navy)", "var(--navy)", "var(--emerald)", "var(--emerald)"}
	sumBD := 0
	for _, v := range breakdownVals {
		sumBD += v
	}
	breakdown := make([]domain.ProfileRepBreakdownRow, 0, len(breakdownVals))
	for i, v := range breakdownVals {
		pct := 0
		if sumBD > 0 {
			pct = (100 * v) / sumBD
		}
		breakdown = append(breakdown, domain.ProfileRepBreakdownRow{
			Label:   breakdownLabels[i],
			Value:   v,
			Percent: pct,
			Color:   breakdownColors[i],
		})
	}

	mastery, err := r.listMastery(ctx, userID, schoolID)
	if err != nil {
		return nil, err
	}

	contributions := []domain.ProfileContributionRow{
		{
			Icon: "zap", IconColor: "var(--navy)", IconBg: "var(--navy-light)",
			Title: "Flashcards shared", Subtitle: "Visible to school or public",
			Value: fmt.Sprintf("%d", cardsShared), ValueColor: "var(--navy)",
		},
		{
			Icon: "message-circle", IconColor: "var(--navy)", IconBg: "var(--navy-light)",
			Title: "Upvotes on your posts", Subtitle: "Total on your feed posts",
			Value: fmt.Sprintf("%d", feedUpvotesSum), ValueColor: "var(--navy)",
		},
		{
			Icon: "users", IconColor: "var(--emerald)", IconBg: "var(--emerald-light)",
			Title: "Tutoring sessions", Subtitle: "Completed requests you took part in",
			Value: fmt.Sprintf("%d", tutoringDone), ValueColor: "var(--emerald)",
		},
		{
			Icon: "shield-check", IconColor: "var(--emerald)", IconBg: "var(--emerald-light)",
			Title: "Feed comments", Subtitle: "Replies you posted",
			Value: fmt.Sprintf("%d", feedComments), ValueColor: "var(--emerald)",
		},
	}

	return &domain.ProfileSummary{
		ReputationTotal: reputationTotal,
		CardsShared:     cardsShared,
		SessionsDone:    sessionsDone,
		Breakdown:       breakdown,
		Mastery:         mastery,
		Contributions:   contributions,
	}, nil
}

func (r *PostgresProfileRepo) listMastery(ctx context.Context, userID, schoolID uuid.UUID) ([]domain.ProfileMasteryRow, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT t.name,
		        LEAST(100, GREATEST(0, ROUND((tm.mastery_score)::numeric * 100)::int)) AS pct
		 FROM topic_mastery tm
		 INNER JOIN topics t ON t.id = tm.topic_id AND t.school_id = tm.school_id
		 WHERE tm.user_id = $1 AND tm.school_id = $2
		 ORDER BY tm.mastery_score DESC
		 LIMIT 8`,
		userID, schoolID,
	)
	if err != nil {
		return nil, fmt.Errorf("listMastery: %w", err)
	}
	defer rows.Close()

	var out []domain.ProfileMasteryRow
	for rows.Next() {
		var name string
		var pct int
		if err := rows.Scan(&name, &pct); err != nil {
			return nil, fmt.Errorf("listMastery row: %w", err)
		}
		band := "red"
		if pct >= 70 {
			band = "green"
		} else if pct >= 40 {
			band = "amber"
		}
		out = append(out, domain.ProfileMasteryRow{Name: name, Mastery: pct, Band: band})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listMastery: %w", err)
	}
	return out, nil
}
