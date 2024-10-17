package rpc

import (
	optionhub "github.com/s21platform/optionhub-proto/optionhub-proto"
	user "github.com/s21platform/user-proto/user-proto"
)

type Server struct {
	optionhub.UnimplementedOptionhubServiceServer
	user.UnimplementedUserServiceServer
	dbRepo DbRepo
	ufrR   UserFriendsRegisterSrv
}

func New(repo DbRepo, ufrR UserFriendsRegisterSrv) *Server {
	return &Server{dbRepo: repo, ufrR: ufrR}
}
