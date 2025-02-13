package main

import (
	"avito-shop/internal/data"
)

func main() {
	data.InitDB()
	defer data.CloseDB()

}
