package locality

import (
	"Sellers/internal/domain"
	"context"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"errors"
)



func TestCreateOk( t *testing.T){

	localityToSave := domain.Locality{

		
		ZipCode: "6700",
		LocalityName: "Lujan",
		ProvinceName: "Buenos Aires",
		CountryName: "Argentina",

	}

	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()
	mock.
	ExpectPrepare("INSERT INTO localities").
	ExpectExec().
	WithArgs(localityToSave.ZipCode, localityToSave.LocalityName, localityToSave.ProvinceName, localityToSave.CountryName).
	WillReturnResult(sqlmock.NewResult(1,1)).
	WillReturnError(nil)

	localityRepository := NewRepository(db)

	actualId, err := localityRepository.Save(context.Background(),localityToSave)
	assert.Nil(t, err)
	assert.Equal(t, 1, actualId)
	assert.NoError(t, mock.ExpectationsWereMet())
}



func TestCreateConflict(t *testing.T){


	localityToSave := domain.Locality{
		ZipCode:      "6700",
		LocalityName: "Lujan",
		ProvinceName: "Buenos Aires",
		CountryName:  "Argentina",
	}

	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()

	mock.
		ExpectPrepare("INSERT INTO localities").
		ExpectExec().
		WithArgs(localityToSave.ZipCode, localityToSave.LocalityName, localityToSave.ProvinceName, localityToSave.CountryName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	localityRepository := NewRepository(db)

	ctx := context.Background()
	actualId, err := localityRepository.Save(ctx, localityToSave)
	assert.Nil(t, err)
	assert.Equal(t, 1, actualId)

	localityCompare := domain.Locality{
		ZipCode:      "6700",
		LocalityName: "Lujan",
		ProvinceName: "Buenos Aires",
		CountryName:  "Argentina",
	}

	mock.
		ExpectPrepare("INSERT INTO localities").
		ExpectExec().
		WithArgs(localityToSave.ZipCode, localityCompare.LocalityName, localityCompare.ProvinceName, localityCompare.CountryName).
		WillReturnError(errors.New("zipCode duplicated"))

	ctx = context.Background()
	actualId, err = localityRepository.Save(ctx, localityCompare)
	assert.NotNil(t, err)
	assert.Equal(t, 0, actualId)
	assert.Equal(t, localityToSave.ZipCode, localityCompare.ZipCode)
	assert.NoError(t, mock.ExpectationsWereMet())

}

func TestDoesExistsZipCode(t *testing.T) {
	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()

	columns := []string{"zip_code"}
	rows := sqlmock.NewRows(columns)

	sellerID := "4000"
	rows.AddRow(sellerID)
	mock.
		ExpectQuery("SELECT zip_code FROM localities WHERE zip_code=\\?;").
		WithArgs(sellerID).
		WillReturnRows(rows)

	localityRepository := NewRepository(db)

	doesExist := localityRepository.Exists(context.Background(), sellerID)
	assert.True(t, doesExist)
	assert.NoError(t, mock.ExpectationsWereMet())
}
func TestDoesntExistsZipCode(t *testing.T) {
	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()

	testID := "5000"

	mock.
		ExpectQuery("SELECT zip_code FROM localities WHERE zip_code=\\?;").
		WithArgs(testID)

	localityRepository := NewRepository(db)

	doesExist := localityRepository.Exists(context.Background(), testID)

	assert.False(t, doesExist)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetLocalityByZipCodeFail(t *testing.T) {
	expectedLocality := domain.Locality{
		ID: 			1,
		ZipCode:      "6700",
		LocalityName: "Lujan",
		ProvinceName: "Buenos Aires",
		CountryName:  "Argentina",
	}

	db, mock, sqlMockErr := sqlmock.New()
	assert.Nil(t, sqlMockErr)
	defer db.Close()

	columns := []string{"id", "zip_code", "locality_name","province_name", "country_name"}
	rows := sqlmock.NewRows(columns)

	sellerId := 1
	rows.AddRow(sellerId, "6700", "Lujan", "Buenos Aires", "Argentina")
	mock.
		ExpectQuery("SELECT \\* FROM localities WHERE zip_code=\\?").
		WithArgs("4000").
		WillReturnRows(rows)

	localityRepository := NewRepository(db)

	actualResult, err := localityRepository.GetByZipCode(context.Background(), "4000")

	assert.Nil(t, err)
	assert.NotNil(t, actualResult)
	assert.NotEqual(t, expectedLocality, actualResult)
	assert.NoError(t, mock.ExpectationsWereMet())
}


// func TestGetSellersSuccess(t *testing.T){
// 	expectedSeller := domain.Seller{
// 		ID:                 1,
// 		CID: 				1,
// 		CompanyName: 		"Meli",
// 		Address:            "Bulnes 10",
// 		Telephone:          "123456",
// 		LocalitiesId: 		1,
// 	}

// 	expectedLocality := domain.Locality{
// 		ID: 			1,
// 		ZipCode:      "6700",
// 		LocalityName: "Lujan",
// 		ProvinceName: "Buenos Aires",
// 		CountryName:  "Argentina",
// 	}

// 	db, mock, sqlMockErr := sqlmock.New()
// 	assert.Nil(t, sqlMockErr)
// 	defer db.Close()

// 	columns := []string{"id", "cid","company_name","address","telephone","localities_id"}
// 	rows := sqlmock.NewRows(columns)

// 	sellerId := 1
// 	rows.AddRow(sellerId, 1, "Meli", "Bulnes 10", "123456", 1)
// 	mock.
// 		ExpectQuery("SELECT \\* FROM sellers WHERE localities_id=\\?").
// 		WithArgs(expectedSeller).
// 		WillReturnRows(rows)

// 	localityRepository := NewRepository(db)

// 	actualResult, err := localityRepository.GetSellers(context.Background(), expectedLocality)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, actualResult)
// 	assert.Equal(t, expectedSeller, actualResult)
// 	assert.NoError(t, mock.ExpectationsWereMet())


// }