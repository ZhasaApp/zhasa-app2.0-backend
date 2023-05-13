package entities

type CreateUserBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type RequestAuthCodeBody struct {
	Phone string `json:"phone"`
}

type RequestAuthCodeResponse struct {
	OtpId int32 `json:"otp_id"`
}

type TryAuthBody struct {
	OtpId int32 `json:"otp_id"`
	Otp   int32 `json:"otp"`
}

type AuthResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Token     string `json:"token"`
}
