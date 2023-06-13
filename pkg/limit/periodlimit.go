package limit

import (
	"context"
	"fmt"
	"github.com/codfrm/cago/pkg/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"

	"github.com/codfrm/cago/pkg/utils"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/redis/go-redis/v9"
)

// TODO: redis lua脚本保证原子性
const script = `
`

// PeriodLimit 周期限流器,redis zet实现滑动窗口
type PeriodLimit struct {
	period, quota int64
	limitStore    *redis.Client
	keyPrefix     string
}

// NewPeriodLimit 创建周期限流器,period单位秒,quota限流数量
func NewPeriodLimit(period, quota int64, limitStore *redis.Client, keyPrefix string) *PeriodLimit {
	return &PeriodLimit{
		period:     period,
		quota:      quota,
		limitStore: limitStore,
		keyPrefix:  keyPrefix,
	}
}

func (p *PeriodLimit) key(key string) string {
	return p.keyPrefix + ":" + key
}

func (p *PeriodLimit) Take(ctx context.Context, key string) (func() error, error) {
	key = p.key(key)
	now := time.Now().Unix()
	cnt, err := p.limitStore.ZCount(ctx, key, strconv.FormatInt(now-p.period, 10), "+inf").Result()
	if err != nil {
		if redis.Nil != err {
			return nil, err
		}
	}
	total, err := p.limitStore.ZCard(ctx, key).Result()
	if err != nil {
		if redis.Nil != err {
			return nil, err
		}
	}
	// 当总数为1000的余数时,删除过期记录
	if total > 1000 && total%1000 == 0 {
		go func() {
			if err := p.limitStore.ZRemRangeByScore(ctx, key, "-inf", strconv.FormatInt(now-p.period*2+60, 10)).Err(); err != nil {
				logger.Ctx(ctx).Error("删除过期记录失败", zap.String("key", key), zap.Error(err))
			}
		}()
	}
	if cnt < p.quota {
		flag := utils.RandString(8, utils.Mix)
		err := p.limitStore.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: flag,
		}).Err()
		if err != nil {
			return nil, err
		}
		if err := p.limitStore.Expire(ctx, key, time.Duration(p.period+60)*time.Second).Err(); err != nil {
			return nil, err
		}
		// 删除本次记录
		return func() error {
			return p.limitStore.ZRem(ctx, key, flag).Err()
		}, nil
	}
	log := fmt.Sprintf("%d秒内产生了太多请求", p.period)
	return nil, httputils.NewError(http.StatusTooManyRequests, -1, log)
}

func (p *PeriodLimit) FuncTake(ctx context.Context, key string, f func() (interface{}, error)) (interface{}, error) {
	cancel, err := p.Take(ctx, key)
	if err != nil {
		return nil, err
	}
	resp, err := f()
	if err != nil {
		if err := cancel(); err != nil {
			return nil, err
		}
		return nil, err
	}
	return resp, nil
}

// Count 获取用量
func (p *PeriodLimit) Count(ctx context.Context, key string, period int64) (int64, error) {
	key = p.key(key)
	now := time.Now().Unix()
	cnt, err := p.limitStore.ZCount(ctx, key, strconv.FormatInt(now-period, 10), "+inf").Result()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
