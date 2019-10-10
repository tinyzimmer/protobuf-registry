import * as React from "react";

import {
  Button,
  Classes,
  Dialog,
  FileInput,
  InputGroup,
  Text,
  Toaster,
} from "@blueprintjs/core";

const toaster = Toaster.create()

function getReadableFileSizeString(fileSizeInBytes) {
    var i = -1;
    var byteUnits = [' kB', ' MB', ' GB', ' TB', 'PB', 'EB', 'ZB', 'YB'];
    do {
        fileSizeInBytes = fileSizeInBytes / 1024;
        i++;
    } while (fileSizeInBytes > 1024);

    return Math.max(fileSizeInBytes, 0.1).toFixed(1) + byteUnits[i];
};

function getBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => {
      let encoded = reader.result.toString().replace(/^data:(.*,)?/, '');
      if ((encoded.length % 4) > 0) {
        encoded += '='.repeat(4 - (encoded.length % 4));
      }
      resolve(encoded);
    };
    reader.onerror = error => reject(error);
  });
}

class UploadForm extends React.Component {
  constructor(props) {
    super(props);
    this.remountFunc = props.remountFunc
    this.postURL = '/api/proto'
    this.state = {
      isOpen: false,
      fileInputText: "Choose file...",
      fileInputCaption: "",
      fileInputWidth: "",
      fileCaptionHidden: true,
      nameIntent: "primary",
      fileIntent: "primary",
      packageName: "",
      packageVersion: "0.0.1",
      packageBody: "",
      wiggling: false,
    }
    this.wiggleTimeoutId = 3
    this.toaster = toaster
    this.clear = this.clear.bind(this)
    this.wiggle = this.wiggle.bind(this)
    this.handleOpen = this.handleOpen.bind(this)
    this.handleClose = this.handleClose.bind(this)
    this.handleFileInput = this.handleFileInput.bind(this)
    this.handleNameChange = this.handleNameChange.bind(this)
    this.handleVersionChange = this.handleVersionChange.bind(this)
    this.handleSubmit = this.handleSubmit.bind(this)
    this.handleSubmitResult = this.handleSubmitResult.bind(this)
  }

  clear() {
    this.setState({
      isOpen: false,
      fileInputText: "Choose file...",
      fileInputCaption: "",
      fileInputWidth: "",
      fileCaptionHidden: true,
      nameIntent: "primary",
      fileIntent: "primary",
      packageName: "",
      packageVersion: "0.0.1",
      packageBody: "",
      wiggling: false,
    })
  }

  handleOpen() {
    this.setState({ isOpen: true });
  }

  handleClose() {
    this.setState({ isOpen: false });
  }

  wiggle() {
    window.clearTimeout(this.wiggleTimeoutId);
    this.setState({ wiggling: true });
    this.wiggleTimeoutId = window.setTimeout(() => this.setState({ wiggling: false }), 300);
  }

  handleSubmitResult(result) {
    if (result.error !== undefined) {
      this.toaster.show(
        {
          message: result.error,
          intent: "danger",
          icon: "warning-sign",
        }
      )
    } else {
      this.toaster.show(
        {
          message: this.state.packageName + ' version ' + this.state.packageVersion + ' uploaded',
          intent: "success",
          icon: "tick"
        }
      );
      this.handleClose()
      this.clear()
      this.remountFunc()
    }
  }

  handleNameChange(e) {
    if (e.target.value === "") {
      this.setState({nameIntent: "danger"})
    } else {
      this.setState({nameIntent: "primary"})
    }
    this.setState({packageName: e.target.value})
  }

  handleVersionChange(e) {
    if (e.target.value === "") {
      this.setState({packageVersion: "0.0.1"})
    } else {
      this.setState({packageVersion: e.target.value})
    }
  }

  handleFileInput(e) {
    var file = e.target.files[0]
    getBase64(file).then( data => {
      this.setState({
        fileInputText: file.name + ' - Size: ' + getReadableFileSizeString(file.size)
      });
      this.setState({
        packageBody: data,
      })
      this.setState({
        fileInputWidth: this.state.fileInputText.length*11 + 'px'
      })
      return ''
    })
    return
  }

  handleSubmit(e) {
    var postBody = {
      name: this.state.packageName,
      version: this.state.packageVersion,
      body: this.state.packageBody,
    }
    if (this.state.packageName === "") {
      this.wiggle()
      this.setState({nameIntent: "danger"})
      return
    } else if (this.state.packageBody === "") {
      this.wiggle()
      this.setState({fileIntent: "danger"})
    }
    fetch(this.postURL, {method: "POST", body: JSON.stringify(postBody)}).then(results => {
      return results.json()
    }).then(result => {
      this.handleSubmitResult(result)
    })
  }

  render() {
    return (
      <Button icon="upload" text="Upload" className="bp3-minimal" onClick={this.handleOpen}>
        <Dialog
          icon="upload"
          title="Upload New Package"
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
            <div className="wrapper">

              <div className="upload-form-container">
                <div className="upload-form">
                  <div className="upload-form-item">
                    <InputGroup
                      round
                      onChange={this.handleNameChange}
                      intent={this.state.nameIntent}
                      placeholder="Package name (required)"
                      leftIcon="bookmark"
                    />
                  </div>
                  <div className="upload-form-item">
                    <InputGroup
                      round
                      onChange={this.handleVersionChange}
                      intent="primary"
                      placeholder="Version: 0.0.1"
                      leftIcon="git-branch"
                    />
                  </div>
                  <br ></br>
                  <div className="upload-form-item">
                    <FileInput
                      style={{width: this.state.fileInputWidth}}
                      intent={this.state.fileIntent}
                      text={this.state.fileInputText}
                      onInputChange={this.handleFileInput}
                    />
                  </div>
                  <div className="upload-form-item" hidden={this.state.fileCaptionHidden}>
                    <Text>{this.state.fileInputCaption}</Text>
                  </div>
                  <br></br>
                  <div className="upload-form-item">
                    <Button
                      intent="primary"
                      onClick={this.handleSubmit}
                      style={{width: "75px"}}
                       className={this.state.wiggling ? "wiggle" : ""}
                    >
                      Upload
                    </Button>
                  </div>
                </div>
              </div>
              <br></br>
              <div className="upload-form-container">
                <div className="upload-form">
                  <div className="upload-form-item">
                    <Button>text</Button>
                  </div>
                </div>
              </div>

            </div>
          </div>
        </Dialog>
      </Button>
    );
  }
}

export default UploadForm;
