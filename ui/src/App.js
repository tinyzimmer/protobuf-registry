import React from 'react';
import './App.css';
import { HashRouter as Router, Switch, Route } from 'react-router-dom';
import RegistryNavBar from './components/registry-navbar.jsx';
import ProtobufTable from './components/protobuf-table.jsx'
import ServerConfig from './components/server-config.jsx';
import APIDoc from './components/api-doc.jsx';
import ProtoBrowser from './components/proto-browser.jsx';


function App() {
  return (
    <div className="App">
      <Router>
        <div>
          <RegistryNavBar />
          <br></br>
          <div style={{paddingLeft: '2em', paddingRight: '2em'}}>
            <Switch>
              <Route exact path="/" component={ProtobufTable} />
              <Route exact path="/config" component={ServerConfig} />
              <Route exact path="/apidoc" component={APIDoc} />
              <Route exact path="/browser" component={ProtoBrowser} />
            </Switch>
          </div>
        </div>
      </Router>
    </div>
  );
}

export default App;
