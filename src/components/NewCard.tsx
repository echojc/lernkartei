import * as React from 'react';

import { Card } from './Card';

import * as styles from './newCard.less';

interface Result {
  Base: string;
  Definitions: string[];
  Forms: string[];
}

const results = [
  {
    Base: 'lesen',
    Definitions: [
      'read',
      'pick',
      'lecture',
      'get',
    ],
    Forms: [
      'liest',
      'las',
      'gelesen haben',
    ],
  },
  {
    Base: 'Hund',
    Definitions: [
      'dog',
      'hound',
      'canid',
      'scoundrel',
    ],
    Forms: [
      'der Hund',
      'die Hunde',
    ],
  },
  {
    Base: 'gut',
    Definitions: [
      'good',
      'well',
      'all right',
      'fine',
    ],
    Forms: [
      'besser',
      'am besten',
    ],
  },
];

function renderResult(r: Result): React.ReactElement<{}> {
  const forms = r.Forms
    ? <span className={styles.extended}>({r.Forms.join(', ')})</span>
    : null;

  const defs = r.Definitions
    ? r.Definitions.join(', ')
    : null;

  return (
    <section className={styles.result} key={r.Base}>
      <dl>
        <dt>{r.Base} {forms}</dt>
        <dd>{defs}</dd>
      </dl>
    </section>
  );
}

export class NewCard extends React.Component<{}, {}> {
  render() {
    return (
      <Card
        front={
          <div className={styles.search}>
            <input className={styles.input} placeholder={'Suchen...'} />
            <div className={styles.results}>
              {results.map(renderResult)}
            </div>
          </div>
        }
        disableFlip={true}
      />
    );
  }
}
