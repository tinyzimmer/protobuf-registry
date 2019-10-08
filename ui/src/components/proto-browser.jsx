import React, { Component } from "react";
import SyntaxHighlighter from 'react-syntax-highlighter';
import { solarizedDark } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import { Classes, Tree, Card } from "@blueprintjs/core";

const Header = () => {
  return (
    <h4 className="font-weight-bold">File Browser (super beta)</h4>
  )
}

function enumerateFiles(nodeData, cb) {
  var files = []
  var directories = []
  var knownDirs = []
  nodeData.rawChildren.map((value, index) => {
    if (nodeData.isDir) {
      value = value.replace(nodeData.label+'/', '')
    }
    var split = value.split('/')
    if (split.length === 1) {
      var file = {
        id: index,
        hasCaret: false,
        icon: "document-open",
        label: value,
        isFile: true,
        parent: nodeData.parent,
        version: nodeData.version,
        fullPath: [nodeData.fullPath, value].join('/'),
      }
      files.push(file)
    } else {
      if (!knownDirs.includes(split[0])) {
        knownDirs.push(split[0])
        var rawChildren = []
        nodeData.rawChildren.map((v, i) => {
          var spl = v.split('/')
          if (spl[0] === split[0]) {
            rawChildren.push(v.replace(split[0] + '/', ''))
          }
          return ''
        })
        var dir = {
          id: index,
          hasCaret: true,
          icon: "folder-close",
          label: split[0],
          isDir: true,
          parent: nodeData.parent,
          version: nodeData.version,
          fullPath: [nodeData.fullPath, split[0]].join('/'),
          rawChildren: rawChildren
        }
        directories.push(dir)
      }
    }
    return ''
  })
  cb(directories, files)
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
    this.handleDirExpand = this.handleDirExpand.bind(this)
    this.handleNodeClick = this.handleNodeClick.bind(this)
    this.handleNodeCollapse = this.handleNodeCollapse.bind(this)
    this.handleNodeExpand = this.handleNodeExpand.bind(this)
    this.handleVersionExpand = this.handleVersionExpand.bind(this)
    this.forEachNode = this.forEachNode.bind(this)
  }

  handleFileClick(nodeData) {
    console.log(nodeData)
    var url = '/api/proto/' + nodeData.parent + '/' + nodeData.version + '/raw' + nodeData.fullPath
    fetch(url)
    .then(results => {
      return results.text()
    }).then(fileText => {
      this.setState({fileText: fileText})
      this.setState({fileViewHidden: false})
    })
  }

  handleDirExpand(nodeData) {
    enumerateFiles(nodeData, (directories, files) => {
      nodeData.childNodes = directories.concat(files)
      this.setState(this.state)
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
    if (nodeData.isDir) {
      nodeData.icon = 'folder-close'
    }
    nodeData.isExpanded = false;
    this.setState(this.state);
  }

  handleVersionExpand(nodeData) {
    fetch('/api/proto/' + nodeData.parent + '/' + nodeData.label)
    .then(results => {
      return results.json()
    }).then(data => {
      nodeData.rawChildren = data.sourceFiles
      nodeData.fullPath = ''
      enumerateFiles(nodeData, (directories, files) => {
        nodeData.childNodes = directories.concat(files)
        this.setState(this.state)
      })
    })
  }

  handleNodeExpand(nodeData) {
    if (nodeData.isVersion) {
      this.handleVersionExpand(nodeData)
    } else if (nodeData.isDir) {
      nodeData.icon = 'folder-open'
      this.handleDirExpand(nodeData)
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
          icon: "globe-network",
          label: value.name,
        }
        var children = []
        value.versions.map((version, i) => {
          children.push({
            id: i,
            hasCaret: true,
            icon: 'git-merge',
            label: version.version,
            isVersion: true,
            version: version.version,
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
      <div align="left" style={{paddingLeft: '5em', paddingRight: '5em'}}>
        <Header />
        <br></br>
        <div className="wrapper">
          <div style={{width: '35%'}}>
            <Tree
              contents={this.state.nodes}
              onNodeClick={this.handleNodeClick}
              onNodeCollapse={this.handleNodeCollapse}
              onNodeExpand={this.handleNodeExpand}
              className={Classes.ELEVATION_2}
            />
          </div>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
          <div hidden={this.state.fileViewHidden} style={{width: '65%'}}>
            <Card elevation="3" className="bp3-dark" style={{width: '100%'}}>
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
