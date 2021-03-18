package terraform

type Opts struct {
	Backend   Backend
	Variables Variables
	Arguments Arguments
}

func NewOpts() Opts {
	return Opts{
		Backend:   Backend{},
		Variables: Variables{},
		Arguments: Arguments{},
	}
}

func (o *Opts) WithBackend(b Backend) *Opts {
	o.Backend = b
	return o
}

func (o *Opts) WithVariables(v Variables) *Opts {
	o.Variables = v
	return o
}

func (o *Opts) WithArguments(a Arguments) *Opts {
	o.Arguments = a
	return o
}
