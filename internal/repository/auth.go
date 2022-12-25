package repository

import (
	"encoding/json"

	"github.com/XVNDEX/blackrocksouth_test/internal/data"
	"github.com/XVNDEX/blackrocksouth_test/internal/entity"
)

func NewCredentials() (*entity.Credentials, error) {
	var logins map[string]string
	err := json.Unmarshal(data.Credentials, &logins)
	if err != nil {
		return nil, err
	}

	return &entity.Credentials{Logins: logins}, nil
}
