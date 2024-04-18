package db

type User struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	TeamId int    `json:"teamId"`
}

type Team struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// Teams 只是一个切片 []Team，用于充当数据库表
var Teams = []Team{
	{Id: 3, Name: "Apple"},
	{Id: 4, Name: "Byte"},
}

// Users 只是一个切片 []User，用于充当数据库表
var Users = []User{
	{Id: 1, Name: "Alice", TeamId: 3},
	{Id: 2, Name: "Bob", TeamId: 4},
}
