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

import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { loadSnippets } from './actions'
import { Col, Grid, PageHeader, Panel, Row } from 'react-bootstrap'

import './style.css'

export class Home extends React.PureComponent {
  componentDidMount () {
    this.props.loadSnippets()
  }
  render () {
    return (
      <Grid>
        <Row>
          <Col>
            <PageHeader> Snippets for the week of Jun 3~10: </PageHeader>
            {
              this.props.snippets.map((u) => (
                <Panel key={u.email} className='snippet-card'>
                  <strong>{u.name}</strong>
                  <ul>
                    {u.snippets.map((s) => (
                      <li key={u.email + s}>{s}</li>
                    ))}
                  </ul>
                </Panel>
              ))
            }
          </Col>
        </Row>
      </Grid>
    )
  }
}

const s2p = state => ({ snippets: state.snippets.snippets })
const d2p = dispatch => bindActionCreators({ loadSnippets }, dispatch)
export default connect(s2p, d2p)(Home)
