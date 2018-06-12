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

package snippet // import "cirello.io/snippetsd/pkg/models/snippet"

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Snippet aggregates all the information of a snippet.
type Snippet struct {
	ID        int64      `db:"id" json:"id"`
	UserID    string     `db:"user_id" json:"user_id"`
	CreatedAt *time.Time `db:"created_at" json:"created_at"`

	Contents []*Content `db:"-" json:"content"`
}

// AddContent add one or more contents to a snippet.
func (s *Snippet) AddContent(contents ...*Content) {
	s.Contents = append(s.Contents, contents...)
}

// DeleteContent remove one content of a snippet. Returns if the content was
// found and removed.
func (s *Snippet) DeleteContent(content *Content) bool {
	var contents []*Content
	var found bool
	for _, c := range s.Contents {
		if c.ID != content.ID {
			contents = append(contents, c)
			found = true
		}
	}
	return found
}

// HasContent checks if the snippet has any content
func (s *Snippet) HasContent() bool {
	return len(s.Contents) > 0
}

// LoadAll load all snippets for the current week
func LoadAll(db *sqlx.DB) ([]*Snippet, error) {
	repo := NewRepository(db)
	return repo.All(WithContent())
}
