package main

import (
	"net/http"
	"worker/database"
	"worker/usecase"

	"github.com/gin-gonic/gin"
)

type BodyRequest struct {
	Filename string `json:"filename"`
}

func main() {
	app := gin.Default()
	database.ConnectDB()

	app.POST("/proccess", func(ctx *gin.Context) {
		var bodyRequest BodyRequest
		err := ctx.ShouldBindJSON(&bodyRequest)

		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"msg": err.Error(),
			})
			return
		}

		data := usecase.ProcessCsv(bodyRequest.Filename)

		if err != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{
				"msg": err.Error(),
			})
			return
		}

		ctx.JSON(200, gin.H{"msg": data})
	})
	app.Run(":8080")
}
