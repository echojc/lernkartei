import * as React from 'react';

import * as styles from './card.less';

interface Props {
  front: string | string[];
  back: string | string[];
}

interface State {
  isFront: boolean;
}

function renderWords(words: string | string[]) {
  const normalized = Array.isArray(words) ? words : [words];
  return normalized.map(word => (<div key={word}>{word}</div>));
}

export class Card extends React.Component<Props, State> {
  state = {
    isFront: true,
  };

  render() {
    return (
      <section className={styles.cardContainer}>
        <article
          className={styles.card}
          onClick={() => this.setState({ isFront: !this.state.isFront })}
        >
          <div>
            {renderWords(this.state.isFront ? this.props.front : this.props.back)}
          </div>
        </article>
      </section>
    );
  }
}
