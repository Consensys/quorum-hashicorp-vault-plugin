package entities

import "time"

type Key struct {
	ID         string            `json:"id"`
	Curve      string            `json:"curve"`
	Algorithm  string            `json:"algorithm"`
	PrivateKey string            `json:"private_key"`
	PublicKey  string            `json:"public_key"`
	Namespace  string            `json:"namespace,omitempty"`
	Tags       map[string]string `json:"tags,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	Version    int               `json:"version"`
}
