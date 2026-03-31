package seed

import "time"

// Canvas-format response structs with full realistic fields.

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
	ID           int        `json:"id"`
	AssignmentID int        `json:"assignment_id"`
	UserID       int        `json:"user_id"`
	Score        *float64   `json:"score"`
	Grade        string     `json:"grade"`
	EnteredGrade string     `json:"entered_grade"`
	EnteredScore *float64   `json:"entered_score"`
	SubmittedAt  *time.Time `json:"submitted_at"`
	GradedAt     *time.Time `json:"graded_at"`
	Attempt      int        `json:"attempt"`
	WorkflowState string   `json:"workflow_state"`
	Late         bool       `json:"late"`
	Missing      bool       `json:"missing"`
	Excused      bool       `json:"excused"`
}

func ptr[T any](v T) *T { return &v }

func t(year, month, day, hour, min int) time.Time {
	return time.Date(year, time.Month(month), day, hour, min, 0, 0, time.UTC)
}

func tp(year, month, day, hour, min int) *time.Time {
	v := t(year, month, day, hour, min)
	return &v
}

// -------------------------------------------------------------------
// Seed data — mapped from frontend hardcoded academic content
// -------------------------------------------------------------------

var CurrentUser = User{
	ID:           5001,
	Name:         "Alex Kim",
	SortableName: "Kim, Alex",
	FirstName:    "Alex",
	LastName:     "Kim",
	ShortName:    "Alex",
	LoginID:      "alex.kim@mit.edu",
	Email:        "alex.kim@mit.edu",
	AvatarURL:    "https://i.pravatar.cc/150?u=alexkim",
	Locale:       "en",
	TimeZone:     "America/New_York",
	Bio:          nil,
}

var semesterStart = tp(2025, 1, 21, 0, 0)
var semesterEnd = tp(2025, 5, 15, 23, 59)

var Courses = []Course{
	{
		ID:               1001,
		Name:             "CS225 Discrete Mathematics",
		CourseCode:        "CS225",
		WorkflowState:    "available",
		AccountID:        1,
		EnrollmentTermID: 202501,
		StartAt:          semesterStart,
		EndAt:            semesterEnd,
		CreatedAt:        t(2024, 12, 1, 0, 0),
		DefaultView:      "feed",
		IsPublic:         false,
		TimeZone:         "America/New_York",
		UUID:             "cs225-uuid-discrete-math-2025",
		HideFinalGrades:  false,
		CourseFormat:     "on_campus",
		TotalStudents:    48,
		SyllabusBody: `<h2>CS225 Discrete Mathematics — Spring 2025</h2>
<h3>Topics</h3>
<ul>
  <li>Weeks 1-2: Logic &amp; Proof Techniques</li>
  <li>Weeks 3-4: Set Theory &amp; Relations</li>
  <li>Weeks 5-6: Mathematical Induction</li>
  <li>Weeks 7-8: Recursion &amp; Recurrences</li>
  <li>Weeks 9-10: Graph Theory</li>
  <li>Weeks 11-12: Combinatorics</li>
</ul>
<h3>Assessments</h3>
<ul>
  <li>Weekly Problem Sets (40%)</li>
  <li>Midterm Exam — March 12 (25%)</li>
  <li>Final Exam — May 8 (35%)</li>
</ul>`,
	},
	{
		ID:               1002,
		Name:             "18.06 Linear Algebra",
		CourseCode:        "18.06",
		WorkflowState:    "available",
		AccountID:        1,
		EnrollmentTermID: 202501,
		StartAt:          semesterStart,
		EndAt:            semesterEnd,
		CreatedAt:        t(2024, 12, 1, 0, 0),
		DefaultView:      "feed",
		IsPublic:         false,
		TimeZone:         "America/New_York",
		UUID:             "1806-uuid-linear-algebra-2025",
		HideFinalGrades:  false,
		CourseFormat:     "on_campus",
		TotalStudents:    62,
		SyllabusBody: `<h2>18.06 Linear Algebra — Spring 2025</h2>
<h3>Topics</h3>
<ul>
  <li>Weeks 1-3: Vector Spaces &amp; Subspaces</li>
  <li>Weeks 4-6: Linear Transformations &amp; Matrices</li>
  <li>Weeks 7-9: Eigenvalues &amp; Eigenvectors</li>
  <li>Weeks 10-12: Orthogonality &amp; Least Squares</li>
</ul>`,
	},
	{
		ID:               1003,
		Name:             "6.006 Introduction to Algorithms",
		CourseCode:        "6.006",
		WorkflowState:    "available",
		AccountID:        1,
		EnrollmentTermID: 202501,
		StartAt:          semesterStart,
		EndAt:            semesterEnd,
		CreatedAt:        t(2024, 12, 1, 0, 0),
		DefaultView:      "feed",
		IsPublic:         false,
		TimeZone:         "America/New_York",
		UUID:             "6006-uuid-algorithms-2025",
		HideFinalGrades:  false,
		CourseFormat:     "on_campus",
		TotalStudents:    55,
		SyllabusBody: `<h2>6.006 Introduction to Algorithms — Spring 2025</h2>
<h3>Topics</h3>
<ul>
  <li>Weeks 1-2: Algorithmic Thinking &amp; Complexity</li>
  <li>Weeks 3-5: Sorting &amp; Searching</li>
  <li>Weeks 6-8: Graph Algorithms</li>
  <li>Weeks 9-11: Dynamic Programming</li>
  <li>Week 12: Advanced Topics</li>
</ul>`,
	},
}

// CoursesByID is a lookup map populated by init().
var CoursesByID map[int]Course

func init() {
	CoursesByID = make(map[int]Course, len(Courses))
	for _, c := range Courses {
		CoursesByID[c.ID] = c
	}
}

// Assignments per course — keyed by course ID.
var Assignments = map[int][]Assignment{
	1001: {
		{
			ID: 2001, Name: "Problem Set 1", CourseID: 1001,
			Description:     "Logic fundamentals and truth tables. Covers propositional and predicate logic.",
			DueAt:           tp(2025, 2, 10, 23, 59),
			PointsPossible:  100, GradingType: "points", AssignmentGroupID: 100,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 1, 22, 9, 0), UpdatedAt: t(2025, 1, 22, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1001/assignments/2001",
		},
		{
			ID: 2002, Name: "Problem Set 2", CourseID: 1001,
			Description:     "Set theory operations, relations, and equivalence classes.",
			DueAt:           tp(2025, 2, 20, 23, 59),
			PointsPossible:  100, GradingType: "points", AssignmentGroupID: 100,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 2, 3, 9, 0), UpdatedAt: t(2025, 2, 3, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1001/assignments/2002",
		},
		{
			ID: 2003, Name: "Problem Set 3", CourseID: 1001,
			Description:     "Mathematical induction and recursion. Covers strong and weak induction, recursive definitions.",
			DueAt:           tp(2025, 2, 27, 23, 59),
			PointsPossible:  100, GradingType: "points", AssignmentGroupID: 100,
			Position: 3, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 2, 14, 9, 0), UpdatedAt: t(2025, 2, 14, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1001/assignments/2003",
		},
		{
			ID: 2004, Name: "Midterm Exam", CourseID: 1001,
			Description:     "Covers all material through Week 8: Logic, Set Theory, Induction, and Recursion.",
			DueAt:           tp(2025, 3, 12, 14, 0),
			PointsPossible:  200, GradingType: "points", AssignmentGroupID: 101,
			Position: 4, SubmissionTypes: []string{"on_paper"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 1, 22, 9, 0), UpdatedAt: t(2025, 1, 22, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1001/assignments/2004",
		},
	},
	1002: {
		{
			ID: 2011, Name: "Problem Set 1", CourseID: 1002,
			Description:     "Vector spaces, linear independence, and spanning sets.",
			DueAt:           tp(2025, 2, 15, 23, 59),
			PointsPossible:  100, GradingType: "points", AssignmentGroupID: 200,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 1, 22, 9, 0), UpdatedAt: t(2025, 1, 22, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1002/assignments/2011",
		},
		{
			ID: 2012, Name: "Problem Set 2", CourseID: 1002,
			Description:     "Matrix operations, determinants, and linear transformations.",
			DueAt:           tp(2025, 3, 5, 23, 59),
			PointsPossible:  100, GradingType: "points", AssignmentGroupID: 200,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 2, 10, 9, 0), UpdatedAt: t(2025, 2, 10, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1002/assignments/2012",
		},
	},
	1003: {
		{
			ID: 2021, Name: "Algorithm Analysis HW1", CourseID: 1003,
			Description:     "Asymptotic analysis, recurrence relations, and master theorem applications.",
			DueAt:           tp(2025, 2, 18, 23, 59),
			PointsPossible:  80, GradingType: "points", AssignmentGroupID: 300,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 1, 22, 9, 0), UpdatedAt: t(2025, 1, 22, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1003/assignments/2021",
		},
	},
}

// Submissions per course — keyed by course ID.
var Submissions = map[int][]Submission{
	1001: {
		{
			ID: 3001, AssignmentID: 2001, UserID: 5001,
			Score: ptr(92.0), Grade: "92", EnteredGrade: "92", EnteredScore: ptr(92.0),
			SubmittedAt: tp(2025, 2, 10, 18, 30), GradedAt: tp(2025, 2, 12, 10, 0),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3002, AssignmentID: 2002, UserID: 5001,
			Score: ptr(85.0), Grade: "85", EnteredGrade: "85", EnteredScore: ptr(85.0),
			SubmittedAt: tp(2025, 2, 20, 21, 15), GradedAt: tp(2025, 2, 22, 14, 0),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3003, AssignmentID: 2003, UserID: 5001,
			Score: nil, Grade: "", EnteredGrade: "", EnteredScore: nil,
			SubmittedAt: nil, GradedAt: nil,
			Attempt: 0, WorkflowState: "unsubmitted", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3004, AssignmentID: 2004, UserID: 5001,
			Score: nil, Grade: "", EnteredGrade: "", EnteredScore: nil,
			SubmittedAt: nil, GradedAt: nil,
			Attempt: 0, WorkflowState: "unsubmitted", Late: false, Missing: false, Excused: false,
		},
	},
	1002: {
		{
			ID: 3011, AssignmentID: 2011, UserID: 5001,
			Score: ptr(88.0), Grade: "88", EnteredGrade: "88", EnteredScore: ptr(88.0),
			SubmittedAt: tp(2025, 2, 15, 20, 0), GradedAt: tp(2025, 2, 17, 11, 0),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3012, AssignmentID: 2012, UserID: 5001,
			Score: nil, Grade: "", EnteredGrade: "", EnteredScore: nil,
			SubmittedAt: nil, GradedAt: nil,
			Attempt: 0, WorkflowState: "unsubmitted", Late: false, Missing: false, Excused: false,
		},
	},
	1003: {
		{
			ID: 3021, AssignmentID: 2021, UserID: 5001,
			Score: ptr(72.0), Grade: "72", EnteredGrade: "72", EnteredScore: ptr(72.0),
			SubmittedAt: tp(2025, 2, 18, 22, 45), GradedAt: tp(2025, 2, 20, 9, 30),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
	},
}
