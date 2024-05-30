package util

import (
	"database/sql/driver"
	"errors"
	"math/big"
)

type BigFloat big.Float

func (b *BigFloat) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return errors.New("failed to scan big float")
	}

	v := big.NewFloat(0)
	v.SetString(str)
	*b = BigFloat(*v)
	return nil
}

func (b BigFloat) Value() (driver.Value, error) {
	v := big.Float(b)
	return v.String(), nil
}

func (b BigFloat) String() string {
	v := big.Float(b)
	return v.String()
}

func (b BigFloat) Float() *big.Float {
	v := big.Float(b)
	return &v
}
