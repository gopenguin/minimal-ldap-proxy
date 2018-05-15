package pkg

import (
	"fmt"
	"github.com/gopenguin/minimal-ldap-proxy/pkg/password"
	"github.com/gopenguin/minimal-ldap-proxy/types"
	sql "github.com/jmoiron/sqlx"
	jww "github.com/spf13/jwalterweatherman"
)

func NewBackend(driver string, connString string, authQuery string, searchQuery string) (types.Backend, error) {
	db, err := sql.Open(driver, connString)
	if err != nil {
		return nil, err
	}

	return &sqlBackend{
		db: db,

		authQuery:   authQuery,
		searchQuery: searchQuery,
	}, nil
}

type sqlBackend struct {
	db *sql.DB

	authQuery   string
	searchQuery string
}

func (b *sqlBackend) Authenticate(user string, pw string) bool {
	row := b.db.QueryRow(b.authQuery, user)

	var passwordHash string
	err := row.Scan(&passwordHash)
	if err != nil {
		jww.WARN.Printf("Error fetching pw: %v", err)
		return false
	}

	return password.Verify(pw, passwordHash)
}

func (b *sqlBackend) Search(user string, attributes []string) *types.Result {
	attrs := make(map[string]interface{})

	rows, err := b.db.Queryx(b.searchQuery, user)
	if err != nil {
		jww.WARN.Printf("Error searching user: %v", err)
		return nil
	}
	defer rows.Close()

	result := &types.Result{
		Attributes: make(map[string][]string),
	}

	for rows.Next() {
		err = rows.MapScan(attrs)
		if err != nil {
			jww.WARN.Printf("Error searching user: %v", err)
			continue
		}

		mapBytesToString(attrs)

		for _, ldapAttr := range attributes {
			result.Attributes[ldapAttr] = append(result.Attributes[ldapAttr], fmt.Sprint(attrs[ldapAttr]))
		}
	}

	deduplicateAttributes(result)

	return result
}

func deduplicateAttributes(result *types.Result) {
	for i := range result.Attributes {
		result.Attributes[i] = deduplicateStringSlice(result.Attributes[i])
	}
}

func deduplicateStringSlice(values []string) []string {
	encountered := make(map[string]bool)
	var result []string

	for _, value := range values {
		if encountered[value] {
			// already processed
		} else {
			encountered[value] = true
			result = append(result, value)
		}
	}

	return result
}

func mapBytesToString(m map[string]interface{}) {
	for k, v := range m {
		if b, ok := v.([]byte); ok {
			m[k] = string(b)
		}
	}
}
