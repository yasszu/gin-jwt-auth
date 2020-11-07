package main

import (
	"gin-jwt-auth/conf"
	"gin-jwt-auth/db"
	"gin-jwt-auth/handler"
	"gin-jwt-auth/jwt"
	"gin-jwt-auth/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load conf
	cnf, err := conf.NewConf()
	if err != nil {
		panic(err.Error())
	}

	// Echo instance
	r := gin.Default()

	// Middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
	}))

	// Establish DB connection
	conn, err := db.NewConn(cnf)
	if err != nil {
		panic(err.Error())
	}

	accountRepository := repository.NewAccountRepository(conn)
	accountHandler := handler.NewAccountHandler(accountRepository, cnf)

	// /..
	r.GET("/", handler.Index)
	accountHandler.RegisterRoot(r)

	// /v1/..
	v1 := r.Group("/v1")
	v1.Use(jwt.AuthMiddleware(jwt.HeaderAuthConfig()))
	accountHandler.RegisterV1(v1)

	// Start server
	r.Run(":" + cnf.Server.Port)
}
