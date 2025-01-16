package middleware

import "time"

type Auth interface {
	CreateToken(username string, duration time.Duration) (string, error)
	CreateTokenWithRoles(username string, roles []string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
