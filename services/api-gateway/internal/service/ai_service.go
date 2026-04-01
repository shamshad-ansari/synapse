package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/ai"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

// AIService runs RAG: embed query, retrieve similar cards, generate candidates, persist accepted cards with async embedding.
type AIService struct {
	Repo   domain.LearningRepository
	Gen    ai.FlashcardGenerator
	Embed  ai.TextEmbedder
	Logger *zap.Logger
}

// NewAIService constructs an AIService. Gen and Embed must be non-nil when used.
func NewAIService(repo domain.LearningRepository, gen ai.FlashcardGenerator, emb ai.TextEmbedder, log *zap.Logger) *AIService {
	return &AIService{
		Repo:   repo,
		Gen:    gen,
		Embed:  emb,
		Logger: log,
	}
}

// GenerateFlashcardsFromNote embeds note text, retrieves similar flashcards, and returns LLM candidates (not persisted).
func (s *AIService) GenerateFlashcardsFromNote(
	ctx context.Context,
	userID, schoolID, courseID uuid.UUID,
	noteContent, courseName string,
	topicID *uuid.UUID,
) ([]ai.GeneratedCard, error) {
	_ = topicID // reserved for future filtering
	if _, err := s.Repo.GetCourse(ctx, courseID, userID, schoolID); err != nil {
		return nil, fmt.Errorf("GenerateFlashcardsFromNote: %w", err)
	}
	queryVec, err := s.Embed.Embed(ctx, ai.EmbedInput{Text: noteContent})
	if err != nil {
		return nil, fmt.Errorf("GenerateFlashcardsFromNote: embed: %w", err)
	}

	similar, err := s.Repo.SearchSimilarFlashcards(ctx, userID, schoolID, queryVec, 5)
	if err != nil {
		return nil, fmt.Errorf("GenerateFlashcardsFromNote: search: %w", err)
	}
	existing := make([]ai.CardContext, 0, len(similar))
	for _, f := range similar {
		existing = append(existing, ai.CardContext{Prompt: f.Prompt, Answer: f.Answer})
	}

	cards, err := s.Gen.GenerateFlashcards(ctx, ai.GenerateFlashcardsInput{
		NoteContent:   noteContent,
		CourseContext: courseName,
		ExistingCards: existing,
		MaxCards:      6,
	})
	if err != nil {
		return nil, fmt.Errorf("GenerateFlashcardsFromNote: generate: %w", err)
	}
	return cards, nil
}

// SaveAcceptedCards persists accepted AI cards and schedules async embedding per card.
func (s *AIService) SaveAcceptedCards(
	ctx context.Context,
	userID, schoolID, courseID uuid.UUID,
	cards []ai.GeneratedCard,
	topicID *uuid.UUID,
) ([]*domain.Flashcard, error) {
	if _, err := s.Repo.GetCourse(ctx, courseID, userID, schoolID); err != nil {
		return nil, fmt.Errorf("SaveAcceptedCards: %w", err)
	}
	if topicID != nil {
		topics, err := s.Repo.ListTopics(ctx, courseID, userID, schoolID)
		if err != nil {
			return nil, fmt.Errorf("SaveAcceptedCards: %w", err)
		}
		ok := false
		for _, t := range topics {
			if t.ID == *topicID {
				ok = true
				break
			}
		}
		if !ok {
			return nil, &domain.ValidationError{Message: "topic not found in this course"}
		}
	}

	domainCards := make([]*domain.Flashcard, 0, len(cards))
	for _, c := range cards {
		if c.Prompt == "" || c.Answer == "" {
			continue
		}
		ct := c.CardType
		if ct == "" {
			ct = "qa"
		}
		domainCards = append(domainCards, &domain.Flashcard{
			SchoolID:   schoolID,
			CourseID:   courseID,
			UserID:     userID,
			TopicID:    topicID,
			CardType:   ct,
			Prompt:     c.Prompt,
			Answer:     c.Answer,
			CreatedBy:  "ai",
			Visibility: "private",
		})
	}
	if len(domainCards) == 0 {
		return nil, &domain.ValidationError{Message: "no valid cards to save"}
	}

	saved, err := s.Repo.BatchCreateFlashcards(ctx, domainCards)
	if err != nil {
		return nil, fmt.Errorf("SaveAcceptedCards: %w", err)
	}

	now := time.Now()
	for _, fc := range saved {
		state := &domain.SchedulerState{
			FlashcardID:  fc.ID,
			UserID:       userID,
			SchoolID:     schoolID,
			EaseFactor:   2.5,
			IntervalDays: 1,
			DueAt:        now,
			LapseCount:   0,
			LastReviewAt: nil,
		}
		if err := s.Repo.UpsertSchedulerState(ctx, state); err != nil {
			return nil, fmt.Errorf("SaveAcceptedCards: scheduler: %w", err)
		}
	}

	for _, card := range saved {
		c := card
		go func() {
			bg, cancel := context.WithTimeout(context.Background(), 90*time.Second)
			defer cancel()
			vec, err := s.Embed.Embed(bg, ai.EmbedInput{Text: c.Prompt + "\n" + c.Answer})
			if err != nil {
				if s.Logger != nil {
					s.Logger.Warn("async flashcard embed failed", zap.String("flashcard_id", c.ID.String()), zap.Error(err))
				}
				return
			}
			if err := s.Repo.UpdateFlashcardEmbedding(bg, c.ID, userID, schoolID, vec); err != nil && s.Logger != nil {
				s.Logger.Warn("async flashcard embedding update failed", zap.String("flashcard_id", c.ID.String()), zap.Error(err))
			}
		}()
	}

	return saved, nil
}
