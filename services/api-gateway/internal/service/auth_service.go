package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"time"
	"unicode"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/config"
	"github.com/shamshad-ansari/synapse/services/api-gateway/internal/domain"
)

type RegisterInput struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	SchoolDomain string `json:"school_domain"`
}

type LoginInput struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	SchoolDomain string `json:"school_domain"`
}

type AuthOutput struct {
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
	User         *domain.UserResponse `json:"user"`
}

// AuthService defines the business operations for authentication.
type AuthService interface {
	Register(ctx context.Context, in RegisterInput) (*AuthOutput, error)
	Login(ctx context.Context, in LoginInput) (*AuthOutput, error)
	GetCurrentUser(ctx context.Context, userID, schoolID uuid.UUID) (*domain.User, error)
}

type authServiceImpl struct {
	repo domain.UserRepository
	cfg  *config.Config
}

func NewAuthService(repo domain.UserRepository, cfg *config.Config) AuthService {
	return &authServiceImpl{repo: repo, cfg: cfg}
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func (s *authServiceImpl) Register(ctx context.Context, in RegisterInput) (*AuthOutput, error) {
	if err := validateRegisterInput(in); err != nil {
		return nil, err
	}

	school, err := s.repo.CreateSchool(ctx, domainToSchoolName(in.SchoolDomain), in.SchoolDomain)
	if err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("Register: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, school.ID, in.Name, in.Email, string(hash))
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			return nil, &domain.ValidationError{Message: "email already registered at this school"}
		}
		return nil, fmt.Errorf("Register: %w", err)
	}

	return s.generateTokens(user)
}

func (s *authServiceImpl) Login(ctx context.Context, in LoginInput) (*AuthOutput, error) {
	school, err := s.repo.FindSchoolByDomain(ctx, in.SchoolDomain)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("Login: %w", domain.ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("Login: %w", err)
	}

	user, err := s.repo.FindUserByEmail(ctx, in.Email, school.ID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("Login: %w", domain.ErrInvalidCredentials)
		}
		return nil, fmt.Errorf("Login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return nil, fmt.Errorf("Login: %w", domain.ErrInvalidCredentials)
	}

	return s.generateTokens(user)
}

func (s *authServiceImpl) GetCurrentUser(ctx context.Context, userID, schoolID uuid.UUID) (*domain.User, error) {
	user, err := s.repo.FindUserByID(ctx, userID, schoolID)
	if err != nil {
		return nil, fmt.Errorf("GetCurrentUser: %w", err)
	}
	return user, nil
}

func (s *authServiceImpl) generateTokens(user *domain.User) (*AuthOutput, error) {
	now := time.Now()

	accessClaims := jwt.MapClaims{
		"user_id":   user.ID.String(),
		"school_id": user.SchoolID.String(),
		"email":     user.Email,
		"type":      "access",
		"exp":       now.Add(time.Duration(s.cfg.JWTExpiryHours) * time.Hour).Unix(),
		"iat":       now.Unix(),
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("generateTokens: access: %w", err)
	}

	refreshClaims := jwt.MapClaims{
		"user_id":   user.ID.String(),
		"school_id": user.SchoolID.String(),
		"type":      "refresh",
		"exp":       now.Add(time.Duration(s.cfg.JWTRefreshExpiryDays) * 24 * time.Hour).Unix(),
		"iat":       now.Unix(),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("generateTokens: refresh: %w", err)
	}

	return &AuthOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	}, nil
}

func validateRegisterInput(in RegisterInput) error {
	if in.Name == "" {
		return &domain.ValidationError{Message: "name is required"}
	}
	if !emailRegex.MatchString(in.Email) {
		return &domain.ValidationError{Message: "invalid email format"}
	}
	if err := validatePassword(in.Password); err != nil {
		return err
	}
	if in.SchoolDomain == "" {
		return &domain.ValidationError{Message: "school_domain is required"}
	}
	if _, err := net.LookupHost(in.SchoolDomain); err != nil {
		// Fallback: check it at least looks like a hostname (contains a dot)
		if !isValidHostname(in.SchoolDomain) {
			return &domain.ValidationError{Message: "school_domain must be a valid hostname"}
		}
	}
	return nil
}

func validatePassword(pw string) error {
	if len(pw) < 8 {
		return &domain.ValidationError{Message: "password must be at least 8 characters"}
	}
	hasDigit := false
	for _, ch := range pw {
		if unicode.IsDigit(ch) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return &domain.ValidationError{Message: "password must contain at least 1 digit"}
	}
	return nil
}

func isValidHostname(h string) bool {
	for i, ch := range h {
		if ch == '.' && i > 0 && i < len(h)-1 {
			return true
		}
	}
	return false
}

// domainToSchoolName produces a default school name from its domain.
func domainToSchoolName(d string) string {
	return d
}
