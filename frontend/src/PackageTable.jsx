import React from 'react';
import TableWithSort from './TableWithSort';

const head = [
    {
        id: 'type',
        label: 'Type'
    },
    {
        id: 'name',
        label: 'Name'
    },
    {
        id: 'state',
        label: 'State',
        element: 'state_card'
    },
    {
        id: 'changed_at',
        label: 'Last update'
    },
    {
        id: 'comment',
        label: 'Comment'
    },
];


const options = [20, 100, 500, 2500];

export default function PackageTable({ packages }) {
    const handler = (event, row) => { };
    return (
        <TableWithSort head={head} rowsPerPageOptions={options} rows={packages} handleClick={handler} defaultOrder="changed_at" />
    );
}
