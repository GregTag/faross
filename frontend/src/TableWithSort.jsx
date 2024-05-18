import React, { useState } from 'react';
import { Box, Table, Typography, FormControl, FormLabel, IconButton, Link, Select, Option } from '@mui/joy';
import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward';
import { visuallyHidden } from '@mui/utils';

function labelDisplayedRows({ from, to, count }) {
    return `${from}–${to} of ${count !== -1 ? count : `more than ${to}`}`;
}

function descendingComparator(a, b, orderBy) {
    if (b[orderBy] < a[orderBy]) {
        return -1;
    }
    if (b[orderBy] > a[orderBy]) {
        return 1;
    }
    return 0;
}

function getComparator(order, orderBy) {
    return order === 'desc'
        ? (a, b) => descendingComparator(a, b, orderBy)
        : (a, b) => -descendingComparator(a, b, orderBy);
}

function EnhancedTableHead({ order, orderBy, onRequestSort, head }) {
    const createSortHandler = (property) => (event) => {
        onRequestSort(event, property);
    };

    return (
        <thead>
            <tr>
                {head.map((headCell) => {
                    const active = orderBy === headCell.id;
                    return (
                        <th
                            key={headCell.id}
                            aria-sort={
                                active ? { asc: 'ascending', desc: 'descending' }[order] : undefined
                            }
                        >
                            <Link
                                underline="none"
                                color="neutral"
                                textColor={active ? 'primary.plainColor' : undefined}
                                component="button"
                                onClick={createSortHandler(headCell.id)}
                                fontWeight="lg"
                                endDecorator={
                                    <ArrowDownwardIcon sx={{ opacity: active ? 1 : 0 }} />
                                }
                                sx={{
                                    '& svg': {
                                        transition: '0.2s',
                                        transform:
                                            active && order === 'desc' ? 'rotate(0deg)' : 'rotate(180deg)'
                                    },
                                    '&:hover': { '& svg': { opacity: 1 } }
                                }}
                            >
                                {headCell.label}
                                {active
                                    ? (
                                        <Box component="span" sx={visuallyHidden}>
                                            {order === 'desc' ? 'sorted descending' : 'sorted ascending'}
                                        </Box>
                                    )
                                    : null}
                            </Link>
                        </th>
                    );
                })}
            </tr>
        </thead >
    );
}

export default function TableWithSort({ head, rows, rowsPerPageOptions, handleClick, defaultOrder }) {
    const [order, setOrder] = useState('asc');
    const [orderBy, setOrderBy] = useState(defaultOrder);
    const [page, setPage] = useState(0);
    const [rowsPerPage, setRowsPerPage] = useState(rowsPerPageOptions[0]);

    const handleRequestSort = (event, property) => {
        const isAsc = orderBy === property && order === 'asc';
        setOrder(isAsc ? 'desc' : 'asc');
        setOrderBy(property);
    };

    const handleChangePage = (newPage) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event, newValue) => {
        setRowsPerPage(parseInt(newValue.toString(), 10));
        setPage(0);
    };

    const getLabelDisplayedRowsTo = () => {
        if (rows.length === -1) {
            return (page + 1) * rowsPerPage;
        }
        return rowsPerPage === -1
            ? rows.length
            : Math.min(rows.length, (page + 1) * rowsPerPage);
    };

    // Avoid a layout jump when reaching the last page with empty rows.
    const emptyRows =
        page > 0 ? Math.max(0, (1 + page) * rowsPerPage - rows.length) : 0;

    return (
        <Table
            aria-labelledby="tableTitle"
            hoverRow
            sx={{
                '--TableCell-headBackground': 'transparent',
                '--TableCell-selectedBackground': (theme) =>
                    theme.vars.palette.success.softBg,
                '& thead th': {
                    width: `calc(100% / ${head.length})`
                }
            }}
        >
            <EnhancedTableHead
                order={order}
                orderBy={orderBy}
                onRequestSort={handleRequestSort}
                head={head}
            />
            <tbody>
                {rows.slice().sort(getComparator(order, orderBy))
                    .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                    .map((row) => {
                        return (
                            <tr
                                key={row.id}
                                onClick={(event) => handleClick(event, row)}
                            >
                                {head.map(({ id, element }, index) => (<th key={index}>{row[element] || row[id]}</th>))}
                            </tr>
                        );
                    })}
                {emptyRows > 0 && (
                    <tr
                        style={{
                            height: `calc(${emptyRows} * 40px)`,
                            '--TableRow-hoverBackground': 'transparent'
                        }}
                    >
                        <td colSpan={6} aria-hidden />
                    </tr>
                )}
            </tbody>
            <tfoot>
                <tr>
                    <td colSpan={6}>
                        <Box
                            sx={{
                                display: 'flex',
                                alignItems: 'center',
                                gap: 2,
                                justifyContent: 'flex-end'
                            }}
                        >
                            <FormControl orientation="horizontal" size="sm">
                                <FormLabel>Rows per page:</FormLabel>
                                <Select onChange={handleChangeRowsPerPage} value={rowsPerPage}>
                                    {rowsPerPageOptions.map((option) => (
                                        <Option key={option} value={option}>{option}</Option>))}
                                </Select>
                            </FormControl>
                            <Typography textAlign="center" sx={{ minWidth: 80 }}>
                                {labelDisplayedRows({
                                    from: rows.length === 0 ? 0 : page * rowsPerPage + 1,
                                    to: getLabelDisplayedRowsTo(),
                                    count: rows.length === -1 ? -1 : rows.length
                                })}
                            </Typography>
                            <Box sx={{ display: 'flex', gap: 1 }}>
                                <IconButton
                                    size="sm"
                                    color="neutral"
                                    variant="outlined"
                                    disabled={page === 0}
                                    onClick={() => handleChangePage(page - 1)}
                                    sx={{ bgcolor: 'background.surface' }}
                                >
                                    <KeyboardArrowLeftIcon />
                                </IconButton>
                                <IconButton
                                    size="sm"
                                    color="neutral"
                                    variant="outlined"
                                    disabled={
                                        rows.length !== -1
                                            ? page >= Math.ceil(rows.length / rowsPerPage) - 1
                                            : false
                                    }
                                    onClick={() => handleChangePage(page + 1)}
                                    sx={{ bgcolor: 'background.surface' }}
                                >
                                    <KeyboardArrowRightIcon />
                                </IconButton>
                            </Box>
                        </Box>
                    </td>
                </tr>
            </tfoot>
        </Table>
    );
}
