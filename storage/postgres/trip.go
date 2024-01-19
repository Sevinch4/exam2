package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"sync"
)

var (
	tripNumberMu      sync.Mutex
	currentTripNumber = 1
)

type tripRepo struct {
	db *sql.DB
}

func NewTripRepo(db *sql.DB) storage.ITripRepo {
	return &tripRepo{
		db: db,
	}
}

func (c *tripRepo) Create(req models.CreateTrip) (string, error) {
	id := uuid.New()

	tripNumberMu.Lock()
	tripNumberID := currentTripNumber
	currentTripNumber++
	tripNumberMu.Unlock()
	req.TripNumberID = fmt.Sprintf("T-%d", tripNumberID)

	query := `INSERT INTO trips (id, trip_number_id, from_city_id, to_city_id, driver_id, price)
                        values($1, $2, $3, $4, $5, $6)`

	if _, err := c.db.Exec(query,
		id,
		req.TripNumberID,
		req.FromCityID,
		req.ToCityID,
		req.DriverID,
		req.Price); err != nil {
		fmt.Println("error is while inserting data to trip", err.Error())
		return "", err
	}

	return req.TripNumberID, nil
}

func (c *tripRepo) Get(id string) (models.Trip, error) {
	trip := models.Trip{}
	query := `SELECT t.id, t.trip_number_id, t.from_city_id, t.to_city_id, t.driver_id, t.price,
						f.id as city_from_id, f.name as city_from_name, f.created_at as city_from_date,
						ct.id as city_to_id, ct.name as city_to_name, ct.created_at as city_to_date,
						d.id as driver_id, d.full_name as driver_name, d.phone as driver_phone, 
						d.from_city_id as driver_from_city_id, d.to_city_id as driver_to_city_id,
						d.created_at, t.created_at,
						fd.id as driver_from_id, fd.name as driver_from_name, fd.created_at as driver_from_date,
						ctd.id as driver_to_id, ctd.name as driver_to_name, ctd.created_at as driver_to_date
						FROM trips AS t
    					LEFT JOIN cities AS f ON t.from_city_id = f.id
    					LEFT JOIN cities AS ct ON t.to_city_id = ct.id
    					LEFT JOIN drivers AS d ON t.driver_id = d.id
						LEFT JOIN cities AS fd ON d.from_city_id = fd.id
    					LEFT JOIN cities AS ctd ON d.to_city_id = ctd.id
						WHERE t.trip_number_id = $1
    			`
	if err := c.db.QueryRow(query, id).Scan(
		&trip.ID, &trip.TripNumberID, &trip.FromCityID, &trip.ToCityID, &trip.DriverID, &trip.Price,
		&trip.FromCityData.ID, &trip.FromCityData.Name, &trip.FromCityData.CreatedAt,
		&trip.ToCityData.ID, &trip.ToCityData.Name, &trip.ToCityData.CreatedAt,
		&trip.DriverData.ID, &trip.DriverData.FullName, &trip.DriverData.Phone, &trip.DriverData.FromCityID, &trip.DriverData.ToCityID,
		&trip.DriverData.CreatedAt, &trip.CreatedAt,
		&trip.DriverData.FromCityData.ID, &trip.DriverData.FromCityData.Name, &trip.DriverData.FromCityData.CreatedAt,
		&trip.DriverData.ToCityData.ID, &trip.DriverData.ToCityData.Name, &trip.DriverData.ToCityData.CreatedAt,
	); err != nil {
		fmt.Println("error is while selecting data", err.Error())
		return models.Trip{}, err
	}
	return trip, nil
}

func (c *tripRepo) GetList(req models.GetListRequest) (models.TripsResponse, error) {
	var (
		page              = req.Page
		offset            = (page - 1) * req.Limit
		trips             = []models.Trip{}
		query, countQuery string
		count             = 0
	)

	countQuery = `SELECT count(1) FROM trips`
	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning count", err.Error())
		return models.TripsResponse{}, err
	}

	query = `SELECT t.id, t.trip_number_id, t.from_city_id, t.to_city_id, t.driver_id, t.price,
						f.id as city_from_id, f.name as city_from_name, f.created_at as city_from_date,
						ct.id as city_to_id, ct.name as city_to_name, ct.created_at as city_to_date,
						d.id as driver_id, d.full_name as driver_name, d.phone as driver_phone, 
						d.from_city_id as driver_from_city_id, d.to_city_id as driver_to_city_id,
						d.created_at, t.created_at,
						fd.id as driver_from_id, fd.name as driver_from_name, fd.created_at as driver_from_date,
						ctd.id as driver_to_id, ctd.name as driver_to_name, ctd.created_at as driver_to_date
						FROM trips AS t
    					LEFT JOIN cities AS f ON t.from_city_id = f.id
    					LEFT JOIN cities AS ct ON t.to_city_id = ct.id
    					LEFT JOIN drivers AS d ON t.driver_id = d.id
						LEFT JOIN cities AS fd ON d.from_city_id = fd.id
    					LEFT JOIN cities AS ctd ON d.to_city_id = ctd.id `
	query += ` LIMIT $1 OFFSET $2`
	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting data", err.Error())
		return models.TripsResponse{}, err
	}

	for rows.Next() {
		trip := models.Trip{}
		if err = rows.Scan(
			&trip.ID, &trip.TripNumberID, &trip.FromCityID, &trip.ToCityID, &trip.DriverID, &trip.Price,
			&trip.FromCityData.ID, &trip.FromCityData.Name, &trip.FromCityData.CreatedAt,
			&trip.ToCityData.ID, &trip.ToCityData.Name, &trip.ToCityData.CreatedAt,
			&trip.DriverData.ID, &trip.DriverData.FullName, &trip.DriverData.Phone, &trip.DriverData.FromCityID, &trip.DriverData.ToCityID,
			&trip.DriverData.CreatedAt, &trip.CreatedAt,
			&trip.DriverData.FromCityData.ID, &trip.DriverData.FromCityData.Name, &trip.DriverData.FromCityData.CreatedAt,
			&trip.DriverData.ToCityData.ID, &trip.DriverData.ToCityData.Name, &trip.DriverData.ToCityData.CreatedAt,
		); err != nil {
			fmt.Println("error is while selecting data", err.Error())
			return models.TripsResponse{}, err
		}
		trips = append(trips, trip)
	}

	return models.TripsResponse{
		Trips: trips,
		Count: count,
	}, nil
}

func (c *tripRepo) Update(req models.Trip) (string, error) {
	query := `UPDATE trips SET from_city_id = $1, to_city_id = $2, driver_id = $3, price = $4
             		WHERE id = $5`

	if _, err := c.db.Exec(query, &req.FromCityID, &req.ToCityID, &req.DriverID, &req.Price, req.ID); err != nil {
		fmt.Println("error is while updating trip", err.Error())
		return "", err
	}

	return req.TripNumberID, nil
}

func (c *tripRepo) Delete(id string) error {
	query := `DELETE FROM trip WHERE id = $1`
	if _, err := c.db.Exec(query, id); err != nil {
		fmt.Println("error is while deleting trip", err.Error())
		return err
	}
	return nil
}
