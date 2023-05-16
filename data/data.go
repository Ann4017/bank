package data

type C_user struct {
	I_seq       int    `json:"seq"`
	S_id        string `json:"id"`
	S_password  string `json:"password"`
	S_name      string `json:"name"`
	S_email     string `json:"email"`
	S_phone_num string `json:"phone_num"`
}

type C_account struct {
	I_user_seq    int    `json:"user_seq"`
	S_account_num string `json:"account_num"`
	I_balance     int    `json:"balance"`
}

type C_transaction struct {
	I_seq         int    `json:"seq"`
	S_account_num string `json:"account_num"`
	S_type        string `json:"type"`
	I_amount      int    `json:"amount"`
	T_date        string `json:"date"`
}

type C_transfer struct {
	S_from_account_num string `json:"from_account_num"`
	S_to_account_num   string `json:"to_account_num"`
	I_amount           int    `json:"amount"`
}
