package service

type (
	UserGetReq struct {
		Id int
	}

	UserGetRes struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		TeamId int    `json:"teamId"`
	}
)

type (
	TeamGetReq struct {
		Id int
	}

	TeamGetRes struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
)

type (
	TeamGetUsersReq struct {
		Id int
	}
	TeamGetUsersRes struct {
		Users []UserGetRes `json:"users"`
	}
)
