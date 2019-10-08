import React, { Component } from "react";
import SyntaxHighlighter from 'react-syntax-highlighter';
import { solarizedDark } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import { Classes, Tree, Card } from "@blueprintjs/core";

const Header = () => {
  return (
    <h4 className="font-weight-bold">File Browser (super beta)</h4>
  )
}

class ProtoBrowser extends Component {
  constructor(props) {
    super(props);
    this.state = {
      nodes: [],
      fileViewHidden: true,
      fileText: "",
      fileTextHeader: "",
    }
    this.handleFileClick = this.handleFileClick.bind(this)
    this.handleNodeClick = this.handleNodeClick.bind(this)
    this.handleNodeCollapse = this.handleNodeCollapse.bind(this)
    this.handleNodeExpand = this.handleNodeExpand.bind(this)
    this.handleVersionExpand = this.handleVersionExpand.bind(this)
    this.forEachNode = this.forEachNode.bind(this)
  }

  handleFileClick(nodeData) {
    fetch('/api/proto/' + nodeData.parent + '/' + nodeData.version + '/raw/' + nodeData.label)
    .then(results => {
      return results.text()
    }).then(fileText => {
      this.setState({fileText: fileText})
      this.setState({fileTextHeader: nodeData.parent + '/' + nodeData.version + '/' + nodeData.label})
      this.setState({fileViewHidden: false})
    })
  }

  handleNodeClick(nodeData, _nodePath, e) {
    if (nodeData.isFile) {
      this.handleFileClick(nodeData)
    }
    const originallySelected = nodeData.isSelected;
    if (!e.shiftKey) {
        this.forEachNode(this.state.nodes, n => (n.isSelected = false));
    }
    nodeData.isSelected = originallySelected == null ? true : !originallySelected;
    this.setState(this.state);
  };

  handleNodeCollapse(nodeData) {
    nodeData.isExpanded = false;
    this.setState(this.state);
  }

  handleVersionExpand(nodeData) {
    fetch('/api/proto/' + nodeData.parent + '/' + nodeData.label)
    .then(results => {
      return results.json()
    }).then(data => {
      var children = []
      data.sourceFiles.map((value, index) => {
        children.push({
          id: index,
          hasCaret: false,
          icon: "document-open",
          label: value,
          isFile: true,
          parent: nodeData.parent,
          version: nodeData.label,
        })
        return ''
      })
      nodeData.childNodes = children
      this.setState(this.state)
    })
  }

  handleNodeExpand(nodeData) {
    if (nodeData.isVersion) {
      this.handleVersionExpand(nodeData)
    }
    nodeData.isExpanded = true;
    this.setState(this.state);
  }

  forEachNode(nodes, callback) {
      if (nodes == null) {
          return;
      }
      for (const node of nodes) {
          callback(node);
          this.forEachNode(node.childNodes, callback);
      }
  }

  componentDidMount() {
    fetch('/api/proto')
    .then(results => {
      return results.json()
    }).then(data => {
      var nodes = []
      data.map((value, index) => {
        var node = {
          id: index,
          hasCaret: true,
          icon: "folder-close",
          label: value.name,
          isVersion: false,
        }
        var children = []
        value.versions.map((version, i) => {
          children.push({
            id: i,
            hasCaret: true,
            icon: 'git-merge',
            label: version.version,
            isVersion: true,
            parent: value.name,
          })
          return ''
        })
        node.childNodes = children
        nodes.push(node)
        return ''
      });
      this.setState({nodes: nodes})
    })
  }

  render() {
    return (
      <div align="left" style={{paddingLeft: '10em', paddingRight: '10em'}}>
        <Header />
        <br></br>
        <div className="wrapper">
          <div style={{width: '25%'}}>
            <Tree
              contents={this.state.nodes}
              onNodeClick={this.handleNodeClick}
              onNodeCollapse={this.handleNodeCollapse}
              onNodeExpand={this.handleNodeExpand}
              className={Classes.ELEVATION_2}
            />
          </div>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
          <div hidden={this.state.fileViewHidden}>
            <Card elevation="3" className="bp3-dark">
              <strong>{this.state.fileTextHeader}</strong>
              <br></br><br></br>
              <SyntaxHighlighter language="protobuf" style={solarizedDark}>
                {this.state.fileText}
              </SyntaxHighlighter>
            </Card>
          </div>
        </div>
      </div>
    );
  }
}

export default ProtoBrowser;
