import * as React from 'react';

import { Card } from 'components/Card';

export class App extends React.Component<{}, {}> {
  render() {
    return (
      <div>
        <main>
          <Card front={'go'} back={['gehen', 'geht', 'gang', 'gegangen']} />
          <Card front={'good'} back={['gut', 'besser', 'am besten']} />
          <Card front={'dog'} back={['der Hund', 'die Hunde']} />
        </main>
      </div>
    );
  }
}
