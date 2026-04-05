package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/http/middleware"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/respond"
)

// AutopilotHandler serves stub autopilot endpoints until the learning service exists.
type AutopilotHandler struct {
	DB *pgxpool.Pool
}

type nbaAction struct {
	Icon       string `json:"icon"`
	Title      string `json:"title"`
	Reason     string `json:"reason"`
	Duration   string `json:"duration"`
	ButtonText string `json:"button_text"`
	Route      string `json:"route"`
}

type contractPayload struct {
	CourseName          string  `json:"course_name"`
	ExamDate            string  `json:"exam_date"`
	DaysUntil           int     `json:"days_until"`
	Status              string  `json:"status"`
	WeeklyHoursBudget   float64 `json:"weekly_hours_budget"`
	HoursDone           float64 `json:"hours_done"`
	Readiness           int     `json:"readiness"`
}

type weakTopicPayload struct {
	Name    string `json:"name"`
	Mastery int    `json:"mastery"`
	Bars    []int  `json:"bars"`
}

type statPayload struct {
	Value string `json:"value"`
	Label string `json:"label"`
	Color string `json:"color"`
}

type deadlineAlertPayload struct {
	Course string `json:"course"`
	Title  string `json:"title"`
	Days   int    `json:"days"`
	Type   string `json:"type"`
}

type todayPayload struct {
	GreetingName  string                `json:"greeting_name"`
	Actions       []nbaAction           `json:"actions"`
	Contract      contractPayload       `json:"contract"`
	WeakTopics    []weakTopicPayload    `json:"weak_topics"`
	Stats         []statPayload         `json:"stats"`
	Streak        int                   `json:"streak"`
	DeadlineAlert deadlineAlertPayload `json:"deadline_alert"`
}

func (h *AutopilotHandler) authIDs(r *http.Request) (userID, schoolID uuid.UUID, ok bool) {
	uid, ok := middleware.UserIDFromCtx(r.Context())
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	sid, ok := middleware.SchoolIDFromCtx(r.Context())
	if !ok {
		return uuid.Nil, uuid.Nil, false
	}
	return uid, sid, true
}

// Today returns API-backed autopilot data for the current user.
func (h *AutopilotHandler) Today(w http.ResponseWriter, r *http.Request) {
	userID, schoolID, ok := h.authIDs(r)
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	greetingName := "Student"
	_ = h.DB.QueryRow(r.Context(), `SELECT name FROM users WHERE id = $1 AND school_id = $2`, userID, schoolID).Scan(&greetingName)
	firstName := strings.Split(strings.TrimSpace(greetingName), " ")[0]
	if firstName == "" {
		firstName = "Student"
	}

	var courseName string
	_ = h.DB.QueryRow(
		r.Context(),
		`SELECT name FROM courses WHERE user_id = $1 AND school_id = $2 ORDER BY created_at DESC LIMIT 1`,
		userID, schoolID,
	).Scan(&courseName)
	if courseName == "" {
		courseName = "No active course"
	}

	var dueCount int
	_ = h.DB.QueryRow(
		r.Context(),
		`SELECT COUNT(*)
		 FROM scheduler_states s
		 JOIN flashcards f ON f.id = s.flashcard_id
		 WHERE s.user_id = $1 AND s.school_id = $2
		   AND f.user_id = $1 AND f.school_id = $2
		   AND s.due_at <= now()`,
		userID, schoolID,
	).Scan(&dueCount)

	var cardsDone int
	var totalMs int64
	var accuracy float64
	_ = h.DB.QueryRow(
		r.Context(),
		`SELECT
		   COUNT(*) FILTER (WHERE ts >= now() - interval '7 days'),
		   COALESCE(SUM(response_time_ms) FILTER (WHERE ts >= now() - interval '7 days'), 0),
		   COALESCE(AVG(CASE WHEN correct THEN 1.0 ELSE 0.0 END) FILTER (WHERE ts >= now() - interval '7 days'), 0)
		 FROM review_events
		 WHERE user_id = $1 AND school_id = $2`,
		userID, schoolID,
	).Scan(&cardsDone, &totalMs, &accuracy)

	accuracyPct := int(accuracy * 100)
	studyHours := float64(totalMs) / 3600000.0

	weakTopics := make([]weakTopicPayload, 0, 3)
	rows, err := h.DB.Query(
		r.Context(),
		`SELECT COALESCE(t.name, 'Unscoped'),
		        LEAST(95, GREATEST(15, 95 - (COUNT(*) * 9)::int)) AS mastery
		 FROM review_events re
		 JOIN flashcards f ON f.id = re.flashcard_id AND f.user_id = $1 AND f.school_id = $2
		 LEFT JOIN topics t ON t.id = f.topic_id AND t.school_id = $2
		 WHERE re.user_id = $1 AND re.school_id = $2
		   AND re.confused = true
		   AND re.ts >= now() - interval '14 days'
		 GROUP BY COALESCE(t.name, 'Unscoped')
		 ORDER BY COUNT(*) DESC
		 LIMIT 3`,
		userID, schoolID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var name string
			var mastery int
			if rows.Scan(&name, &mastery) == nil {
				rawBars := []int{mastery + 16, mastery + 7, mastery, mastery - 5, mastery - 9}
				bars := make([]int, len(rawBars))
				for i, v := range rawBars {
					if v < 0 {
						v = 0
					}
					if v > 100 {
						v = 100
					}
					bars[i] = v
				}
				weakTopics = append(weakTopics, weakTopicPayload{
					Name:    name,
					Mastery: mastery,
					Bars:    bars,
				})
			}
		}
	}
	topWeakTopic := ""
	if len(weakTopics) > 0 {
		topWeakTopic = weakTopics[0].Name
	}

	readiness := 0
	if len(weakTopics) > 0 {
		hotspots := len(weakTopics)
		readiness = 85 - hotspots*7
		if readiness < 40 {
			readiness = 40
		}
	}

	streak := h.computeStreak(r.Context(), userID, schoolID)

	var (
		deadlineCourse string
		deadlineTitle  string
		deadlineAt     *time.Time
	)
	_ = h.DB.QueryRow(
		r.Context(),
		`SELECT COALESCE(lc.lms_course_name, ''), la.title, la.due_at
		 FROM lms_assignments la
		 JOIN lms_courses lc
		   ON lc.lms_course_id = la.lms_course_id
		  AND lc.user_id = $1
		  AND lc.school_id = $2
		 WHERE la.school_id = $2
		   AND la.due_at IS NOT NULL
		   AND la.due_at >= now()
		 ORDER BY la.due_at ASC
		 LIMIT 1`,
		userID, schoolID,
	).Scan(&deadlineCourse, &deadlineTitle, &deadlineAt)

	var deadline deadlineAlertPayload
	if deadlineAt != nil {
		days := int(time.Until(*deadlineAt).Hours() / 24)
		if days < 0 {
			days = 0
		}
		deadline = deadlineAlertPayload{
			Course: deadlineCourse,
			Title:  deadlineTitle,
			Days:   days,
			Type:   "assignment",
		}
	}

	var actions []nbaAction
	if dueCount > 0 {
		actions = append(actions, nbaAction{
			Icon:       "clock",
			Title:      "Review due cards",
			Reason:     "Keep your recall fresh with today's due queue",
			Duration:   "~12 min",
			ButtonText: "Start Review",
			Route:      "/review",
		})
	} else {
		actions = append(actions, nbaAction{
			Icon:       "zap",
			Title:      "Start your review habit",
			Reason:     "Add flashcards or connect Canvas so Synapse can queue spaced repetition for you.",
			Duration:   "—",
			ButtonText: "Open Review",
			Route:      "/review",
		})
	}
	if len(weakTopics) > 0 && topWeakTopic != "" {
		actions = append(actions,
			nbaAction{
				Icon:       "book-open",
				Title:      "Revisit weak topic",
				Reason:     "Recent confusion concentrated in " + topWeakTopic,
				Duration:   "~10 min",
				ButtonText: "Open Notes",
				Route:      "/notes",
			},
			nbaAction{
				Icon:       "users",
				Title:      "Ask for tutor help",
				Reason:     "Peers who list " + topWeakTopic + " can help—match on mastery in Tutoring.",
				Duration:   "~15 min",
				ButtonText: "Find Tutors",
				Route:      "/tutoring",
			},
		)
	}

	weeklyBudget := 8.0
	contractStatus := "on_track"
	if courseName == "No active course" || courseName == "" {
		weeklyBudget = 0
	}
	if weeklyBudget == 0 && readiness == 0 {
		contractStatus = "no_data"
	}

	data := todayPayload{
		GreetingName: firstName,
		Actions:      actions,
		Contract: contractPayload{
			CourseName:        courseName,
			ExamDate:          "TBD",
			DaysUntil:         0,
			Status:            contractStatus,
			WeeklyHoursBudget: weeklyBudget,
			HoursDone:         studyHours,
			Readiness:         readiness,
		},
		WeakTopics: weakTopics,
		Stats: []statPayload{
			{Value: itoa(cardsDone), Label: "Cards done", Color: "var(--navy)"},
			{Value: formatHours(studyHours), Label: "Study time", Color: "var(--navy)"},
			{Value: itoa(accuracyPct) + "%", Label: "Accuracy", Color: "var(--emerald)"},
			{Value: itoa(dueCount), Label: "Due", Color: "var(--emerald)"},
		},
		Streak:        streak,
		DeadlineAlert: deadline,
	}
	respond.JSON(w, http.StatusOK, data)
}

func (h *AutopilotHandler) computeStreak(ctx context.Context, userID, schoolID uuid.UUID) int {
	rows, err := h.DB.Query(
		ctx,
		`SELECT DISTINCT DATE(ts)
		 FROM review_events
		 WHERE user_id = $1 AND school_id = $2
		 ORDER BY DATE(ts) DESC
		 LIMIT 30`,
		userID, schoolID,
	)
	if err != nil {
		return 0
	}
	defer rows.Close()

	dates := make([]time.Time, 0, 30)
	for rows.Next() {
		var d time.Time
		if rows.Scan(&d) == nil {
			dates = append(dates, d)
		}
	}
	if len(dates) == 0 {
		return 0
	}

	streak := 0
	expect := time.Now().UTC().Truncate(24 * time.Hour)
	for _, d := range dates {
		day := d.UTC().Truncate(24 * time.Hour)
		if day.Equal(expect) {
			streak++
			expect = expect.Add(-24 * time.Hour)
			continue
		}
		if streak == 0 && day.Equal(expect.Add(-24*time.Hour)) {
			// tolerate a day rollover near midnight.
			expect = expect.Add(-24 * time.Hour)
			streak++
			expect = expect.Add(-24 * time.Hour)
			continue
		}
		break
	}
	return streak
}

func itoa(v int) string {
	return strconv.Itoa(v)
}

func formatHours(h float64) string {
	if h < 0 {
		h = 0
	}
	return fmt.Sprintf("%.1fh", h)
}
