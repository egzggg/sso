package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"sso/internal/domain/models"
	"sso/internal/storage"
)

type Storage struct {
	db *pgxpool.Pool
}

//func New(storagePath string) (*Storage, error) {
//	const op = "storage.sqlite.New"
//
//	db, err := sql.Open("sqlite3", storagePath)
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	return &Storage{db: db}, nil
//}

func New(ctx context.Context, databaseURL string) (*Storage, error) {
	const op = "storage.postgres.New"
	pool, err := pgxpool.New(ctx, databaseURL) // создаем пул подключения к базе
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := pool.Ping(ctx); err != nil { // проверяем есть ли реальное соединение с базой
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: pool}, nil // возвращаем наше соединение
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	var id int64

	err := s.db.QueryRow(ctx, //QueryRow отправляет скл запрос и возвращает айдишник
		`INSERT INTO users(email, pass_hash) VALUES($1, $2) RETURNING id`,
		email, passHash,
	).Scan(&id) // записывает айд который вернулся из QueryRow в переменную

	if err != nil {
		var pgErr *pgconn.PgError   // переменнай с типом кастомной ошибки постгреса
		if errors.As(err, &pgErr) { // проверяем является ли эта ошбка постгреса
			if pgErr.Code == "23505" { // уникальный email уже есть
				return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err) // если не ошибака постгреса взвращаем любую другую ощибку
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	var user models.User

	err := s.db.QueryRow(ctx,
		`SELECT id, email, pass_hash FROM users WHERE email = $1`,
		email).Scan(&user.ID, &user.Email, &user.PassHash)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	var isAdmin bool

	err := s.db.QueryRow(ctx, `SELECT is_admin FROM users WHERE id = $1`, userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil

}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	//panic("not implemented")
	const op = "storage.postgres.App"

	var app models.App

	err := s.db.QueryRow(ctx, `SELECT id, name, secret FROM apps WHERE id = $1`, appID).Scan(&app.ID, &app.Name, &app.Secret)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}

func (s *Storage) Pool() *pgxpool.Pool {
	return s.db
}
