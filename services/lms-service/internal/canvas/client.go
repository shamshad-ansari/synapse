package canvas

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ResolveServerURL returns the URL to use for server-side calls to a Canvas instance.
// If an internal override is configured (CANVAS_INTERNAL_URL), it replaces the
// scheme+host of institutionURL so browser-facing links still use the public URL
// while server-to-server calls use the internal hostname.
func ResolveServerURL(institutionURL, internalOverride string) string {
	if internalOverride == "" {
		return institutionURL
	}
	return internalOverride
}

// CanvasAPIError is returned when the Canvas API responds with a non-200 status.
type CanvasAPIError struct {
	StatusCode int
	Body       string
}

func (e *CanvasAPIError) Error() string {
	return fmt.Sprintf("canvas api error: status %d: %s", e.StatusCode, e.Body)
}

type CanvasCourse struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	EnrollmentTermID int    `json:"enrollment_term_id"`
}

type CanvasAssignment struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	DueAt             *time.Time `json:"due_at"`
	PointsPossible    float64    `json:"points_possible"`
	AssignmentGroupID int        `json:"assignment_group_id"`
}

type CanvasSubmission struct {
	AssignmentID int        `json:"assignment_id"`
	Score        *float64   `json:"score"`
	Grade        string     `json:"grade"`
	EnteredGrade string     `json:"entered_grade"`
	SubmittedAt  *time.Time `json:"submitted_at"`
	GradedAt     *time.Time `json:"graded_at"`
	WorkflowState string    `json:"workflow_state"`
	Late          bool      `json:"late"`
	Missing       bool      `json:"missing"`
	Excused       bool      `json:"excused"`
}

type CanvasAnnouncement struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	PostedAt  *time.Time `json:"posted_at"`
	HTMLURL   string     `json:"html_url"`
}

type CanvasDiscussionTopic struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	PostedAt  *time.Time `json:"posted_at"`
	HTMLURL   string     `json:"html_url"`
}

type canvasCourseWithSyllabus struct {
	SyllabusBody string `json:"syllabus_body"`
}

// CanvasClient communicates with a Canvas LMS instance.
type CanvasClient struct {
	institutionURL string
	accessToken    string
	httpClient     *http.Client
}

// NewCanvasClient creates a CanvasClient with a 30-second HTTP timeout.
func NewCanvasClient(institutionURL, accessToken string) *CanvasClient {
	return &CanvasClient{
		institutionURL: institutionURL,
		accessToken:    accessToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CanvasClient) doRequest(ctx context.Context, path string, result any) error {
	url := c.institutionURL + path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("doRequest: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("doRequest: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("doRequest: read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return &CanvasAPIError{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("doRequest: json decode: %w", err)
	}

	return nil
}

// GetCourses fetches the authenticated user's active student enrollments.
func (c *CanvasClient) GetCourses(ctx context.Context) ([]CanvasCourse, error) {
	var courses []CanvasCourse
	err := c.doRequest(ctx, "/api/v1/courses?enrollment_type=student&enrollment_state=active&per_page=50", &courses)
	if err != nil {
		return nil, fmt.Errorf("GetCourses: %w", err)
	}
	return courses, nil
}

// GetAssignments fetches all assignments for a given course.
func (c *CanvasClient) GetAssignments(ctx context.Context, courseID string) ([]CanvasAssignment, error) {
	var assignments []CanvasAssignment
	path := fmt.Sprintf("/api/v1/courses/%s/assignments?per_page=100", courseID)
	err := c.doRequest(ctx, path, &assignments)
	if err != nil {
		return nil, fmt.Errorf("GetAssignments: %w", err)
	}
	return assignments, nil
}

// GetSubmissions fetches the authenticated user's submissions for a course.
func (c *CanvasClient) GetSubmissions(ctx context.Context, courseID string) ([]CanvasSubmission, error) {
	var submissions []CanvasSubmission
	path := fmt.Sprintf("/api/v1/courses/%s/students/submissions?student_ids[]=self&per_page=100", courseID)
	err := c.doRequest(ctx, path, &submissions)
	if err != nil {
		return nil, fmt.Errorf("GetSubmissions: %w", err)
	}
	return submissions, nil
}

// GetSyllabusBody fetches the syllabus HTML body for a course.
func (c *CanvasClient) GetSyllabusBody(ctx context.Context, courseID string) (string, error) {
	var course canvasCourseWithSyllabus
	path := fmt.Sprintf("/api/v1/courses/%s?include[]=syllabus_body", courseID)
	err := c.doRequest(ctx, path, &course)
	if err != nil {
		return "", fmt.Errorf("GetSyllabusBody: %w", err)
	}
	return course.SyllabusBody, nil
}

// GetAnnouncements fetches recent announcements for a course.
func (c *CanvasClient) GetAnnouncements(ctx context.Context, courseID string) ([]CanvasAnnouncement, error) {
	var out []CanvasAnnouncement
	path := fmt.Sprintf("/api/v1/courses/%s/announcements?per_page=50", courseID)
	err := c.doRequest(ctx, path, &out)
	if err != nil {
		return nil, fmt.Errorf("GetAnnouncements: %w", err)
	}
	return out, nil
}

// GetDiscussionTopics fetches discussion topics for a course.
func (c *CanvasClient) GetDiscussionTopics(ctx context.Context, courseID string) ([]CanvasDiscussionTopic, error) {
	var out []CanvasDiscussionTopic
	path := fmt.Sprintf("/api/v1/courses/%s/discussion_topics?per_page=50", courseID)
	err := c.doRequest(ctx, path, &out)
	if err != nil {
		return nil, fmt.Errorf("GetDiscussionTopics: %w", err)
	}
	return out, nil
}
