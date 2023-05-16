package handlers

import (
	"bank/data"
	"bank/db"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type C_transaction_handler struct {
	C_db *db.C_db
}

func (c *C_transaction_handler) Account_transfer(w http.ResponseWriter, r *http.Request) {
	transfer := data.C_transfer{}
	err := json.NewDecoder(r.Body).Decode(&transfer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var balance int
	row := c.C_db.PC_sql_db.QueryRow("select balance from account where account_num = ?", transfer.S_from_account_num)
	err = row.Scan(&balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if balance < transfer.I_amount {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := c.C_db.PC_sql_db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("update account set balance = balance - ? where account_num = ?", transfer.I_amount, transfer.S_from_account_num)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("update account set balance = balance + ? where account_num = ?", transfer.I_amount, transfer.S_to_account_num)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("insert into transaction (account_num, type, amount, date) values (?, 'withdraw', ?, ?)", transfer.S_from_account_num, transfer.I_amount, time.Now())
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("insert into transaction (account_num, type, amount, date) values (?, 'deposit', ?, ?)", transfer.S_to_account_num, transfer.I_amount, time.Now())
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *C_transaction_handler) Select_transaction_history(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	num, err := strconv.Atoi(vars["account_num"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := c.C_db.PC_sql_db.Prepare("select account_num, type, amount, date from transaction where account_num = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(num)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	transactions := []data.C_transaction{}

	for rows.Next() {
		transaction := data.C_transaction{}
		err := rows.Scan(&transaction.S_account_num, &transaction.S_type, &transaction.I_amount, &transaction.T_date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}
