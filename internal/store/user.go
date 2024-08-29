package store

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserParams struct {
	Username string `form:"username" binding:"required,alphanum,min=3,max=30"`
	Password string `form:"password" binding:"required,min=3,max=30"`
}
