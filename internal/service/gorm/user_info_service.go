package gorm

import (
	"errors"
	"fmt"
	"go_chat/internal/dao"
	"go_chat/internal/dto/request"
	"go_chat/internal/dto/respond"
	"go_chat/internal/model"
	"go_chat/internal/service/sms"
	"go_chat/pkg/constants"
	"go_chat/pkg/zlog"
	"gorm.io/gorm"
)

type userInfoService struct{}

var UserInfoService = new(userInfoService)

// Login 登录
func (u *userInfoService) Login(loginReq request.LoginRequest) (string, *respond.LoginRespond, int) {
	password := loginReq.Password
	var user model.UserInfo
	res := dao.GormDB.First(&user, "telephone = ?", loginReq.Telephone)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "用户不存在，请注册"
			zlog.Error(message)
			return message, nil, -2
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	if user.Password != password {
		message := "密码不正确，请重试"
		zlog.Error(message)
		return message, nil, -2
	}
	loginRsq := &respond.LoginRespond{
		Uuid:      user.Uuid,
		Telephone: user.Telephone,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		Birthday:  user.Birthday,
		Signature: user.Signature,
		IsAdmin:   user.IsAdmin,
		Status:    user.Status,
	}
	year, month, day := user.CreatedAt.Date()
	loginRsq.CreatedAt = fmt.Sprintf("%d.%d.%d", year, month, day)

	return "登陆成功", loginRsq, 0
}

// SendSmsCode 发送短信验证码 - 验证码登录
func (u *userInfoService) SendSmsCode(telephone string) (string, int) {
	return sms.VerificationCode(telephone)
}
