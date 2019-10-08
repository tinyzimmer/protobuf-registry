import React from 'react';
import './App.css';
import { HashRouter as Router, Switch, Route } from 'react-router-dom';
import { Divider } from '@blueprintjs/core';
import RegistryNavBar from './components/registry-navbar.jsx';
import ProtobufTable from './components/protobuf-table.jsx'
import ServerConfig from './components/server-config.jsx';
import APIDoc from './components/api-doc.jsx';
import ProtoBrowser from './components/proto-browser.jsx';

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      table: 0,
      config: 0,
      doc: 0,
      browser: 0
    }
  }

  remountComponents() {
    this.setState({table: this.state.table + 1})
    this.setState({config: this.state.config + 1})
    this.setState({doc: this.state.doc + 1})
    this.setState({browser: this.state.browser + 1})
  }

  render() {
    return (
      <div className="App">
        <Router>
          <div>
            <RegistryNavBar remountFunc={this.remountComponents} />
            <Divider></Divider>
            <div style={{paddingLeft: '2em', paddingRight: '2em'}}>
              <Switch>
                <Route key={this.state.table} exact path="/" component={ProtobufTable} />
                <Route key={this.state.config} exact path="/config" component={ServerConfig} />
                <Route key={this.state.doc} exact path="/apidoc" component={APIDoc} />
                <Route key={this.state.browser} exact path="/browser" component={ProtoBrowser} />
              </Switch>
            </div>
          </div>
        </Router>
      </div>
    );
  }
}

export default App;
