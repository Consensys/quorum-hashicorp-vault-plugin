package entities

type ZksAccount struct {
	Curve      string `json:"curve"`
	Algorithm  string `json:"algorithm"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Namespace  string `json:"namespace,omitempty"`
}
