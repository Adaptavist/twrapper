package runner

import (
	"log"

	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform/model/hcl"
)

type Runner struct {
	Config    *hcl.Root // Terraform configuration E.G Backends
	Arguments Arguments // Tarraform arguments E.G plan, apply
	Variables Variables // Terraform variables to be places into terraform.tfvars.json
}

// New instance of runner
func New() Runner {
	return Runner{
		Config:    hcl.New(),
		Variables: Variables{},
		Arguments: Arguments{},
	}
}

// WithBackend adds backend configuration to the runner
func (o *Runner) WithBackend(b hcl.TerraformBackend) *Runner {
	o.Config.Terraform.Backend = b
	return o
}

// WithCloud adds Terraform Cloud configuration to the runner
func (o *Runner) WithCloud(c *hcl.Cloud) *Runner {
	o.Config.Terraform.Cloud = c
	return o
}

func (o *Runner) GetCloud() *hcl.Cloud {
	if c := o.Config.Terraform.Cloud; c != nil {
		return o.Config.Terraform.Cloud.(*hcl.Cloud)
	}
	return nil
}

// WithVariables adds terraform variables
func (o *Runner) WithVariables(v Variables) *Runner {
	o.Variables = v
	return o
}

// WithArguments adds terraform CLI arguments to the runner
func (o *Runner) WithArguments(a Arguments) *Runner {
	o.Arguments = a
	return o
}

// writeTerraformFile creates terraform configuration file
func (o *Runner) writeTerraformFile() error {
	data, err := o.Config.ToHCL()

	if err != nil {
		return err
	}

	return writeJSONFile("terraform.tf.json", data)
}

// writeVarFile creates terraform variables file
func (o *Runner) writeVarFile() error {
	if len(o.Variables) > 0 {
		return writeJSONFile("terraform.tfvars.json", o.Variables)
	}
	return nil
}

// Configure writes terraform files
func (o *Runner) Configure() (err error) {
	log.Println("writing terraform.tf.json")
	if err = o.writeTerraformFile(); err != nil {
		return
	}
	err = o.writeVarFile()
	return
}

// Go runs terraform with args
func (r *Runner) Go() error {
	return run("terraform", r.Arguments...)
}
