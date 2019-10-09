import React, { Component } from 'react';
import { Navbar, Alignment, Button, Icon } from '@blueprintjs/core';
import { Link } from 'react-router-dom';

import UploadForm from './upload-form.jsx';

class RegistryNavBar extends Component {
  constructor(props) {
    super(props);
    this.remountFunc = props.remountFunc;
    this.state = {
      active: "Home"
    }
  }

  toggleConfig = () => {
    this.setState({active: 'Config'});
  }

  toggleAPI = () => {
    this.setState({active: 'API'});
  }

  toggleBrowser = () => {
    this.setState({active: 'Browser'});
  }

  toggleHome = () => {
    this.setState({active: 'Home'});
  }

  render() {
    return (
      <Navbar className='bp3-dark'>
        <Navbar.Group align={Alignment.LEFT}>
          <Navbar.Heading><Icon icon="globe-network" intent="primary"></Icon>&nbsp;<strong>Protobuf Registry</strong></Navbar.Heading>
          <Navbar.Divider />
          <Link to="/" style={{ textDecoration: 'none', color: 'white' }}>
            <Button
            onClick={this.toggleHome}
            active={this.state.active === 'Home'}
            className="bp3-minimal"
            icon="home"
            text="Home"
            />
          </Link>
          <Link to="/config" style={{ textDecoration: 'none', color: 'white' }}>
            <Button
            onClick={this.toggleConfig}
            active={this.state.active === 'Config'}
            className="bp3-minimal"
            icon="cog"
            text="Configuration"
            />
          </Link>
          <Navbar.Divider />
          <UploadForm remountFunc={this.remountFunc} />
          <Navbar.Divider />
          <Link to="/apidoc" style={{ textDecoration: 'none', color: 'white' }}>
            <Button
            onClick={this.toggleAPI}
            active={this.state.active === 'API'}
            className="bp3-minimal"
            icon="code"
            text="API"
            />
          </Link>
          <Navbar.Divider />
          <Link to="/browser" style={{ textDecoration: 'none', color: 'white' }}>
            <Button
            onClick={this.toggleBrowser}
            active={this.state.active === 'Browser'}
            className="bp3-minimal"
            icon="folder-open"
            text="Browser"
            />
          </Link>
        </Navbar.Group>
        <Navbar.Group align={Alignment.RIGHT}>
          <a href="https://github.com/tinyzimmer/protobuf-registry" target="_blank" rel="noopener noreferrer">
          <Button
            icon="git-branch"
            text="GitHub"
          />
          </a>
        </Navbar.Group>
      </Navbar>
    )
  }
}

export default RegistryNavBar;
