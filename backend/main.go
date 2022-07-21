package main

import (
	"erp/customers"
	"erp/reports"
	"erp/timetracking"
	"erp/users"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	app := gin.New()
	u := app.Group("/users")
	users.Main(u)

	c := app.Group("/customers")
	customers.Main(c)

	r := app.Group("/reports")
	reports.Main(r)

	t := app.Group("/timetracking")
	timetracking.Main(t)

	app.Run(":3000")
}
