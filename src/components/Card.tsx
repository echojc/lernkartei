import * as classNames from 'classnames';
import * as React from 'react';

import * as styles from './card.less';

interface Props {
  front: React.ReactElement<{}> | React.ReactElement<{}>[];
  back?: React.ReactElement<{}> | React.ReactElement<{}>[];
  disableFlip?: boolean;
  disableEnter?: boolean;
}

interface State {
  isFlipped: boolean;
  isEntered: boolean;
}

export class Card extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    const { disableFlip, disableEnter } = this.props;
    this.state = {
      isFlipped: !disableFlip && !disableEnter,
      isEntered: !!disableEnter,
    };

    if (!disableEnter) {
      setTimeout(() => {
        this.setState({ isFlipped: false, isEntered: true });
      }, 3000);
    }
  }

  render() {
    const { front, back, disableFlip } = this.props;
    const { isFlipped, isEntered } = this.state;

    return (
      <article className={styles.container}>
        <section
          className={classNames(styles.card, { [styles.flipped]: isFlipped })}
          onClick={() => isEntered && !disableFlip && this.setState({ isFlipped: !isFlipped })}
        >
          <div className={styles.front}>{front}</div>
          <div className={styles.back}>{back}</div>
        </section>
      </article>
    );
  }
}
