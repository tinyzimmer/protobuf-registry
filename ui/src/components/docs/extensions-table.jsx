import React, { Component } from "react";
import {  HTMLTable, Pre } from '@blueprintjs/core';

class ExtensionsTable extends Component {
  constructor(props) {
    super(props);
    this.extensions = props.extensions
    console.log(props.extensions)
  }

  render() {
    return (
      <div>
        <div align="center">
          <HTMLTable bordered striped condensed>
            <thead>
              <tr>
                <th><strong>Name</strong></th>
                <th><strong>Description</strong></th>
                <th><strong>Label</strong></th>
                <th><strong>Type</strong></th>
                <th><strong>Containing Type</strong></th>
              </tr>
            </thead>
            <tbody>
              {this.extensions.map((ext, extindex) => {
                return (
                  <tr key={extindex}>
                    <td><Pre small="true">{ext.fullName}</Pre></td>
                    <td>{ext.description}</td>
                    <td>{ext.label}</td>
                    <td><Pre small="true">{ext.fullType}</Pre></td>
                    <td><Pre small="true">{ext.containingFullType}</Pre></td>
                  </tr>
                )
              })}
            </tbody>
          </HTMLTable>
        </div>
      </div>
    );
  }
}

export default ExtensionsTable;
