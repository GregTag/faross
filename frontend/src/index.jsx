import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import store from './logic/store'
import { fetchPackages } from './logic/package';

store.dispatch(fetchPackages());
const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>
);
