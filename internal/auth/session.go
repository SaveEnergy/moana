package auth

// SessionPayload is the signed cookie content.
type SessionPayload struct {
	UserID int64  `json:"uid"`
	Role   string `json:"role"`
	Exp    int64  `json:"exp"`
}

const cookieName = "moana_session"
