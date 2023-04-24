package entities

type Avatar struct {
	Url string
}

type User struct {
	Id     int
	Phone  Phone
	Avatar Avatar
}
