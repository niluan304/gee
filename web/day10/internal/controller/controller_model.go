package controller

import "gee/web/day10/internal/service"

type (
	TeamGetUsersReq struct {
		Id int `uri:"id"`
	}
	TeamGetUsersRes struct {
		*service.TeamGetUsersRes
	}
)
