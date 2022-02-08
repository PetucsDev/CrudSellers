package locality


import (
	"context"
	"database/sql"
	"Sellers/internal/domain"
)

type Repository interface{
	Save(ctx context.Context, l domain.Locality) (int, error)
	GetAll(ctx context.Context) ([]domain.Locality, error)
	Exists(ctx context.Context, id string) bool
}


type repository struct{
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		db: db,
	}
}


func (r *repository) GetAll(ctx context.Context) ([]domain.Locality, error) {
	query := "SELECT * FROM localities"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var localities []domain.Locality

	for rows.Next() {
		s := domain.Locality{}
		_ = rows.Scan(&s.ID,s.ZipCode ,&s.LocalityName, &s.ProvinceName, &s.CountryName)
		localities = append(localities, s)
	}

	return localities, nil
}

func (r *repository) Save(ctx context.Context, l domain.Locality) (int, error) {
	query := "INSERT INTO localities (zip_code, locality_name, province_name, province_name, country_name) VALUES (?,?, ?, ?, ?)"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(l.ZipCode, l.LocalityName, l.ProvinceName, l.CountryName)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *repository) Exists(ctx context.Context, id string) bool {
	query := "SELECT zip_code FROM localities WHERE zip_code=?;"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&id)
	return err == nil
}

