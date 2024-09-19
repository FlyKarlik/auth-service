package queries

// Refresh Tokens Queries
const (
	CreateRefreshTokenQuery = "INSERT INTO refresh_tokens (id,user_id,refresh_hash,updated_at) VALUES ($1,$2,$3,$4)"
	GetRefreshTokenQuery    = "SELECT id,user_id,refresh_hash,updated_at FROM refresh_tokens WHERE id=$1"
	UpdateRefreshTokenQuery = "UPDATE refresh_tokens SET refresh_hash=$1,updated_at=$2 WHERE id=$3"
)
