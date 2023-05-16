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

type BranchResponse struct {
	Id          int32  `json:"id"`
	Description string `json:"description"`
}

type UserProfileResponse struct {
	Id        int32          `json:"id"`
	Avatar    string         `json:"avatar"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Phone     string         `json:"phone"`
	Branch    BranchResponse `json:"branch"`
	Role      string         `json:"role"`
}

type AuthResponse struct {
	UserProfileResponse UserProfileResponse `json:"user_profile_response"`
	Token               string              `json:"token"`
}
