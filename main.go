package main

import (
	"bank/db"
	"bank/handlers"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	c_db := &db.C_db{}

	err := c_db.Load_config("config/config.ini")
	if err != nil {
		fmt.Println(err)
	}

	err = c_db.Connect_db()
	if err != nil {
		fmt.Println(err)
	}

	u := &handlers.C_user_handler{
		C_db: c_db,
	}

	a := &handlers.C_account_handler{
		C_db: c_db,
	}

	t := &handlers.C_transaction_handler{
		C_db: c_db,
	}

	r := mux.NewRouter()

	r.HandleFunc("/user", u.Insert_user).Methods(http.MethodPost)
	r.HandleFunc("/user/{seq}", u.Update_user).Methods(http.MethodPut)
	r.HandleFunc("/user/{seq}", u.Delete_user).Methods(http.MethodDelete)

	r.HandleFunc("/account", a.Insert_account).Methods(http.MethodPost)
	r.HandleFunc("/account/{user_seq}", a.Select_account).Methods(http.MethodGet)

	r.HandleFunc("/transaction", t.Account_transfer).Methods(http.MethodPost)
	r.HandleFunc("/transaction/{account_num}", t.Select_transaction_history).Methods(http.MethodGet)

	http.ListenAndServe(":8000", r)
}
