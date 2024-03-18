package logic

import (
	"beyond/application/user/rpc/user"
	"beyond/pkg/encrypt"
	"beyond/pkg/jwt"
	"context"
	"errors"
	"fmt"
	"strings"

	"beyond/application/applet/internal/svc"
	"beyond/application/applet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

const (
	prefixActivation = "biz#activation#%s"
)

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.RegisterResponse, err error) {
	req.Name = strings.TrimSpace(req.Name)
	if len(req.Name) == 0 {
		return nil, errors.New("名字不能为空")
	}
	req.Password = strings.TrimSpace(req.Password)
	if len(req.Password) == 0 {
		return nil, errors.New("密码不能为空")
	}
	req.Password = encrypt.EncPassword(req.Password)
	req.VerificationCode = strings.TrimSpace(req.VerificationCode)
	if len(req.VerificationCode) == 0 {
		return nil, errors.New("verification code cannot be empty")
	}

	err = l.checkVerificationCode(l.ctx, req.Mobile, req.VerificationCode)
	if err != nil {
		return nil, err
	}

	mobile, err := encrypt.EncMobile(req.Mobile)
	if err != nil {
		logx.Errorf("EncMobile mobile: %s error: %v", req.Mobile, err)
		return nil, err
	}
	fmt.Println("in FindByMobile")
	userRet, err := l.svcCtx.UserRPC.FindByMobile(l.ctx, &user.FindByMobileRequest{
		Mobile: mobile,
	})
	fmt.Println("out FindByMobile")
	if err != nil {
		fmt.Println("FindByMobile error")
		return nil, err
	}

	if userRet != nil && userRet.UserId > 0 {
		return nil, errors.New("this mobile is already registered")
	}
	fmt.Println("in Register")
	regRet, err := l.svcCtx.UserRPC.Register(l.ctx, &user.RegisterRequest{
		Username: req.Name,
		Mobile:   mobile,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("out Register")
	token, err := jwt.BuildTokens(jwt.TokenOptions{
		AccessSecret: l.svcCtx.Config.Auth.AccessSecret,
		AccessExpire: l.svcCtx.Config.Auth.RefreshExpire,
		Fields: map[string]interface{}{
			"userId": regRet.UserId,
		},
	})
	if err != nil {
		return nil, err
	}

	return &types.RegisterResponse{
		UserId: regRet.UserId,
		Token:  token,
	}, nil
	return
}

func (l *RegisterLogic) checkVerificationCode(ctx context.Context, mobile, code string) error {
	cacheCode, err := getActivationCache(mobile, l.svcCtx.BizRedis)
	if err != nil {
		return err
	}
	if cacheCode == "" {
		return errors.New("verification code expired")
	}
	if cacheCode != code {
		return errors.New("verification code failed")
	}

	return nil
}
