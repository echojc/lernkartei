import * as React from 'react';

import * as styles from './app.less';

export class App extends React.Component<{}, {}> {
  render() {
    return (
      <main>
        <article className={styles.card}>WORD</article>
      </main>
    );
  }
}
