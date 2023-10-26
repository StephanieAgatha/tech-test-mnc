package usecase

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestRedisSetGet(t *testing.T) {
	var ctx = context.Background()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", //isi jika ada pw nya
		DB:       0,
	})

	//set some key val
	err := client.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	// get key and val
	value, err := client.Get(ctx, "key").Result()
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}

	//cek jika tidak sama
	if value != "value" {
		t.Fatalf("Value is not the same: got %v, want %v", value, "value")
	}

	fmt.Println("Test passed")
}
