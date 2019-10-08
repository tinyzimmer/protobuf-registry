import React, { Component } from "react";
import SyntaxHighlighter from 'react-syntax-highlighter';
import { atomOneDark } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import { Card, Divider } from '@blueprintjs/core';

const Header = () => {
  return (
    <h4>Server Configuration</h4>
  )
}

class ServerConfig extends Component {
  constructor(props) {
    super(props);
    this.state = {
      configStr: ""
    };
  }

  componentDidMount() {
    fetch('/api/config')
    .then(result => {
      return result.text();
    }).then(data => {
      this.setState({configStr: data})
      console.log(this.state)
    })
  }

  render() {
    return (
      <div align="left" style={{paddingLeft: '10em', paddingRight: '10em'}}>
        <Card className="bp3-dark" style={{width: '100%'}}>
          <Header />
        </Card>
        <Divider></Divider>
        <SyntaxHighlighter language="json" style={atomOneDark}>
          {this.state.configStr}
        </SyntaxHighlighter>
      </div>
    );
  }
}

export default ServerConfig;
