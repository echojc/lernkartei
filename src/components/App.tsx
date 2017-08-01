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

function deserialize(s: string): Card[] {
  try {
    const cards = JSON.parse(s);
    cards.forEach((c: Card) => c.isNew = false);
    return cards;
  } catch (e) {
    // tslint:disable-next-line
    console.error(e);
    return [];
  }
}

export class App extends React.Component<{}, State> {
  constructor(props: {}) {
    super(props);
    const data = localStorage.getItem('cards');
    const cards = data ? deserialize(data) : [];
    this.state = {
      cards,
    };
  }

  add = (front: string, back: string[]) => {
    const { cards } = this.state;
    const newCard = {
      isNew: true,
      front,
      back,
    };

    const existingIndex = findIndex(cards, c => key(c) === key(newCard));
    const newCards = existingIndex < 0
      ? [newCard].concat(cards)
      : [newCard].concat(cards.slice(0, existingIndex)).concat(cards.slice(existingIndex + 1));

    localStorage.setItem('cards', JSON.stringify(newCards));
    this.setState({ cards: newCards });
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
