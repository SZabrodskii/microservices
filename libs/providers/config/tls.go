package config

type TLSParams struct {
	TLSKey 			   string
	TLSCertificate     string
	TLSRootCertificate string
}

type TLSConfig interface {
	GetKey() string
	GetCertificate() string
	GetRootCertificate() string
}

func (t *TLSParams) GetKey() string {
	return t.TLSKey
}

func (t *TLSParams) GetCertificate() string {
	return t.TLSCertificate
}

func (t *TLSParams) GetRootCertificate() string {
	return t.TLSRootCertificate
}
