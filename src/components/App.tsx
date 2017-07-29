import * as React from 'react';

import { WordCard } from 'components/WordCard';

export class App extends React.Component<{}, {}> {
  render() {
    return (
      <div>
        <main>
          <WordCard front={'go'} back={['gehen', 'geht', 'ging', 'gegangen']} />
          <WordCard front={'good'} back={['gut', 'besser', 'am besten']} />
          <WordCard front={'dog'} back={['der Hund', 'die Hunde']} />
        </main>
      </div>
    );
  }
}
