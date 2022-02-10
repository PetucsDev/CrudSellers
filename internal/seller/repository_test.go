package seller

import (
	"Sellers/internal/domain"
	"context"
	"errors"
	//"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)


func TestCreateSellerOk( t *testing.T){

	sellerToSave := domain.Seller{

		
		ID:                 1,
		CID: 				1,
		CompanyName: 		"Meli",
		Address:            "Bulnes 10",
		Telephone:          "123456",
		LocalitiesId: 		1,

	}

	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()
	mock.
	ExpectPrepare("INSERT INTO sellers").
	ExpectExec().
	WithArgs(sellerToSave.CID, sellerToSave.CompanyName, sellerToSave.Address, sellerToSave.Telephone, sellerToSave.LocalitiesId).
	WillReturnResult(sqlmock.NewResult(1,1)).
	WillReturnError(nil)

	sellerRepository := NewRepository(db)

	actualId, err := sellerRepository.Save(context.Background(),sellerToSave)
	assert.Nil(t, err)
	assert.Equal(t, 1, actualId)
	assert.NoError(t, mock.ExpectationsWereMet())
}


func TestCreateConflict(t *testing.T){


	sellerToSave := domain.Seller{

		
		ID:                 1,
		CID: 				1,
		CompanyName: 		"Meli",
		Address:            "Bulnes 10",
		Telephone:          "123456",
		LocalitiesId: 		1,

	}

	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()

	mock.
		ExpectPrepare("INSERT INTO sellers").
		ExpectExec().
		WithArgs(sellerToSave.CID, sellerToSave.CompanyName, sellerToSave.Address, sellerToSave.Telephone, sellerToSave.LocalitiesId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	localityRepository := NewRepository(db)

	ctx := context.Background()
	actualId, err := localityRepository.Save(ctx, sellerToSave)
	assert.Nil(t, err)
	assert.Equal(t, 1, actualId)

	sellerCompare := domain.Seller{
		ID:                 1,
		CID: 				1,
		CompanyName: 		"Meli",
		Address:            "Bulnes 10",
		Telephone:          "123456",
		LocalitiesId: 		1,
	}

	mock.
		ExpectPrepare("INSERT INTO sellers").
		ExpectExec().
		WithArgs(sellerCompare.CID, sellerCompare.CompanyName, sellerCompare.Address, sellerCompare.Telephone, sellerCompare.LocalitiesId).
		WillReturnError(errors.New("cid duplicated"))

	ctx = context.Background()
	actualId, err = localityRepository.Save(ctx, sellerCompare)
	assert.NotNil(t, err)
	assert.Equal(t, 0, actualId)
	assert.Equal(t, sellerToSave.CID, sellerCompare.CID)
	assert.NoError(t, mock.ExpectationsWereMet())

}


// func TestGetAllSellersOk(t *testing.T){
// 	sellerToSave := []domain.Seller{

		
// 		{ 
// 			ID:                 1,
// 			CID: 				1,
// 			CompanyName: 		"Meli",
// 			Address:            "Bulnes 10",
// 			Telephone:          "123456",
// 			LocalitiesId: 		1,
// 		},
// 		{
// 			ID:                 2,
// 			CID: 				2,
// 			CompanyName: 		"Baires Dev",
// 			Address:            "Belgrano 3200",
// 			Telephone:          "3814471789",
// 			LocalitiesId: 		2,
// 		},
// 	}

// 	db, mock, sqlMockErr := sqlmock.New()
// 	assert.Nil(t, sqlMockErr)
// 	defer db.Close()

// 	mock.
// 		ExpectPrepare("SELECT * FROM sellers").
// 		ExpectExec().
// 		WithArgs(sellerToSave).
// 		WillReturnResult(sqlmock.NewResult(1, 1))

// 	sellerRepository := NewRepository(db)

// 	ctx := context.Background()
// 	actualId, err := sellerRepository.GetAll(ctx)
// 	assert.Nil(t, err)
// 	assert.Len(t, actualId, 2)
// }

// func TestGetSellerByIdOk(t *testing.T){
// 	expectedSeller := domain.Seller{
// 			ID:                 1,
// 			CID: 				1,
//  			CompanyName: 		"Meli",
//  			Address:            "Bulnes 10",
//  			Telephone:          "123456",
//  			LocalitiesId: 		1,
// 	}

// 	db, mock, sqlMockErr := sqlmock.New()
// 	assert.Nil(t, sqlMockErr)
// 	defer db.Close()

// 	columns := []string{"id","cid","company_name","address","telephone","localities_id"}
// 	rows := sqlmock.NewRows(columns)

// 	sellerId := 1
// 	rows.AddRow(sellerId,1,"Meli","Bulnes 10","123456",1)
// 	mock.
// 	ExpectPrepare("SELECT * FROM sellers WHERE id=?;").
// 	ExpectQuery().
// 	WithArgs(6).
// 	WillReturnRows(rows)


// 	sellerRepository := NewRepository(db)

// 	actualResult, err := sellerRepository.Get(context.Background(), 1)

// 	actualResult1 := domain.Seller{
// 			ID:                 1,
// 			CID: 				1,
//  			CompanyName: 		"Meli",
//  			Address:            "Bulnes 10",
//  			Telephone:          "123456",
//  			LocalitiesId: 		1,
// 	}
// 	fmt.Println(actualResult)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, actualResult1)
// 	assert.Equal(t, expectedSeller, actualResult1)
// 	assert.NoError(t, mock.ExpectationsWereMet())
// }