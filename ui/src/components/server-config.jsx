import React, { Component } from "react";
import SyntaxHighlighter from 'react-syntax-highlighter';
import { github } from 'react-syntax-highlighter/dist/esm/styles/hljs';

const Header = () => {
  return (
    <h4 className="font-weight-bold">Server Configuration</h4>
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
        <Header />
        <SyntaxHighlighter language="json" style={github}>
          {this.state.configStr}
        </SyntaxHighlighter>
      </div>
    );
  }
}

export default ServerConfig;
