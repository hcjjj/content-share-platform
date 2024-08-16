package web

type LoginSMSReq struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type SignUpReq struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginJWTReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserEditReq struct {
	// 改邮箱，密码，或者能不能改手机号

	Nickname string `json:"nickname"`
	// YYYY-MM-DD
	Birthday string `json:"birthday"`
	AboutMe  string `json:"aboutMe"`
}

type SendSMSCodeReq struct {
	Phone string `json:"phone"`
}
