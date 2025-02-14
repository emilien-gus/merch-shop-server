package services

import "errors"

var items = map[string]int{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}

func GetItem(item string) (int, error) {
	product, ok := items[item]
	if !ok {
		return 0, errors.New("item not found")
	}
	return product, nil
}
