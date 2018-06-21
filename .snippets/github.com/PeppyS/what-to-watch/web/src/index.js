import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux'
import Store from './Services/Store'
import App from './App';
import registerServiceWorker from './registerServiceWorker';
import './Style/index.css';

ReactDOM.render(
    <Provider store={Store}>
        <App />
    </Provider>,
    document.getElementById('root')
);

registerServiceWorker();
