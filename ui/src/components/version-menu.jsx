import * as React from "react";

import { Button, Classes, Dialog, HTMLTable, Tag } from "@blueprintjs/core";
import ProtoDownloadButton from './proto-download-button.jsx';
import DeleteButton from './version-delete.jsx';

class VersionMenu extends React.Component {
  constructor(props) {
    super(props);
    this.meta = props.meta;
    this.state = {
      isOpen: false,
      versions: this.meta.versions
    }
    this.handleOpen = this.handleOpen.bind(this)
    this.handleClose = this.handleClose.bind(this)
    this.handleVersionDeleted = this.handleVersionDeleted.bind(this)
  }

  handleOpen() {
    this.setState({ isOpen: true });
  }

  handleClose() {
    this.setState({ isOpen: false });
  }

  handleVersionDeleted(name, version) {
    var newVersions = [];
    this.state.versions.map((value, index) => {
      if (value.version !== version) {
        newVersions.push(value);
      }
      return ''
    })
    this.setState({versions: newVersions})
  }

  render() {
    var title = "Versions | " + this.meta.name
    return (
      <Button intent="primary" icon="projects" text="Versions" onClick={this.handleOpen}>
        <Dialog
          icon="info-sign"
          title={title}
          isOpen={this.state.isOpen}
          onClose={this.handleClose}
          autoFocus={true}
          canEscapeKeyClose={true}
          canOutsideClickClose={true}
          isCloseButtonShown={false}
          style={{width:'800px'}}
          className="bp3-dark"
        >
          <div className={Classes.DIALOG_BODY}>
            <HTMLTable>
              <thead>
                <tr>
                  <th></th>
                  <th>Version</th>
                  <th>Uploaded</th>
                </tr>
              </thead>
              <tbody>
                {this.state.versions.map((value, index) => {
                  return (
                    <tr key={value.version}>
                      <td><DeleteButton callback={this.handleVersionDeleted} name={value.name} version={value.version}></DeleteButton></td>
                      <td><Tag icon="git-branch" large>{value.version}</Tag></td>
                      <td>{new Date(value.lastUpdated).toString().replace(/\(.*\)/, '')}</td>
                      <td><ProtoDownloadButton buttonText="Download" name={value.name} version={value.version}/></td>
                    </tr>
                  );
                })}
              </tbody>
            </HTMLTable>
          </div>
        </Dialog>
      </Button>
    );
  }
}

export default VersionMenu;
