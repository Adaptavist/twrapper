package hcl

import (
	"log"
)

type Terraform struct {
	Backend ToHCL `json:"backend"`
	Cloud   ToHCL `json:"cloud"`
}

func (t Terraform) ToHCL() (interface{}, error) {
	hclEntity := map[string]interface{}{}

	var backend map[string]interface{}
	var cloud map[string]interface{}
	// we may need to run some checks

	if t.Backend != nil {
		if b, err := t.Backend.ToHCL(); err == nil {
			backend = b.(map[string]interface{})
		} else {
			return nil, err
		}
	}

	if t.Cloud != nil {
		if c, err := t.Cloud.ToHCL(); err == nil {
			cloud = c.(map[string]interface{})
		} else {
			return nil, err
		}
	}

	// Having both a backend and a cloud conflict isn't going to work so lets deal with it.
	if backend != nil && cloud != nil {
		// We need to double check there is a workspace set
		if _, ok := cloud["workspaces"]; ok {
			log.Println("prioritising cloud over backend")
			backend = nil
		}
	}

	if backend != nil {
		hclEntity["backend"] = backend
	}

	if cloud != nil {
		if _, ok := cloud["workspaces"]; ok {
			hclEntity["cloud"] = cloud
		}
	}

	return hclEntity, nil
}
