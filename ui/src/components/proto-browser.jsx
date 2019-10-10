import React, { Component } from "react";
import SyntaxHighlighter from 'react-syntax-highlighter';
import { atomOneDark } from 'react-syntax-highlighter/dist/esm/styles/hljs';
import {
  Classes,
  Tree,
  Card,
  Breadcrumbs,
  Icon,
  Divider,
  Spinner,
  Tag,
  Collapse,
} from "@blueprintjs/core";
import ProtoDocTables from './docs/proto-doc-tables.jsx';

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
  directories.sort((a, b) => (a.label > b.label) ? 1 : -1)
  files.sort((a, b) => (a.label > b.label) ? 1 : -1)
  cb(directories, files)
}


class ProtoBrowser extends Component {
  constructor(props) {
    super(props);
    this.state = {
      nodes: [],
      fileViewHidden: true,
      fileTextExpanded: false,
      docTextExpanded: true,
      fileText: "",
      docContents: {
        files: [],
        scalarValueTypes: [],
      },
      fileTextHeader: "",
      breadcrumbs: [],
      loading: true,
    }
    this.handleFileContentsClick = this.handleFileContentsClick.bind(this)
    this.handleDocContentsClick = this.handleDocContentsClick.bind(this)
    this.handleFileClick = this.handleFileClick.bind(this)
    this.handleDirExpand = this.handleDirExpand.bind(this)
    this.handleNodeClick = this.handleNodeClick.bind(this)
    this.handleNodeCollapse = this.handleNodeCollapse.bind(this)
    this.handleNodeExpand = this.handleNodeExpand.bind(this)
    this.handleVersionExpand = this.handleVersionExpand.bind(this)
    this.forEachNode = this.forEachNode.bind(this)
  }

  handleFileContentsClick() {
    this.setState({fileTextExpanded: !this.state.fileTextExpanded})
  }

  handleDocContentsClick() {
    this.setState({docTextExpanded: !this.state.docTextExpanded})
  }

  handleFileClick(nodeData) {
    var url = '/api/proto/' + nodeData.parent + '/' + nodeData.version + '/raw' + nodeData.fullPath

    // fetch file text
    fetch(url)
    .then(results => {
      return results.json()
    }).then(res => {
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

      if (res.error !== undefined) {
        this.setState({fileText: "// " + nodeData.fullPath.substr(1) + " appears to be a remote import to this package"})
      }  else {
        this.setState({fileText: res.content})
      }
    })

    // fetch file docs
    url = '/api/proto/' + nodeData.parent + '/' + nodeData.version + '/meta' + nodeData.fullPath
    fetch(url)
    .then(results => {
      return results.json()
    }).then(data => {
      this.setState({
        docContents: {
          files: data.files,
          scalarValueTypes: data.scalarValueTypes
        }
      })
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
      nodes.sort((a, b) => (a.label > b.label) ? 1 : -1)
      this.setState({loading: false})
      this.setState({nodes: nodes})
    })
  }

  renderCurrentBreadcrumb({ text, ...restProps }) {
    return <Breadcrumbs {...restProps}>{text} <Icon icon="star" /></Breadcrumbs>;
  };

  render() {
    return (
      <div align="center">
        <Header />
        <div align="center" hidden={!this.state.loading}>
          <Spinner size={Spinner.SIZE_LARGE}></Spinner>
        </div>
        <div hidden={this.state.nodes.length !== 0 || this.state.loading}>
          <Card elevation="3" className="bp3-dark">
            The registry is empty
          </Card>
        </div>
        <div align="left" style={{paddingLeft: '5em', paddingRight: '5em'}}>
          <br></br>
          <div className="wrapper">
            <div
              hidden={this.state.nodes.length === 0 || this.state.loading}
              style={{paddingRight: '2em', height: '700px', overflowX: 'auto', overflowY: 'auto', display: 'flex', flexDirection: "column"}}
            >
                <Tree
                  contents={this.state.nodes}
                  onNodeClick={this.handleNodeClick}
                  onNodeCollapse={this.handleNodeCollapse}
                  onNodeExpand={this.handleNodeExpand}
                  className={Classes.TREE}
                />
            </div>
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
            <div hidden={this.state.fileViewHidden} style={{width: '90%'}}>
              <Card elevation="3" className="bp3-dark" style={{width: '100%'}}>
                <Breadcrumbs currentBreadcumbRenderer={this.renderCurrentBreadcrumb} items={this.state.breadcrumbs} />
                <br></br>
                <Tag interactive icon={this.state.fileTextExpanded ? "caret-down" : "caret-right"} onClick={this.handleFileContentsClick}>File Contents</Tag>
                <Collapse isOpen={this.state.fileTextExpanded}>
                  <SyntaxHighlighter language="protobuf" style={atomOneDark}>
                    {this.state.fileText}
                  </SyntaxHighlighter>
                </Collapse>
                <br></br>
                <br></br>
                <Tag interactive icon={this.state.docTextExpanded ? "caret-down" : "caret-right"} onClick={this.handleDocContentsClick}>Documentation</Tag>
                <br></br>
                <br></br>
                <Collapse isOpen={this.state.docTextExpanded}>
                  <div>
                    {this.state.docContents.files.length && this.state.docContents.files.map((file, index) => {
                      return <ProtoDocTables key={file.name} docs={file} />
                    })}
                  </div>
                </Collapse>
              </Card>
            </div>
          </div>
          <br></br>
        </div>
      </div>
    );
  }
}

export default ProtoBrowser;
