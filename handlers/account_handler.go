package handlers

import (
	"bank/data"
	"bank/db"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type C_account_handler struct {
	C_db *db.C_db
}

const account_num_len = 13

func Generate_unique_account_num() string {
	rand.Seed(time.Now().UnixNano())

	set_num := "0123456789"
	set_num_len := len(set_num)

	result := make([]byte, account_num_len)
	for i := 0; i < account_num_len; i++ {
		random_index := rand.Intn(set_num_len)
		result[i] = set_num[random_index]
	}

	return string(result)
}

func (c *C_account_handler) Insert_account(w http.ResponseWriter, r *http.Request) {
	account := data.C_account{}
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	account_num := Generate_unique_account_num()

	stmt, err := c.C_db.PC_sql_db.Prepare("update account set account_num = ? where user_seq = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(account_num, account.I_user_seq)
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
