package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vave-tool/internal/domain"
	"google.golang.org/api/idtoken"
)

type authService struct {
	userRepo       domain.UserRepository
	jwtSecret      []byte
	googleClientID string
}

func NewAuthService(userRepo domain.UserRepository, jwtSecret string, googleClientID string) domain.AuthService {
	return &authService{
		userRepo:       userRepo,
		jwtSecret:      []byte(jwtSecret),
		googleClientID: googleClientID,
	}
}

type jwtClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	GoogleID string `json:"google_id"`
}

func (s *authService) AuthenticateWithGoogle(ctx context.Context, googleIDToken string) (string, *domain.User, error) {
	payload, err := idtoken.Validate(ctx, googleIDToken, s.googleClientID)
	if err != nil {
		return "", nil, fmt.Errorf("%w: invalid google token", domain.ErrUnauthorized)
	}

	googleID := payload.Subject
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)
	picture, _ := payload.Claims["picture"].(string)

	if email == "" {
		return "", nil, fmt.Errorf("%w: email not provided by google", domain.ErrUnauthorized)
	}

	user, err := s.userRepo.GetByGoogleID(ctx, googleID)
	if err != nil {
		if err == domain.ErrNotFound {
			user, err = s.userRepo.Create(ctx, &domain.User{
				GoogleID: googleID,
				Email:    email,
				Name:     name,
				Picture:  picture,
			})
			if err != nil {
				return "", nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			return "", nil, fmt.Errorf("failed to find user: %w", err)
		}
	} else if user.Name != name || user.Picture != picture {
		user.Name = name
		user.Picture = picture
		if updated, updateErr := s.userRepo.Update(ctx, user); updateErr == nil {
			user = updated
		}
	}

	if !user.IsActive {
		return "", nil, fmt.Errorf("%w: user account is inactive", domain.ErrForbidden)
	}

	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
		UserID:   user.ID,
		Email:    user.Email,
		Name:     user.Name,
		GoogleID: user.GoogleID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, user, nil
}

func (s *authService) ValidateJWT(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrUnauthorized, err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrUnauthorized
	}

	return &domain.JWTClaims{
		UserID:   claims.UserID,
		Email:    claims.Email,
		Name:     claims.Name,
		GoogleID: claims.GoogleID,
	}, nil
}
