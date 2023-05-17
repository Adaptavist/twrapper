package runner

type Variables map[string]interface{}

func (v Variables) HasVars() bool {
	return len(v) > 0
}
