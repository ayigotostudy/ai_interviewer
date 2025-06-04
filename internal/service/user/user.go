package userService

import (
	"ai_jianli_go/internal/dao"
	"ai_jianli_go/pkg/utils"
	"ai_jianli_go/types/model"
	"ai_jianli_go/types/req"
	"ai_jianli_go/types/resp"
	"ai_jianli_go/types/resp/common"
)

type UserService struct {
	dao *dao.UserDAO
}

func NewUserService(dao *dao.UserDAO) *UserService {
	return &UserService{dao: dao}
}

func (s *UserService) Register(request *req.RegisterReq) int64 {
	_, err := s.dao.GetUserByEmail(request.Email)
	if err == nil {
		return common.CodeUserExist
	}
	encPwd := utils.Encrypt(request.Password)
	user := &model.User{
		Email:    request.Email,
		PassWord: encPwd,
	}
	err = s.dao.CreateUser(user)
	if err != nil {
		return common.CodeCreateUserFail
	}
	return common.CodeSuccess
}

func (s *UserService) Login(request *req.LoginReq) (any, int64) {
	user, err := s.dao.GetUserByEmail(request.Email)
	res := resp.LoginResp{
		Token: "",
	}
	if err != nil {
		return res, common.CodeUserNotExist
	}
	if user.PassWord != utils.Encrypt(request.Password) {
		return res, common.CodeInvalidPassword
	}
	token, err := utils.GetToken(user.ID, user.Role)
	if err != nil {
		return res, common.CodeServerBusy
	}

	res.Token = token
	return res, common.CodeSuccess
}
