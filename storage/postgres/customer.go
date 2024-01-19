package postgres

import (
	"city2city/api/models"
	"city2city/storage"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type customerRepo struct {
	db *sql.DB
}

func NewCustomerRepo(db *sql.DB) storage.ICustomerRepo {
	return customerRepo{
		db,
	}
}

func (c customerRepo) Create(customer models.CreateCustomer) (string, error) {
	id := uuid.New()
	query := `INSERT INTO customers(id, full_name, phone, email) values($1, $2, $3, $4)`

	if _, err := c.db.Exec(query, id, customer.FullName, customer.Phone, customer.Email); err != nil {
		fmt.Println("error is while inserting customer", err.Error())
		return "", err
	}
	return id.String(), nil
}

func (c customerRepo) Get(id string) (models.Customer, error) {
	customer := models.Customer{}
	query := `SELECT id, full_name, phone, email, created_at  FROM customers WHERE id = $1`

	if err := c.db.QueryRow(query, id).Scan(
		&customer.ID,
		&customer.FullName,
		&customer.Phone,
		&customer.Email,
		&customer.CreatedAt); err != nil {
		fmt.Println("error is while selecting customer by id", err.Error())
		return models.Customer{}, err
	}
	return customer, nil
}

func (c customerRepo) GetList(req models.GetListRequest) (models.CustomersResponse, error) {
	var (
		page              = req.Page
		offset            = (page - 1) * req.Limit
		query, countQuery string
		count             = 0
		customers         = []models.Customer{}
	)

	countQuery = `select count(1) from customers`

	if err := c.db.QueryRow(countQuery).Scan(&count); err != nil {
		fmt.Println("error is while selecting count", err.Error())
		return models.CustomersResponse{}, err
	}

	query = `SELECT id, full_name, phone, email, created_at FROM customers `

	query += ` LIMIT $1 OFFSET $2`

	rows, err := c.db.Query(query, req.Limit, offset)
	if err != nil {
		fmt.Println("error is while selecting customers", err.Error())
		return models.CustomersResponse{}, err
	}

	for rows.Next() {
		c := models.Customer{}
		if err = rows.Scan(&c.ID, &c.FullName, &c.Phone, &c.Email, &c.CreatedAt); err != nil {
			fmt.Println("error is while selecting customers", err.Error())
			return models.CustomersResponse{}, err
		}
		customers = append(customers, c)
	}
	return models.CustomersResponse{
		Customers: customers,
		Count:     count,
	}, nil
}

func (c customerRepo) Update(customer models.Customer) (string, error) {
	query := `UPDATE customers SET full_name = $1, phone = $2, email = $3 WHERE id = $4`
	if _, err := c.db.Exec(query, &customer.FullName, &customer.Phone, &customer.Email, &customer.ID); err != nil {
		fmt.Println("error is while updating customer", err.Error())
		return "", err
	}
	return customer.ID, nil
}

func (c customerRepo) Delete(id string) error {
	query := `DELETE FROM customers WHERE id = $1`

	if _, err := c.db.Exec(query, id); err != nil {
		fmt.Println("error is while deleting customer", err.Error())
		return err
	}

	return nil
}
