package middleware

type JWTError struct {
	Type    string
	Message string
}

func (e *JWTError) Error() string {
	return e.Message
}

var (
	ErrNoHeader      = &JWTError{Type: "no_header", Message: "no authorization header"}
	ErrInvalidFormat = &JWTError{Type: "invalid_format", Message: "invalid authorization header format"}
	ErrTokenExpired  = &JWTError{Type: "expired", Message: "token has expired"}
	ErrInvalidToken  = &JWTError{Type: "invalid", Message: "invalid token"}
)
