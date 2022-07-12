package main

import (
	"erp/users"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app := gin.New()
	r := app.Group("/users")
	users.Main(r)

	app.Run(":3000")
}
