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

// import { push } from 'react-router-redux'

export function loadSnippetsByUser () {
  return (dispatch) => {
    fetch('http://localhost:5100/snippetsByUser', {
      credentials: 'include'
    })
    .then(res => res.json())
    .catch((e) => {
      console.log('cannot load state:', e)
    })
    .then((snippets) => {
      dispatch({
        type: 'snippets/USER_SNIPPETS_LOADED',
        snippets: snippets
      })
    })
  }
}

export function saveSnippet (contents) {
  return (dispatch) => {
    fetch('http://localhost:5100/storeSnippet', {
      credentials: 'include',
      method: 'POST',
      body: JSON.stringify({contents})
    })
    .then(res => res.json())
    .catch((e) => {
      console.log('cannot store snippet:', e)
    })
    .then(() => {
      loadSnippetsByUser()(dispatch)
    })
  }
}