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
		err := client.Bind("username", "password")
		assert.Nil(t, err)
		assert.Equal(t, "username", backend.username)
		assert.Equal(t, "password", backend.password)

		backend.bindResult = false
		err = client.Bind("username", "password")
		assert.EqualError(t, err, "LDAP Result Code 49 \"Invalid Credentials\": ")
	})
}

func TestFrontend_handleUserSearch(t *testing.T) {
	withLdapServerAndClient(t, map[string]string{"attr1": "a1", "attr2": "a2", "attr3": "a3"}, func(t *testing.T, backend *testBackend, client *ldap.Conn) {
		result, err := client.Search(&ldap.SearchRequest{
			Filter: "(objectClass=*)",
		})

		assert.EqualError(t, err, "LDAP Result Code 16 \"No Such Attribute\": ")

		backend.searchResult = []types.Result{
			{
				Rdn: "attr1=abc",
				Attributes: []types.Attribute{
					{Name: "attr2", Value: "def"},
					{Name: "attr3", Value: "ghi"},
				},
			},
		}

		result, err = client.Search(&ldap.SearchRequest{
			Attributes: []string{"attr2", "attr3"},
			Filter:     "(attr1=abc)",
		})

		assert.Nil(t, err)
		assert.Equal(t, "abc", backend.username)
		assert.Equal(t, map[string]string{"attr2": "a2", "attr3": "a3"}, backend.attributes)
		assert.Len(t, result.Entries, 1)
		assert.Equal(t, "attr1=abc", result.Entries[0].DN)
		assert.Len(t, result.Entries[0].Attributes, 2)
	})
}

func withLdapServerAndClient(t *testing.T, attrs map[string]string, inner func(t *testing.T, backend *testBackend, client *ldap.Conn)) {
	backend := &testBackend{}
	frontend := NewFrontend("127.0.0.1:0", attrs, backend)
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
