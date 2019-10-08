import * as React from "react";

import VersionMenu from './version-menu.jsx';
import ProtoDownloadButton from './proto-download-button.jsx';

class ProtobufMenu extends React.Component {
  constructor(props) {
    super(props);
    this.meta = props.meta;
  }

  render() {
    return (
      <div className="wrapper">
        <VersionMenu meta={this.meta}></VersionMenu>
        &nbsp;&nbsp;
        <ProtoDownloadButton buttonText="Download Latest" name={this.meta.name} version={this.meta.latest} />,
      </div>
    );
  }
}

export default ProtobufMenu;
