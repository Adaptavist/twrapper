package hcl

type ToHCL interface {
	// Return entity ready for marshalling into terraform configuration.
	ToHCL() (interface{}, error)
}

type TerraformBackend interface {
	ToHCL
	// Validate the entity is in a desired state
	Validate() error
	// Get the workspace key
	GetWorkspace() string
	// Set the workspace key
	SetWorkspace(workspace string)
}
