package pkg

import (
	"fmt"
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

func (b *sqlBackend) Authenticate(user string, password string) bool {
	row := b.db.QueryRow(b.authQuery, user)

	var passwordHash string
	err := row.Scan(&passwordHash)
	if err != nil {
		jww.WARN.Printf("Error fetching password: %v", err)
		return false
	}

	return Verify(password, passwordHash)
}

func (b *sqlBackend) Search(user string, attributes []string) []types.Result {
	attrs := make(map[string]interface{})

	rows, err := b.db.Queryx(b.searchQuery, user)
	if err != nil {
		jww.WARN.Printf("Error searching user: %v", err)
		return nil
	}
	defer rows.Close()

	var results []types.Result

	for rows.Next() {
		err = rows.MapScan(attrs)
		if err != nil {
			jww.WARN.Printf("Error searching user: %v", err)
			continue
		}

		mapBytesToString(attrs)

		result := types.Result{
			Attributes: make(map[string]string),
		}
		for _, ldapAttr := range attributes {
			result.Attributes[ldapAttr] = fmt.Sprint(attrs[ldapAttr])
		}

		results = append(results, result)
	}

	return results
}

func mapBytesToString(m map[string]interface{}) {
	for k, v := range m {
		if b, ok := v.([]byte); ok {
			m[k] = string(b)
		}
	}
}
