package seed

import "time"

// ───────────────────────────────────────────────────────────────────
// Canvas-format response structs with full realistic fields.
// ───────────────────────────────────────────────────────────────────

type User struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	SortableName string  `json:"sortable_name"`
	FirstName    string  `json:"first_name"`
	LastName     string  `json:"last_name"`
	ShortName    string  `json:"short_name"`
	LoginID      string  `json:"login_id"`
	Email        string  `json:"email"`
	AvatarURL    string  `json:"avatar_url"`
	Locale       string  `json:"locale"`
	TimeZone     string  `json:"time_zone"`
	Bio          *string `json:"bio"`
}

type Course struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	CourseCode       string     `json:"course_code"`
	WorkflowState    string     `json:"workflow_state"`
	AccountID        int        `json:"account_id"`
	EnrollmentTermID int        `json:"enrollment_term_id"`
	StartAt          *time.Time `json:"start_at"`
	EndAt            *time.Time `json:"end_at"`
	CreatedAt        time.Time  `json:"created_at"`
	DefaultView      string     `json:"default_view"`
	SyllabusBody     string     `json:"syllabus_body,omitempty"`
	IsPublic         bool       `json:"is_public"`
	TimeZone         string     `json:"time_zone"`
	UUID             string     `json:"uuid"`
	HideFinalGrades  bool       `json:"hide_final_grades"`
	CourseFormat     string     `json:"course_format"`
	TotalStudents    int        `json:"total_students"`
}

type Assignment struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	DueAt             *time.Time `json:"due_at"`
	LockAt            *time.Time `json:"lock_at"`
	UnlockAt          *time.Time `json:"unlock_at"`
	CourseID          int        `json:"course_id"`
	PointsPossible    float64    `json:"points_possible"`
	GradingType       string     `json:"grading_type"`
	AssignmentGroupID int        `json:"assignment_group_id"`
	Position          int        `json:"position"`
	SubmissionTypes   []string   `json:"submission_types"`
	Published         bool       `json:"published"`
	HasSubmitted      bool       `json:"has_submitted_submissions"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	HTMLURL           string     `json:"html_url"`
}

type Submission struct {
	ID            int        `json:"id"`
	AssignmentID  int        `json:"assignment_id"`
	UserID        int        `json:"user_id"`
	Score         *float64   `json:"score"`
	Grade         string     `json:"grade"`
	EnteredGrade  string     `json:"entered_grade"`
	EnteredScore  *float64   `json:"entered_score"`
	SubmittedAt   *time.Time `json:"submitted_at"`
	GradedAt      *time.Time `json:"graded_at"`
	Attempt       int        `json:"attempt"`
	WorkflowState string     `json:"workflow_state"`
	Late          bool       `json:"late"`
	Missing       bool       `json:"missing"`
	Excused       bool       `json:"excused"`
}

type Announcement struct {
	ID       int        `json:"id"`
	Title    string     `json:"title"`
	Message  string     `json:"message"`
	PostedAt *time.Time `json:"posted_at"`
	HTMLURL  string     `json:"html_url"`
}

type DiscussionTopic struct {
	ID       int        `json:"id"`
	Title    string     `json:"title"`
	Message  string     `json:"message"`
	PostedAt *time.Time `json:"posted_at"`
	HTMLURL  string     `json:"html_url"`
}

// Canvas Planner Items — the primary API developers use for planner apps
type PlannerSubmission struct {
	Submitted    bool `json:"submitted"`
	Excused      bool `json:"excused"`
	Graded       bool `json:"graded"`
	Late         bool `json:"late"`
	Missing      bool `json:"missing"`
	NeedsGrading bool `json:"needs_grading"`
	HasFeedback  bool `json:"has_feedback"`
}

type Plannable struct {
	ID             int        `json:"id"`
	Title          string     `json:"title"`
	DueAt          *time.Time `json:"due_at"`
	PointsPossible float64    `json:"points_possible"`
	CourseID       int        `json:"course_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type PlannerItem struct {
	ContextType   string             `json:"context_type"`
	CourseID      int                `json:"course_id,omitempty"`
	PlannableDate time.Time          `json:"plannable_date"`
	Plannable     Plannable          `json:"plannable"`
	PlannableType string             `json:"plannable_type"`
	NewActivity   bool               `json:"new_activity"`
	Submissions   *PlannerSubmission `json:"submissions"`
	HTMLURL       string             `json:"html_url"`
	ContextName   string             `json:"context_name"`
}

type CalendarEvent struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	StartAt       *time.Time `json:"start_at"`
	EndAt         *time.Time `json:"end_at"`
	LocationName  string     `json:"location_name"`
	ContextCode   string     `json:"context_code"`
	WorkflowState string     `json:"workflow_state"`
	AllDay        bool       `json:"all_day"`
	Description   string     `json:"description,omitempty"`
}

func ptr[T any](v T) *T { return &v }

// now returns the current reference time (today at midnight UTC).
func now() time.Time {
	t := time.Now().UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// day returns now() + offset days at the given hour:minute.
func day(offset, hour, min int) time.Time {
	n := now()
	return n.Add(time.Duration(offset)*24*time.Hour +
		time.Duration(hour)*time.Hour +
		time.Duration(min)*time.Minute)
}

func dayp(offset, hour, min int) *time.Time {
	v := day(offset, hour, min)
	return &v
}

// -------------------------------------------------------------------
// Seed data — time-relative so data always looks current
// -------------------------------------------------------------------

var CurrentUser = User{
	ID:           5001,
	Name:         "Alex Kim",
	SortableName: "Kim, Alex",
	FirstName:    "Alex",
	LastName:     "Kim",
	ShortName:    "Alex",
	LoginID:      "alex@mit.edu",
	Email:        "alex@mit.edu",
	AvatarURL:    "https://i.pravatar.cc/150?u=alexkim",
	Locale:       "en",
	TimeZone:     "America/New_York",
	Bio:          nil,
}

// Semester runs from ~8 weeks ago to ~8 weeks from now.
func semesterStart() *time.Time { v := day(-56, 0, 0); return &v }
func semesterEnd() *time.Time   { v := day(56, 23, 59); return &v }

// Courses uses init() to populate with time-relative dates.
var Courses []Course
var CoursesByID map[int]Course

func init() {
	ss := semesterStart()
	se := semesterEnd()
	created := day(-60, 0, 0)

	Courses = []Course{
		{
			ID: 1001, Name: "African History", CourseCode: "HIS-180-ONL.2026SP",
			WorkflowState: "available", AccountID: 1, EnrollmentTermID: 202601,
			StartAt: ss, EndAt: se, CreatedAt: created,
			DefaultView: "feed", IsPublic: false, TimeZone: "America/New_York",
			UUID: "his-180-uuid", HideFinalGrades: false, TotalStudents: 40,
		},
		{
			ID: 1002, Name: "Algorithms & Systems", CourseCode: "CSCI-390ASD-01.2026SP",
			WorkflowState: "available", AccountID: 1, EnrollmentTermID: 202601,
			StartAt: ss, EndAt: se, CreatedAt: created,
			DefaultView: "feed", IsPublic: false, TimeZone: "America/New_York",
			UUID: "csci-390asd-uuid", HideFinalGrades: false, TotalStudents: 60,
		},
		{
			ID: 1003, Name: "Careers in Tech", CourseCode: "CSCI-390CIT-01.2026SP",
			WorkflowState: "available", AccountID: 1, EnrollmentTermID: 202601,
			StartAt: ss, EndAt: se, CreatedAt: created,
			DefaultView: "feed", IsPublic: false, TimeZone: "America/New_York",
			UUID: "csci-390cit-uuid", HideFinalGrades: false, TotalStudents: 50,
		},
		{
			ID: 1004, Name: "Composition II", CourseCode: "CORE-160-07.2026SP",
			WorkflowState: "available", AccountID: 1, EnrollmentTermID: 202601,
			StartAt: ss, EndAt: se, CreatedAt: created,
			DefaultView: "feed", IsPublic: false, TimeZone: "America/New_York",
			UUID: "core-160-uuid", HideFinalGrades: false, TotalStudents: 25,
		},
		{
			ID: 1005, Name: "Intro to Software Engineering", CourseCode: "CSCI-390ISE-01.2026SP",
			WorkflowState: "available", AccountID: 1, EnrollmentTermID: 202601,
			StartAt: ss, EndAt: se, CreatedAt: created,
			DefaultView: "feed", IsPublic: false, TimeZone: "America/New_York",
			UUID: "csci-390ise-uuid", HideFinalGrades: false, TotalStudents: 55,
		},
		{
			ID: 1006, Name: "Junior Seminar", CourseCode: "CSCI-310-01.2026SP",
			WorkflowState: "available", AccountID: 1, EnrollmentTermID: 202601,
			StartAt: ss, EndAt: se, CreatedAt: created,
			DefaultView: "feed", IsPublic: false, TimeZone: "America/New_York",
			UUID: "csci-310-uuid", HideFinalGrades: false, TotalStudents: 30,
		},
		{
			ID: 1007, Name: "Machine Learning", CourseCode: "CSCI-380-01.2026SP",
			WorkflowState: "available", AccountID: 1, EnrollmentTermID: 202601,
			StartAt: ss, EndAt: se, CreatedAt: created,
			DefaultView: "feed", IsPublic: false, TimeZone: "America/New_York",
			UUID: "csci-380-uuid", HideFinalGrades: false, TotalStudents: 70,
		},
	}

	CoursesByID = make(map[int]Course, len(Courses))
	for _, c := range Courses {
		CoursesByID[c.ID] = c
	}
}

var Assignments = map[int][]Assignment{
	1001: {
		{
			ID: 2001, Name: "Historical Connections: Brittania", CourseID: 1001,
			Description:    "Midterm reading coverage.",
			DueAt:          dayp(-5, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 100,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: day(-49, 9, 0), UpdatedAt: day(-49, 9, 0),
		},
		{
			ID: 2002, Name: "Historical Connections: Southern Africa", CourseID: 1001,
			Description:    "Historical comparison.",
			DueAt:          dayp(-2, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 100,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: day(-35, 9, 0), UpdatedAt: day(-35, 9, 0),
		},
		{
			ID: 2003, Name: "Historical Connection Summary", CourseID: 1001,
			Description:    "Short summary.",
			DueAt:          dayp(0, 23, 30),
			PointsPossible: 50, GradingType: "points", AssignmentGroupID: 100,
			Position: 3, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: day(-21, 9, 0), UpdatedAt: day(-21, 9, 0),
		},
		{
			ID: 2004, Name: "Historical Connections Paper", CourseID: 1001,
			Description:    "End of month paper.",
			DueAt:          dayp(2, 23, 59),
			PointsPossible: 200, GradingType: "points", AssignmentGroupID: 100,
			Position: 4, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: day(-7, 9, 0), UpdatedAt: day(-7, 9, 0),
		},
	},
	1002: {
		{
			ID: 2011, Name: "Homework 06", CourseID: 1002,
			Description:    "Dynamic programming and arrays.",
			DueAt:          dayp(-25, 12, 0),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 200,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: day(-42, 9, 0), UpdatedAt: day(-42, 9, 0),
		},
		{
			ID: 2012, Name: "Homework 07", CourseID: 1002,
			Description:    "Graph representations.",
			DueAt:          dayp(-17, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 200,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: day(-28, 9, 0), UpdatedAt: day(-28, 9, 0),
		},
		{
			ID: 2013, Name: "Homework 10", CourseID: 1002,
			Description:    "Trees and balances.",
			DueAt:          dayp(4, 12, 0),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 200,
			Position: 3, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: day(-7, 9, 0), UpdatedAt: day(-7, 9, 0),
		},
	},
	1004: {
		{
			ID: 2021, Name: "Paper 2", CourseID: 1004,
			Description:    "Comparison of rhetoric.",
			DueAt:          dayp(-22, 23, 59),
			PointsPossible: 150, GradingType: "points", AssignmentGroupID: 300,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: day(-44, 9, 0), UpdatedAt: day(-44, 9, 0),
		},
	},
	1006: {
		{
			ID: 2031, Name: "Progress Report 1", CourseID: 1006,
			Description:    "Seminar presentation draft.",
			DueAt:          dayp(-20, 23, 59),
			PointsPossible: 50, GradingType: "points", AssignmentGroupID: 400,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: day(-49, 9, 0), UpdatedAt: day(-49, 9, 0),
		},
		{
			ID: 2032, Name: "Demo", CourseID: 1006,
			Description:    "System demonstration in class.",
			DueAt:          dayp(8, 14, 0),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 401,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: day(-35, 9, 0), UpdatedAt: day(-35, 9, 0),
		},
		{
			ID: 2033, Name: "Final Creative Project Proposal", CourseID: 1006,
			Description:    "Due end of term.",
			DueAt:          dayp(11, 23, 59),
			PointsPossible: 200, GradingType: "points", AssignmentGroupID: 400,
			Position: 3, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: day(-7, 9, 0), UpdatedAt: day(-7, 9, 0),
		},
	},
	1007: {
		{
			ID: 2041, Name: "Quiz 1 - Requires Respondus", CourseID: 1007,
			Description:    "Baseline linear algebra and probability.",
			DueAt:          dayp(-29, 15, 42),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 500,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: day(-44, 9, 0), UpdatedAt: day(-44, 9, 0),
		},
		{
			ID: 2042, Name: "Midterm Exam", CourseID: 1007,
			Description:    "Models 1 through 4.",
			DueAt:          dayp(-5, 15, 59),
			PointsPossible: 200, GradingType: "points", AssignmentGroupID: 501,
			Position: 2, SubmissionTypes: []string{"online_quiz"}, Published: true, HasSubmitted: true,
			CreatedAt: day(-28, 9, 0), UpdatedAt: day(-28, 9, 0),
		},
		{
			ID: 2043, Name: "Week 7", CourseID: 1007,
			Description:    "Weekly module.",
			DueAt:          dayp(-22, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 500,
			Position: 3, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: day(-28, 9, 0), UpdatedAt: day(-28, 9, 0),
		},
		{
			ID: 2044, Name: "Week 8", CourseID: 1007,
			Description:    "Weekly module.",
			DueAt:          dayp(-15, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 500,
			Position: 4, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: day(-28, 9, 0), UpdatedAt: day(-28, 9, 0),
		},
		{
			ID: 2045, Name: "Week 10", CourseID: 1007,
			Description:    "Weekly module.",
			DueAt:          dayp(-1, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 500,
			Position: 5, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: day(-28, 9, 0), UpdatedAt: day(-28, 9, 0),
		},
		{
			ID: 2046, Name: "Week 11", CourseID: 1007,
			Description:    "Weekly module.",
			DueAt:          dayp(6, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 500,
			Position: 6, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: day(-28, 9, 0), UpdatedAt: day(-28, 9, 0),
		},
		{
			ID: 2047, Name: "Week 12", CourseID: 1007,
			Description:    "Weekly module.",
			DueAt:          dayp(13, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 500,
			Position: 7, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: day(-28, 9, 0), UpdatedAt: day(-28, 9, 0),
		},
	},
}

var Submissions = map[int][]Submission{
	1001: {
		{
			ID: 3001, AssignmentID: 2001, UserID: 5001,
			Score: ptr(90.0), Grade: "90", EnteredGrade: "90", EnteredScore: ptr(90.0),
			SubmittedAt: dayp(-6, 18, 30), GradedAt: dayp(-4, 10, 0),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3002, AssignmentID: 2002, UserID: 5001,
			Score: ptr(85.0), Grade: "85", EnteredGrade: "85", EnteredScore: ptr(85.0),
			SubmittedAt: dayp(-3, 21, 15), GradedAt: dayp(-1, 14, 0),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
	},
	1002: {
		{
			ID: 3011, AssignmentID: 2011, UserID: 5001,
			Score: ptr(88.0), Grade: "88", EnteredGrade: "88", EnteredScore: ptr(88.0),
			SubmittedAt: dayp(-26, 20, 0), GradedAt: dayp(-24, 11, 0),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3012, AssignmentID: 2012, UserID: 5001,
			Score: ptr(91.0), Grade: "91", EnteredGrade: "91", EnteredScore: ptr(91.0),
			SubmittedAt: dayp(-18, 21, 0), GradedAt: dayp(-15, 14, 0),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
	},
	1004: {
		{
			ID: 3021, AssignmentID: 2021, UserID: 5001,
			Score: ptr(95.0), Grade: "95", EnteredGrade: "95", EnteredScore: ptr(95.0),
			SubmittedAt: dayp(-23, 22, 45), GradedAt: dayp(-20, 9, 30),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
	},
	1006: {
		{
			ID: 3031, AssignmentID: 2031, UserID: 5001,
			Score: ptr(47.0), Grade: "47", EnteredGrade: "47", EnteredScore: ptr(47.0),
			SubmittedAt: dayp(-21, 22, 31), GradedAt: dayp(-19, 10, 5),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
	},
	1007: {
		{
			ID: 3041, AssignmentID: 2041, UserID: 5001,
			Score: ptr(82.0), Grade: "82", EnteredGrade: "82", EnteredScore: ptr(82.0),
			SubmittedAt: dayp(-29, 20, 8), GradedAt: dayp(-27, 11, 20),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3042, AssignmentID: 2042, UserID: 5001,
			Score: ptr(91.0), Grade: "91", EnteredGrade: "91", EnteredScore: ptr(91.0),
			SubmittedAt: dayp(-5, 20, 8), GradedAt: dayp(-2, 11, 20),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3043, AssignmentID: 2043, UserID: 5001,
			Score: ptr(100.0), Grade: "100", EnteredGrade: "100", EnteredScore: ptr(100.0),
			SubmittedAt: dayp(-22, 20, 8), GradedAt: dayp(-20, 11, 20),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3044, AssignmentID: 2044, UserID: 5001,
			Score: ptr(100.0), Grade: "100", EnteredGrade: "100", EnteredScore: ptr(100.0),
			SubmittedAt: dayp(-15, 20, 8), GradedAt: dayp(-10, 11, 20),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
	},
}

// ───────────────────────────────────────────────────────────────────
// Announcements & Discussion Topics
// ───────────────────────────────────────────────────────────────────

var Announcements = map[int][]Announcement{
	1001: {
		{ID: 9001, Title: "Midterm Review Session", Message: "Live review this Friday at 5 PM in room 4-145.", PostedAt: dayp(-3, 14, 30), HTMLURL: "https://mock-canvas.local/courses/1001/announcements/9001"},
		{ID: 9002, Title: "Problem Set 4 Clarification", Message: "Q4 expects a closed-form solution via master theorem.", PostedAt: dayp(-2, 10, 0), HTMLURL: "https://mock-canvas.local/courses/1001/announcements/9002"},
	},
	1003: {
		{ID: 9011, Title: "Graph Unit Starts Monday", Message: "Please read chapter 12 before lecture.", PostedAt: dayp(-1, 9, 0), HTMLURL: "https://mock-canvas.local/courses/1003/announcements/9011"},
	},
}

var DiscussionTopics = map[int][]DiscussionTopic{
	1001: {
		{ID: 9501, Title: "Induction step intuition", Message: "Share how you explain induction to a classmate.", PostedAt: dayp(-5, 18, 0), HTMLURL: "https://mock-canvas.local/courses/1001/discussion_topics/9501"},
		{ID: 9502, Title: "Recurrence pitfalls", Message: "Common mistakes when expanding recurrence trees.", PostedAt: dayp(-4, 11, 0), HTMLURL: "https://mock-canvas.local/courses/1001/discussion_topics/9502"},
	},
	1003: {
		{ID: 9511, Title: "Dijkstra vs BFS", Message: "When weighted edges change your approach.", PostedAt: dayp(-3, 17, 45), HTMLURL: "https://mock-canvas.local/courses/1003/discussion_topics/9511"},
	},
}

// ───────────────────────────────────────────────────────────────────
// Canvas Planner Items — /api/v1/planner/items
// Generated from assignments with time-relative due dates
// ───────────────────────────────────────────────────────────────────

// GetPlannerItems builds planner items dynamically from assignments.
func GetPlannerItems() []PlannerItem {
	var items []PlannerItem
	for courseID, assignments := range Assignments {
		courseName := ""
		if c, ok := CoursesByID[courseID]; ok {
			courseName = c.Name
		}
		subs := Submissions[courseID]
		subByAssignment := make(map[int]Submission, len(subs))
		for _, s := range subs {
			subByAssignment[s.AssignmentID] = s
		}
		for _, a := range assignments {
			if a.DueAt == nil {
				continue
			}
			sub, hasSub := subByAssignment[a.ID]
			ps := &PlannerSubmission{}
			if hasSub {
				ps.Submitted = sub.WorkflowState == "graded" || sub.WorkflowState == "submitted"
				ps.Graded = sub.WorkflowState == "graded"
				ps.Late = sub.Late
				ps.Missing = sub.Missing
				ps.Excused = sub.Excused
				ps.HasFeedback = sub.Score != nil
			}
			items = append(items, PlannerItem{
				ContextType:   "Course",
				CourseID:      courseID,
				PlannableDate: *a.DueAt,
				Plannable: Plannable{
					ID:             a.ID,
					Title:          a.Name,
					DueAt:          a.DueAt,
					PointsPossible: a.PointsPossible,
					CourseID:       courseID,
					CreatedAt:      a.CreatedAt,
					UpdatedAt:      a.UpdatedAt,
				},
				PlannableType: "assignment",
				NewActivity:   !ps.Submitted && a.DueAt.After(now()),
				Submissions:   ps,
				HTMLURL:       a.HTMLURL,
				ContextName:   courseName,
			})
		}
	}
	return items
}

// ───────────────────────────────────────────────────────────────────
// Calendar Events — /api/v1/calendar_events
// Recurring class meetings for each course
// ───────────────────────────────────────────────────────────────────

// GetCalendarEvents builds recurring lecture/recitation events for the current week.
func GetCalendarEvents() []CalendarEvent {
	var events []CalendarEvent
	n := now()
	// Find Monday of current week
	weekday := int(n.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := n.Add(time.Duration(-(weekday - 1)) * 24 * time.Hour)

	id := 8001
	// CS225 — MWF 10:00-10:50
	for _, doff := range []int{0, 2, 4} { // Mon, Wed, Fri
		start := monday.Add(time.Duration(doff)*24*time.Hour + 10*time.Hour)
		end := start.Add(50 * time.Minute)
		events = append(events, CalendarEvent{
			ID: id, Title: "CS225 Lecture", StartAt: &start, EndAt: &end,
			LocationName: "Room 4-145", ContextCode: "course_1001",
			WorkflowState: "active", AllDay: false,
		})
		id++
	}
	// 18.06 — TR 11:00-12:15
	for _, doff := range []int{1, 3} { // Tue, Thu
		start := monday.Add(time.Duration(doff)*24*time.Hour + 11*time.Hour)
		end := start.Add(75 * time.Minute)
		events = append(events, CalendarEvent{
			ID: id, Title: "18.06 Lecture", StartAt: &start, EndAt: &end,
			LocationName: "Room 10-250", ContextCode: "course_1002",
			WorkflowState: "active", AllDay: false,
		})
		id++
	}
	// 6.006 — MWF 13:00-13:50
	for _, doff := range []int{0, 2, 4} {
		start := monday.Add(time.Duration(doff)*24*time.Hour + 13*time.Hour)
		end := start.Add(50 * time.Minute)
		events = append(events, CalendarEvent{
			ID: id, Title: "6.006 Lecture", StartAt: &start, EndAt: &end,
			LocationName: "Room 26-100", ContextCode: "course_1003",
			WorkflowState: "active", AllDay: false,
		})
		id++
	}
	// 8.01 — TR 14:00-15:15
	for _, doff := range []int{1, 3} {
		start := monday.Add(time.Duration(doff)*24*time.Hour + 14*time.Hour)
		end := start.Add(75 * time.Minute)
		events = append(events, CalendarEvent{
			ID: id, Title: "8.01 Lecture", StartAt: &start, EndAt: &end,
			LocationName: "Room 6-120", ContextCode: "course_1004",
			WorkflowState: "active", AllDay: false,
		})
		id++
	}
	// 6.042J — MWF 15:00-15:50
	for _, doff := range []int{0, 2, 4} {
		start := monday.Add(time.Duration(doff)*24*time.Hour + 15*time.Hour)
		end := start.Add(50 * time.Minute)
		events = append(events, CalendarEvent{
			ID: id, Title: "6.042J Lecture", StartAt: &start, EndAt: &end,
			LocationName: "Room 32-123", ContextCode: "course_1005",
			WorkflowState: "active", AllDay: false,
		})
		id++
	}
	return events
}
