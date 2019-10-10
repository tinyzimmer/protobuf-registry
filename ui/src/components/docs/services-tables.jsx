import React, { Component } from "react";
import {  HTMLTable, Divider, Pre, Tag } from '@blueprintjs/core';

class ServicesTables extends Component {
  constructor(props) {
    super(props);
    this.services = props.services
  }

  render() {
    return (
      <div>
        {this.services.map((svc, svcindex) => {
          return (
            <div key={svcindex}>
              <div align="center">
                <strong>{svc.name}</strong>
              </div>
              <br></br>
              <div>
                <Tag>{svc.fullName}</Tag>&nbsp;&nbsp;
                <i>{svc.description}</i>
              </div>
              <br></br>
              <div align="center">
                <i>Methods</i>
                <br></br>
                <HTMLTable bordered striped condensed>
                  <thead>
                    <tr>
                      <th><strong>Name</strong></th>
                      <th><strong>Description</strong></th>
                      <th><strong>Request Type</strong></th>
                      <th><strong>Response Type</strong></th>
                    </tr>
                  </thead>
                  <tbody>
                    {svc.methods.map((method, methodindex) => {
                      return (
                        <tr key={methodindex}>
                          <td><Pre small="true">{method.name}</Pre></td>
                          <td>{method.description}</td>
                          <td><Pre>{method.requestFullType}</Pre></td>
                          <td><Pre small="true">{method.responseFullType}</Pre></td>
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

export default ServicesTables;
