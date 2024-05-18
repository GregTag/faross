import { Chip } from '@mui/joy';
import { createAsyncThunk, createSlice } from '@reduxjs/toolkit';
import { PackageURL } from 'packageurl-js';

// const test_packages = [
//     {
//         "changed_at": "2024-05-17T00:54:08.932360442+03:00",
//         "comment": "",
//         "purl": "pkg:pypi/numpy",
//         "score": 10,
//         "state": "healthy"
//     },
//     {
//         "changed_at": "2024-05-17T00:55:40.210506992+03:00",
//         "comment": "",
//         "purl": "pkg:pypi/scipy",
//         "score": 0,
//         "state": "pending"
//     },
//     {
//         "changed_at": "2024-05-18T01:52:05.787423096+03:00",
//         "comment": "Some very long comment that should not fit in edges Some very long comment that should not fit in edges",
//         "purl": "pkg:pypi/advancedpythonmalware",
//         "score": 10,
//         "state": "healthy"
//     }
// ]

const request_url = "/api/status"

export const fetchPackages = createAsyncThunk('package/fetchPackages', async (_, thunkAPI) => {
    const response = await fetch(request_url);
    const data = await response.json();
    // const data = test_packages;
    return data
})

function mapColor(state) {
    switch (state) {
        case "pending":
            return "warning"
        case "healthy":
            return "success";
        case "quarantined":
            return "danger";
        case "unquarantined":
            return "primary";
        default:
            return "neutral";
    }
}

function mapper(pkg) {
    const purl = PackageURL.fromString(pkg.purl)
    const type = purl.type
    const name = (purl.namespace ? purl.namespace + '/' : '') + purl.name + (purl.version ? ' ' + purl.version : '')
    const changed_at = pkg.changed_at.slice(0, 19)
    const state_card = (<Chip color={mapColor(pkg.state)} >{pkg.state}</Chip>);
    return { type, name, changed_at, state_card, state: pkg.state, comment: pkg.comment }
}

const packageSlice = createSlice({
    name: 'package',
    initialState: {
        packages: [],
    },
    extraReducers: (builder) => {
        builder.addCase(fetchPackages.fulfilled, (state, action) => {
            state.packages = action.payload.map(mapper);
        })
    }
})

export const selectPackages = (state) => state.package.packages
export default packageSlice.reducer
