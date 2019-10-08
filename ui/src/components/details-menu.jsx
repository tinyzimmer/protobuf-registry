import * as React from "react";

import { Button, Dialog, Classes, HTMLTable, Tag, Pre, Collapse } from "@blueprintjs/core";

class FieldsCollapse extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isOpen: false
    }
    var fieldsStr = '';
    Object.keys(props.fields).map(function(keyName, keyIndex) {
      fieldsStr = fieldsStr + keyName + ' | ' + props.fields[keyName] + '\n'
      return ''
    })
    this.fieldsStr = fieldsStr;
  }

  handleClick = () => this.setState({ isOpen: !this.state.isOpen });

  render() {

    return (
      <div>
        <Button intent="primary" onClick={this.handleClick}>{this.state.isOpen ? "Hide" : "Show"} fields</Button>
        <Collapse isOpen={this.state.isOpen} keepChildrenMounted={true}>
            <Pre>{this.fieldsStr}</Pre>
        </Collapse>
      </div>
    );
  }
}

class DetailsMenu extends React.Component {
  constructor(props) {
    super(props);
    this.meta = props.meta;
    this.apiURL = props.apiURL;
    this.state = {
      isOpen: false,
      messages: [],
    };
    this.dataURL = this.apiURL + '/api/proto/' + this.meta.name + '/' + this.meta.latest;
    this.handleOpen = this.handleOpen.bind(this);
    this.handleClose = this.handleClose.bind(this);
  }

  handleOpen() {
    //this.fetchData();
    this.setState({ isOpen: true });
  }

  handleClose() {
    this.setState({ isOpen: false });
  }

  componentDidMount() {
    fetch(this.dataURL)
    .then(results => {
      return results.json()
    }).then(data => {
      var messages = [];
      var sourceFiles = [];
      data.messages.map((value, index) => {
        messages.push({
          name: value.name,
          fields: value.fields,
        })
        return ''
      })
      data.sourceFiles.map((value, index) => {
        sourceFiles.push(<Tag key={index}>{value}</Tag>)
        return ''
      })
      this.setState({messages: messages})
      this.setState({sourceFiles: sourceFiles})
    })
  }

  render() {
    var title = "Details | " + this.meta.name
    return (
      <Button icon="more" text="Details" onClick={this.handleOpen}>
        <Dialog
          icon="info-sign"
          title={title}
          isOpen={this.state.isOpen}
          onClose={this.handleClose}
          autoFocus={true}
          canEscapeKeyClose={true}
          canOutsideClickClose={true}
          isCloseButtonShown={false}
          style={{width:'700px'}}
        >
          <div className={Classes.DIALOG_BODY}>
            <HTMLTable interactive striped bordered>
              <thead>
                <tr>
                  <th>Message</th>
                  <th>Fields</th>
                </tr>
              </thead>
              <tbody>
                {this.state.messages.map((value, index) => {
                  return (
                    <tr key={value.name}>
                      <td><strong>{value.name}</strong></td>
                      <td><FieldsCollapse fields={value.fields} /></td>
                    </tr>
                  );
                })}
              </tbody>
            </HTMLTable>
            <br></br>
            <p><strong>Source Files: </strong>{this.state.sourceFiles}</p>
          </div>
        </Dialog>
      </Button>
    );
  }
}

export default DetailsMenu;
