package utils

type ToolConf struct {
	Cert CertConf
}

type CertConf struct {
	KeyType string

	CommonName   string
	Organization []string
	Country      []string

	NotBefore int
	NotAfter  int

	DNSNames       []string
	EmailAddresses []string
	IPAddress      []string

	PermittedDNSDomains     []string
	CRLDistributionPoints   []string
	SM2SignatureAlgorithm   int
	ECDSASignatureAlgorithm int
}
