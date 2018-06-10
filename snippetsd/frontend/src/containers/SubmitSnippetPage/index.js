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
import {
  Button,
  Col,
  Form,
  Grid,
  PageHeader,
  Row,
  FormControl,
  FormGroup,
  InputGroup
 } from 'react-bootstrap'

import './style.css'

const SubmitSnippetPage = () => (
  <Grid>
    <Row>
      <Col>
        <PageHeader> What did you do past week? </PageHeader>

        <Form>
          <FormGroup bsSize='large'>
            <InputGroup>
              <InputGroup.Addon className='snippet-field-left-addon'>•</InputGroup.Addon>
              <FormControl type='text' />
            </InputGroup>
          </FormGroup>
          <FormGroup bsSize='large'>
            <InputGroup>
              <InputGroup.Addon className='snippet-field-left-addon'>•</InputGroup.Addon>
              <FormControl type='text' />
            </InputGroup>
          </FormGroup>
          <FormGroup bsSize='large'>
            <InputGroup>
              <InputGroup.Addon className='snippet-field-left-addon'>•</InputGroup.Addon>
              <FormControl type='text' />
            </InputGroup>
          </FormGroup>
          <FormGroup bsSize='large'>
            <InputGroup>
              <InputGroup.Addon className='snippet-field-left-addon'>•</InputGroup.Addon>
              <FormControl type='text' />
            </InputGroup>
          </FormGroup>
          <div className='snippet-submit'><Button>submit</Button></div>
        </Form>
      </Col>
    </Row>
  </Grid>
)

export default SubmitSnippetPage
