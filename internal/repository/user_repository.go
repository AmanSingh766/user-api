package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	db "github.com/dinesh/user-api/db/sqlc"
)

// ErrNotFound is returned when a record doesn't exist
var ErrNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, name string, dob time.Time) (db.User, error)
	GetByID(ctx context.Context, id int32) (db.User, error)
	Update(ctx context.Context, id int32, name string, dob time.Time) (db.User, error)
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context, limit, offset int32) ([]db.User, error)
	Count(ctx context.Context) (int64, error)
}

type userRepo struct {
	q *db.Queries
}

func NewUserRepository(database *sql.DB) UserRepository {
	return &userRepo{q: db.New(database)}
}

func (r *userRepo) Create(ctx context.Context, name string, dob time.Time) (db.User, error) {
	return r.q.CreateUser(ctx, name, dob)
}

func (r *userRepo) GetByID(ctx context.Context, id int32) (db.User, error) {
	u, err := r.q.GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return db.User{}, ErrNotFound
	}
	return u, err
}

func (r *userRepo) Update(ctx context.Context, id int32, name string, dob time.Time) (db.User, error) {
	u, err := r.q.UpdateUser(ctx, name, dob, id)
	if errors.Is(err, sql.ErrNoRows) {
		return db.User{}, ErrNotFound
	}
	return u, err
}

func (r *userRepo) Delete(ctx context.Context, id int32) error {
	// Check existence first so we can distinguish 404 from actual errors
	_, err := r.q.GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	return r.q.DeleteUser(ctx, id)
}

func (r *userRepo) List(ctx context.Context, limit, offset int32) ([]db.User, error) {
	return r.q.ListUsers(ctx, limit, offset)
}

func (r *userRepo) Count(ctx context.Context) (int64, error) {
	return r.q.CountUsers(ctx)
}
