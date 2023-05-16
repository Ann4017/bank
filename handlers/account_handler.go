package handlers

import (
	"bank/data"
	"bank/db"
	"bank/util"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type C_account_handler struct {
	C_db *db.C_db
}

func (c *C_account_handler) Insert_account(w http.ResponseWriter, r *http.Request) {
	account := data.C_account{}
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account_num := util.Generate_unique_account_num()

	stmt, err := c.C_db.PC_sql_db.Prepare("insert into account (user_seq, account_num) value (?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.I_user_seq, account_num)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	account.S_account_num = account_num

	json.NewEncoder(w).Encode(account)
}

func (c *C_account_handler) Select_account(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seq, err := strconv.Atoi(vars["user_seq"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account := data.C_account{}

	stmt, err := c.C_db.PC_sql_db.Prepare("select account_num, balance from account where user_seq = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	err = stmt.QueryRow(seq).Scan(&account.S_account_num, &account.I_balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(account)
}
