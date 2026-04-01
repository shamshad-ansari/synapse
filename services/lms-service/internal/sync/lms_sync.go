package sync

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shamshad-ansari/synapse/services/lms-service/internal/canvas"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/crypto"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/domain"
	"github.com/shamshad-ansari/synapse/services/lms-service/internal/oauth"
)

// Syncer pulls Canvas courses, assignments, and graded submissions into Postgres.
type Syncer struct {
	Repo   domain.LMSRepository
	Logger *zap.Logger
	encKey []byte
}

// NewSyncer returns a Syncer that decrypts OAuth tokens with encryptionKey.
func NewSyncer(repo domain.LMSRepository, logger *zap.Logger, encryptionKey []byte) *Syncer {
	return &Syncer{Repo: repo, Logger: logger, encKey: encryptionKey}
}

type SyncOptions struct {
	InternalURL  string
	CanvasClientID string
	CanvasClientSecret string
}

// SyncUser fetches Canvas data for the connection and upserts into lms_* tables.
// internalURL is cfg.CanvasInternalURL when set; otherwise conn.InstitutionURL is used via ResolveServerURL.
func (s *Syncer) SyncUser(ctx context.Context, conn *domain.LMSConnection, opts SyncOptions) error {
	serverURL := canvas.ResolveServerURL(conn.InstitutionURL, opts.InternalURL)
	runID, runErr := s.Repo.StartSyncRun(ctx, conn.UserID, conn.SchoolID)
	if runErr != nil {
		s.Logger.Warn("sync run start failed", zap.Error(runErr))
	}
	finishRun := func(status string, courses int, errMsg *string) {
		if runID == uuid.Nil {
			return
		}
		if err := s.Repo.FinishSyncRun(ctx, runID, status, courses, errMsg); err != nil {
			s.Logger.Warn("sync run finish failed", zap.Error(err))
		}
	}

	if err := s.ensureFreshToken(ctx, conn, serverURL, opts.CanvasClientID, opts.CanvasClientSecret); err != nil {
		msg := err.Error()
		finishRun("error", 0, &msg)
		return fmt.Errorf("SyncUser: ensureFreshToken: %w", err)
	}
	plain, err := crypto.Decrypt(conn.AccessToken, s.encKey)
	if err != nil {
		msg := err.Error()
		finishRun("error", 0, &msg)
		return fmt.Errorf("SyncUser: decrypt token: %w", err)
	}
	token := string(plain)

	client := canvas.NewCanvasClient(serverURL, token)

	courses, err := client.GetCourses(ctx)
	if err != nil {
		msg := err.Error()
		finishRun("error", 0, &msg)
		return fmt.Errorf("SyncUser: GetCourses: %w", err)
	}

	now := time.Now().UTC()
	userID := conn.UserID
	schoolID := conn.SchoolID

	for _, course := range courses {
		courseIDStr := strconv.Itoa(course.ID)
		lastSync := now
		lmsCourse := &domain.LMSCourse{
			UserID:         userID,
			SchoolID:       schoolID,
			LMSCourseID:    courseIDStr,
			LMSCourseName:  course.Name,
			LMSTerm:        strconv.Itoa(course.EnrollmentTermID),
			EnrollmentType: "student",
			LastSyncedAt:   &lastSync,
		}
		if err := s.Repo.UpsertCourse(ctx, lmsCourse); err != nil {
			msg := err.Error()
			finishRun("error", len(courses), &msg)
			return fmt.Errorf("SyncUser: UpsertCourse %s: %w", courseIDStr, err)
		}

		assignments, err := client.GetAssignments(ctx, courseIDStr)
		if err != nil {
			msg := err.Error()
			finishRun("error", len(courses), &msg)
			return fmt.Errorf("SyncUser: GetAssignments %s: %w", courseIDStr, err)
		}

		pointsByAssignmentID := make(map[int]float64, len(assignments))
		for _, a := range assignments {
			pointsByAssignmentID[a.ID] = a.PointsPossible

			ptsPossible := a.PointsPossible
			lmsA := &domain.LMSAssignment{
				SchoolID:        schoolID,
				LMSAssignmentID: strconv.Itoa(a.ID),
				LMSCourseID:     courseIDStr,
				Title:           a.Name,
				DueAt:           a.DueAt,
				PointsPossible:  &ptsPossible,
				AssignmentGroup: strconv.Itoa(a.AssignmentGroupID),
				LastSyncedAt:    now,
			}
			if err := s.Repo.UpsertAssignment(ctx, lmsA); err != nil {
				msg := err.Error()
				finishRun("error", len(courses), &msg)
				return fmt.Errorf("SyncUser: UpsertAssignment %d: %w", a.ID, err)
			}
		}

		submissions, err := client.GetSubmissions(ctx, courseIDStr)
		if err != nil {
			msg := err.Error()
			finishRun("error", len(courses), &msg)
			return fmt.Errorf("SyncUser: GetSubmissions %s: %w", courseIDStr, err)
		}

		for _, sub := range submissions {
			ptsVal := pointsByAssignmentID[sub.AssignmentID]
			state := &domain.LMSSubmissionState{
				UserID:          userID,
				SchoolID:        schoolID,
				LMSAssignmentID: strconv.Itoa(sub.AssignmentID),
				LMSCourseID:     courseIDStr,
				WorkflowState:   sub.WorkflowState,
				Missing:         sub.Missing,
				Late:            sub.Late,
				Excused:         sub.Excused,
				SubmittedAt:     sub.SubmittedAt,
				GradedAt:        sub.GradedAt,
				Score:           sub.Score,
				PointsPossible:  &ptsVal,
				SyncedAt:        now,
			}
			if err := s.Repo.UpsertSubmissionState(ctx, state); err != nil {
				msg := err.Error()
				finishRun("error", len(courses), &msg)
				return fmt.Errorf("SyncUser: UpsertSubmissionState assignment %d: %w", sub.AssignmentID, err)
			}

			if sub.Score == nil {
				continue
			}
			score := *sub.Score

			ev := &domain.LMSGradeEvent{
				UserID:          userID,
				SchoolID:        schoolID,
				LMSAssignmentID: strconv.Itoa(sub.AssignmentID),
				LMSCourseID:     courseIDStr,
				Score:           &score,
				PointsPossible:  &ptsVal,
				SubmittedAt:     sub.SubmittedAt,
				GradedAt:        sub.GradedAt,
				GradeType:       "assignment",
				SyncedAt:        now,
			}
			if err := s.Repo.UpsertGradeEvent(ctx, ev); err != nil {
				msg := err.Error()
				finishRun("error", len(courses), &msg)
				return fmt.Errorf("SyncUser: UpsertGradeEvent assignment %d: %w", sub.AssignmentID, err)
			}
		}

		announcements, err := client.GetAnnouncements(ctx, courseIDStr)
		if err == nil {
			for _, a := range announcements {
				item := &domain.LMSAnnouncement{
					UserID:            userID,
					SchoolID:          schoolID,
					LMSCourseID:       courseIDStr,
					LMSAnnouncementID: strconv.Itoa(a.ID),
					Title:             a.Title,
					Message:           a.Message,
					PostedAt:          a.PostedAt,
					HTMLURL:           a.HTMLURL,
					LastSyncedAt:      now,
				}
				if upErr := s.Repo.UpsertAnnouncement(ctx, item); upErr != nil {
					msg := upErr.Error()
					finishRun("error", len(courses), &msg)
					return fmt.Errorf("SyncUser: UpsertAnnouncement %d: %w", a.ID, upErr)
				}
			}
		}

		discussions, err := client.GetDiscussionTopics(ctx, courseIDStr)
		if err == nil {
			for _, d := range discussions {
				item := &domain.LMSDiscussionTopic{
					UserID:       userID,
					SchoolID:     schoolID,
					LMSCourseID:  courseIDStr,
					LMSTopicID:   strconv.Itoa(d.ID),
					Title:        d.Title,
					Message:      d.Message,
					PostedAt:     d.PostedAt,
					HTMLURL:      d.HTMLURL,
					LastSyncedAt: now,
				}
				if upErr := s.Repo.UpsertDiscussionTopic(ctx, item); upErr != nil {
					msg := upErr.Error()
					finishRun("error", len(courses), &msg)
					return fmt.Errorf("SyncUser: UpsertDiscussionTopic %d: %w", d.ID, upErr)
				}
			}
		}
	}

	if err := s.Repo.UpdateConnectionSyncStatus(ctx, userID, schoolID, "active", &now); err != nil {
		msg := err.Error()
		finishRun("error", len(courses), &msg)
		return fmt.Errorf("SyncUser: UpdateConnectionSyncStatus: %w", err)
	}

	s.Logger.Info("lms sync completed",
		zap.String("user_id", userID.String()),
		zap.Int("courses_synced", len(courses)),
	)
	finishRun("success", len(courses), nil)

	return nil
}

func (s *Syncer) ensureFreshToken(ctx context.Context, conn *domain.LMSConnection, serverURL, clientID, clientSecret string) error {
	if conn.TokenExpiresAt.After(time.Now().Add(5 * time.Minute)) {
		return nil
	}
	if clientID == "" || clientSecret == "" {
		return nil
	}

	refreshPlain, err := crypto.Decrypt(conn.RefreshToken, s.encKey)
	if err != nil {
		return fmt.Errorf("decrypt refresh token: %w", err)
	}
	refreshToken := string(refreshPlain)
	// Personal token mode has no refresh grant.
	if refreshToken == "" || refreshToken == "personal_token_no_refresh" {
		return nil
	}

	tokenResp, err := oauth.RefreshToken(ctx, serverURL, clientID, clientSecret, refreshToken)
	if err != nil {
		return fmt.Errorf("refresh token: %w", err)
	}

	encAccessToken, err := crypto.Encrypt([]byte(tokenResp.AccessToken), s.encKey)
	if err != nil {
		return fmt.Errorf("encrypt access token: %w", err)
	}
	newRefresh := refreshToken
	if tokenResp.RefreshToken != "" {
		newRefresh = tokenResp.RefreshToken
	}
	encRefreshToken, err := crypto.Encrypt([]byte(newRefresh), s.encKey)
	if err != nil {
		return fmt.Errorf("encrypt refresh token: %w", err)
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	conn.AccessToken = encAccessToken
	conn.RefreshToken = encRefreshToken
	conn.TokenExpiresAt = expiresAt
	if err := s.Repo.UpsertConnection(ctx, conn); err != nil {
		return fmt.Errorf("persist refreshed connection: %w", err)
	}
	return nil
}
