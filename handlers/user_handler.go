package handlers

import (
	"bank/data"
	"bank/db"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type C_user_handler struct {
	C_db *db.C_db
}

func (c *C_user_handler) Insert_user(w http.ResponseWriter, r *http.Request) {
	user := data.C_user{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.S_id == "" || user.S_password == "" || user.S_name == "" || user.S_email == "" || user.S_phone_num == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := c.C_db.PC_sql_db.Prepare("insert into user (id, password, name, email, phone_num) value (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	row, err := stmt.Exec(user.S_id, user.S_password, user.S_name, user.S_email, user.S_phone_num)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	seq, err := row.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.I_seq = int(seq)
	json.NewEncoder(w).Encode(user)
}

func (c *C_user_handler) Update_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seq, err := strconv.Atoi(vars["seq"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := data.C_user{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := c.C_db.PC_sql_db.Prepare("update users set password = ?, name = ?, email = ?, phone_num = ? where seq = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.S_password, user.S_name, user.S_email, user.S_phone_num, seq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (c *C_user_handler) Delete_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	seq, err := strconv.Atoi(vars["seq"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stmt, err := c.C_db.PC_sql_db.Prepare("delete from users where seq = ?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(seq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
