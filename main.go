package main

import (
	"promanage/backend/db"
	"promanage/backend/routers"
	"promanage/backend/utils"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("未找到.env文件，或加载失败")
	}

	if err := db.InitDB(); err != nil {
		panic("failed to connect database: " + err.Error())
	}

	r := gin.Default()
	r.Use(utils.CORSMiddleware())

	routers.RegisterRoutes(r)

	r.Run(":8080")
}
