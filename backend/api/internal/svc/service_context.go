// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"
	"fmt"
	"time"

	"dmh/api/internal/config"
	"dmh/api/internal/middleware"
	"dmh/api/internal/service"
	"dmh/common/poster"
	"dmh/common/wechatpay"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config               config.Config
	DB                   *gorm.DB
	PasswordService      *service.PasswordService
	AuditService         *service.AuditService
	SessionService       *service.SessionService
	PosterService        *poster.Service
	PosterRateLimiter    middleware.RateLimiter
	DefaultRateLimiter   middleware.RateLimiter
	WeChatPayService     *wechatpay.Service
	PermissionMiddleware *middleware.PermissionMiddleware
}

type redisAdapter struct {
	client *redis.Client
}

func (r *redisAdapter) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *redisAdapter) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisAdapter) Expire(ctx context.Context, key string, seconds int) error {
	return r.client.Expire(ctx, key, time.Duration(seconds)*time.Second).Err()
}

func (r *redisAdapter) TTL(ctx context.Context, key string) (string, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", int(ttl.Seconds())), nil
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化GORM数据库连接
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		logx.Errorf("连接数据库失败: %v", err)
		// 不中断服务，继续运行（使用mock数据）
	} else {
		// 测试连接
		sqlDB, err := db.DB()
		if err != nil {
			logx.Errorf("获取数据库实例失败: %v", err)
		} else if err := sqlDB.Ping(); err != nil {
			logx.Errorf("数据库ping失败: %v", err)
		} else {
			logx.Info("数据库连接成功")
		}
	}

	passwordService := service.NewPasswordService(db)
	auditService := service.NewAuditService(db)
	sessionService := service.NewSessionService(db, passwordService)
	posterService := poster.NewService("/opt/data/posters", "http://localhost:8889/api/v1")

	wechatPayConfig := &wechatpay.Config{
		AppID:           c.WeChatPay.AppID,
		MchID:           c.WeChatPay.MchID,
		APIKey:          c.WeChatPay.APIKey,
		APIKeyV3:        c.WeChatPay.APIKeyV3,
		APIClientCert:   c.WeChatPay.APIClientCert,
		APIClientKey:    c.WeChatPay.APIClientKey,
		NotifyURL:       c.WeChatPay.NotifyURL,
		RefundNotifyURL: c.WeChatPay.RefundNotifyURL,
		Sandbox:         c.WeChatPay.Sandbox,
		CacheTTL:        c.WeChatPay.CacheTTL,
		MockEnabled:     c.WeChatPay.MockEnabled,
		UnifiedOrderURL: c.WeChatPay.UnifiedOrderURL,
		HTTPTimeoutMs:   c.WeChatPay.HTTPTimeoutMs,
	}
	wechatPayService := wechatpay.NewService(wechatPayConfig)
	logx.Infof("微信支付配置: AppID=%s, MchID=%s, Sandbox=%v, MockEnabled=%v, CacheTTL=%ds",
		c.WeChatPay.AppID, c.WeChatPay.MchID, c.WeChatPay.Sandbox, c.WeChatPay.MockEnabled, c.WeChatPay.CacheTTL)

	ctx := context.Background()
	var redisAdapterClient middleware.RedisClient
	useRedis := c.RateLimit.PosterGenerate.Storage == "redis" || c.RateLimit.Default.Storage == "redis"
	if useRedis {
		if c.Redis.Host == "" {
			logx.Errorf("Redis未配置，已退回内存存储")
		} else {
			redisClient := redis.NewClient(&redis.Options{
				Addr:     c.Redis.Host,
				Password: c.Redis.Pass,
				DB:       0,
				PoolSize: 10,
			})
			if err := redisClient.Ping(ctx).Err(); err != nil {
				logx.Errorf("Redis连接失败: %v", err)
			} else {
				redisAdapterClient = &redisAdapter{client: redisClient}
			}
		}
	} else {
		logx.Info("使用内存存储")
	}

	createRateLimiter := func(storageType string, redisClient middleware.RedisClient, maxRequests int, duration int, prefix string) middleware.RateLimiter {
		if storageType == "redis" && redisClient != nil {
			return middleware.NewRedisRateLimiter(redisClient, prefix, maxRequests, time.Duration(duration)*time.Second)
		}
		return middleware.NewMemoryRateLimiter(maxRequests, time.Duration(duration)*time.Second)
	}

	posterRateLimiter := createRateLimiter(c.RateLimit.PosterGenerate.Storage, redisAdapterClient, c.RateLimit.PosterGenerate.MaxRequests, c.RateLimit.PosterGenerate.WindowDuration, "poster")
	defaultRateLimiter := createRateLimiter(c.RateLimit.Default.Storage, redisAdapterClient, c.RateLimit.Default.MaxRequests, c.RateLimit.Default.WindowDuration, "default")

	// 初始化权限中间件
	var permissionMiddleware *middleware.PermissionMiddleware
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			logx.Errorf("获取底层SQL DB失败: %v", err)
		}
		permissionMiddleware = middleware.NewPermissionMiddleware(sqlDB)
	} else {
		permissionMiddleware = middleware.NewPermissionMiddleware(nil)
	}

	return &ServiceContext{
		Config:               c,
		DB:                   db,
		PasswordService:      passwordService,
		AuditService:         auditService,
		SessionService:       sessionService,
		PosterService:        posterService,
		PosterRateLimiter:    posterRateLimiter,
		DefaultRateLimiter:   defaultRateLimiter,
		WeChatPayService:     wechatPayService,
		PermissionMiddleware: permissionMiddleware,
	}
}
