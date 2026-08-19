package domain_models

type User struct {
	ID       int64
	Email    string
	PassHash []byte
}
