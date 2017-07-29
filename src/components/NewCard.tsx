import * as React from 'react';

import { Card } from './Card';

import * as styles from './newCard.less';

export class NewCard extends React.Component<{}, {}> {
  render() {
    return (
      <Card
        front={
          <input className={styles.input} />
        }
        disableFlip={true}
      />
    );
  }
}
