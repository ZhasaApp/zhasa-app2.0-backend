package apiadmin

import (
	"github.com/gin-gonic/gin"
	"zhasa2.0/user/repository"
)

type Server struct {
	getUserByPhoneFunc    repository.GetUserByPhoneFunc
	createUserFunc        repository.CreateUserFunc
	makeUserAsManagerFunc repository.MakeUserAsManagerFunc
}

func NewServer(getUserByPhoneFunc repository.GetUserByPhoneFunc, createUserFunc repository.CreateUserFunc, makeManagerAsUserFunc repository.MakeUserAsManagerFunc) *Server {
	return &Server{
		getUserByPhoneFunc:    getUserByPhoneFunc,
		createUserFunc:        createUserFunc,
		makeUserAsManagerFunc: makeManagerAsUserFunc,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
