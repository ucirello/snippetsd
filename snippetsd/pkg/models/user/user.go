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

package user // import "cirello.io/snippetsd/pkg/models/user"
import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// User aggregates all the information of a snippet user.
type User struct {
	ID    int64  `db:"id" json:"id"`
	Email string `db:"email" json:"email"`
}

func (u *User) String() string {
	return fmt.Sprintf("%v", u.Email)
}

// Add inserts a user into the repository.
func Add(db *sqlx.DB, u *User) (*User, error) {
	return NewRepository(db).Insert(u)
}

// NewFromEmail creates a user from a given email.
func NewFromEmail(email string) (*User, error) {
	// TODO: validate email
	return &User{
		Email: email,
	}, nil
}
