package handler

import (
	"Sellers/internal/domain"
	"Sellers/internal/seller"
	"Sellers/pkg/web"
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
	"net/http"
)

// type request struct{
// 	ID      int `json:"id"`
// 	CID     int `json:"cid"`
// 	CompanyName string `json:"company_name"`
// 	Address string `json:"address"`
// 	Telephone string `json:"telephone"`
// }

type Seller struct {
	// sellerService seller.Service

	service seller.Service
}

const tokenCompare string = "1234"


func NewSeller(p seller.Service) *Seller {
	return &Seller{
		// sellerService: s,
		service: p,
	}
}



func (s *Seller) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
	

		if token != tokenCompare {
			c.JSON(401, web.NewResponse(401, nil, "Token inválido"))
			return
		}


		id, err := strconv.ParseInt(c.Param("id"),10, 64)
		if err != nil  {
			c.JSON(400, gin.H{ "error":  "invalid ID"})
			return
		}

		p,err := s.service.Get(c.Request.Context(), int(id))

		if err != nil{
			c.JSON(404, web.NewResponse(404, nil, "No se encuentra el seller con el id ingresado"))
			return
		}

		c.JSON(200, p)


	}
}


func (s *Seller) GetAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token != tokenCompare {
			c.JSON(401, web.NewResponse(401, nil, "Token inválido"))
			return
		}

		p, err := s.service.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		if len(p) == 0 {
			c.JSON(404, web.NewResponse(404, nil, "No hay sellers"))
			return
		}
		c.JSON(200, web.NewResponse(200, p, ""))

	}
}



func (s *Seller) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")

		if token != tokenCompare {
			c.JSON(401, web.NewResponse(401, nil, "Token inválido"))
			return
		}

		var req domain.Seller

		if err := c.Bind(&req); err != nil {
			c.JSON(400, web.NewResponse(400, nil, err.Error()))
			return
		}

		if s.service.Exists(c.Request.Context(), req.CID) {
			web.Error(c, 409, "%s", "Ya existe un seller con ese CID")
			return
		}

		if req.CID == 0 {
			c.JSON(422, web.NewResponse(422, nil, "El cid es requerido"))
			return
		}

		if req.Address == "" {
			c.JSON(422, web.NewResponse(422, nil, "El domicilio es requerido"))
			return
		}

		if req.CompanyName == "" {
			c.JSON(422, web.NewResponse(422, nil, "El nombre de la compañia es requerido"))
			return
		}

		if req.Telephone == "" {
			c.JSON(422, web.NewResponse(422, nil, "El telefono es requerido"))
			return
		}
		if req.LocalitiesId == 0 {
			c.JSON(422, web.NewResponse(422, nil, "El localitiesId es requerido"))
			return
		}

		p, err := s.service.Save(c.Request.Context(), req)
		if err != nil {
			c.JSON(422, web.NewResponse(422, nil, err.Error()))
			return
		}
		c.JSON(201, web.NewResponse(201, p, ""))
	}

}

func (s *Seller) Update() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("token")
		if token != tokenCompare {
			c.JSON(401, gin.H{ "error": "token inválido" })
			return
		}
		
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "%s", "El id no es valido ")
			return
		}

		req := domain.Seller{}

		if err := c.ShouldBindJSON(&req); err != nil{
			web.Error(c, http.StatusBadRequest, "Error en la peticion: %s", err)
			return
		}
		lastS, err := s.service.Get(c.Request.Context(), int(id))
		if err != nil{
			web.Error(c, http.StatusNotFound, "No se encuentra el seller con el id ingresado")
			return
		}

		newS := updateSellerFields(lastS, req, int(id))
		err = s.service.Update(c.Request.Context(), newS)
		if err != nil{
			web.Error(c, http.StatusInternalServerError, "Ocurrio un error al actualizar el seller")
			return
		}

		web.Success(c, http.StatusOK, newS)
	
	}
	
}


func updateSellerFields(lastSeller domain.Seller, newSeller domain.Seller, id int) domain.Seller {
	if newSeller.Address != lastSeller.Address && newSeller.Address != "" {
		lastSeller.Address = newSeller.Address
	}
	if newSeller.Telephone != lastSeller.Telephone && newSeller.Telephone != "" {
		lastSeller.Telephone = newSeller.Telephone
	}
	if newSeller.CID != lastSeller.CID && newSeller.CID != 0 {
		lastSeller.CID = newSeller.CID
	}
	if newSeller.CompanyName != lastSeller.CompanyName && newSeller.CompanyName != "" {
		lastSeller.CompanyName = newSeller.CompanyName
	}
	
	lastSeller.ID = id
	return lastSeller
}


func (s *Seller) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token != tokenCompare {	
			c.JSON(401, gin.H{ "error": "token inválido" })
			return
		}

		id, err := strconv.ParseInt(c.Param("id"),10, 64)
		if err != nil {
			web.Error(c, http.StatusBadRequest, "%s", "El id ingresado no es valido")
			return
		}
		seller, err := s.service.Get(c.Request.Context(), int(id))
		if seller.ID == 0 {
			web.Error(c, http.StatusNotFound, "No se encuentra el seller con id: %d", id)
			return
		}


		err = s.service.Delete(c.Request.Context(),int(id))
		if err != nil {
			web.Error(c, http.StatusInternalServerError, "Ocurrio un error al eliminar el warehouse")
			return
		}
		c.JSON(200, gin.H{ "data": fmt.Sprintf("El seller %d ha sido eliminado", id) })


	}
}
