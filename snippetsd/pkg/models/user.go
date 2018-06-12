// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package models // import "cirello.io/snippetsd/pkg/models"

import (
	"cirello.io/snippetsd/pkg/errors"
	"github.com/jmoiron/sqlx"
)

// User stores th user of snippets
type User struct {
	ID    int64  `db:"id" json:"id"`
	Email string `db:"email" json:"email"`
}

// UsersDAO provides DB persistence to snippets ssers.
type UsersDAO struct {
	db *sqlx.DB
}

// NewUsersDAO instanties a UsersDAO
func NewUsersDAO(db *sqlx.DB) *UsersDAO {
	return &UsersDAO{db: db}
}

// Bootstrap creates table if missing.
func (b *UsersDAO) Bootstrap() error {
	cmds := []string{
		`create table if not exists users (
			id integer primary key autoincrement,
			email varchar(255)
		);
		`,
		`create index if not exists users_email on users (email)`,
	}

	for _, cmd := range cmds {
		_, err := b.db.Exec(cmd)
		if err != nil {
			return errors.E(err)
		}
	}

	return nil
}

// GetByEmail loads a user by email.
func (b *UsersDAO) GetByEmail(email string) ([]*User, error) {
	var Users []*User
	err := b.db.Select(&Users, "SELECT * FROM users WHERE email = $1", email)
	return Users, errors.E(err)
}

// Insert a user.
func (b *UsersDAO) Insert(user *User) (*User, error) {
	_, err := b.db.NamedExec(`
		INSERT INTO users
		(email)
		VALUES (:email)
	`, user)
	if err != nil {
		return nil, errors.E(err)
	}

	err = b.db.Get(user, `
		SELECT
			*
		FROM
			users
		WHERE
			id = last_insert_rowid()
	`)
	if err != nil {
		return nil, errors.E(err)
	}
	return user, errors.E(err)
}

// Update a user.
func (b *UsersDAO) Update(user *User) error {
	_, err := b.db.NamedExec(`
		UPDATE Users
		SET
			email = :email
		WHERE
			id = :id
	`, user)
	return errors.E(err)
}
