package terraform

import "reflect"

type Root struct {
	Terraform Terraform `json:"terraform"`
}

type Terraform struct {
	Backend Backend `json:"backend"`
}

type Backend map[string]map[string]interface{}

func (b Backend) IsEmpty() bool {
	return reflect.DeepEqual(b, Backend{})
}

type Variables map[string]interface{}

func (v Variables) IsEmpty() bool {
	return reflect.DeepEqual(v, Variables{})
}

type Arguments []string

func (a Arguments) IsEmpty() bool {
	return reflect.DeepEqual(a, Arguments{})
}
