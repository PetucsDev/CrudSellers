package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"Sellers/internal/seller"
	"Sellers/cmd/server/handler"
	"Sellers/internal/locality"
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
	r.buildLocalityRoutes()

}

func (r *router) setGroup() {
	r.rg = r.r.Group("/api/v1")
}

func (r *router) buildSellerRoutes() {
	// Example
	repo := seller.NewRepository(r.db)
	service := seller.NewService(repo)
	handler := handler.NewSeller(service)
	r.r.GET("/sellers", handler.GetAll())
	r.r.GET("/sellers/:id", handler.Get())
	r.r.POST("/sellers", handler.Create())
	r.r.DELETE("/sellers/:id", handler.Delete())
	r.r.PATCH("/sellers/:id", handler.Update())
}

func (r *router) buildLocalityRoutes() {
	// Example
	repo := locality.NewRepository(r.db)
	service := locality.NewService(repo)
	handler := handler.NewLocality(service)
	r.r.POST("/localities", handler.Create())
	r.r.GET("/localities/reportSellers", handler.Get())
	
}

