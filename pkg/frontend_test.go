package pkg

import (
	"github.com/go-ldap/ldap"
	"github.com/gopenguin/minimal-ldap-proxy/types"
	"github.com/stretchr/testify/assert"
	"github.com/vjeantet/ldapserver"
	"testing"
	"time"
)

var _ types.Backend = (*testBackend)(nil)

type testBackend struct {
	username   string
	password   string
	bindResult bool

	attributes   map[string]string
	searchResult []types.Result
}

func (t *testBackend) Authenticate(username string, password string) bool {
	t.username = username
	t.password = password

	return t.bindResult
}

func (t *testBackend) Search(user string, attributes map[string]string) []types.Result {
	t.username = user
	t.attributes = attributes

	return t.searchResult
}

func TestFrontend_handleBind(t *testing.T) {
	withLdapServerAndClient(t, map[string]string{}, func(t *testing.T, backend *testBackend, client *ldap.Conn) {
		backend.bindResult = true
		err := client.Bind("cn=username,ou=People,dc=example,dc=com", "password")
		assert.NoError(t, err)
		assert.Equal(t, "username", backend.username)
		assert.Equal(t, "password", backend.password)

		backend.bindResult = false
		err = client.Bind("cn=username,ou=People,dc=example,dc=com", "password")
		assert.EqualError(t, err, "LDAP Result Code 49 \"Invalid Credentials\": ")
	})
}

func TestFrontend_handleUserSearch(t *testing.T) {
	withLdapServerAndClient(t, map[string]string{"attr1": "a1", "attr2": "a2", "attr3": "a3"}, func(t *testing.T, backend *testBackend, client *ldap.Conn) {
		result, err := client.Search(&ldap.SearchRequest{
			BaseDN: "ou=People,dc=example,dc=com",
			Filter: "(objectClass=*)",
		})

		assert.EqualError(t, err, "LDAP Result Code 16 \"No Such Attribute\": ")

		backend.searchResult = []types.Result{
			{
				Attributes: map[string]string{
					"cn":    "abc",
					"attr2": "def",
					"attr3": "ghi",
				},
			},
		}

		result, err = client.Search(&ldap.SearchRequest{
			BaseDN:     "ou=People,dc=example,dc=com",
			Attributes: []string{"attr2", "attr3"},
			Filter:     "(cn=abc)",
		})

		assert.NoError(t, err)
		assert.Equal(t, "abc", backend.username)
		assert.Equal(t, map[string]string{"attr2": "a2", "attr3": "a3"}, backend.attributes)
		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "cn=abc,ou=People,dc=example,dc=com", result.Entries[0].DN)
		assert.Len(t, result.Entries[0].Attributes, 3)
	})
}

func TestFrontend_userFromDn(t *testing.T) {
	f := &Frontend{
		rDn:    "cn",
		baseDn: "ou=People,dc=example,dc=com",
	}
	errorMsg := "dn must have a prefix of 'cn=' and suffix of ',ou=People,dc=example,dc=com'"

	tests := []struct {
		name     string
		value    string
		result   string
		errorMsg string
	}{
		{
			name:   "Full match",
			value:  "cn=user1,ou=People,dc=example,dc=com",
			result: "user1",
		},
		{
			name:     "Prefix missing",
			value:    "user1,ou=People,dc=example,dc=com",
			errorMsg: errorMsg,
		},
		{
			name:     "Suffix missing",
			value:    "cn=user1",
			errorMsg: errorMsg,
		},
		{
			name:     "Prefix and suffix missing",
			value:    "user1",
			errorMsg: errorMsg,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			result, err := f.userFromDn(test.value)

			if test.errorMsg != "" {
				assert.EqualError(t, err, test.errorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.result, result)
			}
		})
	}
}

func withLdapServerAndClient(t *testing.T, attrs map[string]string, inner func(t *testing.T, backend *testBackend, client *ldap.Conn)) {
	backend := &testBackend{}
	frontend := NewFrontend("127.0.0.1:0", "ou=People,dc=example,dc=com", "cn", attrs, backend)
	frontend.Serve()
	defer frontend.Stop()

	if !waitListenerReady(frontend.server, 2*time.Second) {
		t.Errorf("server not ready after 2 seconds")
		return
	}

	client, err := ldap.Dial("tcp", frontend.server.Listener.Addr().String())
	assert.Nil(t, err)
	defer client.Close()

	inner(t, backend, client)
}

func waitListenerReady(server *ldapserver.Server, d time.Duration) bool {
	success := make(chan bool)
	cancle := make(chan bool)

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)

			select {
			case <-cancle:
				return
			default:
			}

			func() {

				defer func() {
					recover()
				}()

				server.Listener.Addr()
				close(success)
			}()
		}
	}()

	select {
	case <-success:
		return true
	case <-(time.NewTimer(d)).C:
	}

	close(cancle)

	return false
}
