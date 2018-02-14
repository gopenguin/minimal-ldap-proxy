package types

type CmdConfig struct {
	ServerAddress string

	Driver string
	Conn   string

	AuthQuery   string
	SearchQuery string
	BaseDn      string
	Attributes  map[string]string
	Rdn         string
}

type Result struct {
	Rdn        string
	Attributes []Attribute
}

type Attribute struct {
	Name  string
	Value string
}

type Backend interface {
	Authenticate(username string, password string) bool
	Search(user string, attributes map[string]string) []Result
}
