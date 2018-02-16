package pkg

import (
	"testing"

	"github.com/gopenguin/minimal-ldap-proxy/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestSqlBackend_Authenticate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error during db setup: %v", err)
	}

	defer db.Close()

	backend := &sqlBackend{
		db:        db,
		authQuery: "SELECT password FROM user WHERE name = ?",
	}

	mock.ExpectQuery("SELECT password FROM user WHERE name = ?").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"password"}).AddRow("{SSHA}RrAeHR4zMHdNUfvtEibV9yTbtmMY7nF/"))

	result := backend.Authenticate("username", "test123")

	assert.True(t, result)
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestSqlBackend_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error during db setup: %v", err)
	}

	defer db.Close()

	backend := &sqlBackend{
		db:          db,
		searchQuery: "SELECT %s FROM user WHERE name = ?",
	}

	mock.ExpectQuery("SELECT attr1, attr3 FROM user WHERE name = ?").WithArgs("username").WillReturnRows(sqlmock.NewRows([]string{"attr1", "attr3"}).AddRow("a", "b"))

	result := backend.Search("username", map[string]string{"ldap1": "attr1", "ldap2": "attr3"})

	assert.EqualValues(t, []types.Result{
		{
			Attributes: map[string]string{
				"ldap1": "a",
				"ldap2": "b",
			},
		},
	}, result)

	assert.Nil(t, mock.ExpectationsWereMet())
}
