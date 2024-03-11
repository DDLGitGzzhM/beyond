package logic

import (
	"beyond/application/applet/internal/svc"
	"beyond/application/applet/internal/types"
	"context"

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
	// todo: add your logic here and delete this line
	return
}
