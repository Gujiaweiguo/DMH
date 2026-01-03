package syncadapter

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
)

// SyncTask 同步任务
type SyncTask struct {
	TaskId    string    `json:"task_id"`
	Type      string    `json:"type"` // order, reward
	OrderId   int64     `json:"order_id"`
	RewardId  int64     `json:"reward_id,omitempty"`
	Attempts  int       `json:"attempts"`
	CreatedAt time.Time `json:"created_at"`
}

// SyncQueue 同步队列（基于Redis）
type SyncQueue struct {
	redis  *redis.Client
	key    string
	logger logx.Logger
}

// NewSyncQueue 创建同步队列
func NewSyncQueue(redisClient *redis.Client, queueKey string) *SyncQueue {
	return &SyncQueue{
		redis:  redisClient,
		key:    queueKey,
		logger: logx.WithContext(context.Background()),
	}
}

// Enqueue 将任务加入队列
func (q *SyncQueue) Enqueue(task *SyncTask) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return q.redis.RPush(ctx, q.key, data).Err()
}

// Dequeue 从队列取出任务（阻塞式）
func (q *SyncQueue) Dequeue(timeout time.Duration) (*SyncTask, error) {
	ctx := context.Background()
	result, err := q.redis.BLPop(ctx, timeout, q.key).Result()
	if err != nil {
		return nil, err
	}

	if len(result) < 2 {
		return nil, nil
	}

	var task SyncTask
	if err := json.Unmarshal([]byte(result[1]), &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// Length 获取队列长度
func (q *SyncQueue) Length() (int64, error) {
	ctx := context.Background()
	return q.redis.LLen(ctx, q.key).Result()
}

// Clear 清空队列
func (q *SyncQueue) Clear() error {
	ctx := context.Background()
	return q.redis.Del(ctx, q.key).Err()
}
