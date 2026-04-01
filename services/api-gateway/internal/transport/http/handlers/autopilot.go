package handlers

import (
	"net/http"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/transport/respond"
)

// AutopilotHandler serves stub autopilot endpoints until the learning service exists.
type AutopilotHandler struct{}

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

// Today returns stub data matching the legacy today-tab hardcoded UI.
// TODO: replace with real NBA engine in learning service phase
func (h *AutopilotHandler) Today(w http.ResponseWriter, r *http.Request) {
	_ = r // RequireAuth applied by router
	data := todayPayload{
		GreetingName: "Alex",
		Actions: []nbaAction{
			{Icon: "clock", Title: "Review 12 due cards", Reason: "8 cards overdue · forgetting risk rising on Induction", Duration: "~14 min", ButtonText: "Start Review", Route: "/review"},
			{Icon: "book-open", Title: `Relearn "Induction step" section`, Reason: "Confusion hotspot · 3 confused marks in last session", Duration: "~8 min", ButtonText: "Open Section", Route: "/notes"},
			{Icon: "users", Title: "Request tutor on Recursion base cases", Reason: "Confusion plateau detected · 2 peers available now", Duration: "~15 min", ButtonText: "Find Tutors", Route: "/tutoring"},
		},
		Contract: contractPayload{
			CourseName:        "CS225 · Discrete Mathematics",
			ExamDate:          "Mar 12",
			DaysUntil:         16,
			Status:            "on_track",
			WeeklyHoursBudget: 8,
			HoursDone:         4.5,
			Readiness:         73,
		},
		WeakTopics: []weakTopicPayload{
			{Name: "Recursion", Mastery: 29, Bars: []int{45, 32, 29, 25, 20}},
			{Name: "Induction", Mastery: 51, Bars: []int{62, 55, 51, 48, 45}},
			{Name: "Set Theory", Mastery: 78, Bars: []int{88, 82, 78, 76, 74}},
		},
		Stats: []statPayload{
			{Value: "47", Label: "Cards done", Color: "var(--navy)"},
			{Value: "3.2h", Label: "Study time", Color: "var(--navy)"},
			{Value: "82%", Label: "Accuracy", Color: "var(--emerald)"},
			{Value: "0 Due", Label: "Due", Color: "var(--emerald)"},
		},
		Streak: 7,
		DeadlineAlert: deadlineAlertPayload{
			Course: "CS225",
			Title:  "Problem Set 3",
			Days:   2,
			Type:   "assignment",
		},
	}
	respond.JSON(w, http.StatusOK, data)
}
