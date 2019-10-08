import React, { Component } from 'react';
import { PaginationOptions, TableColumns } from './table-options.jsx';

import BootstrapTable from 'react-bootstrap-table-next';
import paginationFactory from 'react-bootstrap-table2-paginator';
import ToolkitProvider, { Search } from 'react-bootstrap-table2-toolkit';
import { Tag } from "@blueprintjs/core";

import ProtobufMenu from './protobuf-menu.jsx';
import ProtoDownloadButton from './proto-download-button.jsx';
import DeleteButton from './version-delete.jsx';

const { SearchBar } = Search;

const nameTag = (parent, data) => {
  return (
    <div className="wrapper">
      <DeleteButton apiURL={parent.apiURL} name={data.name} version="*" callback={parent.deleteCallback} />
      &nbsp;&nbsp;
      <Tag icon="comment" large>{data.name}</Tag>
  </div>
  )
}

var weekday=new Array(7);
weekday[0]="Monday";
weekday[1]="Tuesday";
weekday[2]="Wednesday";
weekday[3]="Thursday";
weekday[4]="Friday";
weekday[5]="Saturday";
weekday[6]="Sunday";


class ProtobufTable extends Component {
  constructor(props) {
    super(props)
    this.apiURL = props.apiURL;
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
    fetch(this.apiURL + '/api/proto')
    .then(results => {
      return results.json()
    }).then(data => {
      var newRows = []
      data.map((value, index) => {
        newRows.push({
          name: nameTag(this, value),
          latest: <Tag icon="git-branch" large>{value.latest}</Tag>,
          latestUploaded: new Date(value.latestUploaded).toString().replace(/\(.*\)/, ''),
          details: <ProtobufMenu apiURL={this.apiURL} meta={value} />,
          download: <ProtoDownloadButton apiURL={this.apiURL} buttonText="Download Latest" name={value.name} version={value.latest} />,
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
