package credentials

type TerraformCredential struct {
	Token string `json:"token"`
}

type TerraformCredentials struct {
	Credentials map[string]TerraformCredential `json:"credentials"`
}

func (c *TerraformCredentials) Get(name string) (TerraformCredential, bool) {
	cred, ok := c.Credentials[name]
	return cred, ok
}

func New() TerraformCredentials {
	return TerraformCredentials{}
}
