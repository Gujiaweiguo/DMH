package syncadapter

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"dmh/model"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// SyncWorker 同步Worker
type SyncWorker struct {
	adapter *SyncAdapter
	queue   *SyncQueue
	db      *gorm.DB
	logger  logx.Logger
	stop    chan struct{}
}

// NewSyncWorker 创建同步Worker
func NewSyncWorker(adapter *SyncAdapter, queue *SyncQueue, db *gorm.DB) *SyncWorker {
	return &SyncWorker{
		adapter: adapter,
		queue:   queue,
		db:      db,
		logger:  logx.WithContext(context.Background()),
		stop:    make(chan struct{}),
	}
}

// Start 启动Worker
func (w *SyncWorker) Start() {
	w.logger.Info("SyncWorker started")

	for {
		select {
		case <-w.stop:
			w.logger.Info("SyncWorker stopped")
			return
		default:
			w.processTask()
		}
	}
}

// Stop 停止Worker
func (w *SyncWorker) Stop() {
	close(w.stop)
}

// processTask 处理任务
func (w *SyncWorker) processTask() {
	// 1. 从队列取任务（阻塞1秒）
	task, err := w.queue.Dequeue(1 * time.Second)
	if err != nil {
		if err.Error() != "redis: nil" {
			w.logger.Error("dequeue failed:", err)
		}
		return
	}

	if task == nil {
		return
	}

	w.logger.Infof("processing task: %s, type: %s, order_id: %d", task.TaskId, task.Type, task.OrderId)

	// 2. 执行同步
	var syncErr error
	switch task.Type {
	case "order":
		syncErr = w.syncOrder(task.OrderId)
	case "reward":
		syncErr = w.syncReward(task.RewardId)
	default:
		w.logger.Errorf("unknown task type: %s", task.Type)
		return
	}

	// 3. 更新同步状态
	if syncErr != nil {
		w.logger.Errorf("sync failed: %v", syncErr)
		w.updateSyncStatus(task.OrderId, "failed", syncErr.Error())

		// 重试逻辑（最多3次）
		if task.Attempts < 3 {
			task.Attempts++
			w.logger.Infof("retry task: %s, attempt: %d", task.TaskId, task.Attempts)
			w.queue.Enqueue(task) // 重新入队
		} else {
			w.logger.Errorf("task failed after 3 attempts: %s", task.TaskId)
		}
	} else {
		w.updateSyncStatus(task.OrderId, "synced", "")
		w.logger.Infof("task completed: %s", task.TaskId)
	}
}

// syncOrder 同步订单
func (w *SyncWorker) syncOrder(orderId int64) error {
	// 1. 查询订单数据
	var order model.Order
	if err := w.db.Where("id = ?", orderId).First(&order).Error; err != nil {
		return fmt.Errorf("get order failed: %w", err)
	}

	// 2. 解析表单数据
	var formData map[string]interface{}
	if err := json.Unmarshal([]byte(order.FormData), &formData); err != nil {
		formData = make(map[string]interface{})
	}

	// 3. 构建同步数据
	syncData := &SyncOrderData{
		OrderId:    order.Id,
		CampaignId: order.CampaignId,
		MemberId:   order.MemberId,   // 会员ID
		UnionID:    order.UnionID,    // 微信 unionid
		Phone:      order.Phone,
		FormData:   formData,
		Amount:     order.Amount,
		PayStatus:  order.PayStatus,
		CreatedAt:  order.CreatedAt,
	}

	// 4. 执行同步
	return w.adapter.SyncOrder(context.Background(), syncData)
}

// syncReward 同步奖励
func (w *SyncWorker) syncReward(rewardId int64) error {
	// 1. 查询奖励数据
	var reward model.Reward
	if err := w.db.Where("id = ?", rewardId).First(&reward).Error; err != nil {
		return fmt.Errorf("get reward failed: %w", err)
	}

	// 2. 构建同步数据
	syncData := &SyncRewardData{
		RewardId: reward.Id,
		UserId:   reward.UserId,
		MemberId: reward.MemberId, // 会员ID
		OrderId:  reward.OrderId,
		Amount:   reward.Amount,
		Status:   reward.Status,
	}

	if reward.SettledAt.Valid {
		syncData.SettledAt = reward.SettledAt.Time
	}

	// 3. 执行同步
	return w.adapter.SyncReward(context.Background(), syncData)
}

// updateSyncStatus 更新同步状态
func (w *SyncWorker) updateSyncStatus(orderId int64, status, errorMsg string) {
	var syncedAt sql.NullTime
	if status == "synced" {
		syncedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	// 创建或更新sync_logs记录
	log := model.SyncLog{
		OrderId:    orderId,
		SyncType:   "order",
		SyncStatus: status,
		ErrorMsg:   errorMsg,
		SyncedAt:   syncedAt,
		UpdatedAt:  time.Now(),
	}

	// 尝试查找现有记录
	var existing model.SyncLog
	result := w.db.Where("order_id = ? AND sync_type = ?", orderId, "order").First(&existing)

	if result.Error == nil {
		// 更新现有记录
		log.Id = existing.Id
		log.Attempts = existing.Attempts + 1
		w.db.Save(&log)
	} else {
		// 创建新记录
		log.Attempts = 1
		log.CreatedAt = time.Now()
		w.db.Create(&log)
	}

	// 同时更新订单表的sync_status字段
	w.db.Model(&model.Order{}).Where("id = ?", orderId).Update("sync_status", status)
}

// AsyncSyncOrder 异步同步订单（将任务加入队列）
func AsyncSyncOrder(queue *SyncQueue, orderId int64) error {
	task := &SyncTask{
		TaskId:    uuid.New().String(),
		Type:      "order",
		OrderId:   orderId,
		Attempts:  0,
		CreatedAt: time.Now(),
	}

	return queue.Enqueue(task)
}

// AsyncSyncReward 异步同步奖励（将任务加入队列）
func AsyncSyncReward(queue *SyncQueue, rewardId int64) error {
	task := &SyncTask{
		TaskId:    uuid.New().String(),
		Type:      "reward",
		RewardId:  rewardId,
		Attempts:  0,
		CreatedAt: time.Now(),
	}

	return queue.Enqueue(task)
}
