import React, { Component } from "react";
import { Button, Card, Elevation, Pre, Collapse, Divider, Spinner } from "@blueprintjs/core";

const Header = () => {
  return (
    <div>
      <Card  elevation="4" className="bp3-dark">
        <h4 className="font-weight-bold">API Index</h4>
      </Card>
      <Divider></Divider>
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
      <Card className="bp3-dark" elevation={Elevation.THREE}>
        <div className="wrapper">
          <Button intent={intent} icon={icon} text={this.state.data.method} onClick={this.handleClick}></Button>
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
    this.state = {
      routes: [],
      loading: true,
    };
  }

  componentDidMount() {
    fetch('/api')
    .then(result => {
      return result.json();
    }).then(data => {
      this.setState({routes: data.routes})
      this.setState({loading: false})
      console.log(this.state)
    })
  }

  render() {
    return (
      <div align="center">
        <Header />
        <div align="center" hidden={!this.state.loading}>
          <Spinner size={Spinner.SIZE_LARGE}></Spinner>
        </div>
        <div align="left" style={{paddingLeft: '10em', paddingRight: '10em'}}>
          {this.state.routes.map((route, index) => {
            return (
              <div>
                <DocCollapse key={index} data={route} />
                <br></br>
              </div>
            );
          })}
        </div>
      </div>
    );
  }
}

export default APIDoc;
