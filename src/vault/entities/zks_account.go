package entities

const (
	ZksCurveBN256     = "bn256"
	ZksAlgorithmEDDSA = "eddsa"
)

type ZksAccount struct {
	Curve      string   `json:"curve"`
	Algorithm  string   `json:"algorithm"`
	PrivateKey string   `json:"privateKey"`
	PublicKey  string   `json:"publicKey"`
	Namespace  string   `json:"namespace,omitempty"`
}
