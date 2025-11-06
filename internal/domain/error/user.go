package error

var (
	ErrUserNotFound      = New("user not found")
	ErrUserAlreadyExists = New("user already exists").SetCode(409)
)
