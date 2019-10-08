import * as React from "react";

import VersionMenu from './version-menu.jsx';
import DetailsMenu from './details-menu.jsx';


class ProtobufMenu extends React.Component {
  constructor(props) {
    super(props);
    this.meta = props.meta;
    this.apiURL = props.apiURL;
  }

  render() {
    return (
      <div className="wrapper">
        <VersionMenu meta={this.meta} apiURL={this.apiURL}></VersionMenu>
        &nbsp;&nbsp;
        <DetailsMenu meta={this.meta} apiURL={this.apiURL}></DetailsMenu>
      </div>
    );
  }
}

export default ProtobufMenu;
