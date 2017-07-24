import * as React from 'react';
import * as ReactDOM from 'react-dom';

import './index.less';

const rootEl = document.getElementById('root');

// -- HOT RELOAD SUPPORT -- //
/* tslint:disable */
declare const __DEV__: boolean;
declare const module: { hot: any };
declare const require: { (path: string): any; };
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
  const App = require('./components/App').App;
/* tslint:enable */
// -- HOT RELOAD SUPPORT -- //
  ReactDOM.render(<App />, rootEl);
}
