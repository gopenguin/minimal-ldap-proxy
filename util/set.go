package util

type StringSet map[string]*struct{}

func NewStringSet(values []string) {
	var res = make(map[string]*struct{})
	for _, value := range values {
		res[value] = nil
	}
}

func (ss StringSet) Present(value string) bool {
	_, ok := ss[value]
	return ok
}
