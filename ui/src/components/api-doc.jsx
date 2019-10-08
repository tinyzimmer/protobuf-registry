import React, { Component } from "react";
import { Button, Card, Elevation, Pre, Collapse } from "@blueprintjs/core";

const Header = () => {
  return (
    <div>
      <h4 className="font-weight-bold">API Index</h4>
      <br></br>
    </div>
  )
}

class DocCollapse extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isOpen: false,
      data: props.data,
    }
  }

  handleClick = () => this.setState({ isOpen: !this.state.isOpen });

  render() {
    var intent = "primary"
    var icon =""
    if (this.state.data.method === "GET") {
      intent = "success"
      icon = "info-sign"
    } else if (this.state.data.method === "POST") {
      intent = "warning"
      icon = "cloud-upload"
    } else if (this.state.data.method === "DELETE") {
      intent = "danger"
      icon = "cross"
    }
    return (
      <Card className="bp3-dark" interactive={true} elevation={Elevation.THREE} onClick={this.handleClick}>
        <div className="wrapper">
          <Button small={true} intent={intent} icon={icon} text={this.state.data.method}></Button>
          &nbsp;&nbsp;<strong>{this.state.data.path}</strong>
        </div>
        <Collapse isOpen={this.state.isOpen} keepChildrenMounted={true}>
            <Pre>{this.state.data.description}</Pre>
        </Collapse>
      </Card>
    );
  }
}

class APIDoc extends Component {
  constructor(props) {
    super(props);
    this.apiURL = props.apiURL
    this.state = {
      routes: [],
    };
  }

  componentDidMount() {
    fetch(this.apiURL + '/api')
    .then(result => {
      return result.json();
    }).then(data => {
      this.setState({routes: data.routes})
      console.log(this.state)
    })
  }

  render() {
    return (
      <div align="left" style={{paddingLeft: '10em', paddingRight: '10em'}}>
        <Header />
        {this.state.routes.map((route, index) => {
          return (
            <div>
              <DocCollapse key={index} data={route} />
              <br></br>
            </div>
          );
        })}
      </div>
    );
  }
}

export default APIDoc;
