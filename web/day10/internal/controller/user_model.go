package controller

type (
	UserUpsertReq struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		Job  string `json:"job"`
	}

	UserUpsertRes struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		Job  string `json:"job"`
	}
)
