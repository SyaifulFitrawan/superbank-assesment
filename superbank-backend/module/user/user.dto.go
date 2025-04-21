package user

type UserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}
