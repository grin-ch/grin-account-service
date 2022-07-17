package service

import (
	"context"

	"github.com/grin-ch/grin-account-service/pkg/auth"
	"github.com/grin-ch/grin-account-service/pkg/model"
	"github.com/grin-ch/grin-account-service/pkg/util"
	"github.com/grin-ch/grin-api/api/account"
	"github.com/grin-ch/grin-api/api/captcha"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService struct {
	pv            *model.Provider
	captchaClient captcha.CaptchaServiceClient
}

func NewUserService(pv *model.Provider, captchaClient captcha.CaptchaServiceClient) account.UserServiceServer {
	return &userService{
		pv,
		captchaClient,
	}
}

// AuthFuncOverride 覆盖authFunc
// 跳过部分接口的权限验证
func (*userService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	return ctx, nil
}

// 注册
func (s *userService) SignUp(ctx context.Context, req *account.SignUpReq) (*account.SignUpRsp, error) {
	// 参数校验
	if len(req.Username) < 3 || len(req.Username) > 24 {
		return nil, status.Errorf(codes.InvalidArgument, "length of username must between 3 to 12")
	}
	var phoneNumber, email string
	if util.ValidatePhoneNumberFormat(req.Contact) {
		phoneNumber = req.Contact
	} else if util.ValidateEmailFormat(req.Contact) {
		email = req.Contact
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "contact invalid")
	}
	if len(req.Password) < 8 {
		return nil, status.Errorf(codes.InvalidArgument, "password length less 8")
	}
	encodePwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "params invalid")
	}

	// 验证码校验
	rsp, err := s.captchaClient.Verify(ctx, &captcha.VerifyReq{
		Key:     req.Contact,
		Value:   req.Captcha,
		Purpose: captcha.Purpose_SIGN_UP,
	})
	if err != nil {
		log.Errorf("captcha verify error:%s", err.Error())
		return nil, status.Errorf(codes.Internal, "captcha verify error")
	}
	if !rsp.Success {
		return nil, status.Errorf(codes.InvalidArgument, "captcha invalid")
	}

	// 创建帐号
	err = s.pv.CreateUser(model.User{
		Username:    req.Username,
		PhoneNumber: phoneNumber,
		Email:       email,
		Password:    string(encodePwd),
	})
	if err != nil {
		return nil, err
	}
	return &account.SignUpRsp{
		Success: true,
		Message: "Sign Up!",
	}, nil
}

// 登入
func (s *userService) SignIn(ctx context.Context, req *account.SignInReq) (*account.SignInRsp, error) {
	// 验证码校验
	rsp, err := s.captchaClient.Verify(ctx, &captcha.VerifyReq{
		Key:     req.Key,
		Value:   req.Value,
		Purpose: captcha.Purpose_SIGN_IN,
	})
	if err != nil {
		log.Errorf("captcha verify error:%s", err.Error())
		return nil, status.Errorf(codes.Internal, "captcha verify error")
	}
	if !rsp.Success {
		return nil, status.Errorf(codes.InvalidArgument, "captcha invalid")
	}

	// 密码校验
	var phoneNumber, email string
	if util.ValidatePhoneNumberFormat(req.Contact) {
		phoneNumber = req.Contact
	} else if util.ValidateEmailFormat(req.Contact) {
		email = req.Contact
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "contact invalid")
	}

	user, err := s.pv.PickUser(model.User{PhoneNumber: phoneNumber, Email: email})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "account not found")
	}
	token, err := auth.Generate(user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "generate token error")
	}

	return &account.SignInRsp{
		Token: token,
	}, nil
}

// 重置密码
func (s *userService) ResetPasswd(ctx context.Context, req *account.ResetPasswdReq) (*account.ResetPasswdRsp, error) {
	var phoneNumber, email string
	if util.ValidatePhoneNumberFormat(req.Contact) {
		phoneNumber = req.Contact
	} else if util.ValidateEmailFormat(req.Contact) {
		email = req.Contact
	}
	encodePwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "params invalid")
	}
	// TODO
	// 重置密码
	s.pv.CreateUser(model.User{
		PhoneNumber: phoneNumber,
		Email:       email,
		Password:    string(encodePwd),
	})
	return nil, nil
}
