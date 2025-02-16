package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetItem_Success(t *testing.T) {
	price, err := GetItem("t-shirt")
	assert.NoError(t, err)
	assert.Equal(t, 80, price)

	price, err = GetItem("cup")
	assert.NoError(t, err)
	assert.Equal(t, 20, price)

	price, err = GetItem("pink-hoody")
	assert.NoError(t, err)
	assert.Equal(t, 500, price)
}

func TestGetItem_NotFound(t *testing.T) {
	price, err := GetItem("non-existent-item")
	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
	assert.Equal(t, 0, price)
}
