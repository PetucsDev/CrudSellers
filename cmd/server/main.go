package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"Sellers/cmd/server/routes"
)

func main() {

	//db, _ := sql.Open("mysql", "meli_sprint_user:Meli_Sprint#123@/melisprint") // me falta el host: puerto, nombre bd
	db, _ := sql.Open("mysql", "root:peti96_cnc@tcp(127.0.0.1:3306)/sellerbd")
	r := gin.Default()

	router := routes.NewRouter(r, db)
	router.MapRoutes()

	if err := r.Run(":3001"); err != nil {
		panic(err)
	}
}