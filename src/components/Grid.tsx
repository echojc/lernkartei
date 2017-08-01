import * as React from 'react';
import { SpringGrid, makeResponsive } from 'react-stonecutter';

import { WordCard } from 'components/WordCard';

import * as styles from './grid.less';

interface Card {
  front: string;
  back: string[];
  isNew: boolean;
}

interface Props {
  cards: Card[];
}

function key(card: Card): string {
  return card.front + card.back.join();
}

const ResponsiveGrid = makeResponsive(SpringGrid, {
  maxWidth: 1920,
  defaultColumns: 1,
});

export class Grid extends React.Component<Props, {}> {
  render() {
    const { cards } = this.props;

    return (
      <ResponsiveGrid
        className={styles.grid}
        columnWidth={280}
        itemHeight={160}
        gutterWidth={20}
        gutterHeight={20}
        springConfig={{ stiffness: 170, damping: 22 }}
      >
        {cards.map(card => (
          <div key={key(card)}>
            <WordCard
              front={card.front}
              back={card.back}
              disableAnimate={!card.isNew}
            />
          </div>
        ))}
      </ResponsiveGrid>
    );
  }
}
