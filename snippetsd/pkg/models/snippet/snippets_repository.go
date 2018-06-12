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
	"cirello.io/snippetsd/pkg/errors"
	"github.com/jmoiron/sqlx"
)

// Repository provides a repository of Snippets.
type Repository struct {
	db *sqlx.DB

	contentsRepository *contentsRepository
}

// NewRepository instanties a Repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db:                 db,
		contentsRepository: newContentsRepository(db),
	}
}

// Bootstrap creates table if missing.
func (b *Repository) Bootstrap() error {
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

	if err := b.contentsRepository.Bootstrap(); err != nil {
		return errors.E(err, "cannot create table for contents")
	}

	return nil
}

// RepositoryOption allows to modify the repository calls as needed.
type RepositoryOption func(*Repository, []*Snippet) error

// WithContent will plugin the snippet content.
func WithContent() RepositoryOption {
	return func(b *Repository, snippets []*Snippet) error {
		for _, s := range snippets {
			contents, err := b.contentsRepository.GetBySnippetID(s.ID)
			if err != nil {
				return errors.E(err, "cannot load snippets content")
			}
			s.Contents = contents
		}
		return nil
	}
}

func applyRepositoryOptions(b *Repository, snippets []*Snippet, opts []RepositoryOption) error {
	for _, opt := range opts {
		if err := opt(b, snippets); err != nil {
			return errors.E(err, "failed to  apply repository option")
		}
	}
	return nil
}

// All returns all known snippets.
func (b *Repository) All(opts ...RepositoryOption) ([]*Snippet, error) {
	var snippets []*Snippet
	err := b.db.Select(&snippets, "SELECT * FROM snippets")
	applyRepositoryOptions(b, snippets, opts)
	return snippets, errors.E(err)
}

// Current returns the current week snippets.
func (b *Repository) Current(opts ...RepositoryOption) ([]*Snippet, error) {
	var snippets []*Snippet
	// TODO: convert 7 to variable representing the number of days
	err := b.db.Select(&snippets, "SELECT * FROM snippets WHERE when > 7")
	applyRepositoryOptions(b, snippets, opts)
	return snippets, errors.E(err)
}

// Insert one snippet entry.
func (b *Repository) Insert(snippet *Snippet) (*Snippet, error) {
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

	for _, c := range snippet.Contents {
		c, err = b.contentsRepository.Insert(c)
		if err != nil {
			return snippet, errors.E(err, "cannot insert snippet content")
		}
	}
	return snippet, nil
}

// Update one snippet.
func (b *Repository) Update(snippet *Snippet) error {
	_, err := b.db.NamedExec(`
		UPDATE snippets
		SET
			user_id = :user_id,
			when = :when
		WHERE
			id = :id
	`, snippet)
	if err != nil {
		return errors.E(err)
	}
	for _, c := range snippet.Contents {
		if err := b.contentsRepository.Update(c); err != nil {
			return errors.E(err, "cannot update snippet content")
		}
	}
	return nil
}