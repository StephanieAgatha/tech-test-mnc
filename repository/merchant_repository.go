package repository

import (
	"database/sql"
	"fmt"
	"mnc-test/model"
)

type MerchantRepository interface {
	CreateNewMerchant(merchant model.Merchant) error
	FindAllMerchant() ([]model.Merchant, error)
}

type merchantRepository struct {
	db *sql.DB
}

func (m *merchantRepository) CreateNewMerchant(merchant model.Merchant) error {
	//TODO implement me
	query := "insert into merchants (name) values ($1)"

	_, err := m.db.Exec(query, merchant.Name)
	if err != nil {
		return fmt.Errorf("Failed to exec query %v", err.Error())
	}

	return nil
}

func (m *merchantRepository) FindAllMerchant() ([]model.Merchant, error) {
	//TODO implement me
	var merchants []model.Merchant

	query := "select * from merchants"

	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("Failed to find merchants %v", err.Error())
	}

	for rows.Next() {
		var merchant model.Merchant
		if err = rows.Scan(&merchant.ID, &merchant.Name, &merchant.CreatedAt, &merchant.UpdatedAt); err != nil {
			return nil, fmt.Errorf("Failed to retrieve merchants %v", err.Error())
		}
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

func NewMerchantRepository(db *sql.DB) MerchantRepository {
	return &merchantRepository{
		db: db,
	}
}
