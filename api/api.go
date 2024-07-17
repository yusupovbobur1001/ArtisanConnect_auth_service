package api

import (
	"auth_service/api/handler"
	"database/sql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
	_ "auth_service/api/docs"
)

// @title ReserveDesk API
// @version 1.0
// @description Connecting artists and customers program

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /auth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRouter(db *sql.DB) *gin.Engine {
	h := handler.NewHadler(db)


	router := gin.Default()


	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	user := router.Group("/auth")
	user.POST("/register", h.Register)
	user.POST("/login", h.Login)
	user.POST("/logout", h.Logout)
	user.GET("/refreshtoken", h.RefreshToken)
	user.POST("/passwordrecovery", h.Passwordrecovery)

	return router
}