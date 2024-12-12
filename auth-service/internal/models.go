package internal

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Models struct {
	db   *pgxpool.Pool
	User User
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		db:   db,
		User: User{},
	}
}

type User struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Password  string     `json:"-"`
	Active    int        `json:"active"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (m *Models) PasswordMatches(password string) (bool, error) {
	if password == "" {
		return false, nil
	}

	//TODO: implement password check, for now just return true
	return true, nil

	if err := bcrypt.CompareHashAndPassword([]byte(m.User.Password), []byte(password)); err != nil {
		slog.Info("Passwords do not match")
		return false, fmt.Errorf("passwords do not match")
	}

	return true, nil
}

func (m *Models) GeneratePassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password is required")
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passHash), nil
}

func (m *Models) GetAll(ctx context.Context) ([]*User, error) {
	query, args, err := squirrel.Select("id", "email", "first_name", "last_name", "password_hash", "user_active", "created_at", "updated_at").
		From("users").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := m.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var users []*User
	for rows.Next() {
		var user User
		if err = rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.Active, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (m *Models) GetByEmail(ctx context.Context, email string) (*User, error) {
	query, args, err := squirrel.Select("id", "email", "first_name", "last_name", "password_hash", "user_active", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user User
	err = m.db.QueryRow(ctx, query, args...).
		Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.Active, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		slog.Error("failed to get user by email", "ERROR", err)
		return nil, err
	}

	return &user, nil
}

func (m *Models) GetByID(ctx context.Context, id int) (*User, error) {
	query, args, err := squirrel.Select("id", "email", "first_name", "last_name", "password_hash", "user_active", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user User
	err = m.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Password, &user.Active, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *Models) CreateUser(ctx context.Context, user User) error {
	query, args, err := squirrel.Insert("users").
		Columns("email", "first_name", "last_name", "password_hash", "user_active").
		Values(user.Email, user.FirstName, user.LastName, user.Password, user.Active).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = m.db.Exec(ctx, query, args...)
	return err
}

func (m *Models) UpdateUser(ctx context.Context, user User) error {
	query, args, err := squirrel.Update("users").SetMap(map[string]any{
		"email":         user.Email,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"password_hash": user.Password,
		"user_active":   user.Active,
	}).Where(squirrel.Eq{"id": user.ID}).ToSql()
	if err != nil {
		return err
	}

	_, err = m.db.Exec(ctx, query, args...)
	return err
}

func (m *Models) DeleteUser(ctx context.Context, id int) error {
	query, args, err := squirrel.Delete("users").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return err
	}

	_, err = m.db.Exec(ctx, query, args...)
	return err
}
