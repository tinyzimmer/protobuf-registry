import * as React from "react";
import { AnchorButton } from "@blueprintjs/core";

const formats = [
  "raw",
  "go",
  "python",
  "java",
  "cpp",
  "csharp",
  "objc",
  "php",
  "ruby",
]

class ProtoDownloadButton extends React.Component {
  constructor(props) {
    super(props);
    this.buttonText = props.buttonText;
    this.downloadRoot = '/api/proto/' + props.name + '/' + props.version + '/'
    this.state = {
      downloadURL: this.downloadRoot + 'raw'
    }
    this.handleFormatChange = this.handleFormatChange.bind(this);
  }

  handleFormatChange(e) {
    this.setState({ downloadURL: this.downloadRoot + formats[e.target.value] });
  }

  render() {
    return (
      <div className="wrapper">
        <AnchorButton intent="primary" icon="download" href={this.state.downloadURL}>{this.buttonText}</AnchorButton>
        &nbsp;&nbsp;
        <div className="bp3-select">
          <select
            style={{width: '100px'}}
            onChange={this.handleFormatChange}>
            {formats.map((value, index) => {
              return <option value={index} key={index}>{value}</option>
            })}
          </select>
        </div>
      </div>
    );
  }
}

export default ProtoDownloadButton;
