package domain

import (
	"context"

	"github.com/google/uuid"
)

// ProfileSummary is returned by GET /v1/profile/summary.
type ProfileSummary struct {
	ReputationTotal int                      `json:"reputation_total"`
	CardsShared     int                      `json:"cards_shared"`
	SessionsDone    int                      `json:"sessions_done"`
	Breakdown       []ProfileRepBreakdownRow `json:"reputation_breakdown"`
	Mastery         []ProfileMasteryRow      `json:"mastery"`
	Contributions   []ProfileContributionRow `json:"contributions"`
}

// ProfileRepBreakdownRow drives reputation bars (values are raw counts / points).
type ProfileRepBreakdownRow struct {
	Label   string `json:"label"`
	Value   int    `json:"value"`
	Percent int    `json:"percent"`
	Color   string `json:"color"`
}

// ProfileMasteryRow is one topic row for the mastery list.
type ProfileMasteryRow struct {
	Name    string `json:"name"`
	Mastery int    `json:"mastery"`
	Band    string `json:"band"` // green | amber | red
}

// ProfileContributionRow is one line in the contributions list.
type ProfileContributionRow struct {
	Icon       string `json:"icon"`
	IconColor  string `json:"icon_color"`
	IconBg     string `json:"icon_bg"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Value      string `json:"value"`
	ValueColor string `json:"value_color"`
}

// ProfileRepository loads aggregated profile data for a user.
type ProfileRepository interface {
	GetSummary(ctx context.Context, userID, schoolID uuid.UUID) (*ProfileSummary, error)
}
