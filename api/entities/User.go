package entities

type CreateUserBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
}

type RequestAuthCodeBody struct {
	Phone string `json:"phone_number"`
}

type RequestAuthCodeResponse struct {
	OtpId int32 `json:"otp_id"`
}

type TryAuthBody struct {
	OtpId int32 `json:"otp_id"`
	Otp   int32 `json:"code"`
}

type BranchResponse struct {
	Id          int32   `json:"id"`
	Description string  `json:"description"`
	Brands      []Brand `json:"brands"`
}

type UserProfileResponse struct {
	Id       int32             `json:"id"`
	Avatar   *string           `json:"avatar"`
	FullName string            `json:"full_name"`
	Phone    string            `json:"phone"`
	Branch   *BranchResponse   `json:"branch"`
	Role     string            `json:"role"`
	Branches *[]BranchResponse `json:"branches"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
