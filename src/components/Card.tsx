import * as classNames from 'classnames';
import * as React from 'react';

import * as styles from './card.less';

interface Props {
  front: React.ReactElement<{}> | React.ReactElement<{}>[];
  back: React.ReactElement<{}> | React.ReactElement<{}>[];
}

interface State {
  isFlipped: boolean;
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
          <div className={styles.front}>{front}</div>
          <div className={styles.back}>{back}</div>
        </section>
      </article>
    );
  }
}
