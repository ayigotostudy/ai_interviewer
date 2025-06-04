package req

type LoginReq struct {
	Email    string `json:"email" `
	Password string `json:"password" `
}

type RegisterReq struct {
	Email    string `json:"email" `
	Password string `json:"password" `
}
