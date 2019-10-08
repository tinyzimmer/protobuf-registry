import React from 'react';
import './App.css';
import { HashRouter as Router, Switch, Route } from 'react-router-dom';
import RegistryNavBar from './components/registry-navbar.jsx';
import ProtobufTable from './components/protobuf-table.jsx'
import ServerConfig from './components/server-config.jsx';
import APIDoc from './components/api-doc.jsx';
import ProtoBrowser from './components/proto-browser.jsx';

// uncomment the below to point API server at a seperate instance when using `npm run start`
//const apiURL = 'http://localhost:8080';
const apiURL = '';

function App() {
  return (
    <div className="App">
      <Router>
        <div>
          <RegistryNavBar apiURL={apiURL} />
          <br></br>
          <div style={{paddingLeft: '3em', paddingRight: '3em'}}>
            <Switch>
              <Route exact path="/">
                <ProtobufTable apiURL={apiURL} />
              </Route>
              <Route exact path="/config">
                <ServerConfig apiURL={apiURL} />
              </Route>
              <Route exact path="/apidoc">
                <APIDoc apiURL={apiURL} />
              </Route>
              <Route exact path="/browser">
                <ProtoBrowser apiURL={apiURL} />
              </Route>
            </Switch>
          </div>
        </div>
      </Router>
    </div>
  );
}

export default App;
