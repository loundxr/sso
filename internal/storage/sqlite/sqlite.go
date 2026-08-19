package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	domain_models "sso/internal/domain/models"
	"sso/internal/storage"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.sqlite.SaveUser"
	query := `
	INSERT INTO users (email, pass_hash)
	VALUES(?,?)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteError sqlite3.Error

		if errors.As(err, &sqliteError) && sqliteError.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserAlreadyExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (domain_models.User, error) {
	const op = "storage.sqlite.User"
	query := `
	SELECT id, email, pass_hash FROM users
	WHERE email=?`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return domain_models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)
	var user domain_models.User

	if err := row.Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain_models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return domain_models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	query := `
	SELECT is_admin FROM users
	WHERE id=?`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)
	var isAdmin bool

	if err := row.Scan(&isAdmin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (domain_models.App, error) {
	const op = "storage.sqlite.App"
	query := `
	SELECT id, name, secret FROM apps
	WHERE id=?`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return domain_models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, appID)
	var app domain_models.App

	if err := row.Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain_models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return domain_models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
