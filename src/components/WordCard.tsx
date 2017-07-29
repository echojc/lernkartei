import * as React from 'react';

import { Card } from './Card';

interface Props {
  front: string | string[];
  back: string | string[];
}

function renderWords(words: string | string[]): React.ReactElement<{}>[] {
  const normalized = Array.isArray(words) ? words : [words];
  return normalized.map(word => (<div key={word}>{word}</div>));
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
