package pkg

import (
	"database/sql"
	"fmt"
	"github.com/gopenguin/minimal-ldap-proxy/types"
	_ "github.com/lib/pq"
	jww "github.com/spf13/jwalterweatherman"
	"strings"
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
	jww.INFO.Printf("Authenticating %s\n", user)

	row := b.db.QueryRow(b.authQuery, user)

	var passwordHash string
	err := row.Scan(&passwordHash)
	if err != nil {
		jww.WARN.Printf("Error fetching password: %v", err)
		return false
	}

	return password == passwordHash
}

func (b *sqlBackend) Search(user string, attributes map[string]string) []types.Result {
	var ldapAttrs, sqlAttrs []string

	for ldapAttr, sqlAttr := range attributes {
		ldapAttrs = append(ldapAttrs, ldapAttr)
		sqlAttrs = append(sqlAttrs, sqlAttr)
	}

	jww.INFO.Printf("searching for %s with %s", user, strings.Join(ldapAttrs, ", "))

	attr := make([]string, len(attributes))
	attrP := make([]interface{}, len(attributes))
	for i := range attr {
		attrP[i] = &attr[i]
	}

	rows, err := b.db.Query(fmt.Sprintf(b.searchQuery, strings.Join(sqlAttrs, ", ")), user)
	if err != nil {
		jww.WARN.Printf("Error searching user: %v", err)
		return nil
	}
	defer rows.Close()

	var results []types.Result

	for rows.Next() {
		err = rows.Scan(attrP...)
		if err != nil {
			jww.WARN.Printf("Error searching user: %v", err)
			continue
		}

		result := types.Result{
			Attributes: make([]types.Attribute, len(attributes)),
		}
		for i := range attr {
			result.Attributes[i] = types.Attribute{
				Name:  ldapAttrs[i],
				Value: attr[i],
			}
		}

		results = append(results, result)
	}

	return results
}
