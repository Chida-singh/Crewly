package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"crewly/backend/internal/platform/db"
	"crewly/backend/internal/platform/ids"
)

var ErrNotFound = errors.New("not found")

const (
	RoleTalent = "talent"
	RoleClient = "client"
)

// User is the authentication identity.
type User struct {
	ID         string
	Phone      string
	Role       string
	IsVerified bool
}

// OTPRecord is a stored one-time password.
type OTPRecord struct {
	CodeHash  string
	ExpiresAt time.Time
}

type Repo interface {
	CountRecentOTPs(ctx context.Context, phone string) (int64, error)
	CreateOTP(ctx context.Context, phone, codeHash string, expiresAt time.Time) error
	LatestOTP(ctx context.Context, phone string) (OTPRecord, error)
	DeleteOTPs(ctx context.Context, phone string) error
	UserByPhone(ctx context.Context, phone string) (User, error)
	CreateUser(ctx context.Context, phone, role string) (User, error)
}

type PGRepo struct {
	q *db.Queries
}

func NewPGRepo(pool *pgxpool.Pool) *PGRepo {
	return &PGRepo{q: db.New(pool)}
}

func (p *PGRepo) CountRecentOTPs(ctx context.Context, phone string) (int64, error) {
	return p.q.CountRecentOTPs(ctx, phone)
}

func (p *PGRepo) CreateOTP(ctx context.Context, phone, codeHash string, expiresAt time.Time) error {
	_, err := p.q.CreateOTP(ctx, db.CreateOTPParams{
		Phone:     phone,
		CodeHash:  codeHash,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	return err
}

func (p *PGRepo) LatestOTP(ctx context.Context, phone string) (OTPRecord, error) {
	row, err := p.q.GetLatestOTP(ctx, phone)
	if errors.Is(err, pgx.ErrNoRows) {
		return OTPRecord{}, ErrNotFound
	}
	if err != nil {
		return OTPRecord{}, err
	}
	if !row.ExpiresAt.Valid {
		return OTPRecord{}, ErrNotFound
	}
	return OTPRecord{CodeHash: row.CodeHash, ExpiresAt: row.ExpiresAt.Time}, nil
}

func (p *PGRepo) DeleteOTPs(ctx context.Context, phone string) error {
	return p.q.DeleteOTPsByPhone(ctx, phone)
}

func (p *PGRepo) UserByPhone(ctx context.Context, phone string) (User, error) {
	row, err := p.q.GetUserByPhone(ctx, phone)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrNotFound
	}
	if err != nil {
		return User{}, fmt.Errorf("get user by phone: %w", err)
	}
	return toUser(row), nil
}

func (p *PGRepo) CreateUser(ctx context.Context, phone, role string) (User, error) {
	row, err := p.q.CreateUser(ctx, db.CreateUserParams{Phone: phone, Role: role})
	if err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return toUser(row), nil
}

func toUser(u db.User) User {
	return User{
		ID:         ids.FromPG(u.ID),
		Phone:      u.Phone,
		Role:       u.Role,
		IsVerified: u.IsVerified,
	}
}
