import React, { useEffect } from 'react';
import { Stack } from '@mui/joy';

import PackageTable from './PackageTable'
import { useSelector } from 'react-redux';
import { selectPackages } from './logic/package';

export default function Root() {
    const pkgs = useSelector(selectPackages)
    return (
        <Stack direction="column" sx={{ minHeight: '100dvh' }}>

            <h1>Faross analysed packages</h1>
            <PackageTable type={null} packages={pkgs} />

        </Stack >
    );
}
