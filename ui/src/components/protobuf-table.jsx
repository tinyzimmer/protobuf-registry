import React, { Component } from 'react';
import { PaginationOptions, TableColumns } from './table-options.jsx';

import BootstrapTable from 'react-bootstrap-table-next';
import paginationFactory from 'react-bootstrap-table2-paginator';
import ToolkitProvider, { Search } from 'react-bootstrap-table2-toolkit';
import { Tag } from "@blueprintjs/core";

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
    this.state = {data: []}
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
          latestUploaded: new Date(value.latestUploaded).toString().replace(/\(.*\)/, ''),
          details: <ProtobufMenu meta={value} />,
          rawName: value.name,
        });
        return ''
      });
      this.setState({data: newRows})
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
            <strong><p align="center">Type to filter specs by name</p></strong>
            <SearchBar align="center" { ...props.searchProps } />
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
