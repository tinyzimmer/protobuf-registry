import React, { Component } from 'react';
import { PaginationOptions, TableColumns } from './table-options.jsx';

import BootstrapTable from 'react-bootstrap-table-next';
import paginationFactory from 'react-bootstrap-table2-paginator';
import ToolkitProvider, { Search } from 'react-bootstrap-table2-toolkit';
import { Tag, Card, Divider, Spinner } from "@blueprintjs/core";

import ProtobufMenu from './protobuf-menu.jsx';
import DeleteButton from './version-delete.jsx';

const { SearchBar } = Search;

const nameTag = (parent, data) => {
  return (
    <div className="wrapper">
      <DeleteButton name={data.name} version="*" callback={parent.deleteCallback} />
      &nbsp;&nbsp;
      <Tag icon="comment" large>{data.name}</Tag>
  </div>
  )
}

class ProtobufTable extends Component {
  constructor(props) {
    super(props)
    this.state = {data: [], loading: true}
    this.deleteCallback = this.deleteCallback.bind(this)
  }

  deleteCallback(name, version) {
    var newItems = [];
    this.state.data.map((value, index) => {
      if (value.rawName !== name) {
        newItems.push(value);
      }
      return ''
    })
    this.setState({data: newItems})
  }

  componentDidMount() {
    fetch('/api/proto')
    .then(results => {
      return results.json()
    }).then(data => {
      var newRows = []
      data.map((value, index) => {
        newRows.push({
          name: nameTag(this, value),
          latest: <Tag icon="git-branch" large>{value.latest}</Tag>,
          latestUploaded: <p style={{color: 'white'}}>{new Date(value.latestUploaded).toString().replace(/\(.*\)/, '')}</p>,
          details: <ProtobufMenu meta={value} />,
          rawName: value.name,
        });
        return ''
      });
      this.setState({data: newRows})
      this.setState({loading: false})
      console.log(this.state);
    })
  }

  render() {
    return (
      <ToolkitProvider
        data={ this.state.data }
        columns={ TableColumns }
        hover={true}
        bootstrap4={true}
        keyField="rawName"
        search
      >
      {
        props =>
          <div>
            <Card elevation="4" className="bp3-dark">
              <div className="wrapper" align="center">
                <strong>Filter protocols by name&nbsp;&nbsp;&nbsp;&nbsp;</strong>
                <SearchBar { ...props.searchProps } />
              </div>
            </Card>
            <div align="center" hidden={!this.state.loading}>
              <Spinner size={Spinner.SIZE_LARGE}></Spinner>
            </div>
            <Divider></Divider>
            <BootstrapTable
              { ...props.baseProps }
              bordered={false}
              pagination={ paginationFactory(PaginationOptions) }
            />
          </div>
      }
      </ToolkitProvider>
    )
  }
}

export default ProtobufTable
