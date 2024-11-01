package entity

import (
	"encoding"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var (
	_ encoding.BinaryMarshaler = (*ProductCategory)(nil)
	_ encoding.BinaryMarshaler = (*Product)(nil)
)

type ProductCategory int

const (
	ProductCategoryElectronics ProductCategory = iota + 1
	ProductCategoryClothing
	ProductCategoryAccessories
	ProductCategoryHome
	ProductCategoryLife
)

func (p ProductCategory) MarshalBinary() (data []byte, err error) {
	return []byte(strconv.Itoa(int(p))), nil
}

type Product struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CategoryID  ProductCategory `json:"category_id"`
	Price       decimal.Decimal `json:"price"`
	Rank        int             `json:"rank"`
	CreatedAt   int64           `json:"created_at"`
	UpdatedAt   int64           `json:"updated_at"`
}

func (p *Product) MarshalBinary() (data []byte, err error) {
	buf, err := json.Marshal(p)
	if err != nil {
		return nil, errors.Errorf("marshal product, err: %+v", err)
	}

	return buf, nil
}

func (p *Product) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}
