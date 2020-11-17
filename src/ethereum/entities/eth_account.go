package entities

type ETHAccount struct {
	Address             string `json:"address"`
	PrivateKey          string `json:"privateKey"`
	PublicKey           string `json:"publicKey"`
	CompressedPublicKey string `json:"compressedPublicKey"`
	Namespace           string `json:"namespace,omitempty"`
}
