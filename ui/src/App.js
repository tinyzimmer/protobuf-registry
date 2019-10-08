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
          <div style={{paddingLeft: '3em', paddingRight: '3em'}}>
            <Switch>
              <Route exact path="/">
                <ProtobufTable />
              </Route>
              <Route exact path="/config">
                <ServerConfig />
              </Route>
              <Route exact path="/apidoc">
                <APIDoc />
              </Route>
              <Route exact path="/browser">
                <ProtoBrowser />
              </Route>
            </Switch>
          </div>
        </div>
      </Router>
    </div>
  );
}

export default App;
