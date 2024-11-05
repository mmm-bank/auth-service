package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mmm-bank/auth-service/models"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var _ UserRepo = UserPostgres{}

type UserRepo interface {
	AddUser(user *models.User) error
	AuthUser(user *models.User) error
}

type UserPostgres struct {
	db *pgxpool.Pool
}

func NewUserPostgres(connString string) UserPostgres {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	return UserPostgres{pool}
}

func (u UserPostgres) AddUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	query := "INSERT INTO users (user_id, username, hashed_password) VALUES ($1, $2, $3)"
	_, err = u.db.Exec(context.Background(), query, uuid.New(), user.Username, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to add user: %v", err)
	}
	return nil
}

func (u UserPostgres) AuthUser(user *models.User) error {
	query := "SELECT user_id, hashed_password FROM users WHERE username = $1"

	var hashedPassword string
	err := u.db.QueryRow(context.Background(), query, user.Username).Scan(&user.ID, &hashedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to authenticate user: %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		return fmt.Errorf("invalid credentials")
	}
	return nil
}
