package main

import (
	"bank/db"
)

func main() {
	db := db.C_db{}

	db.Load_config("config/config.ini")
}
