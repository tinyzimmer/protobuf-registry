import React, { Component } from "react";
import {  HTMLTable, Divider, Pre, Tag } from '@blueprintjs/core';

class MessageTables extends Component {
  constructor(props) {
    super(props);
    this.messages = props.messages
  }

  render() {
    return (
      <div>
        {this.messages.map((msg, msgindex) => {
          return (
            <div key={msgindex}>
              <div align="center">
                <strong>{msg.name}</strong>
              </div>
              <br></br>
              <div>
                <Tag>{msg.fullName}</Tag>&nbsp;&nbsp;
                <i>{msg.description}</i>
              </div>
              <br></br>
              <div align="center">
                <i>Fields</i>
                <br></br>
                <HTMLTable bordered striped condensed>
                  <thead>
                    <tr>
                      <th><strong>Name</strong></th>
                      <th><strong>Description</strong></th>
                      <th><strong>Label</strong></th>
                      <th><strong>Type</strong></th>
                    </tr>
                  </thead>
                  <tbody>
                    {msg.fields.map((field, fieldindex) => {
                      return (
                        <tr key={fieldindex}>
                          <td><Pre small="true">{field.name}</Pre></td>
                          <td>
                            {field.description}
                            <br></br>
                            {field.defaultValue !== "" ? "Default: " + field.defaultValue : ""}
                        </td>
                          <td>{field.label}</td>
                          <td><Pre small="true">{field.fullType}</Pre></td>
                        </tr>
                      )
                    })}
                  </tbody>
                </HTMLTable>
              </div>
              <br></br>
              <Divider></Divider>
              <br></br>
            </div>
          )
        })}
      </div>
    );
  }
}

export default MessageTables;
