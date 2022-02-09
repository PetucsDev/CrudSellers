package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"Sellers/internal/domain"
	"Sellers/internal/locality"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	ErrNotFound = errors.New("locality not found")
)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) Save(ctx context.Context, loc domain.Locality) (int, error) {
	args := m.Called(ctx, loc)
	return args.Int(0), args.Error(1)
}

func (m *ServiceMock) Exists(ctx context.Context, zipC string) bool {
	args := m.Called(ctx, zipC)
	return args.Bool(0)
}

func (r *ServiceMock) GetByZipCode(ctx context.Context, zipCode string) (domain.Locality, error) {
	args := r.Called(ctx, zipCode)
	return args.Get(0).(domain.Locality), args.Error(1)
}
func (r *ServiceMock) GetSellers(ctx context.Context, l domain.Locality) ([]domain.Seller, error) {
	args := r.Called(ctx, l)
	return args.Get(0).([]domain.Seller), args.Error(1)
}

func (r *ServiceMock) GetAll(ctx context.Context) ([]domain.Locality, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Locality), args.Error(1)
}

func createServerLocalities(s *Locality) *gin.Engine {
	r := gin.Default()
	locality := r.Group("/api/v1/localities")
	{
		locality.GET("/reportSellers", s.Get())
		locality.POST("/", s.Create())
	
	}
	return r
}


func createRequestTestLocalities(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", "1234")
	return req, httptest.NewRecorder()
}


func TestCreateLocalityOK(t *testing.T) {

	expectedResponse := 
	`{
		"zip_code": "6700",
  		"locality_name": "Lujan",
  		"province_name": "Buenos Aires",
  		"country_name": "Argentina"

	}`


	s := new(ServiceMock)
	s.On("Exists", mock.Anything, mock.Anything).Return(false)
	s.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	service := locality.NewService(s)
	w := NewLocality(service)
	r := createServerLocalities(w)
	req, res := createRequestTestLocalities(http.MethodPost, `/api/v1/localities/`, expectedResponse)
	r.ServeHTTP(res, req)
	
	assert.Equal(t, 201, res.Code, res.Result())
}

func TestCreateLocalityFail(t *testing.T) {
	s := new(ServiceMock)
	s.On("Exists", mock.Anything, mock.Anything).Return(false)
	s.On("Save", mock.Anything, mock.Anything).Return(1, nil)
	service := locality.NewService(s)
	w := NewLocality(service)
	r := createServerLocalities(w)
	body := `
		{
			
			"locality_name": "Lujan",
  			"province_name": "Buenos Aires",
  			"country_name": "Argentina"
		}`
	req, res := createRequestTestLocalities(http.MethodPost, `/api/v1/localities/`, body)
	r.ServeHTTP(res, req)
	assert.Equal(t, 422, res.Code, res.Result())
	
}


func TestCreateLocalityConflict(t *testing.T) {
	s := new(ServiceMock)
	s.On("Exists", mock.Anything, mock.Anything).Return(true)
	s.On("Save", mock.Anything, mock.AnythingOfType("domain.Locality")).Return(1,ErrNotFound)
	service := locality.NewService(s)
	w := NewLocality(service)
	r := createServerLocalities(w)
	body := `
		{
			"zip_code": "6700",
  			"locality_name": "Lujan",
  			"province_name": "Buenos Aires",
  			"country_name": "Argentina"
		}`
	req, res := createRequestTestLocalities(http.MethodPost, "/api/v1/localities/", body)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusConflict, res.Code)
}

