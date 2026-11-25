package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"crewly/backend/internal/platform/sms"
)

var (
	ErrInvalidPhone = errors.New("invalid phone number")
	ErrRateLimited  = errors.New("too many otp requests, try again later")
	ErrInvalidCode  = errors.New("invalid or expired code")
	ErrInvalidRole  = errors.New("role must be talent or client")
)

const (
	otpTTL          = 5 * time.Minute
	maxOTPsPer10Min = 3
	otpDigits       = 6
)

var phonePattern = regexp.MustCompile(`^\+91[6-9]\d{9}$`)

// Service holds the OTP and login business rules.
type Service struct {
	repo   Repo
	sms    sms.Sender
	tokens *Tokens
	now    func() time.Time
}

func NewService(repo Repo, sender sms.Sender, tokens *Tokens) *Service {
	return &Service{repo: repo, sms: sender, tokens: tokens, now: time.Now}
}

// RequestOTP generates a 6 digit code, stores only its hash, and sends it.
func (s *Service) RequestOTP(ctx context.Context, phone string) error {
	normalized, err := NormalizePhone(phone)
	if err != nil {
		return err
	}

	recent, err := s.repo.CountRecentOTPs(ctx, normalized)
	if err != nil {
		return fmt.Errorf("count recent otps: %w", err)
	}
	if recent >= maxOTPsPer10Min {
		return ErrRateLimited
	}

	code, err := generateCode()
	if err != nil {
		return fmt.Errorf("generate otp: %w", err)
	}

	expiresAt := s.now().Add(otpTTL)
	if err := s.repo.CreateOTP(ctx, normalized, s.hashCode(normalized, code), expiresAt); err != nil {
		return fmt.Errorf("store otp: %w", err)
	}
	if err := s.sms.SendOTP(ctx, normalized, code); err != nil {
		return fmt.Errorf("send otp: %w", err)
	}
	return nil
}

// VerifyOTP checks the latest code, creates the user on first login, and issues a JWT.
// role is only used when the phone has no user yet.
func (s *Service) VerifyOTP(ctx context.Context, phone, code, role string) (User, string, error) {
	normalized, err := NormalizePhone(phone)
	if err != nil {
		return User{}, "", err
	}
	if !isSixDigits(code) {
		return User{}, "", ErrInvalidCode
	}

	record, err := s.repo.LatestOTP(ctx, normalized)
	if errors.Is(err, ErrNotFound) {
		return User{}, "", ErrInvalidCode
	}
	if err != nil {
		return User{}, "", fmt.Errorf("load otp: %w", err)
	}
	if s.now().After(record.ExpiresAt) {
		return User{}, "", ErrInvalidCode
	}
	if subtle.ConstantTimeCompare([]byte(record.CodeHash), []byte(s.hashCode(normalized, code))) != 1 {
		return User{}, "", ErrInvalidCode
	}

	if err := s.repo.DeleteOTPs(ctx, normalized); err != nil {
		return User{}, "", fmt.Errorf("consume otp: %w", err)
	}

	user, err := s.repo.UserByPhone(ctx, normalized)
	if errors.Is(err, ErrNotFound) {
		role = strings.ToLower(strings.TrimSpace(role))
		if role != RoleTalent && role != RoleClient {
			return User{}, "", ErrInvalidRole
		}
		user, err = s.repo.CreateUser(ctx, normalized, role)
	}
	if err != nil {
		return User{}, "", fmt.Errorf("resolve user: %w", err)
	}

	token, err := s.tokens.Issue(user.ID, user.Role)
	if err != nil {
		return User{}, "", err
	}
	return user, token, nil
}

func (s *Service) hashCode(phone, code string) string {
	sum := sha256.Sum256([]byte(phone + ":" + code + ":" + s.tokens.pepper()))
	return hex.EncodeToString(sum[:])
}

func generateCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", otpDigits, n.Int64()), nil
}

func isSixDigits(code string) bool {
	if len(code) != otpDigits {
		return false
	}
	for _, r := range code {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// NormalizePhone accepts Indian numbers in common forms and returns +91XXXXXXXXXX.
func NormalizePhone(phone string) (string, error) {
	trimmed := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "").Replace(strings.TrimSpace(phone))

	switch {
	case strings.HasPrefix(trimmed, "+91"):
		// already E.164
	case strings.HasPrefix(trimmed, "91") && len(trimmed) == 12:
		trimmed = "+" + trimmed
	case strings.HasPrefix(trimmed, "0") && len(trimmed) == 11:
		trimmed = "+91" + trimmed[1:]
	case len(trimmed) == 10:
		trimmed = "+91" + trimmed
	default:
		return "", ErrInvalidPhone
	}

	if !phonePattern.MatchString(trimmed) {
		return "", ErrInvalidPhone
	}
	return trimmed, nil
}
