import { findIndex, pick } from 'lodash';
import * as React from 'react';

import { Grid } from 'components/Grid';
import { Search } from 'components/Search';

interface Card {
  front: string;
  back: string[];
  isNew?: boolean;
}

interface State {
  cards: Card[];
}

function key(card: Card): string {
  return card.front + card.back.join();
}

function loadCards(): Card[] {
  try {
    const data = localStorage.getItem('cards');
    const decoded = data ? JSON.parse(data) : [];
    return Array.isArray(decoded) ? decoded : [];
  } catch (e) {
    // tslint:disable-next-line
    console.error(e);
    return [];
  }
}

function saveCards(cards: Card[]): void {
  const stripped = cards.map(c => pick(c, 'front', 'back'));
  console.dir(stripped);
  localStorage.setItem('cards', JSON.stringify(stripped));
}

export class App extends React.Component<{}, State> {
  constructor(props: {}) {
    super(props);
    this.state = {
      cards: loadCards(),
    };
  }

  add = (front: string, back: string[]) => {
    const { cards } = this.state;
    const newCard: Card = {
      front,
      back,
      isNew: true,
    };

    const existingIndex = findIndex(cards, c => key(c) === key(newCard));
    const newCards = existingIndex < 0
      ? [newCard].concat(cards)
      : [newCard].concat(cards.slice(0, existingIndex), cards.slice(existingIndex + 1));

    saveCards(newCards);
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
