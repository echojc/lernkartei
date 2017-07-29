import * as React from 'react';

import { Card } from './Card';

import * as styles from './wordCard.less';

interface Props {
  front: string | string[];
  back: string | string[];
}

function renderWords(words: string | string[]): React.ReactElement<{}> {
  const normalized = Array.isArray(words) ? words : [words];
  return (
    <div className={styles.words}>
      <div>
        {normalized.map(word => (<div key={word}>{word}</div>))}
      </div>
    </div>
  );
}

export class WordCard extends React.Component<Props, {}> {
  render() {
    const { front, back } = this.props;
    return (
      <Card
        front={renderWords(front)}
        back={renderWords(back)}
      />
    );
  }
}
