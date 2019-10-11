import React, { Component } from 'react';

import { Tag, InputGroup, Button, MenuItem } from '@blueprintjs/core';
import { Suggest } from "@blueprintjs/select";

class RemoteDepSelect extends Component {
  constructor(props) {
    super(props);
    this.state = {
      deps: [],
      knownRemotes: [],
    }
    this.handleAddDep = this.handleAddDep.bind(this);
    this.handleRevChange = this.handleRevChange.bind(this);
    this.deleteDep = this.deleteDep.bind(this);
    this.handleItemSelect = this.handleItemSelect.bind(this);
    this.renderItem = this.renderItem.bind(this);
    this.renderNewItemFromQuery = this.renderNewItemFromQuery.bind(this);
  }

  componentDidMount() {
    var knownRemotes = []
    fetch('/api/remotes')
    .then(results => {
      return results.json()
    }).then(remotes => {
      knownRemotes = remotes
      this.setState({knownRemotes: knownRemotes})
    })
  }

  handleItemSelect(item, e, index) {
    fetch('/api/remotes', {method: "PUT", body: JSON.stringify({url: item})})
    .then(results => {
      return results.json()
    }).then(result => {
      console.log(result)
    })
    var deps = this.state.deps.slice()
    deps[index] = {
      url: item,
      revision: deps[index].revision
    }
    this.setState({deps: deps})
  }

  handleAddDep(e) {
    var deps = this.state.deps.slice()
    deps.push({
      url: "",
      revision: ""
    })
    this.setState({deps: deps})
  }

  handleRevChange(e, index) {
    var deps = this.state.deps.slice()
    deps[index] = {
      url: deps[index].url,
      revision: e.target.value
    }
    this.setState({deps: deps})
  }

  deleteDep(index) {
    var deps = this.state.deps.slice()
    deps.splice(index, 1)
    this.setState({deps: deps})
  }

  createNewItemFromQuery(query) {
    return query
  }

  renderNewItemFromQuery(item, active, handleClick) {
    return (
      <MenuItem
        active={active}
        disabled={false}
        onClick={handleClick}
        key={item}
        text={`Fetch ${item}...`}
      />
    );
  }

  renderItem(item, props) {
    return (
      <MenuItem
        active={props.modifiers.active}
        disabled={props.modifiers.disabled}
        onClick={props.handleClick}
        key={this.state.knownRemotes.indexOf(item)}
        text={item}
      />
    );
  }

  renderInput(item) {
    return item
  }

  render() {
    return (
      <div>
        {this.state.deps.map((dep, index) => {
          return (
            <div key={index}>
              <div key={index} className="wrapper">
                <Suggest
                  createNewItemFromQuery={this.createNewItemFromQuery}
                  createNewItemRenderer={this.renderNewItemFromQuery}
                  noResults={<MenuItem disabled={true} text="No results." />}
                  allowCreate={true}
                  resetOnSelect={true}
                  onItemSelect={(item, ev) => this.handleItemSelect(item, ev, index)}
                  items={this.state.knownRemotes}
                  itemRenderer={this.renderItem}
                  inputValueRenderer={this.renderInput}
                  value={this.state.deps[index].revision}
                  inputProps={{
                    round: true,
                    intent: 'primary',
                    leftIcon: 'git-repo',
                    placeholder: 'Git URL',
                    value: dep.url
                  }}
                />
                <InputGroup
                  round
                  onChange={(ev) => this.handleRevChange(ev, index)}
                  intent="primary"
                  placeholder="master"
                  leftIcon="git-branch"
                  value={this.state.deps[index].revision}
                />
                <Button onClick={() => this.deleteDep(index)} icon="delete" small></Button>
              </div>
              <br></br>
            </div>
          );
        })}
        <Tag
          fill={true}
          icon="add"
          interactive={true}
          round={true}
          onClick={this.handleAddDep}
        >
        Add remote import paths
        </Tag>
      </div>
    );
  }
}

export default RemoteDepSelect;
