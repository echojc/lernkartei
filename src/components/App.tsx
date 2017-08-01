import { findIndex } from 'lodash';
import * as React from 'react';

import { Grid } from 'components/Grid';
import { Search } from 'components/Search';

interface Card {
  front: string;
  back: string[];
  isNew: boolean;
}

interface State {
  cards: Card[];
}

function key(card: Card): string {
  return card.front + card.back.join();
}

export class App extends React.Component<{}, {}> {
  state: State = {
    cards: [
      { isNew: false, front: 'go',   back: ['gehen', 'geht', 'ging', 'gegangen sein'] },
      { isNew: false, front: 'good', back: ['gut', 'besser', 'am besten'] },
      { isNew: false, front: 'dog',  back: ['der Hund', 'die Hunde'] },
    ],
  };

  add = (front: string, back: string[]) => {
    const newCard = {
      isNew: true,
      front,
      back,
    };

    const { cards } = this.state;
    const existingIndex = findIndex(cards, c => key(c) === key(newCard));
    if (existingIndex < 0) {
      this.setState({ cards: [newCard].concat(cards) });
    } else {
      const before = cards.slice(0, existingIndex);
      const after = cards.slice(existingIndex + 1);
      this.setState({ cards: [newCard].concat(before).concat(after) });
    }
  }

  render() {
    const { cards } = this.state;
    return (
      <div>
        <main>
          <Search add={this.add} />
          <Grid cards={cards} />
        </main>
      </div>
    );
  }
}
