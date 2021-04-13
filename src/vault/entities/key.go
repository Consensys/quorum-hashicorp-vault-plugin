package entities

import "time"

type Key struct {
	ID         string            `json:"id"`
	Curve      string            `json:"curve"`
	Algorithm  string            `json:"algorithm"`
	PrivateKey string            `json:"privateKey"`
	PublicKey  string            `json:"publicKey"`
	Namespace  string            `json:"namespace,omitempty"`
	Tags       map[string]string `json:"tags,omitempty"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
	Version    int               `json:"version"`
}
