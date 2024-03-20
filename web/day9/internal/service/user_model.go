package service

type (
	UserAddReq struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
		Job  string `form:"job"`
	}

	UserAddRes struct{}
)

type (
	UserGetReq struct {
		Name string `json:"name" form:"name"`
	}

	UserGetRes struct {
		Name string
		Age  int
		Job  string
	}
)
