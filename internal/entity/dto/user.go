package dto

// INPUT DTO (REQUEST BODY)
type UserInput struct {
	Username  string `json:"username" valid:"username"`
	FirstName string `json:"firstName" valid:"firstName"`
	LastName  string `json:"lastName" valid:"lastName"`
}

// OUTPUT DTO (RESPONSE BODY)
type UserOutPut struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	CreatedAt string `json:"createdAt"`
}
