package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type cityRepo struct {
	db *sql.DB
}

func NewCityRepo(db *sql.DB) storage.ICityRepo {
	return cityRepo{
		db,
	}
}

func (c cityRepo) Create(city models.CreateCity) (string, error) {
	id := uuid.New()
	query := `INSERT INTO cities (id, name) values ($1, $2)`

	if _, err := c.db.Exec(query, id, city.Name); err != nil {
		fmt.Println("error is while inserting city", err.Error())
		return "", err
	}

	return id.String(), nil
}

func (c cityRepo) Get(id string) (models.City, error) {
	city := models.City{}
	query := `SELECT * FROM cities WHERE id = $1`

	if err := c.db.QueryRow(query, id).Scan(&city.ID, &city.Name, &city.CreatedAt); err != nil {
		fmt.Println("error is while scanning city", err.Error())
		return models.City{}, err
	}

	return city, nil
}

func (c cityRepo) GetList(req models.GetListRequest) (models.CitiesResponse, error) {
	var (
		page              = req.Page
		offset            = (page - 1) * req.Limit
		countQuery, query string
		count             = 0
		cities            = []models.City{}
	)
	countQuery = `SELECT count(1) FROM cities`

	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while scanning count")
		return models.CitiesResponse{}, err
	}

	query = `SELECT id, name, created_at FROM cities `

	query += ` LIMIT $1 OFFSET $2`
	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting city", err.Error())
		return models.CitiesResponse{}, err
	}

	for rows.Next() {
		city := models.City{}

		if err = rows.Scan(&city.ID, &city.Name, &city.CreatedAt); err != nil {
			fmt.Println("error is while scanning city", err.Error())
			return models.CitiesResponse{}, err
		}

		cities = append(cities, city)
	}

	return models.CitiesResponse{
		Cities: cities,
		Count:  count,
	}, nil
}

func (c cityRepo) Update(city models.City) (string, error) {
	query := `UPDATE cities SET name = $1 WHERE id = $2`

	if _, err := c.db.Exec(query, &city.Name, &city.ID); err != nil {
		fmt.Println("error is while updating city", err.Error())
		return "", err
	}

	return city.ID, nil
}

func (c cityRepo) Delete(id string) error {
	query := `DELETE cities WHERE id = $1`

	if _, err := c.db.Exec(query, id); err != nil {
		fmt.Println("error is while deleting city", err.Error())
		return err
	}
	return nil
}
