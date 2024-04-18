package service

type (
	UserGetReq struct {
		Id int `uri:"id"`
	}

	UserGetRes struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		TeamId int    `json:"teamId"`
	}
)

type (
	TeamGetReq struct {
		Id int `uri:"id"`
	}

	TeamGetRes struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
)
