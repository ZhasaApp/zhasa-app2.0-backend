package apiadmin

import (
	"github.com/gin-gonic/gin"
	"zhasa2.0/user/repository"
)

type Server struct {
	getUserByPhoneFunc repository.GetUserByPhoneFunc
	createUserFunc     repository.CreateUserFunc
}

func NewServer(getUserByPhoneFunc repository.GetUserByPhoneFunc, createUserFunc repository.CreateUserFunc) *Server {
	return &Server{
		getUserByPhoneFunc: getUserByPhoneFunc,
		createUserFunc:     createUserFunc,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
