package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	*redis.Client
}

func NewClient(addr string, passwd string, db int) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})
	return &Cache{client}
}

func (c *Cache) Ping(ctx context.Context) error {
	return c.Client.Ping(ctx).Err()
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	return c.Client.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *Cache) SetStudentsWithScore(ctx context.Context, students []Student) {
	for _, student := range students {
		c.Client.ZIncrBy(ctx, "student_rank", float64(student.Total), student.Name)
	}
}

func (c *Cache) GetStudentsByRank(ctx context.Context) ([]string, error) {
	var result []string
	rankings, err := c.Client.ZRevRangeByScore(ctx, "student_rank", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()

	if err != nil {
		return result, err
	}

	for _, student := range rankings {
		result = append(result, student)
	}
	return result, nil
}
