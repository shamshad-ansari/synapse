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
		CourseCode:       "CS225",
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
		CourseCode:       "18.06",
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
		CourseCode:       "6.006",
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
	{
		ID:               1004,
		Name:             "8.01 Physics I: Mechanics",
		CourseCode:       "8.01",
		WorkflowState:    "available",
		AccountID:        1,
		EnrollmentTermID: 202501,
		StartAt:          semesterStart,
		EndAt:            semesterEnd,
		CreatedAt:        t(2024, 12, 2, 0, 0),
		DefaultView:      "modules",
		IsPublic:         false,
		TimeZone:         "America/New_York",
		UUID:             "801-uuid-physics-mechanics-2025",
		HideFinalGrades:  false,
		CourseFormat:     "on_campus",
		TotalStudents:    131,
		SyllabusBody: `<h2>8.01 Physics I: Mechanics — Spring 2025</h2>
<h3>Topics</h3>
<ul>
  <li>Kinematics and Newton's Laws</li>
  <li>Energy, Momentum, and Rotational Motion</li>
  <li>Oscillations, Gravitation, and Applications</li>
</ul>`,
	},
	{
		ID:               1005,
		Name:             "6.042J Mathematics for Computer Science",
		CourseCode:       "6.042J",
		WorkflowState:    "available",
		AccountID:        1,
		EnrollmentTermID: 202501,
		StartAt:          semesterStart,
		EndAt:            semesterEnd,
		CreatedAt:        t(2024, 12, 2, 0, 0),
		DefaultView:      "feed",
		IsPublic:         false,
		TimeZone:         "America/New_York",
		UUID:             "6042j-uuid-math-for-cs-2025",
		HideFinalGrades:  false,
		CourseFormat:     "on_campus",
		TotalStudents:    89,
		SyllabusBody: `<h2>6.042J Mathematics for Computer Science — Spring 2025</h2>
<h3>Topics</h3>
<ul>
  <li>Proofs and Invariants</li>
  <li>Graphs, Trees, and Number Theory</li>
  <li>Combinatorics and Probability</li>
</ul>`,
	},
	{
		ID:               1006,
		Name:             "6.100A Introduction to CS Programming in Python",
		CourseCode:       "6.100A",
		WorkflowState:    "available",
		AccountID:        1,
		EnrollmentTermID: 202501,
		StartAt:          semesterStart,
		EndAt:            semesterEnd,
		CreatedAt:        t(2024, 12, 3, 0, 0),
		DefaultView:      "feed",
		IsPublic:         false,
		TimeZone:         "America/New_York",
		UUID:             "6100a-uuid-python-programming-2025",
		HideFinalGrades:  false,
		CourseFormat:     "on_campus",
		TotalStudents:    210,
		SyllabusBody: `<h2>6.100A Intro Programming in Python — Spring 2025</h2>
<h3>Topics</h3>
<ul>
  <li>Control flow, functions, recursion</li>
  <li>Data structures and algorithmic problem solving</li>
  <li>Testing, debugging, and mini-projects</li>
</ul>`,
	},
	{
		ID:               1007,
		Name:             "21W.789 Communicating with Data",
		CourseCode:       "21W.789",
		WorkflowState:    "available",
		AccountID:        1,
		EnrollmentTermID: 202501,
		StartAt:          semesterStart,
		EndAt:            semesterEnd,
		CreatedAt:        t(2024, 12, 4, 0, 0),
		DefaultView:      "feed",
		IsPublic:         false,
		TimeZone:         "America/New_York",
		UUID:             "21789-uuid-communicating-with-data-2025",
		HideFinalGrades:  false,
		CourseFormat:     "seminar",
		TotalStudents:    24,
		SyllabusBody: `<h2>21W.789 Communicating with Data — Spring 2025</h2>
<h3>Topics</h3>
<ul>
  <li>Narrative structure for technical writing</li>
  <li>Data-supported arguments and visual rhetoric</li>
  <li>Peer review and revision workshops</li>
</ul>`,
	},
	{
		ID:               1008,
		Name:             "15.053 Optimization Methods in Management Science",
		CourseCode:       "15.053",
		WorkflowState:    "available",
		AccountID:        1,
		EnrollmentTermID: 202501,
		StartAt:          semesterStart,
		EndAt:            semesterEnd,
		CreatedAt:        t(2024, 12, 5, 0, 0),
		DefaultView:      "feed",
		IsPublic:         false,
		TimeZone:         "America/New_York",
		UUID:             "15053-uuid-optimization-methods-2025",
		HideFinalGrades:  false,
		CourseFormat:     "on_campus",
		TotalStudents:    53,
		SyllabusBody: `<h2>15.053 Optimization Methods — Spring 2025</h2>
<h3>Topics</h3>
<ul>
  <li>Linear programming and duality</li>
  <li>Network models and integer optimization</li>
  <li>Decision analysis applications</li>
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
			Description:    "Logic fundamentals and truth tables. Covers propositional and predicate logic.",
			DueAt:          tp(2025, 2, 10, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 100,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 1, 22, 9, 0), UpdatedAt: t(2025, 1, 22, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1001/assignments/2001",
		},
		{
			ID: 2002, Name: "Problem Set 2", CourseID: 1001,
			Description:    "Set theory operations, relations, and equivalence classes.",
			DueAt:          tp(2025, 2, 20, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 100,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 2, 3, 9, 0), UpdatedAt: t(2025, 2, 3, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1001/assignments/2002",
		},
		{
			ID: 2003, Name: "Problem Set 3", CourseID: 1001,
			Description:    "Mathematical induction and recursion. Covers strong and weak induction, recursive definitions.",
			DueAt:          tp(2025, 2, 27, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 100,
			Position: 3, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 2, 14, 9, 0), UpdatedAt: t(2025, 2, 14, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1001/assignments/2003",
		},
		{
			ID: 2004, Name: "Midterm Exam", CourseID: 1001,
			Description:    "Covers all material through Week 8: Logic, Set Theory, Induction, and Recursion.",
			DueAt:          tp(2025, 3, 12, 14, 0),
			PointsPossible: 200, GradingType: "points", AssignmentGroupID: 101,
			Position: 4, SubmissionTypes: []string{"on_paper"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 1, 22, 9, 0), UpdatedAt: t(2025, 1, 22, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1001/assignments/2004",
		},
	},
	1002: {
		{
			ID: 2011, Name: "Problem Set 1", CourseID: 1002,
			Description:    "Vector spaces, linear independence, and spanning sets.",
			DueAt:          tp(2025, 2, 15, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 200,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 1, 22, 9, 0), UpdatedAt: t(2025, 1, 22, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1002/assignments/2011",
		},
		{
			ID: 2012, Name: "Problem Set 2", CourseID: 1002,
			Description:    "Matrix operations, determinants, and linear transformations.",
			DueAt:          tp(2025, 3, 5, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 200,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 2, 10, 9, 0), UpdatedAt: t(2025, 2, 10, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1002/assignments/2012",
		},
		{
			ID: 2013, Name: "Eigenvalues Quiz", CourseID: 1002,
			Description:    "Timed quiz on characteristic polynomials and eigenbases.",
			DueAt:          tp(2025, 3, 12, 21, 0),
			PointsPossible: 40, GradingType: "points", AssignmentGroupID: 201,
			Position: 3, SubmissionTypes: []string{"online_quiz"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 2, 25, 9, 0), UpdatedAt: t(2025, 2, 25, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1002/assignments/2013",
		},
		{
			ID: 2014, Name: "Midterm I", CourseID: 1002,
			Description:    "Covers vector spaces, matrix algebra, and linear maps.",
			DueAt:          tp(2025, 3, 20, 13, 0),
			PointsPossible: 120, GradingType: "points", AssignmentGroupID: 202,
			Position: 4, SubmissionTypes: []string{"on_paper"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 2, 1, 9, 0), UpdatedAt: t(2025, 2, 1, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1002/assignments/2014",
		},
	},
	1003: {
		{
			ID: 2021, Name: "Algorithm Analysis HW1", CourseID: 1003,
			Description:    "Asymptotic analysis, recurrence relations, and master theorem applications.",
			DueAt:          tp(2025, 2, 18, 23, 59),
			PointsPossible: 80, GradingType: "points", AssignmentGroupID: 300,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 1, 22, 9, 0), UpdatedAt: t(2025, 1, 22, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1003/assignments/2021",
		},
		{
			ID: 2022, Name: "Sorting Lab", CourseID: 1003,
			Description:    "Implement and benchmark merge sort, quicksort, and heapsort.",
			DueAt:          tp(2025, 2, 28, 23, 59),
			PointsPossible: 60, GradingType: "points", AssignmentGroupID: 300,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 2, 12, 10, 0), UpdatedAt: t(2025, 2, 12, 10, 0),
			HTMLURL: "https://mock-canvas.local/courses/1003/assignments/2022",
		},
		{
			ID: 2023, Name: "Graph Algorithms Pset", CourseID: 1003,
			Description:    "BFS/DFS shortest paths and minimum spanning trees.",
			DueAt:          tp(2025, 3, 18, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 301,
			Position: 3, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 3, 1, 10, 0), UpdatedAt: t(2025, 3, 1, 10, 0),
			HTMLURL: "https://mock-canvas.local/courses/1003/assignments/2023",
		},
		{
			ID: 2024, Name: "Dynamic Programming Quiz", CourseID: 1003,
			Description:    "Short quiz covering memoization and tabulation tradeoffs.",
			DueAt:          tp(2025, 3, 22, 20, 30),
			PointsPossible: 25, GradingType: "points", AssignmentGroupID: 302,
			Position: 4, SubmissionTypes: []string{"online_quiz"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 3, 5, 8, 0), UpdatedAt: t(2025, 3, 5, 8, 0),
			HTMLURL: "https://mock-canvas.local/courses/1003/assignments/2024",
		},
	},
	1004: {
		{
			ID: 2031, Name: "Kinematics Worksheet", CourseID: 1004,
			Description:    "Uniform acceleration and projectile motion problems.",
			DueAt:          tp(2025, 2, 11, 23, 59),
			PointsPossible: 50, GradingType: "points", AssignmentGroupID: 400,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 1, 25, 9, 0), UpdatedAt: t(2025, 1, 25, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1004/assignments/2031",
		},
		{
			ID: 2032, Name: "Forces & Free-Body Diagram Quiz", CourseID: 1004,
			Description:    "Interpret force diagrams and apply Newton's second law.",
			DueAt:          tp(2025, 2, 21, 22, 0),
			PointsPossible: 30, GradingType: "points", AssignmentGroupID: 401,
			Position: 2, SubmissionTypes: []string{"online_quiz"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 2, 8, 9, 0), UpdatedAt: t(2025, 2, 8, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1004/assignments/2032",
		},
		{
			ID: 2033, Name: "Midterm Mechanics", CourseID: 1004,
			Description:    "In-person mechanics midterm.",
			DueAt:          tp(2025, 3, 16, 14, 0),
			PointsPossible: 150, GradingType: "points", AssignmentGroupID: 402,
			Position: 3, SubmissionTypes: []string{"on_paper"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 2, 1, 9, 0), UpdatedAt: t(2025, 2, 1, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1004/assignments/2033",
		},
	},
	1005: {
		{
			ID: 2041, Name: "Proof Techniques Set", CourseID: 1005,
			Description:    "Direct proof, contradiction, and induction exercises.",
			DueAt:          tp(2025, 2, 13, 23, 59),
			PointsPossible: 90, GradingType: "points", AssignmentGroupID: 500,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 1, 29, 9, 0), UpdatedAt: t(2025, 1, 29, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1005/assignments/2041",
		},
		{
			ID: 2042, Name: "Combinatorics Challenge", CourseID: 1005,
			Description:    "Counting, permutations/combinations, and pigeonhole principle.",
			DueAt:          tp(2025, 3, 2, 23, 59),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 500,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 2, 14, 9, 0), UpdatedAt: t(2025, 2, 14, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1005/assignments/2042",
		},
		{
			ID: 2043, Name: "Probability Quiz I", CourseID: 1005,
			Description:    "Conditional probability and Bayes rule.",
			DueAt:          tp(2025, 3, 11, 21, 30),
			PointsPossible: 35, GradingType: "points", AssignmentGroupID: 501,
			Position: 3, SubmissionTypes: []string{"online_quiz"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 2, 24, 8, 0), UpdatedAt: t(2025, 2, 24, 8, 0),
			HTMLURL: "https://mock-canvas.local/courses/1005/assignments/2043",
		},
	},
	1006: {
		{
			ID: 2051, Name: "Python Functions Lab", CourseID: 1006,
			Description:    "Functional decomposition, unit tests, and code readability.",
			DueAt:          tp(2025, 2, 12, 23, 59),
			PointsPossible: 60, GradingType: "points", AssignmentGroupID: 600,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 1, 30, 9, 0), UpdatedAt: t(2025, 1, 30, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1006/assignments/2051",
		},
		{
			ID: 2052, Name: "Recursion Practice", CourseID: 1006,
			Description:    "Implement recursive strategies and analyze stack traces.",
			DueAt:          tp(2025, 2, 26, 23, 59),
			PointsPossible: 75, GradingType: "points", AssignmentGroupID: 600,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 2, 9, 9, 0), UpdatedAt: t(2025, 2, 9, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1006/assignments/2052",
		},
		{
			ID: 2053, Name: "Mini Project: Study Planner", CourseID: 1006,
			Description:    "Build a CLI study planner using files and dictionaries.",
			DueAt:          tp(2025, 3, 19, 23, 59),
			PointsPossible: 120, GradingType: "points", AssignmentGroupID: 601,
			Position: 3, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 2, 27, 9, 0), UpdatedAt: t(2025, 2, 27, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1006/assignments/2053",
		},
	},
	1007: {
		{
			ID: 2061, Name: "Essay Draft 1", CourseID: 1007,
			Description:    "First draft: Explain a data story for a non-technical audience.",
			DueAt:          tp(2025, 2, 22, 18, 0),
			PointsPossible: 50, GradingType: "points", AssignmentGroupID: 700,
			Position: 1, SubmissionTypes: []string{"online_text_entry"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 2, 1, 10, 0), UpdatedAt: t(2025, 2, 1, 10, 0),
			HTMLURL: "https://mock-canvas.local/courses/1007/assignments/2061",
		},
		{
			ID: 2062, Name: "Peer Review Memo", CourseID: 1007,
			Description:    "Structured peer feedback with evidence-backed revisions.",
			DueAt:          tp(2025, 3, 6, 18, 0),
			PointsPossible: 40, GradingType: "points", AssignmentGroupID: 700,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 2, 20, 10, 0), UpdatedAt: t(2025, 2, 20, 10, 0),
			HTMLURL: "https://mock-canvas.local/courses/1007/assignments/2062",
		},
		{
			ID: 2063, Name: "Final Narrative Brief", CourseID: 1007,
			Description:    "Final revised communication brief with appendix visuals.",
			DueAt:          tp(2025, 4, 2, 18, 0),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 701,
			Position: 3, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 3, 10, 10, 0), UpdatedAt: t(2025, 3, 10, 10, 0),
			HTMLURL: "https://mock-canvas.local/courses/1007/assignments/2063",
		},
	},
	1008: {
		{
			ID: 2071, Name: "LP Modeling Homework", CourseID: 1008,
			Description:    "Model production planning and sensitivity analysis in LP.",
			DueAt:          tp(2025, 2, 17, 23, 0),
			PointsPossible: 100, GradingType: "points", AssignmentGroupID: 800,
			Position: 1, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: true,
			CreatedAt: t(2025, 2, 1, 9, 0), UpdatedAt: t(2025, 2, 1, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1008/assignments/2071",
		},
		{
			ID: 2072, Name: "Network Flow Case", CourseID: 1008,
			Description:    "Apply max-flow/min-cut to logistics planning scenario.",
			DueAt:          tp(2025, 3, 9, 23, 0),
			PointsPossible: 90, GradingType: "points", AssignmentGroupID: 800,
			Position: 2, SubmissionTypes: []string{"online_upload"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 2, 21, 9, 0), UpdatedAt: t(2025, 2, 21, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1008/assignments/2072",
		},
		{
			ID: 2073, Name: "Integer Programming Quiz", CourseID: 1008,
			Description:    "Mixed integer model formulation and branch-and-bound basics.",
			DueAt:          tp(2025, 3, 21, 20, 0),
			PointsPossible: 30, GradingType: "points", AssignmentGroupID: 801,
			Position: 3, SubmissionTypes: []string{"online_quiz"}, Published: true, HasSubmitted: false,
			CreatedAt: t(2025, 3, 2, 9, 0), UpdatedAt: t(2025, 3, 2, 9, 0),
			HTMLURL: "https://mock-canvas.local/courses/1008/assignments/2073",
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
		{
			ID: 3013, AssignmentID: 2013, UserID: 5001,
			Score: ptr(35.0), Grade: "35", EnteredGrade: "35", EnteredScore: ptr(35.0),
			SubmittedAt: tp(2025, 3, 12, 19, 54), GradedAt: tp(2025, 3, 13, 9, 10),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3014, AssignmentID: 2014, UserID: 5001,
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
		{
			ID: 3022, AssignmentID: 2022, UserID: 5001,
			Score: ptr(56.0), Grade: "56", EnteredGrade: "56", EnteredScore: ptr(56.0),
			SubmittedAt: tp(2025, 2, 28, 23, 3), GradedAt: tp(2025, 3, 2, 8, 45),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3023, AssignmentID: 2023, UserID: 5001,
			Score: nil, Grade: "", EnteredGrade: "", EnteredScore: nil,
			SubmittedAt: nil, GradedAt: nil,
			Attempt: 0, WorkflowState: "unsubmitted", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3024, AssignmentID: 2024, UserID: 5001,
			Score: nil, Grade: "", EnteredGrade: "", EnteredScore: nil,
			SubmittedAt: nil, GradedAt: nil,
			Attempt: 0, WorkflowState: "unsubmitted", Late: false, Missing: false, Excused: false,
		},
	},
	1004: {
		{
			ID: 3031, AssignmentID: 2031, UserID: 5001,
			Score: ptr(47.0), Grade: "47", EnteredGrade: "47", EnteredScore: ptr(47.0),
			SubmittedAt: tp(2025, 2, 11, 22, 31), GradedAt: tp(2025, 2, 13, 10, 5),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3032, AssignmentID: 2032, UserID: 5001,
			Score: ptr(24.0), Grade: "24", EnteredGrade: "24", EnteredScore: ptr(24.0),
			SubmittedAt: tp(2025, 2, 21, 21, 42), GradedAt: tp(2025, 2, 22, 12, 1),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3033, AssignmentID: 2033, UserID: 5001,
			Score: nil, Grade: "", EnteredGrade: "", EnteredScore: nil,
			SubmittedAt: nil, GradedAt: nil,
			Attempt: 0, WorkflowState: "unsubmitted", Late: false, Missing: false, Excused: false,
		},
	},
	1005: {
		{
			ID: 3041, AssignmentID: 2041, UserID: 5001,
			Score: ptr(78.0), Grade: "78", EnteredGrade: "78", EnteredScore: ptr(78.0),
			SubmittedAt: tp(2025, 2, 13, 20, 8), GradedAt: tp(2025, 2, 15, 11, 20),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3042, AssignmentID: 2042, UserID: 5001,
			Score: nil, Grade: "", EnteredGrade: "", EnteredScore: nil,
			SubmittedAt: nil, GradedAt: nil,
			Attempt: 0, WorkflowState: "unsubmitted", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3043, AssignmentID: 2043, UserID: 5001,
			Score: ptr(28.0), Grade: "28", EnteredGrade: "28", EnteredScore: ptr(28.0),
			SubmittedAt: tp(2025, 3, 11, 19, 11), GradedAt: tp(2025, 3, 12, 9, 20),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
	},
	1006: {
		{
			ID: 3051, AssignmentID: 2051, UserID: 5001,
			Score: ptr(54.0), Grade: "54", EnteredGrade: "54", EnteredScore: ptr(54.0),
			SubmittedAt: tp(2025, 2, 12, 22, 5), GradedAt: tp(2025, 2, 14, 10, 17),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3052, AssignmentID: 2052, UserID: 5001,
			Score: ptr(67.0), Grade: "67", EnteredGrade: "67", EnteredScore: ptr(67.0),
			SubmittedAt: tp(2025, 2, 26, 22, 14), GradedAt: tp(2025, 2, 28, 9, 41),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3053, AssignmentID: 2053, UserID: 5001,
			Score: nil, Grade: "", EnteredGrade: "", EnteredScore: nil,
			SubmittedAt: nil, GradedAt: nil,
			Attempt: 0, WorkflowState: "unsubmitted", Late: false, Missing: false, Excused: false,
		},
	},
	1007: {
		{
			ID: 3061, AssignmentID: 2061, UserID: 5001,
			Score: ptr(44.0), Grade: "44", EnteredGrade: "44", EnteredScore: ptr(44.0),
			SubmittedAt: tp(2025, 2, 22, 16, 40), GradedAt: tp(2025, 2, 24, 13, 28),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3062, AssignmentID: 2062, UserID: 5001,
			Score: ptr(36.0), Grade: "36", EnteredGrade: "36", EnteredScore: ptr(36.0),
			SubmittedAt: tp(2025, 3, 6, 17, 12), GradedAt: tp(2025, 3, 8, 11, 2),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3063, AssignmentID: 2063, UserID: 5001,
			Score: nil, Grade: "", EnteredGrade: "", EnteredScore: nil,
			SubmittedAt: nil, GradedAt: nil,
			Attempt: 0, WorkflowState: "unsubmitted", Late: false, Missing: false, Excused: false,
		},
	},
	1008: {
		{
			ID: 3071, AssignmentID: 2071, UserID: 5001,
			Score: ptr(84.0), Grade: "84", EnteredGrade: "84", EnteredScore: ptr(84.0),
			SubmittedAt: tp(2025, 2, 17, 21, 55), GradedAt: tp(2025, 2, 19, 9, 49),
			Attempt: 1, WorkflowState: "graded", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3072, AssignmentID: 2072, UserID: 5001,
			Score: nil, Grade: "", EnteredGrade: "", EnteredScore: nil,
			SubmittedAt: nil, GradedAt: nil,
			Attempt: 0, WorkflowState: "unsubmitted", Late: false, Missing: false, Excused: false,
		},
		{
			ID: 3073, AssignmentID: 2073, UserID: 5001,
			Score: nil, Grade: "", EnteredGrade: "", EnteredScore: nil,
			SubmittedAt: nil, GradedAt: nil,
			Attempt: 0, WorkflowState: "unsubmitted", Late: false, Missing: false, Excused: false,
		},
	},
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

var Announcements = map[int][]Announcement{
	1001: {
		{ID: 9001, Title: "Midterm Review Session", Message: "Live review this Friday at 5 PM in room 4-145.", PostedAt: tp(2025, 3, 8, 14, 30), HTMLURL: "https://mock-canvas.local/courses/1001/announcements/9001"},
		{ID: 9002, Title: "Problem Set 3 Clarification", Message: "Q4 expects an inductive proof with explicit hypothesis.", PostedAt: tp(2025, 3, 9, 10, 0), HTMLURL: "https://mock-canvas.local/courses/1001/announcements/9002"},
	},
	1003: {
		{ID: 9011, Title: "Graph Unit Starts Monday", Message: "Please read chapter 12 before lecture.", PostedAt: tp(2025, 3, 10, 9, 0), HTMLURL: "https://mock-canvas.local/courses/1003/announcements/9011"},
	},
}

var DiscussionTopics = map[int][]DiscussionTopic{
	1001: {
		{ID: 9501, Title: "Induction step intuition", Message: "Share how you explain induction to a classmate.", PostedAt: tp(2025, 3, 6, 18, 0), HTMLURL: "https://mock-canvas.local/courses/1001/discussion_topics/9501"},
		{ID: 9502, Title: "Recurrence pitfalls", Message: "Common mistakes when expanding recurrence trees.", PostedAt: tp(2025, 3, 7, 11, 0), HTMLURL: "https://mock-canvas.local/courses/1001/discussion_topics/9502"},
	},
	1003: {
		{ID: 9511, Title: "Dijkstra vs BFS", Message: "When weighted edges change your approach.", PostedAt: tp(2025, 3, 8, 17, 45), HTMLURL: "https://mock-canvas.local/courses/1003/discussion_topics/9511"},
	},
}
