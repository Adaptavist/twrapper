package hcl

import "encoding/json"

type Root struct {
	Terraform Terraform `json:"terraform" yaml:"terraform" ommitempty:"true"`
}

func New() *Root {
	return &Root{
		Terraform: Terraform{},
	}
}

func (r *Root) ToHCL() (interface{}, error) {
	terraform, err := r.Terraform.ToHCL()

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"terraform": terraform,
	}, nil
}

func (r *Root) ToJSON() ([]byte, error) {
	hcl, err := r.ToHCL()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(hcl, "", "\t")
}
