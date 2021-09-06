package fputil

type Rule interface {
	Validate(value interface{}) bool
}

type RuleFunc func(value interface{}) bool

func (r RuleFunc) Validate(value interface{}) bool {
	return r(value)
}
