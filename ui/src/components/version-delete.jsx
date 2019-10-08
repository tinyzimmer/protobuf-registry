import * as React from "react";

import { Alert, Button, Intent, Toaster } from "@blueprintjs/core";

const toaster = Toaster.create()

class DeleteButton extends React.Component {
  constructor(props) {
    super(props);
    this.apiURL = props.apiURL;
    this.name = props.name;
    this.version = props.version;
    this.callback = props.callback;
    this.deleteURL = this.apiURL + '/api/proto/' + this.name
    if (this.version !== '*') {
      this.deleteURL = this.deleteURL + '/' + this.version
    }
    this.toaster = toaster
    this.state = {
      isOpen: false
    }
    this.handleDeleteConfirm = this.handleDeleteConfirm.bind(this)
    this.toast = this.toast.bind(this)
  }

  handleDeleteOpen = () => this.setState({ isOpen: true });
  handleDeleteCancel = () => this.setState({ isOpen: false });

  toast(msg) {
    this.toaster.show(
      {
        message: msg,
        intent: 'danger',
        icon: 'delete'
      }
    );
  }

  handleDeleteConfirm = () => {
    fetch(this.deleteURL, {method: 'DELETE'})
    .then(results => {
      return results.json()
    }).then(result => {
      this.setState({ isOpen: false });
      this.callback(this.name, this.version)
      this.toast(result.result)
    })
  };

  getDeleteMessage() {
    if (this.version === '*') {
      return (
        <p>Are you sure you want to delete <b>{this.name}?</b><br /><br /> <b>All versions</b> will be removed permanently.</p>
      );
    } else {
      return (
        <p>Are you sure you want to delete <b>{this.name}</b> version <b>{this.version}</b>? It will be removed permanently.</p>
      );
    }
  }

  render() {
    return (
      <div>
        <Button icon="trash" intent={Intent.DANGER} onClick={this.handleDeleteOpen} text="" />
        <Alert
            cancelButtonText="Cancel"
            confirmButtonText="Delete"
            icon="trash"
            intent={Intent.DANGER}
            isOpen={this.state.isOpen}
            onCancel={this.handleDeleteCancel}
            onConfirm={this.handleDeleteConfirm}
        >
        {this.getDeleteMessage()}
        </Alert>
      </div>
    );
  }
}

export default DeleteButton;
