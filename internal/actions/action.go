package actions

type Action interface {
}

var actions = map[string]Action{}

func Exists(name string) bool {
	return actions[name] != nil
}
