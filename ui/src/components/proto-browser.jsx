import React, { Component } from "react";
import SyntaxHighlighter from 'react-syntax-highlighter';
import { atomOneDark } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import { Classes, Tree, Card, Breadcrumbs, Icon, Divider } from "@blueprintjs/core";

const Header = () => {
  return (
    <div>
      <Card elevation="4" className="bp3-dark" >
        <h4 className="font-weight-bold">Protocol Browser (super beta)</h4>
      </Card>
      <Divider></Divider>
    </div>
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
        className: 'tree-node',
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
          className: 'tree-node',
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
      breadcrumbs: [],
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
      var crumbs = [
        { icon: 'globe-network', text: nodeData.parent },
        { icon: 'git-merge', text: nodeData.version },
      ]
      nodeData.fullPath.split('/').map((value, index) => {
        if (value === nodeData.label) {
          return ''
        } else if (value !== "") {
          crumbs.push({ icon: 'folder-open', text: value })
        }
        return ''
      })
      crumbs.push({ icon: 'document-open', text: nodeData.label })
      this.setState({breadcrumbs: crumbs})
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
          className: 'tree-node',
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
            className: 'tree-node',
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

  renderCurrentBreadcrumb({ text, ...restProps }) {
    return <Breadcrumbs {...restProps}>{text} <Icon icon="star" /></Breadcrumbs>;
  };

  render() {
    return (
      <div align="left" style={{paddingLeft: '5em', paddingRight: '5em'}}>
        <Header />
        <br></br>
        <div className="wrapper">
          <div style={{width: '35%'}}>
            <Card elevation="3" className="bp3-dark">
            <Tree
              contents={this.state.nodes}
              onNodeClick={this.handleNodeClick}
              onNodeCollapse={this.handleNodeCollapse}
              onNodeExpand={this.handleNodeExpand}
              className={Classes.TREE}
            />
            </Card>
          </div>
          &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
          <div hidden={this.state.fileViewHidden} style={{width: '65%'}}>
            <Card elevation="3" className="bp3-dark" style={{width: '100%'}}>
              <Breadcrumbs currentBreadcumbRenderer={this.renderCurrentBreadcrumb} items={this.state.breadcrumbs} />
              <br></br>
              <SyntaxHighlighter language="protobuf" style={atomOneDark}>
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
