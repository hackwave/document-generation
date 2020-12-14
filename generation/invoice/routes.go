package server

import (
	"net/http"
	"strconv"

	gin "github.com/gin-gonic/gin"
)

func (self *Server) Route() *Server {

	self.Router.GET("/invoice/by-index/:index", func(context *gin.Context) {
		index, _ := strconv.Atoi(context.Param("index"))
		context.Writer.WriteHeader(http.StatusOK)
		context.Writer.Write([]byte(self.Merchant.Invoice(self.Merchant.OrderByIndex(index))))

	})

	self.Router.GET("/", func(context *gin.Context) {
		context.Writer.WriteHeader(http.StatusOK)
		context.Writer.Write([]byte(self.Merchant.YearReport(2019)))
	})

	return self
}
