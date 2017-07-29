import * as React from 'react';

import { NewCard } from 'components/NewCard';
import { WordCard } from 'components/WordCard';

interface Card {
  front: string;
  back: string[];
}

interface State {
  cards: Card[];
}

export class App extends React.Component<{}, {}> {
  state: State = {
    cards: [
      { front: 'go',   back: ['gehen', 'geht', 'ging', 'gegangen'] },
      { front: 'good', back: ['gut', 'besser', 'am besten'] },
      { front: 'dog',  back: ['der Hund', 'die Hunde'] },
    ],
  };

  add = (front: string, back: string[]) => {
    const { cards } = this.state;
    this.setState({ cards: [{ front, back }].concat(cards) });
  }

  render() {
    const { cards } = this.state;
    return (
      <div>
        <main>
          <NewCard add={this.add} />
          {cards.map(card => (
            <WordCard key={card.back.join()} front={card.front} back={card.back} />
          ))}
        </main>
      </div>
    );
  }
}
