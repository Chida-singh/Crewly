package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type fakeRepo struct {
	recent      int64
	otps        []storedOTP
	latest      map[string]OTPRecord
	users       map[string]User
	createdUser []User
	deleted     []string
}

type storedOTP struct {
	phone     string
	codeHash  string
	expiresAt time.Time
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{latest: map[string]OTPRecord{}, users: map[string]User{}}
}

func (f *fakeRepo) CountRecentOTPs(context.Context, string) (int64, error) {
	return f.recent, nil
}

func (f *fakeRepo) CreateOTP(_ context.Context, phone, codeHash string, expiresAt time.Time) error {
	f.otps = append(f.otps, storedOTP{phone: phone, codeHash: codeHash, expiresAt: expiresAt})
	f.recent++
	return nil
}

func (f *fakeRepo) LatestOTP(_ context.Context, phone string) (OTPRecord, error) {
	rec, ok := f.latest[phone]
	if !ok {
		return OTPRecord{}, ErrNotFound
	}
	return rec, nil
}

func (f *fakeRepo) DeleteOTPs(_ context.Context, phone string) error {
	f.deleted = append(f.deleted, phone)
	return nil
}

func (f *fakeRepo) UserByPhone(_ context.Context, phone string) (User, error) {
	u, ok := f.users[phone]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func (f *fakeRepo) CreateUser(_ context.Context, phone, role string) (User, error) {
	u := User{ID: "user-" + phone, Phone: phone, Role: role}
	f.users[phone] = u
	f.createdUser = append(f.createdUser, u)
	return u, nil
}

type fakeSMS struct {
	phone string
	code  string
	sent  bool
}

func (f *fakeSMS) SendOTP(_ context.Context, phone, code string) error {
	f.phone, f.code, f.sent = phone, code, true
	return nil
}

func newTestService(repo Repo, sender *fakeSMS) *Service {
	svc := NewService(repo, sender, NewTokens("test-secret"))
	fixed := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	svc.now = func() time.Time { return fixed }
	return svc
}

func TestRequestOTPStoresHashAndSendsCode(t *testing.T) {
	repo := newFakeRepo()
	sender := &fakeSMS{}
	svc := newTestService(repo, sender)

	if err := svc.RequestOTP(context.Background(), "9876543210"); err != nil {
		t.Fatalf("RequestOTP: %v", err)
	}

	if !sender.sent || sender.phone != "+919876543210" {
		t.Fatalf("expected sms to +919876543210, got %q (sent=%v)", sender.phone, sender.sent)
	}
	if len(sender.code) != otpDigits {
		t.Fatalf("expected %d digit code, got %q", otpDigits, sender.code)
	}
	if len(repo.otps) != 1 {
		t.Fatalf("expected 1 stored otp, got %d", len(repo.otps))
	}
	stored := repo.otps[0]
	if stored.codeHash == sender.code || strings.Contains(stored.codeHash, sender.code) {
		t.Fatal("otp must be stored hashed")
	}
	if stored.expiresAt.Sub(svc.now()) != otpTTL {
		t.Fatalf("expected ttl %v, got %v", otpTTL, stored.expiresAt.Sub(svc.now()))
	}
}

func TestRequestOTPRateLimited(t *testing.T) {
	repo := newFakeRepo()
	repo.recent = maxOTPsPer10Min
	svc := newTestService(repo, &fakeSMS{})

	err := svc.RequestOTP(context.Background(), "9876543210")
	if !errors.Is(err, ErrRateLimited) {
		t.Fatalf("expected ErrRateLimited, got %v", err)
	}
	if len(repo.otps) != 0 {
		t.Fatal("rate limited request must not store an otp")
	}
}

func TestRequestOTPRejectsInvalidPhone(t *testing.T) {
	svc := newTestService(newFakeRepo(), &fakeSMS{})
	if err := svc.RequestOTP(context.Background(), "123"); !errors.Is(err, ErrInvalidPhone) {
		t.Fatalf("expected ErrInvalidPhone, got %v", err)
	}
}

func TestVerifyCreatesNewUserAndIssuesToken(t *testing.T) {
	repo := newFakeRepo()
	svc := newTestService(repo, &fakeSMS{})
	phone := "+919876543210"
	repo.latest[phone] = OTPRecord{CodeHash: svc.hashCode(phone, "123456"), ExpiresAt: svc.now().Add(time.Minute)}

	user, token, err := svc.VerifyOTP(context.Background(), "9876543210", "123456", "talent")
	if err != nil {
		t.Fatalf("VerifyOTP: %v", err)
	}
	if user.Role != RoleTalent {
		t.Fatalf("expected talent user, got %q", user.Role)
	}
	if len(repo.createdUser) != 1 {
		t.Fatalf("expected user creation, got %d", len(repo.createdUser))
	}
	if len(repo.deleted) != 1 {
		t.Fatal("otp should be consumed after a successful verify")
	}

	parsed, err := svc.tokens.Parse(token)
	if err != nil {
		t.Fatalf("parse issued token: %v", err)
	}
	if parsed.ID != user.ID || parsed.Role != RoleTalent {
		t.Fatalf("token claims mismatch: %+v", parsed)
	}
}

func TestVerifyRejectsWrongAndExpiredCode(t *testing.T) {
	phone := "+919876543210"

	tests := []struct {
		name      string
		wrongHash bool
		expired   bool
	}{
		{name: "wrong code", wrongHash: true},
		{name: "expired code", expired: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := newFakeRepo()
			svc := newTestService(repo, &fakeSMS{})

			record := OTPRecord{
				CodeHash:  svc.hashCode(phone, "123456"),
				ExpiresAt: svc.now().Add(time.Minute),
			}
			if tc.wrongHash {
				record.CodeHash = "deadbeef"
			}
			if tc.expired {
				record.ExpiresAt = svc.now().Add(-time.Minute)
			}
			repo.latest[phone] = record

			_, _, err := svc.VerifyOTP(context.Background(), phone, "123456", "talent")
			if !errors.Is(err, ErrInvalidCode) {
				t.Fatalf("expected ErrInvalidCode, got %v", err)
			}
			if len(repo.createdUser) != 0 {
				t.Fatal("must not create a user for an invalid code")
			}
		})
	}
}

func TestVerifyRequiresValidRoleForNewUser(t *testing.T) {
	repo := newFakeRepo()
	svc := newTestService(repo, &fakeSMS{})
	phone := "+919876543210"
	repo.latest[phone] = OTPRecord{CodeHash: svc.hashCode(phone, "123456"), ExpiresAt: svc.now().Add(time.Minute)}

	if _, _, err := svc.VerifyOTP(context.Background(), phone, "123456", "admin"); !errors.Is(err, ErrInvalidRole) {
		t.Fatalf("expected ErrInvalidRole, got %v", err)
	}
}

func TestVerifyKeepsExistingRole(t *testing.T) {
	repo := newFakeRepo()
	svc := newTestService(repo, &fakeSMS{})
	phone := "+919876543210"
	repo.users[phone] = User{ID: "existing", Phone: phone, Role: RoleClient}
	repo.latest[phone] = OTPRecord{CodeHash: svc.hashCode(phone, "654321"), ExpiresAt: svc.now().Add(time.Minute)}

	user, _, err := svc.VerifyOTP(context.Background(), phone, "654321", "talent")
	if err != nil {
		t.Fatalf("VerifyOTP: %v", err)
	}
	if user.Role != RoleClient || user.ID != "existing" {
		t.Fatalf("existing user should be returned unchanged, got %+v", user)
	}
	if len(repo.createdUser) != 0 {
		t.Fatal("existing user must not be recreated")
	}
}

func TestNormalizePhone(t *testing.T) {
	cases := map[string]string{
		"9876543210":      "+919876543210",
		"09876543210":     "+919876543210",
		"919876543210":    "+919876543210",
		"+91 98765 43210": "+919876543210",
	}
	for input, want := range cases {
		got, err := NormalizePhone(input)
		if err != nil || got != want {
			t.Fatalf("NormalizePhone(%q) = %q, %v; want %q", input, got, err, want)
		}
	}
	for _, bad := range []string{"", "123", "5876543210", "+91987654321"} {
		if _, err := NormalizePhone(bad); !errors.Is(err, ErrInvalidPhone) {
			t.Fatalf("expected ErrInvalidPhone for %q, got %v", bad, err)
		}
	}
}
