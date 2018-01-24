package pkg

import (
	"database/sql"
	"fmt"
	"github.com/gopenguin/minimal-ldap-proxy/types"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

func NewBackend(driver string, connString string, authQuery string, searchQuery string) types.Backend {
	return &sqlBackend{
		driver:     driver,
		connString: connString,

		authQuery:   authQuery,
		searchQuery: searchQuery,
	}
}

type sqlBackend struct {
	driver     string
	connString string

	authQuery   string
	searchQuery string
}

func (b *sqlBackend) Authenticate(user string, password string) bool {
	db, err := sql.Open(b.driver, b.connString)
	if err != nil {
		return false
	}

	defer db.Close()

	row := db.QueryRow(b.authQuery, user)

	var passwordHash string
	row.Scan(&passwordHash)

	fmt.Printf("Authenticating %s", user)

	return password == passwordHash
}

func (b *sqlBackend) Search(user string, attributes map[string]string) []types.Result {
	db, err := sql.Open(b.driver, b.connString)
	if err != nil {
		return nil
	}

	defer db.Close()

	var ldapAttrs, sqlAttrs []string

	for ldapAttr, sqlAttr := range attributes {
		ldapAttrs = append(ldapAttrs, ldapAttr)
		sqlAttrs = append(sqlAttrs, sqlAttr)
	}

	attr := make([]string, len(attributes))
	attrP := make([]interface{}, len(attributes))
	for i := range attr {
		attrP[i] = &attr[i]
	}

	rows, err := db.Query(fmt.Sprintf(b.searchQuery, strings.Join(sqlAttrs, ", ")), user)
	if err != nil {
		return nil
	}
	defer rows.Close()

	var results []types.Result

	for rows.Next() {
		err = rows.Scan(attrP...)
		if err != nil {
			continue
		}

		result := types.Result{
			Attributes: make([]types.Attribute, len(attributes)),
		}
		for i := range attr {
			result.Attributes[i] = types.Attribute{
				Name:  ldapAttrs[i],
				Value: (attr[i]),
			}
		}

		results = append(results, result)
	}

	return results
}
