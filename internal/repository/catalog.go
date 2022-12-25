package repository

import (
	"encoding/json"

	"github.com/XVNDEX/blackrocksouth_test/internal/data"
	"github.com/XVNDEX/blackrocksouth_test/internal/entity"
)

func NewCatalog() (*entity.Catalog, error) {
	products := new(entity.Catalog)
	err := json.Unmarshal(data.Catalog, products)
	if err != nil {
		return nil, err
	}

	return products, nil
}
