import React, { Component } from "react";
import { Tag, Tabs, Tab } from '@blueprintjs/core';

import MessageTables from './message-tables.jsx';
import ServicesTables from './services-tables.jsx'
import ExtensionsTable from './extensions-table.jsx';
import EnumsTable from './enums-table.jsx';

class ProtoDocTables extends Component {
  constructor(props) {
    super(props);
    this.docs = props.docs
    var selectedTabId = ""
    if (this.docs.hasMessages) {
      selectedTabId = "msgs"
    } else if (this.docs.hasServices) {
      selectedTabId = "svcs"
    } else if (this.docs.hasExtensions) {
      selectedTabId = "exts"
    } else if (this.docs.enums.length) {
      selectedTabId = "enms"
    }
    this.state = {
      selectedTabId: selectedTabId
    }
    this.handleTabChange = this.handleTabChange.bind(this)
  }

  handleTabChange(tabId) {
    this.setState({selectedTabId: tabId})
  }

  render() {
    return (
      <div>
        <div align="center">
          <h4>{this.docs.name}</h4>
        </div>
        <div>
          <i>{this.docs.description}</i>
          <br></br>
          <br></br>
          {this.docs.package && (this.docs.package !== "") ? <div><strong>Package:  </strong><Tag minimal={true}>{this.docs.package}</Tag></div> : ""}
        </div>
        <br></br>
        <Tabs animate={true} key="horizontal" renderActiveTabPanelOnly={true} vertical={false} onChange={this.handleTabChange} selectedTabId={this.state.selectedTabId}>
          <Tab id="msgs" title="Messages" disabled={!this.docs.hasMessages} panel={<MessageTables messages={this.docs.messages} />} />
          <Tab id="svcs" title="Services" disabled={!this.docs.hasServices} panel={<ServicesTables services={this.docs.services} />} />
          <Tab id="exts" title="Extensions" disabled={!this.docs.hasExtensions} panel={<ExtensionsTable extensions={this.docs.extensions} />} />
          <Tab id="enms" title="Enums" disabled={this.docs.enums.length === 0} panel={<EnumsTable enums={this.docs.enums} />} />
        </Tabs>
      </div>
    );
  }
}

export default ProtoDocTables;
