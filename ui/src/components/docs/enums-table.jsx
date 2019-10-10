import React, { Component } from "react";
import {  HTMLTable, Pre } from '@blueprintjs/core';

class EnumsTable extends Component {
  constructor(props) {
    super(props);
    this.enums = props.enums
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
                <th><strong>Values</strong></th>
              </tr>
            </thead>
            <tbody>
              {this.enums.map((enm, enumindex) => {
                return (
                  <tr key={enumindex}>
                    <td><Pre small="true">{enm.fullName}</Pre></td>
                    <td>{enm.description}</td>
                    <td>
                      {enm.values.map((value, index) => {
                        return (
                            <Pre>{value.name}  |  {value.description}</Pre>
                        );
                      })}
                    </td>
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

export default EnumsTable;
