package dto

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	RegistrationID string `json:"registrationID"`
	Message        string `json:"message"`
}

type VerfiyUserRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type VerfiyUserResponse struct {
	UserID  string `json:"userID"`
	Message string `json:"message"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	SessionID string `json:"sessionID"`
	Message   string `json:"message"`
}
