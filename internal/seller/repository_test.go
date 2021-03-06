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

func TestGetAllSellersOk(t *testing.T) {
	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()

	columns := []string{"id", "cid", "company_name", "address", "telephone", "localities_id"}
	rows := sqlmock.NewRows(columns)
	rows.AddRow(1, 1, "Meli", "Bulnes 10", "123456", 1)
	rows.AddRow(2, 2, "Baires Dev", "Belgrano 3200", "3814471789", 2)

	// TODO
	mock.
		ExpectQuery("SELECT \\* FROM sellers").
		WillReturnRows(rows)

	sellerRepository := NewRepository(db)

	ctx := context.Background()
	sellers, err := sellerRepository.GetAll(ctx)
	assert.Nil(t, err)
	assert.Len(t, sellers, 2)
}

func TestGetSellerByIdOk(t *testing.T) {
	expectedSeller := domain.Seller{
		ID:           1,
		CID:          1,
		CompanyName:  "Meli",
		Address:      "Bulnes 10",
		Telephone:    "123456",
		LocalitiesId: 1,
	}

	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()

	columns := []string{"id", "cid", "company_name", "address", "telephone", "localities_id"}
	rows := sqlmock.NewRows(columns)

	sellerId := 1
	rows.AddRow(sellerId, 1, "Meli", "Bulnes 10", "123456", 1)
	mock.
		ExpectQuery("SELECT \\* FROM sellers WHERE id=\\?;").
		WithArgs(1).
		WillReturnRows(rows)

	sellerRepository := NewRepository(db)

	actualResult, err := sellerRepository.Get(context.Background(), 1)

	assert.Nil(t, err)
	assert.NotNil(t, actualResult)
	assert.Equal(t, expectedSeller, actualResult)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDoesExists(t *testing.T) {
	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()

	columns := []string{"cid"}
	rows := sqlmock.NewRows(columns)

	sellerID := 1
	rows.AddRow(sellerID)
	mock.
		ExpectQuery("SELECT cid FROM sellers WHERE cid=\\?;").
		WithArgs(sellerID).
		WillReturnRows(rows)

	sellerRepository := NewRepository(db)

	doesExist := sellerRepository.Exists(context.Background(), sellerID)
	assert.True(t, doesExist)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDoesntExists(t *testing.T) {
	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()

	testID := 2

	mock.
		ExpectQuery("SELECT cid FROM sellers WHERE cid=\\?;").
		WithArgs(testID)

	sellerRepository := NewRepository(db)

	doesExist := sellerRepository.Exists(context.Background(), testID)

	assert.False(t, doesExist)
	assert.NoError(t, mock.ExpectationsWereMet())
}


func TestUpdateSuccess(t *testing.T) {
	sellerToUpdate := domain.Seller{
		ID:           1,
		CID:          1,
		CompanyName:  "Meli",
		Address:      "Bulnes 10",
		Telephone:    "123456",
		LocalitiesId: 1,
	}

	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()
	mock.
		ExpectPrepare("UPDATE sellers SET cid=\\?, company_name=\\?, address=\\?, telephone=\\? WHERE id=\\?").
		ExpectExec().
		WithArgs(sellerToUpdate.CID, sellerToUpdate.CompanyName, sellerToUpdate.Address, sellerToUpdate.Telephone, sellerToUpdate.ID).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	sellerRepository := NewRepository(db)

	err := sellerRepository.Update(context.Background(), sellerToUpdate)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSuccess(t *testing.T) {
	sellerToDelete := domain.Seller{
		ID:           1,
		CID:          1,
		CompanyName:  "Meli",
		Address:      "Bulnes 10",
		Telephone:    "123456",
		LocalitiesId: 1,
	}

	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()
	mock.
		ExpectPrepare("DELETE FROM sellers WHERE id=\\?").
		ExpectExec().
		WithArgs(sellerToDelete.ID).
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	sellerRepository := NewRepository(db)

	err := sellerRepository.Delete(context.Background(), 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

