import * as React from 'react';
import * as ReactDOM from 'react-dom';

import { App } from './components/App';

import './index.less';

const rootEl = document.getElementById('root');

// -- HOT RELOAD SUPPORT -- //
/* tslint:disable */
declare var __DEV__: boolean;
declare var module: { hot: any };
declare var require: { (path: string): any; };
if (__DEV__) {
  const { AppContainer } = require('react-hot-loader');
  const renderApp = () => {
    const NextApp = require('./components/App').App;
    ReactDOM.render(<AppContainer><NextApp /></AppContainer>, rootEl);
  }
  renderApp();
  module.hot.accept('./components/App', renderApp);

  const MobXDevTools = require('mobx-react-devtools').default;
  ReactDOM.render(<MobXDevTools />, document.getElementById('dev'));

  // const { appState } = require('./state/state');
  // window.state = appState;
} else {
/* tslint:enable */
// -- HOT RELOAD SUPPORT -- //
  ReactDOM.render(<App />, rootEl);
}
