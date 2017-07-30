import { uniqBy } from 'lodash';
import * as React from 'react';

import { Search } from 'components/Search';
import { WordCard } from 'components/WordCard';

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
      { isNew: false, front: 'go',   back: ['gehen', 'geht', 'ging', 'gegangen'] },
      { isNew: false, front: 'good', back: ['gut', 'besser', 'am besten'] },
      { isNew: false, front: 'dog',  back: ['der Hund', 'die Hunde'] },
    ],
  };

  add = (front: string, back: string[]) => {
    const { cards } = this.state;
    this.setState({ cards: [{ isNew: true, front, back }].concat(cards) });
  }

  render() {
    const { cards } = this.state;
    return (
      <div>
        <main>
          <Search add={this.add} />
          {uniqBy(cards, key).map(card => (
            <WordCard key={key(card)} front={card.front} back={card.back} disableAnimate={!card.isNew} />
          ))}
        </main>
      </div>
    );
  }
}
