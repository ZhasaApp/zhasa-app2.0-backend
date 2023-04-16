package entities

type Avatar struct {
	Url string
}

type User struct {
	Id       int
	Email    Email
	Password Password
	Avatar   Avatar
}
