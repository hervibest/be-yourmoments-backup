package entity

type Auth struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email" db:"email"`
	PhoneNumber string `json:"phone_number"`
	Similarity  uint   `json:"similarity"`
}
