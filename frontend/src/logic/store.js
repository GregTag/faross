import { configureStore } from "@reduxjs/toolkit";
import packageReducer from './package'


export default configureStore({
    reducer: {
        package: packageReducer
    }
})
