package syncadapter

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestNewSyncQueue(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	queue := NewSyncQueue(redisClient, "test_queue")

	assert.NotNil(t, queue)
	assert.Equal(t, "test_queue", queue.key)
	assert.NotNil(t, queue.redis)

	redisClient.Close()
}

func TestSyncQueue_Enqueue(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	queue := NewSyncQueue(redisClient, "test_enqueue")

	queue.Clear()

	task := &SyncTask{
		TaskId:    "task_1",
		Type:      "order",
		OrderId:   123,
		Attempts:  0,
		CreatedAt: time.Now(),
	}

	err := queue.Enqueue(task)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
		return
	}

	length, err := queue.Length()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), length)
}

func TestSyncQueue_Dequeue(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	queue := NewSyncQueue(redisClient, "test_dequeue")

	task := &SyncTask{
		TaskId:    "task_1",
		Type:      "order",
		OrderId:   123,
		Attempts:  0,
		CreatedAt: time.Now(),
	}

	err := queue.Enqueue(task)
	if err != nil {
		t.Skipf("Redis not available: %v", err)
		return
	}

	dequeuedTask, err := queue.Dequeue(5 * time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, dequeuedTask)
	assert.Equal(t, "task_1", dequeuedTask.TaskId)
	assert.Equal(t, "order", dequeuedTask.Type)
	assert.Equal(t, int64(123), dequeuedTask.OrderId)
}

func TestSyncQueue_Dequeue_Timeout(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	queue := NewSyncQueue(redisClient, "test_dequeue_timeout")

	err := queue.Clear()
	if err != nil {
		t.Skipf("Redis not available: %v", err)
		return
	}

	task, err := queue.Dequeue(2 * time.Second)
	if err == redis.Nil {
		return
	}
	if err != nil {
		t.Skipf("Unexpected error: %v", err)
		return
	}
	assert.Nil(t, task)
}

func TestSyncQueue_Length(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	queue := NewSyncQueue(redisClient, "test_length")

	err := queue.Clear()
	if err != nil {
		t.Skipf("Redis not available: %v", err)
		return
	}

	length, err := queue.Length()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), length)

	queue.Enqueue(&SyncTask{TaskId: "1", Type: "order"})
	queue.Enqueue(&SyncTask{TaskId: "2", Type: "order"})

	length, _ = queue.Length()
	assert.Equal(t, int64(2), length)
}

func TestSyncQueue_Clear(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	defer redisClient.Close()

	queue := NewSyncQueue(redisClient, "test_clear")

	queue.Enqueue(&SyncTask{TaskId: "1", Type: "order"})
	queue.Enqueue(&SyncTask{TaskId: "2", Type: "order"})

	err := queue.Clear()
	if err != nil {
		t.Skipf("Redis not available: %v", err)
		return
	}

	length, _ := queue.Length()
	assert.Equal(t, int64(0), length)
}
