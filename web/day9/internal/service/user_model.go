package service

type (
	UserAddReq struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		Job  string `json:"job"`
	}

	UserAddRes struct{}
)

type (
	UserGetReq struct {
		Name string `form:"name"`
	}

	UserGetRes struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
		Job  string `json:"job"`
	}
)
