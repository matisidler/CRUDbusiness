package main

import (
	"fmt"

	"github.com/matisidler/CRUDbusiness/pkg/products"
	"github.com/matisidler/CRUDbusiness/storage"
)

func main() {
	driver := storage.MySQL
	storage.NewConnection(driver)

	storageSale, err := storage.DAOProduct(driver)
	if err != nil {
		fmt.Println(err)
	} else {
		err = storageSale.Update(&products.Model{
			ID:    4,
			Price: 150,
			Name:  "Bombilla",
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}
