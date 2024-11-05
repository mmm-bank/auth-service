package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/mmm-bank/auth-service/models"
	"log"
	"time"
)

var _ SessionRepo = SessionRedis{}

type SessionRepo interface {
	AddSession(session models.Session) error
}

type SessionRedis struct {
	client *redis.Client
}

func NewSessionRedis(redisAddr string) SessionRedis {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}

	return SessionRedis{client}
}

func (s SessionRedis) AddSession(session models.Session) error {
	ctx := context.Background()

	ttl := session.ExpiresAt.Sub(time.Now())
	if ttl <= 0 {
		return fmt.Errorf("invalid session expiration time")
	}

	err := s.client.Set(ctx, session.Token, session.UserID, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to add session to redis: %v", err)
	}

	return nil
}
