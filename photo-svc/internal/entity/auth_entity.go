package entity

type Auth struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email" db:"email"`
	PhoneNumber   string `json:"phone_number"`
	UserProfileID string `json:"user_profile_id"`
	Similarity    uint   `json:"similarity"`
}
