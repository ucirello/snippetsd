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
	"time"

	"cirello.io/snippetsd/pkg/errors"
	"github.com/jmoiron/sqlx"
)

// Snippet stores the basic information of a snippet.
type Snippet struct {
	ID     int64      `db:"id" json:"id"`
	UserID string     `db:"user_id" json:"user_id"`
	When   *time.Time `db:"when" json:"when"`
}

// SnippetDAO provides DB persistence to Snippets.
type SnippetDAO struct {
	db *sqlx.DB
}

// NewSnippetDAO instanties a snippetDAO
func NewSnippetDAO(db *sqlx.DB) *SnippetDAO {
	return &SnippetDAO{db: db}
}

// Bootstrap creates table if missing.
func (b *SnippetDAO) Bootstrap() error {
	cmds := []string{
		`create table if not exists snippets (
			id integer primary key autoincrement,
			user_id bigint,
			when datetime
		);
		`,
		`create index if not exists snippets_user_id  on snippets (user_id)`,
		`create index if not exists snippets_when on snippets (when)`,
	}

	for _, cmd := range cmds {
		_, err := b.db.Exec(cmd)
		if err != nil {
			return errors.E(err)
		}
	}

	return nil
}

// All returns all known snippets.
func (b *SnippetDAO) All() ([]*Snippet, error) {
	var snippets []*Snippet
	err := b.db.Select(&snippets, "SELECT * FROM snippets")
	return snippets, errors.E(err)
}

// Current returns the current week snippets.
func (b *SnippetDAO) Current() ([]*Snippet, error) {
	var snippets []*Snippet
	// TODO: convert 7 to variable representing the number of days
	err := b.db.Select(&snippets, "SELECT * FROM snippets WHERE when > 7")
	return snippets, errors.E(err)
}

// Insert one snippet entry.
func (b *SnippetDAO) Insert(snippet *Snippet) (*Snippet, error) {
	_, err := b.db.NamedExec(`
		INSERT INTO snippets
		(user_id, when)
		VALUES (:user_id, :when)
	`, snippet)
	if err != nil {
		return nil, errors.E(err)
	}

	err = b.db.Get(snippet, `
		SELECT
			*
		FROM
			snippets
		WHERE
			id = last_insert_rowid()
	`)
	if err != nil {
		return nil, errors.E(err)
	}
	return snippet, errors.E(err)
}

// Update one snippet.
func (b *SnippetDAO) Update(snippet *Snippet) error {
	_, err := b.db.NamedExec(`
		UPDATE snippets
		SET
			user_id = :user_id,
			when = :when
		WHERE
			id = :id
	`, snippet)
	return errors.E(err)
}
