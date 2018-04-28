package types

type CmdConfig struct {
	ServerAddress string
	Cert          string
	Key           string

	Driver string
	Conn   string

	AuthQuery   string
	SearchQuery string
	BaseDn      string
	Attributes  []string
	Rdn         string
}

type Result struct {
	Rdn        string
	Attributes map[string][]string
}

type Backend interface {
	Authenticate(username string, password string) bool
	Search(user string, attributes []string) *Result
}
