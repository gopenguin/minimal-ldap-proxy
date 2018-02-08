package pkg

import (
	"fmt"
	"github.com/gopenguin/minimal-ldap-proxy/types"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/vjeantet/goldap/message"
	ldap "github.com/vjeantet/ldapserver"
)

type Frontend struct {
	serverAddr string
	attributes map[string]string

	server  *ldap.Server
	backend types.Backend
}

func init() {
	ldap.Logger = jww.INFO
}

func NewFrontend(serverAddr string, attributes map[string]string, backend types.Backend) (frontend *Frontend) {
	frontend = &Frontend{
		serverAddr: serverAddr,
		attributes: attributes,
		server:     ldap.NewServer(),
		backend:    backend,
	}

	router := ldap.NewRouteMux()
	router.Bind(frontend.handleBind)
	router.Search(frontend.handleUserSearch).
		BaseDn("")

	frontend.server.Handle(router)

	return frontend
}

func (f *Frontend) handleBind(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetBindRequest()
	res := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
	if r.AuthenticationChoice() == "simple" {
		user := string(r.Name())
		password := string(r.AuthenticationSimple())

		if f.backend.Authenticate(user, password) {
			res.SetResultCode(ldap.LDAPResultSuccess)
		}
	}

	w.Write(res)
}

func (f *Frontend) handleUserSearch(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetSearchRequest()

	user, err := f.extractUser(r.Filter())
	if err != nil {
		jww.WARN.Printf("extract user: %v", err)
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultNoSuchAttribute)
		w.Write(res)
		return
	}

	filteredAttributes := f.filterAttributes(r.Attributes())

	results := f.backend.Search(user, filteredAttributes)

	for _, result := range results {
		entry := ldap.NewSearchResultEntry(result.Rdn + string(r.BaseObject()))

		for _, attr := range result.Attributes {
			entry.AddAttribute(message.AttributeDescription(attr.Name), message.AttributeValue(attr.Value))
		}

		w.Write(entry)
	}

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func (f *Frontend) Serve() {
	go func() {
		err := f.server.ListenAndServe(f.serverAddr)
		jww.ERROR.Println(err)
	}()
}

func (f *Frontend) Stop() {
	f.server.Stop()
}

func (f *Frontend) filterAttributes(attributes message.AttributeSelection) map[string]string {
	filtered := make(map[string]string)

	for _, attr := range attributes {
		value, ok := f.attributes[string(attr)]
		if ok {
			filtered[string(attr)] = value
		}
	}

	return filtered
}

func (f *Frontend) extractUser(filter message.Filter) (user string, err error) {
	switch filter.(type) {
	case message.FilterEqualityMatch:
		eq := filter.(message.FilterEqualityMatch)
		return string(eq.AssertionValue()), nil
	default:
		return "", fmt.Errorf("filter of type %T not supported", filter)
	}
}
