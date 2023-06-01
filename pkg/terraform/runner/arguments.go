package runner

type Arguments []string

func (a Arguments) HasArgs() bool {
	return len(a) > 0
}
