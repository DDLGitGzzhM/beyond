package logic

import (
	"beyond/application/applet/internal/svc"
	"beyond/application/applet/internal/types"
	"beyond/application/user/rpc/user"
	"beyond/pkg/util"
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	prefixVerificationCount = "biz#verification#count#%s"
	//一天可以请求的次数
	verificationLimitPerDay = 10
	expireActivation        = 60 * 30
)

type VerificationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerificationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerificationLogic {
	return &VerificationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerificationLogic) Verification(req *types.VerificationRequest) (resp *types.VerificationResponse, err error) {
	count, err := l.getVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("无效的手机号 : %s error : %v", req.Mobile, err)
		return nil, err
	}

	if count > verificationLimitPerDay {
		return nil, errors.New("已经超过一天请求 验证码的 次数")
	}

	code, err := getActivationCache(req.Mobile, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("getActivationCache mobile: %s error: %v", req.Mobile, err)
	}

	if len(code) == 0 {
		code = util.RandomNumeric(6)
	}

	_, err = l.svcCtx.UserRPC.SendSms(l.ctx, &user.SendSmsRequest{
		Mobile: req.Mobile,
	})

	if err != nil {
		logx.Errorf("sendSms mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}

	err = saveActivationCache(req.Mobile, code, l.svcCtx.BizRedis)
	if err != nil {
		logx.Errorf("saveActivationCache mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}

	err = l.incrVerificationCount(req.Mobile)
	if err != nil {
		logx.Errorf("incrVerificationCount mobile: %s error: %v", req.Mobile, err)
	}

	return &types.VerificationResponse{}, nil
}

func (l *VerificationLogic) getVerificationCount(mobile string) (int, error) {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	val, err := l.svcCtx.BizRedis.Get(key)
	if err != nil {
		return 0, err
	}
	if len(val) == 0 {
		return 0, nil
	}
	return strconv.Atoi(val)
}

func (l *VerificationLogic) incrVerificationCount(mobile string) error {
	key := fmt.Sprintf(prefixVerificationCount, mobile)
	_, err := l.svcCtx.BizRedis.Incr(key)
	if err != nil {
		return err
	}

	return l.svcCtx.BizRedis.Expireat(key, util.EndOfDay(time.Now()).Unix())
}

func getActivationCache(mobile string, rds *redis.Redis) (string, error) {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Get(key)
}

func saveActivationCache(mobile, code string, rds *redis.Redis) error {
	key := fmt.Sprintf(prefixActivation, mobile)
	return rds.Setex(key, code, expireActivation)
}
