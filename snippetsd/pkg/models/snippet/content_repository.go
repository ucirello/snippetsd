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

// contentsRepository provides a repository of snippets contents.
type contentsRepository struct {
	db *sqlx.DB
}

// newContentsRepository instanties a contentsRepository
func newContentsRepository(db *sqlx.DB) *contentsRepository {
	return &contentsRepository{db: db}
}

// Bootstrap creates table if missing.
func (b *contentsRepository) Bootstrap() error {
	cmds := []string{
		`create table if not exists contents (
			id integer primary key autoincrement,
			snippet_id bigint,
			content text
		);
		`,
		`create index if not exists contents_snippet_id  on contents (snippet_id)`,
	}

	for _, cmd := range cmds {
		_, err := b.db.Exec(cmd)
		if err != nil {
			return errors.E(err)
		}
	}

	return nil
}

// GetBySnippetID loads a snippet's content.
func (b *contentsRepository) GetBySnippetID(id int64) ([]*Content, error) {
	var contents []*Content
	err := b.db.Select(&contents, "SELECT * FROM contents")
	return contents, errors.E(err)
}

// Insert one snippet entry content.
func (b *contentsRepository) Insert(content *Content) (*Content, error) {
	_, err := b.db.NamedExec(`
		INSERT INTO contents
		(snippet_id, content)
		VALUES (:snippet_id, :content)
	`, content)
	if err != nil {
		return nil, errors.E(err)
	}

	err = b.db.Get(content, `
		SELECT
			*
		FROM
			contents
		WHERE
			id = last_insert_rowid()
	`)
	if err != nil {
		return nil, errors.E(err)
	}
	return content, errors.E(err)
}

// Update one content.
func (b *contentsRepository) Update(content *Content) error {
	_, err := b.db.NamedExec(`
		UPDATE contents
		SET
			content = :content
		WHERE
			id = :id
	`, content)
	return errors.E(err)
}
