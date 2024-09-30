package rpc

import (
	user "github.com/s21platform/user-proto/user-proto"
)

type Server struct {
	user.UnimplementedUserServiceServer
	dbRepo DbRepo
	ufrR   UserFriendsRegisterSrv
}

func New(repo DbRepo, ufrR UserFriendsRegisterSrv) *Server {
	return &Server{dbRepo: repo, ufrR: ufrR}
}
