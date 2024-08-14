//go:generate easyjson -all user.go
package hw10programoptimization

//easyjson:json
type User struct {
	//nolint:tagliatelle
	ID int `json:"Id"`
	//nolint:tagliatelle
	Name string `json:"Name"`
	//nolint:tagliatelle
	Username string `json:"Username"`
	//nolint:tagliatelle
	Email string `json:"Email"`
	//nolint:tagliatelle
	Phone string `json:"Phone"`
	//nolint:tagliatelle
	Password string `json:"Password"`
	//nolint:tagliatelle
	Address string `json:"Address"`
}
