package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"Sellers/internal/seller"
	"Sellers/cmd/server/handler"
)

type Router interface {
	MapRoutes()
}

type router struct {
	r  *gin.Engine
	rg *gin.RouterGroup
	db *sql.DB
}

func NewRouter(r *gin.Engine, db *sql.DB) Router {
	return &router{r: r, db: db}
}

func (r *router) MapRoutes() {
	r.setGroup()

	r.buildSellerRoutes()

}

func (r *router) setGroup() {
	r.rg = r.r.Group("/api/v1")
}

func (r *router) buildSellerRoutes() {
	// Example
	repo := seller.NewRepository(r.db)
	service := seller.NewService(repo)
	handler := handler.NewSeller(service)
	r.r.GET("/seller", handler.GetAll())
	r.r.GET("/seller/:id", handler.Get())
	r.r.POST("/seller/create", handler.Create())
	r.r.DELETE("/seller/:id", handler.Delete())
	r.r.PATCH("/seller/:id", handler.Update())
}

