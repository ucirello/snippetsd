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

package snippet

import (
	"time"

	"cirello.io/snippetsd/pkg/errors"
	"cirello.io/snippetsd/pkg/models/user"
	"github.com/jmoiron/sqlx"
)

// Repository provides a repository of Snippets.
type Repository struct {
	db *sqlx.DB
}

// NewRepository instanties a Repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Bootstrap creates table if missing.
func (b *Repository) Bootstrap() error {
	cmds := []string{
		`create table if not exists snippets (
			id integer primary key autoincrement,
			user_id bigint,
			week_start datetime,
			contents bigtext
		);
		`,
		`create index if not exists snippets_user_id on snippets (user_id)`,
		`create index if not exists snippets_week_start on snippets (week_start)`,
	}

	for _, cmd := range cmds {
		_, err := b.db.Exec(cmd)
		if err != nil {
			return errors.E(err, cmd)
		}
	}

	return nil
}

// RepositoryOption allows to modify the repository calls as needed.
type RepositoryOption func(*Repository, *[]*Snippet) error

// WithUser will plug the snippet user.
func WithUser() RepositoryOption {
	return func(b *Repository, snippets *[]*Snippet) error {
		repo := user.NewRepository(b.db)
		for i, s := range *snippets {
			u, err := repo.GetByID(s.ID)
			if err != nil {
				return errors.E(err, "cannot load snippets user")
			}
			s.User = u
			(*snippets)[i] = s
		}
		return nil
	}
}

func applyRepositoryOptions(b *Repository, snippets *[]*Snippet, opts []RepositoryOption) error {
	for _, opt := range opts {
		if err := opt(b, snippets); err != nil {
			return errors.E(err, "failed to apply repository option")
		}
	}
	return nil
}

// All returns all known snippets.
func (b *Repository) All(opts ...RepositoryOption) ([]*Snippet, error) {
	var snippets []*Snippet
	err := b.db.Select(&snippets, "SELECT * FROM snippets")
	if err := applyRepositoryOptions(b, &snippets, opts); err != nil {
		return snippets, errors.E(err)
	}
	return snippets, errors.E(err)
}

// Current returns the current week snippets.
func (b *Repository) Current(opts ...RepositoryOption) ([]*Snippet, error) {
	var snippets []*Snippet

	weekStart := 7 * 24 * time.Hour // TODO: calculate correct week start
	err := b.db.Select(&snippets,
		"SELECT * FROM snippets WHERE week_start >= $1", weekStart)
	if err := applyRepositoryOptions(b, &snippets, opts); err != nil {
		return snippets, errors.E(err)
	}
	return snippets, errors.E(err)
}

// Insert one snippet entry.
func (b *Repository) Insert(snippet *Snippet) (*Snippet, error) {
	_, err := b.db.NamedExec(`
		INSERT INTO snippets
		(user_id, week_start, contents)
		VALUES (:user_id, :week_start, :contents)
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

	return snippet, nil
}

// Update one snippet.
func (b *Repository) Update(snippet *Snippet) error {
	_, err := b.db.NamedExec(`
		UPDATE snippets
		SET
			user_id = :user_id,
			week_start = :week_start,
			contents = :contents
		WHERE
			id = :id
	`, snippet)
	if err != nil {
		return errors.E(err)
	}

	return nil
}
