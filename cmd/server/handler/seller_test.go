package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"Sellers/internal/domain"
	"Sellers/internal/seller"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ServiceM struct {
	mock.Mock
}



func (m *ServiceM) Save(ctx context.Context, s domain.Seller) (int, error) {
	args := m.Called(ctx, s)
	return args.Int(0), args.Error(1)
}

func (m *ServiceM) Get(ctx context.Context, id int) (domain.Seller, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Seller), args.Error(1)
}

func (m *ServiceM) GetAll(ctx context.Context) ([]domain.Seller, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Seller), args.Error(1)
}

func (m *ServiceM) Update(ctx context.Context, s domain.Seller) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *ServiceM) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ServiceM) Exists(ctx context.Context, cid int) bool {
	args := m.Called(ctx, cid)
	return args.Bool(0)
}



func createServer(s *Seller) *gin.Engine {
	r := gin.Default()
	seller := r.Group("/api/v1/sellers")
	{
		seller.GET("/", s.GetAll())
		seller.GET("/:id", s.Get())
		seller.POST("/", s.Create())
		seller.PATCH("/:id", s.Update())
		seller.DELETE("/:id", s.Delete())
	}
	return r
}

func createRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	return req, httptest.NewRecorder()
}


// TESTS CREATE


func TestCreateSellerOK(t *testing.T) {

	expectedResponse := 
	`{
		"cid": 1,
		"address": "Bulnes 10",
		"telephone": "123456",
		"company_name": "Meli",
		"localities_id":1
	}`


	s := new(ServiceM)
	s.On("Exists", mock.Anything, mock.Anything).Return(false)
	s.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodPost, `/api/v1/sellers/`, expectedResponse)
	r.ServeHTTP(res, req)
	
	assert.Equal(t, 201, res.Code, res.Result())
}



func TestCreateSellerFail(t *testing.T) {
	s := new(ServiceM)
	s.On("Exists", mock.Anything, mock.Anything).Return(false)
	s.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	body := `
		{
			
			"company_name": "Pedidos Ya",
			"address": "Av Aconquija 1100",
			"telephone": "987654",
			"localities_id":1
		}`
	req, res := createRequestTest(http.MethodPost, `/api/v1/sellers/`, body)
	r.ServeHTTP(res, req)
	assert.Equal(t, 422, res.Code, res.Result())
	

}
func TestCreateSellerFailCompanyName(t *testing.T) {
	s := new(ServiceM)
	s.On("Exists", mock.Anything, mock.Anything).Return(false)
	s.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	body := `
		{
			"cid": 1,
			"address": "Av Aconquija 1100",
			"telephone": "987654",
			"localities_id":1
		}`
	req, res := createRequestTest(http.MethodPost, `/api/v1/sellers/`, body)
	r.ServeHTTP(res, req)
	assert.Equal(t, 422, res.Code, res.Result())
	

}

func TestCreateSellerFailAddress(t *testing.T) {
	s := new(ServiceM)
	s.On("Exists", mock.Anything, mock.Anything).Return(false)
	s.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	body := `
		{
			"cid": 1,
			"company_name": "Pedidos Ya",
			"telephone": "987654",
			"localities_id":1
		}`
	req, res := createRequestTest(http.MethodPost, `/api/v1/sellers/`, body)
	r.ServeHTTP(res, req)
	assert.Equal(t, 422, res.Code, res.Result())
	

}


func TestCreateSellerFailTelephone(t *testing.T) {
	s := new(ServiceM)
	s.On("Exists", mock.Anything, mock.Anything).Return(false)
	s.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	body := `
		{
			"cid": 1,
			"company_name": "Pedidos Ya",
			"address": "Av Aconquija 1100",
			"localities_id":1
		}`
	req, res := createRequestTest(http.MethodPost, `/api/v1/sellers/`, body)
	r.ServeHTTP(res, req)
	assert.Equal(t, 422, res.Code, res.Result())
	
}

func TestCreateSellerFailLocalities(t *testing.T) {
	s := new(ServiceM)
	s.On("Exists", mock.Anything, mock.Anything).Return(false)
	s.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	body := `
		{
			"cid": 1,
			"company_name": "Pedidos Ya",
			"address": "Av Aconquija 1100",
			"telephone": "987654"
			
		}`
	req, res := createRequestTest(http.MethodPost, `/api/v1/sellers/`, body)
	r.ServeHTTP(res, req)
	assert.Equal(t, 422, res.Code, res.Result())
	
}



func TestCreateSellerConflict(t *testing.T) {
	s := new(ServiceM)
	s.On("Exists", mock.Anything, mock.Anything).Return(true)
	s.On("Save", mock.Anything, mock.AnythingOfType("domain.Seller")).Return(1, errors.New("Ya existe un seller con ese numero de cid"))
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	body := `
		{
			"cid": 1,
			"company_name": "Demo",
			"address": "Bulnes 10",
			"telephone": "123456",
			"localities_id":1
		}`
	req, res := createRequestTest(http.MethodPost, "/api/v1/sellers/", body)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusConflict, res.Code)
}


// TESTS GET

func TestFindAllSellers(t *testing.T) {

	mockResponse := []domain.Seller{}
	mockResponse = append(mockResponse, domain.Seller{
		ID:                 1,
		CID: 				1,
		CompanyName: 		"Meli",
		Address:            "Bulnes 10",
		Telephone:          "123456",
		LocalitiesId:        1,
		
	})
	mockResponse = append(mockResponse, domain.Seller{
			ID:          2,
			CID:         2,
			CompanyName: "Baires Dev",
			Address: " Av Belgrano 1200",
			Telephone: "4833269",
			LocalitiesId:       2,
	})
	s := new(ServiceM)
	s.On("GetAll", mock.Anything, mock.Anything).Return(mockResponse, nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodGet, "/api/v1/sellers/", "")
	r.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)

}
func TestFindAllSellersFailService(t *testing.T) {

	
	s := new(ServiceM)
	s.On("GetAll", mock.Anything, mock.Anything).Return([]domain.Seller{}, ErrNotFound)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodGet, "/api/v1/sellers/", "")
	r.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code)

}

func TestFindAllSellersFailService2(t *testing.T) {

	s := new(ServiceM)
	s.On("GetAll", mock.Anything, mock.Anything).Return([]domain.Seller{}, nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodGet, "/api/v1/sellers/", "")
	r.ServeHTTP(res, req)

	assert.Equal(t, 404, res.Code)

}





func TestFindSellerByIdNonExistent(t *testing.T) {
	expectedResult := domain.Seller{}
	expectedError := "sql: no rows in result set"
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(expectedResult, errors.New(expectedError))
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	req, err := createRequestTest(http.MethodGet, "/api/v1/sellers/100", "")
	r.ServeHTTP(err, req)
	assert.Equal(t, http.StatusNotFound, err.Code)
}



func TestFindSellerByIdExistent(t *testing.T) {
		
	
	mockResponse := domain.Seller{
		ID:                 1,
		CID: 				1,
		CompanyName: 		"Meli",
		Address:            "Bulnes 10",
		Telephone:          "123456",
		LocalitiesId: 		1,
	}

	expectedResponse := `{"id":1,"cid":1,"company_name":"Meli","address":"Bulnes 10","telephone":"123456","localities_id":1}`
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(mockResponse, nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodGet, "/api/v1/sellers/1", "")
	r.ServeHTTP(res, req)
	actualResponse := res.Body.String()
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, expectedResponse, actualResponse)
}


// TESTS UPDATE
func TestUpdateSellerOK(t *testing.T) {
	mockResponse := domain.Seller{
		ID:                 1,
		CID: 				1,
		CompanyName: 		"Meli",
		Address:            "Bulnes 10",
		Telephone:          "123456",
	}
	expectedResponse := `{"id":1,"cid":1,"company_name":"Meli","address":"Av Belgrano 3200","telephone":"654321"}`
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(mockResponse, nil)
	s.On("Update", mock.Anything, mock.Anything).Return(nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	objAc := `{"id":1,"cid":1,"company_name":"Meli","address":"Av Belgrano 3200","telephone":"654321"}`
	req, res := createRequestTest(http.MethodPatch, `/api/v1/sellers/1`, objAc)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, expectedResponse, objAc)

}


func TestUpdateNonSeller(t *testing.T) {

	expectedResult := domain.Seller{}
	expectedError := "sql: no rows in result set"
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(expectedResult, errors.New(expectedError))
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)

	body := `
		{
			"cid": 1,
			"company_name": "Demo",
			"address": "Bulnes 10",
			"telephone": "123456"
		}`
	req, res := createRequestTest(http.MethodPatch, `/api/v1/sellers/2`, body)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusNotFound, res.Code)
}



func TestUpdateSellerFail(t *testing.T) {
	mockResponse := domain.Seller{
			ID:          2,
			CID:         2,
			CompanyName: "Baires Dev",
			Address: " Av Belgrano 1200",
			Telephone: "4833269",
	}
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(mockResponse, nil)
	s.On("Update", mock.Anything, mock.Anything).Return(nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodPatch, `/api/v1/sellers/2`,
		`{
		"id": "34",
		"cid": 1,
		"address": "Av Sarmiento 123",
		"telephone": "395716",
		
		
		}`)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusBadRequest, res.Code)
}


// TESTS DELETE

func TestDeleteNonExistSeller(t *testing.T) {
	expectedError := "sql: no rows in result set"
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(domain.Seller{}, errors.New(expectedError))
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodDelete, `/api/v1/sellers/2`, "{}")
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusNotFound, res.Code)
	assert.Equal(t, 404, res.Code, res.Result())
}


func TestDeleteSellerOK(t *testing.T) {
	mockResponse := domain.Seller{
		ID:                 1,
		CID: 				1,
		CompanyName: 		"Meli",
		Address:            "Bulnes 10",
		Telephone:          "123456",
	}
	s := new(ServiceM)
	s.On("Get", mock.Anything, mock.Anything).Return(mockResponse, nil)
	s.On("Delete", mock.Anything, mock.Anything).Return(nil)
	service := seller.NewService(s)
	w := NewSeller(service)
	r := createServer(w)
	req, res := createRequestTest(http.MethodDelete, `/api/v1/sellers/1`, "")
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)
}




































// import (
// 	"Sellers/internal/domain"
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// )


// func createSellerServer() *gin.Engine {
// 	sellerService := NewSellerServiceMock()
// 	seller := NewSeller(sellerService)
// 	r := gin.Default()

// 	sellersGroup := r.Group("/sellers")
// 	{
// 		sellersGroup.GET("/", seller.GetAll())
// 		sellersGroup.GET("/:id", seller.Get())
// 		sellersGroup.POST("/", seller.Create())
// 		sellersGroup.PATCH("/:id", seller.Update())
// 		sellersGroup.DELETE("/:id", seller.Delete())
// 	}
// 	return r
// }
// func createSellerRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
// 	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Header.Add("token", "1234")
// 	return req, httptest.NewRecorder()
// }

// type SellerServiceMock struct {
// 	Data []byte
// }

// func NewSellerServiceMock() *SellerServiceMock {
// 	datosDePrueba := []domain.Seller{
// 		{
// 			ID:          1,
// 			CID:         1,
// 			CompanyName: "Meli",
// 			Address: "Bulnes 10",
// 			Telephone: "123456",
// 		},
// 		{
// 			ID:          2,
// 			CID:         2,
// 			CompanyName: "Baires Dev",
// 			Address: " Av Belgrano 1200",
// 			Telephone: "4833269",
// 		},
// 		{
// 			ID:          3,
// 			CID:         3,
// 			CompanyName: "Uala",
// 			Address: " Av Mate de Luna 2200",
// 			Telephone: "3814471789",
// 		},
// 	}
// 	datos, _ := json.Marshal(datosDePrueba)
// 	return &SellerServiceMock{
// 		Data: datos,
// 	}
// }



// func (s *SellerServiceMock) GetAll(ctx context.Context) ([]domain.Seller, error) {
// 	var datos []domain.Seller
// 	err := json.Unmarshal(s.Data, &datos)
// 	if err != nil {
// 		return []domain.Seller{}, err
// 	}
// 	return datos, nil
// }


// func (s *SellerServiceMock) Get(ctx context.Context, id int) (domain.Seller, error) {
// 	var datos []domain.Seller
// 	err := json.Unmarshal(s.Data, &datos)
// 	if err != nil {
// 		return domain.Seller{}, err
// 	}
// 	for _, seller := range datos {
// 		if seller.ID == id {
// 			return seller, nil
// 		}
// 	}
// 	return domain.Seller{}, fmt.Errorf("Seller not found")
// }


// func (s *SellerServiceMock) Save(ctx context.Context, se domain.Seller) (domain.Seller, error) {
// 	var datos []domain.Seller
// 	err := json.Unmarshal(s.Data, &datos)
// 	if err != nil {
// 		return domain.Seller{}, err
// 	}
// 	for _, sel := range datos {
// 		if sel.CID == se.CID {
// 			return domain.Seller{}, fmt.Errorf("Seller already exists")
// 		}
// 	}
// 	se.ID = datos[len(datos)-1].ID
// 	datos = append(datos, se)
// 	return se, nil
// }

// func (s *SellerServiceMock) Update(ctx context.Context, se domain.Seller) (domain.Seller, error) {
// 	var datos []domain.Seller
// 	err := json.Unmarshal(s.Data, &datos)
// 	if err != nil {
// 		return domain.Seller{}, err
// 	}
// 	for _, sel := range datos {
// 		if sel.ID == se.ID {
// 			sel.CID = se.CID
// 			sel.CompanyName = se.CompanyName
// 			sel.Address = se.Address
// 			sel.Telephone = se.Telephone
// 			return sel, nil
// 		}
// 	}
// 	return domain.Seller{}, fmt.Errorf("Employee does not exist")
// }

// func (s *SellerServiceMock) Delete(ctx context.Context, id int) error {
// 	var datos []domain.Seller
// 	err := json.Unmarshal(s.Data, &datos)
// 	if err != nil {
// 		return err
// 	}
// 	for _, seller := range datos {
// 		if seller.ID == id {
// 			return nil
// 		}
// 	}
// 	return fmt.Errorf("Seller not found")
// }


// func Test_create_ok(t *testing.T) {
// 	r := createSellerServer()
	
// 	body := `
// 		{
// 			"cid": 4,
// 			"company_name": "Pedidos Ya",
// 			"address": "Av Aconquija 1100",
// 			"telephone": "987654"
// 		}`
// 	req, rr := createSellerRequestTest(http.MethodPost, "/sellers/", body)

// 	r.ServeHTTP(rr, req)

// 	assert.Equal(t, 201, rr.Code, rr.Result())
// }

// func Test_create_fail(t *testing.T) {
// 	r := createSellerServer()
// 	body := `
// 		{
			
// 			"company_name": "Pedidos Ya",
// 			"address": "Av Aconquija 1100",
// 			"telephone": "987654"
// 		}`
// 	req, rr := createSellerRequestTest(http.MethodPost, "/sellers/", body)

// 	r.ServeHTTP(rr, req)

// 	assert.Equal(t, 422, rr.Code, rr.Result())
// }


// func Test_create_conflict(t *testing.T) {
// 	r := createSellerServer()
			
// 	body := `
// 		{
// 			"cid": 1,
// 			"company_name": "Demo",
// 			"address": "Bulnes 10",
// 			"telephone": "123456"
// 		}`
// 	req, rr := createSellerRequestTest(http.MethodPost, "/sellers/", body)

// 	r.ServeHTTP(rr, req)

// 	assert.Equal(t, 400, rr.Code, rr.Result())
// }


// func Test_find_all(t *testing.T) {
// 	r := createSellerServer()
// 	req, rr := createSellerRequestTest(http.MethodGet, "/sellers/", "{}")

// 	r.ServeHTTP(rr, req)

// 	assert.Equal(t, 200, rr.Code, rr.Result())
// }


// func Test_find_by_id_non_existent(t *testing.T) {
// 	r := createSellerServer()
// 	req, rr := createSellerRequestTest(http.MethodGet, "/sellers/4", "{}")
// 	r.ServeHTTP(rr, req)
// 	assert.Equal(t, 401, rr.Code, rr.Result())
// }


// func Test_find_by_id_existent(t *testing.T) {
// 	r := createSellerServer()
// 	req, rr := createSellerRequestTest(http.MethodGet, "/sellers/2", "{}")
// 	r.ServeHTTP(rr, req)
// 	assert.Equal(t, 200, rr.Code, rr.Result())
// }


// func Test_update_ok(t *testing.T) {
// 	r := createSellerServer()
	
// 	body := `
// 		{
// 			"id": 2,
// 			"cid": 2,
// 			"company_name": "Baires Dev",
// 			"address": "Av Aconquija 1200",
// 			"telephone": "4833269"
// 		}`
// 	req, rr := createSellerRequestTest(http.MethodPatch, "/sellers/2", body)

// 	r.ServeHTTP(rr, req)


// 	assert.Equal(t, 200, rr.Code)
// }


// func Test_update_non_existent(t *testing.T) {
// 	r := createSellerServer()
// 	body := `
// 		{
// 			"id": 2,
// 			"cid": 2,
// 			"company_name": "Baires Dev",
// 			"address": "Av Aconquija 1200",
// 			"telephone": "4833269"
// 		}`
// 	req, rr := createSellerRequestTest(http.MethodPatch, "/sellers/5", body)

// 	r.ServeHTTP(rr, req)

// 	assert.Equal(t, 400, rr.Code)
// }

// func Test_delete_non_existent(t *testing.T) {
// 	r := createSellerServer()
// 	req, rr := createSellerRequestTest(http.MethodDelete, "/sellers/4", "{}")

// 	r.ServeHTTP(rr, req)

// 	assert.Equal(t, 404, rr.Code, rr.Result())
// }


// func Test_delete_ok(t *testing.T) {
// 	r := createSellerServer()
// 	req, rr := createSellerRequestTest(http.MethodDelete, "/sellers/1", "{}")

// 	r.ServeHTTP(rr, req)

// 	assert.Equal(t, 200, rr.Code, rr.Result())
// }








