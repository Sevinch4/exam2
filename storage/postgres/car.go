package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type carRepo struct {
	db *sql.DB
}

func NewCarRepo(db *sql.DB) storage.ICarRepo {
	return carRepo{
		db,
	}
}

func (c carRepo) Create(car models.CreateCar) (string, error) {
	id := uuid.New()

	query := `INSERT INTO cars(id, model, brand, number) values ($1, $2, $3, $4)`

	if _, err := c.db.Exec(query, id, car.Model, car.Brand); err != nil {
		fmt.Println("error is while inserting car", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (c carRepo) Get(id string) (models.Car, error) {
	car := models.Car{}
	query := `SELECT c.id, c.model, c.brand, c.number, c.driver_id,
       			d.id as driver_id, d.full_name as driver_name, d.phone as driver_phone,
       			d.from_city_id as driver_from, d.to_city_id as driver_to, 
       			d.created_at as driver_date, c.created_at,
				fc.id as from_city_id, fc.name as from_city_name, fc.created_at as from_city_date,
				tc.id as to_city_id, tc.name as to_city_name, tc.created_at as to_city_date
				FROM cars AS c 
				    LEFT JOIN drivers as d ON c.driver_id = d.id 
					LEFT JOIN cities as fc ON d.from_city_id = fc.id
					LEFT JOIN cities as tc ON d.to_city_id= tc.id
				WHERE c.id = $1`

	if err := c.db.QueryRow(query, id).Scan(
		&car.ID,
		&car.Model,
		&car.Brand,
		&car.Number,
		&car.DriverID,
		&car.DriverData.ID,
		&car.DriverData.FullName,
		&car.DriverData.Phone,
		&car.DriverData.FromCityID,
		&car.DriverData.ToCityData.ID,
		&car.DriverData.CreatedAt,
		&car.CreatedAt,
		&car.DriverData.FromCityData.ID,
		&car.DriverData.FromCityData.Name,
		&car.DriverData.FromCityData.CreatedAt,
		&car.DriverData.ToCityData.Name,
		&car.DriverData.ToCityData.CreatedAt,
		&car.DriverData.ToCityID,
	); err != nil {
		fmt.Println("error is while scanning car by id", err.Error())
		return models.Car{}, err
	}
	return car, nil
}

func (c carRepo) GetList(req models.GetListRequest) (models.CarsResponse, error) {
	//BITMADI bitta data ciqarvotdi
	var (
		page              = req.Page
		offset            = (page - 1) * req.Limit
		query, countQuery string
		count             = 0
		cars              = []models.Car{}
	)

	countQuery = `SELECT count(1) FROM cars`
	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count")
		return models.CarsResponse{}, err
	}

	query = `SELECT c.id, c.model, c.brand, c.number, c.driver_id,
       				d.id as driver_id, d.full_name as driver_name, d.phone as driver_phone,
       				d.from_city_id as driver_from, d.to_city_id as driver_to, d.created_at, c.created_at,
       				f.id as from_city_id, f.name as from_city_name, f.created_at as from_city_date,
       				t.id as to_city_id, t.name as to_city_name, t.created_at as to_city_date
					FROM cars AS c 
					    LEFT JOIN drivers as d ON c.driver_id = d.id
						LEFT JOIN cities as f ON d.from_city_id = f.id
						LEFT JOIN cities as t ON d.to_city_id = t.id `
	query += ` LIMIT $1 OFFSET $2`

	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting cars", err.Error())
		return models.CarsResponse{}, err
	}

	for rows.Next() {
		cs := models.Car{}
		if err = rows.Scan(
			&cs.ID,
			&cs.Model,
			&cs.Brand,
			&cs.Number,
			&cs.DriverID,
			&cs.DriverData.ID,
			&cs.DriverData.FullName,
			&cs.DriverData.Phone,
			&cs.DriverData.FromCityID,
			&cs.DriverData.ToCityID,
			&cs.DriverData.CreatedAt,
			&cs.CreatedAt,
			&cs.DriverData.FromCityData.ID,
			&cs.DriverData.FromCityData.Name,
			&cs.DriverData.FromCityData.CreatedAt,
			&cs.DriverData.ToCityData.ID,
			&cs.DriverData.ToCityData.Name,
			&cs.DriverData.ToCityData.CreatedAt); err != nil {
			fmt.Println("error is while scanning cars", err.Error())
			return models.CarsResponse{}, err
		}
		cars = append(cars, cs)
	}

	return models.CarsResponse{
		Cars:  cars,
		Count: count,
	}, nil
}

func (c carRepo) Update(car models.Car) (string, error) {
	query := `UPDATE cars SET model = $1, brand = $2, number = &3, driver_id = $4 WHERE id = $5`

	if _, err := c.db.Exec(query, &car.Model, &car.Brand, &car.Number, &car.DriverID, &car.ID); err != nil {
		fmt.Println("error is while updating car", err.Error())
		return "", err
	}
	return car.ID, nil
}

func (c carRepo) Delete(id string) error {
	query := `DELETE FROM cars WHERE id = $1`

	if _, err := c.db.Exec(query, id); err != nil {
		fmt.Println("error is while deleting car", err.Error())
		return err
	}
	return nil
}

func (c carRepo) UpdateCarRoute(models.UpdateCarRoute) error {
	route := models.UpdateCarRoute{
		DepartureTime: time.Now(),
	}
	query := `UPDATE drivers SET from_city_id = $1, to_city_id = $2
                         FROM cars 
               WHERE cars.driver_id = drivers.id AND cars.id = $3`
	if _, err := c.db.Exec(query,
		&route.FromCityID,
		&route.ToCityID,
		&route.CarID); err != nil {
		fmt.Println("error is while updating car route", err.Error())
		return err
	}
	return nil
}
func (c carRepo) UpdateCarStatus(status models.UpdateCarStatus) error {
	query := `UPDATE cars SET status = $1 WHERE id = $2`

	if _, err := c.db.Exec(query, &status.Status, &status.ID); err != nil {
		fmt.Println("error is while updating cars", err.Error())
		return err
	}

	return nil
}
