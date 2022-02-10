package locality

import (
	"context"
	"Sellers/internal/domain"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


type repoM struct{
	mock.Mock
}

func (r *repoM) Exists(ctx context.Context, zipC string) bool {
	args := r.Called(ctx, zipC)
	return args.Bool(0)
}

func (r *repoM) Save(ctx context.Context, loc domain.Locality) (int, error) {
	args := r.Called(ctx, loc)
	return args.Int(0), args.Error(1)
}

func (r *repoM) GetAll(ctx context.Context) ([]domain.Locality, error) {
	args := r.Called(ctx)
	return args.Get(0).([]domain.Locality), args.Error(1)
}

func (r *repoM) GetByZipCode(ctx context.Context, zipCode string) (domain.Locality, error) {
	args := r.Called(ctx, zipCode)
	return args.Get(0).(domain.Locality), args.Error(1)
}

func (r *repoM) GetSellers(ctx context.Context, l domain.Locality) ([]domain.Seller, error) {
	args := r.Called(ctx, l)
	return args.Get(0).([]domain.Seller), args.Error(1)
}


func TestCreateLocalitiesOK(t *testing.T) {
	
	repo := new(repoM) 
	
	repo.On("Save", mock.Anything, mock.Anything).Return(1, nil) 
	
	repo.On("Exists", mock.Anything, mock.Anything).Return(false)

	s := NewService(repo) 


	objetoAPersistir := domain.Locality{
	ZipCode: "4000",
	LocalityName: "San Miguel de Tucuman",
	ProvinceName: "Tucuman",
	CountryName: "Argentina",
	}

	
	ns, err := s.Save(context.Background(), objetoAPersistir)

	
	assert.NoError(t, err)
	assert.Equal(t, 1, ns.ID)
	assert.Equal(t, "4000", ns.ZipCode)
}


func TestCreateLocalitiesConflict(t *testing.T){
	repo := new(repoM)
	repo.On("Exists", mock.Anything, mock.Anything).Return(true)
	repo.On("Save", mock.Anything, mock.Anything).Return(domain.Locality{}, ErrNotFound) 
	serviceT := NewService(repo)
	ctx := context.Background()
	result := serviceT.Exists(ctx, domain.Locality{}.ZipCode)
	assert.True(t, result)


}


func TestGetByZipCodeOk(t *testing.T){

	repo := new(repoM)
	repo.On("GetByZipCode", mock.Anything, mock.Anything).Return(domain.Locality{
		
		ID: 1,
		ZipCode: "4000",
		LocalityName: "San Miguel de Tucuman",
		ProvinceName: "Tucuman",
		CountryName: "Argentina",
		
	}, nil)
	
	  obj := domain.Locality {
		ID: 1,
		ZipCode: "4000",
		LocalityName: "San Miguel de Tucuman",
		ProvinceName: "Tucuman",
		CountryName: "Argentina",
	}
	s := NewService(repo)

	objetoRecuperado, err := s.GetByZipCode(context.Background(),"4000")

	assert.NoError(t, err)
	assert.Equal(t,obj, objetoRecuperado, "Los objetos no coinciden")

}


func TestGetByZipCodeFail(t *testing.T){


	repo := new(repoM)
	repo.On("GetByZipCode", mock.Anything, mock.Anything).Return(domain.Locality{
		
		ID: 1,
		ZipCode: "4000",
		LocalityName: "San Miguel de Tucuman",
		ProvinceName: "Tucuman",
		CountryName: "Argentina",
		
	}, nil)
	
	  obj := domain.Locality {
		ID: 1,
		ZipCode: "4000",
		LocalityName: "San Miguel de Tucuman",
		ProvinceName: "Tucuman",
		CountryName: "Argentina",
	}
	s := NewService(repo)

	objetoRecuperado, err := s.GetByZipCode(context.Background(),"6700")

	assert.NoError(t, err)
	assert.Equal(t,obj, objetoRecuperado, "Los objetos no coinciden")
}


func TestGetSellersOk(t *testing.T){
	repo := new(repoM)
	obj := domain.Locality {
		ID: 1,
		ZipCode: "4000",
		LocalityName: "San Miguel de Tucuman",
		ProvinceName: "Tucuman",
		CountryName: "Argentina",
	}
	repo.On("GetSellers", mock.Anything, mock.Anything).Return([]domain.Seller{
		{
			ID:          1,
			CID:         123,
			CompanyName: "Meli",
			Address: "Bulnes 10",
			Telephone: "123456",

		},
		{
			ID:          2,
			CID:         456,
			CompanyName: "Baires Dev",
			Address: " Av Belgrano 1200",
			Telephone: "4833269",
		},
	}, nil)

	
	s:= NewService(repo)

	objetosRecuperados, err := s.GetSellers(context.Background(),obj)

	assert.NoError(t,err)

	assert.Len(t, objetosRecuperados,2)
}

func TestGetSellersFail(t *testing.T){

	repo := new(repoM)
	obj := domain.Locality {
		ID: 1,
		ZipCode: "4000",
		LocalityName: "San Miguel de Tucuman",
		ProvinceName: "Tucuman",
		CountryName: "Argentina",
	}
	repo.On("GetSellers", mock.Anything, mock.Anything).Return([]domain.Seller{}, nil)

	
	s:= NewService(repo)

	objetosRecuperados, err := s.GetSellers(context.Background(),obj)

	assert.NoError(t,err)

	assert.Len(t, objetosRecuperados,0)

}

func TestExist(t *testing.T){
	repo := new(repoM)
	repo.On("Exists", mock.Anything, mock.Anything).Return(false)
	s :=NewService(repo)
	ctx := context.Background()
	err := s.Exists(ctx, "4000")
	assert.Equal(t,false, err)
}
