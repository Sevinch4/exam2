package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type driverRepo struct {
	DB *sql.DB
}

func NewDriverRepo(db *sql.DB) storage.IDriverRepo {
	return driverRepo{
		DB: db,
	}
}

func (d driverRepo) Create(driver models.CreateDriver) (string, error) {
	id := uuid.New()

	query := `INSERT INTO drivers(id, full_name, phone, from_city_id, to_city_id) values($1, $2, $3, $4, $5)`
	if _, err := d.DB.Exec(query, id, driver.FullName, driver.Phone, driver.FromCityID, driver.ToCityID); err != nil {
		fmt.Println("error is while inserting driver", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (d driverRepo) Get(id string) (models.Driver, error) { //dodelat
	driver := models.Driver{}
	query := `SELECT d.id, d.full_name, d.phone, d.from_city_id , fc.id as from_city_id, fc.name as from_city_name,
       fc.created_at as from_city_created_at, 
       d.to_city_id, tc.id as to_city_id, tc.name as to_city_name, tc.created_at as to_city_created_at,
       d.created_at
					FROM drivers AS d 
					    LEFT JOIN cities AS fc on  d.from_city_id = fc.id 
						LEFT JOIN cities as tc on d.to_city_id = tc.id
					WHERE d.id = $1`

	var (
		fromCityID, toCityID     string
		fromCityData, toCityData models.City
	)
	if err := d.DB.QueryRow(query, id).Scan(
		&driver.ID,
		&driver.FullName,
		&driver.Phone,
		&fromCityID,
		&fromCityData.ID,
		&fromCityData.Name,
		&fromCityData.CreatedAt,
		&toCityID,
		&toCityData.ID,
		&toCityData.Name,
		&toCityData.CreatedAt,
		&driver.CreatedAt); err != nil {
		fmt.Println("error is while selecting by id", err.Error())
		return models.Driver{}, err
	}

	driver.FromCityData = fromCityData
	driver.FromCityID = fromCityID
	driver.ToCityData = toCityData
	driver.ToCityID = toCityID

	return driver, nil
}

func (d driverRepo) GetList(req models.GetListRequest) (models.DriversResponse, error) {
	var (
		page              = req.Page
		offset            = (page - 1) * req.Limit
		drivers           = []models.Driver{}
		query, countQuery string
		count             = 0
	)

	countQuery = `SELECT count(1) FROM drivers`
	if err := d.DB.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning count", err.Error())
		return models.DriversResponse{}, err
	}

	query = `SELECT id, full_name, phone, from_city_id, to_city_id, created_at FROM drivers LIMIT $1 OFFSET $2`

	rows, err := d.DB.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting drivers", err.Error())
		return models.DriversResponse{}, err
	}

	for rows.Next() {
		driver := models.Driver{}
		if err = rows.Scan(&driver.ID, &driver.FullName, &driver.Phone, &driver.FromCityID, &driver.ToCityID, &driver.CreatedAt); err != nil {
			fmt.Println("error is while scanning drivers", err.Error())
			return models.DriversResponse{}, err
		}
		drivers = append(drivers, driver)
	}

	return models.DriversResponse{
		Drivers: drivers,
		Count:   count,
	}, nil
}

func (d driverRepo) Update(driver models.Driver) (string, error) {
	query := `UPDATE drivers SET full_name = $1, phone = $2, from_city_id = $3, to_city_id = $4 WHERE id = $5`
	if _, err := d.DB.Exec(query,
		&driver.FullName,
		&driver.Phone,
		&driver.FromCityID,
		&driver.ToCityID,
		&driver.ID); err != nil {
		fmt.Println("error is while updating driver", err.Error())
		return "", err
	}
	return driver.ID, nil
}

func (d driverRepo) Delete(id string) error {
	query := `DELETE FROM drivers WHERE id = $1`
	if _, err := d.DB.Exec(query, id); err != nil {
		fmt.Println("error is while deleting driver", err.Error())
		return err
	}

	return nil
}

func (d driverRepo) GetDriver() {

}
