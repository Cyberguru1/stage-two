package handlers

type registerReq struct {
	Firstname string `json:"firstName"`
	Lastname string `json:"lastName"`
	Email string `json:"email"`
	Password string `json:"password"`
	Phone string `json:"phone"`
}

type orgRegisterReq struct {
	Name string `json:"name"`
	Description string `json:"description"`
}

type orgReq struct {
	Name string `json:"name"`
	OrgId string `json:"orgId"`
	Description string `json:"description"`
}


type orgAddUserReq struct {
	UserId string `json:"userId"`
}

type loginReq struct {
	Email string `json:"email"`
	Password string `json:"password"`
}