package rpc

import (
	user "github.com/s21platform/user-proto/user-proto"
)

type Server struct {
	user.UnimplementedUserServiceServer
	dbRepo     DbRepo
	ufrR       UserFriendsRegisterSrv
	optionhubS OptionhubS
}

func New(repo DbRepo, ufrR UserFriendsRegisterSrv, optionhubService OptionhubS) *Server {
	return &Server{
		dbRepo:     repo,
		ufrR:       ufrR,
		optionhubS: optionhubService,
	}
}
