package pkg

import (
	"fmt"
	"github.com/gopenguin/minimal-ldap-proxy/types"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/vjeantet/goldap/message"
	ldap "github.com/vjeantet/ldapserver"
	"strings"
)

type Frontend struct {
	serverAddr    string
	attributes    []string
	attributesMap map[string]bool

	baseDn string
	rDn    string

	server  *ldap.Server
	backend types.Backend
}

func init() {
	ldap.Logger = jww.INFO
}

func NewFrontend(serverAddr string, baseDn string, rDn string, attributes []string, backend types.Backend) (frontend *Frontend) {
	frontend = &Frontend{
		serverAddr:    serverAddr,
		baseDn:        baseDn,
		rDn:           rDn,
		attributes:    attributes,
		attributesMap: make(map[string]bool),
		server:        ldap.NewServer(),
		backend:       backend,
	}

	for _, attr := range attributes {
		frontend.attributesMap[attr] = true
	}

	router := ldap.NewRouteMux()
	router.Bind(frontend.handleBind)
	router.Search(frontend.handleSearchUser).
		BaseDn(frontend.baseDn)
	router.Search(frontend.handleSearchGeneric)

	frontend.server.Handle(router)

	return frontend
}

func (f *Frontend) handleBind(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetBindRequest()
	res := ldap.NewBindResponse(ldap.LDAPResultInvalidCredentials)
	defer func() { w.Write(res) }()

	if r.AuthenticationChoice() == "simple" {
		dn := string(r.Name())

		user, err := f.userFromDn(dn)
		if err != nil {
			jww.WARN.Printf("Unable to get DN: %v", err)
			return
		}

		jww.INFO.Printf("Authenticating %s\n", user)

		password := string(r.AuthenticationSimple())

		if f.backend.Authenticate(user, password) {
			res.SetResultCode(ldap.LDAPResultSuccess)
		}
	} else {
		jww.INFO.Printf("Unsupported authentication type %s", r.AuthenticationChoice())
	}
}

func (f *Frontend) handleSearchUser(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetSearchRequest()

	filteredAttributes := f.filterAttributes(r.Attributes())

	jww.INFO.Printf("Searching on %s for %s with %s", r.BaseObject(), r.FilterString(), strings.Join(filteredAttributes, ", "))

	user, err := f.userFromFilter(r.Filter())
	if err != nil {
		jww.WARN.Printf("extract user: %v", err)
		res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultNoSuchAttribute)
		w.Write(res)
		return
	}

	results := f.backend.Search(user, filteredAttributes)

	for _, result := range results {
		entry := ldap.NewSearchResultEntry(fmt.Sprintf("%s=%s,%s", f.rDn, result.Attributes[f.rDn], string(r.BaseObject())))

		for key, value := range result.Attributes {
			entry.AddAttribute(message.AttributeDescription(key), message.AttributeValue(value))
		}

		w.Write(entry)
	}

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultSuccess)
	w.Write(res)
}

func (f *Frontend) handleSearchGeneric(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetSearchRequest()

	jww.INFO.Printf("Unhandled search request: %s", r.BaseObject())

	res := ldap.NewSearchResultDoneResponse(ldap.LDAPResultNoSuchObject)
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

func (f *Frontend) filterAttributes(attributes message.AttributeSelection) []string {
	// if no attributes are selected, return all attributes by default
	if len(attributes) == 0 {
		return f.attributes
	}

	filtered := []string{f.rDn}

	for _, attr := range attributes {
		_, ok := f.attributesMap[string(attr)]
		if ok {
			filtered = append(filtered, string(attr))
		}
	}

	return filtered
}

func (f *Frontend) userFromDn(dn string) (user string, err error) {
	prefix := f.rDn + "="
	suffix := "," + f.baseDn

	if strings.HasPrefix(dn, prefix) && strings.HasSuffix(dn, suffix) {
		return strings.TrimSuffix(strings.TrimPrefix(dn, prefix), suffix), nil
	}

	return "", fmt.Errorf("dn must have a prefix of '%s' and suffix of '%s'", prefix, suffix)
}

func (f *Frontend) userFromFilter(filter message.Filter) (user string, err error) {
	switch filter.(type) {
	case message.FilterEqualityMatch:
		eq := filter.(message.FilterEqualityMatch)
		if string(eq.AttributeDesc()) != f.rDn {
			return "", fmt.Errorf("invalid rdn '%s', should be '%s'", string(eq.AttributeDesc()), f.rDn)
		}

		return string(eq.AssertionValue()), nil
	default:
		return "", fmt.Errorf("filter '%T' not supported", filter)
	}
}
