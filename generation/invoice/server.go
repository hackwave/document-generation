package server

import (
	"fmt"

	merchant "./merchant"

	gin "github.com/gin-gonic/gin"
)

type ServerState int

const (
	On ServerState = iota
	Off
	Restarting
	Starting
	Stoping
)

// TODO: Build a state machine later with control
//       over state shifting.

type Server struct {
	Merchant *merchant.Merchant
	Router   *gin.Engine
	State    ServerState
}

func Init(merchant *merchant.Merchant) *Server {
	fmt.Println("Initializing the [docgen]==>[http server]:8080")
	server := &Server{
		Merchant: merchant,
		Router:   gin.Default(),
		State:    Off,
	}
	//gin.SetMode(gin.ReleaseMode)
	server.Router.Use(gin.Recovery())
	server.Route()

	return server
}

func (self *Server) Start() *Server {
	fmt.Println("Starting http @ localhost:8080")
	self.Router.Run("0.0.0.0:8080")
	return self
}

func (self *Server) Stop() *Server {
	return self
}

func (self *Server) Restart() *Server {
	self.Stop()
	self.Start()
	return self
}

func (self *Server) Status() *Server {
	fmt.Println("<===================>")
	fmt.Println("                     ")
	fmt.Println("<===================>")
	return self
}
