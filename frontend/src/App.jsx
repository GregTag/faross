import React from 'react';

import { CssVarsProvider } from '@mui/joy/styles';
import { CssBaseline } from '@mui/joy';
import { Provider } from 'react-redux';
import store from "./logic/store";

import Root from './Root';

export default function App() {
    return (
        <Provider store={store}>
            <CssVarsProvider defaultMode="dark" disableNestedContext>
                <CssBaseline />
                <Root />
            </CssVarsProvider>
        </Provider>
    )
}
