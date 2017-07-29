import * as classNames from 'classnames';
import * as React from 'react';

import * as styles from './card.less';

interface Props {
  front: string | string[];
  back: string | string[];
}

interface State {
  isFlipped: boolean;
}

function renderWords(words: string | string[]) {
  const normalized = Array.isArray(words) ? words : [words];
  return normalized.map(word => (<div key={word}>{word}</div>));
}

export class Card extends React.Component<Props, State> {
  state = {
    isFlipped: false,
  };

  render() {
    const { front, back } = this.props;
    const { isFlipped } = this.state;

    return (
      <article className={styles.container}>
        <section
          className={classNames(styles.card, { [styles.flipped]: isFlipped })}
          onClick={() => this.setState({ isFlipped: !isFlipped })}
        >
          <div className={styles.front}>{renderWords(front)}</div>
          <div className={styles.back}>{renderWords(back)}</div>
        </section>
      </article>
    );
  }
}
