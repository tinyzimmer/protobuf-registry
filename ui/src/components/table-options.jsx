import React from 'react';

export const CustomTotal = (from, to, size) => (
  <span className="react-bootstrap-table-pagination-total">
    &nbsp;&nbsp;Showing { from } to { to } of { size } Protobuf Specs
  </span>
);

export const PaginationOptions = {
  paginationSize: 5,
  pageStartIndex: 1,
  // alwaysShowAllBtns: true, // Always show next and previous button
  // withFirstAndLast: false, // Hide the going to First and Last page button
  // hideSizePerPage: true, // Hide the sizePerPage dropdown always
  // hidePageListOnlyOnePage: true, // Hide the pagination list when only one page
  firstPageText: 'First',
  prePageText: 'Back',
  nextPageText: 'Next',
  lastPageText: 'Last',
  // nextPageTitle: 'First page',
  // prePageTitle: 'Pre page',
  // firstPageTitle: 'Next page',
  // lastPageTitle: 'Last page',
  showTotal: true,
  paginationTotalRenderer: CustomTotal,
  sizePerPageList: [
    {
      text: '15', value: 15
    },
    {
      text: '25', value: 25
    },
    {
      text: '50', value: 50
    }
  ]
};

export const TableColumns = [
  {
    text: 'Name',
    dataField: 'name',
    sort: true,
    sortValue: (cell, row) => row.rawName,
    align: 'left'
  },
  {text: 'Latest', dataField: 'latest' },
  {text: 'Uploaded', dataField: 'latestUploaded'},
  {text: '', dataField: 'details'},
  {text: '', dataField: 'rawName', hidden: true}
]
