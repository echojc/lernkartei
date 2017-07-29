import * as classNames from 'classnames';
import { debounce } from 'lodash';
import * as React from 'react';

import { Card } from './Card';

import * as api from 'api/word';

import * as styles from './newCard.less';

function renderResult(r: api.Result): React.ReactElement<{}> {
  const forms = r.Forms.length > 0
    ? `(${r.Forms.join(', ')})`
    : null;

  const defs = r.Definitions
    ? r.Definitions.join(', ')
    : null;

  const pos = r.PartOfSpeech
    ? `[${api.partOfSpeech(r.PartOfSpeech)}]`
    : null;

  return (
    <section className={styles.result} key={r.Base + r.PartOfSpeech}>
      <dl>
        <dt>{r.Base} <span className={styles.extended}>{pos} {forms}</span></dt>
        <dd>{defs}</dd>
      </dl>
    </section>
  );
}

interface State {
  pending: boolean;
  results: api.Result[] | null;
}

export class NewCard extends React.Component<{}, State> {
  inflight: Promise<void>;

  state: State = {
    pending: false,
    results: null,
  };

  search = debounce((word: string) => {
    this.setState({ pending: true });
    const current = api.lookup(word)
      .then(results => {
        // ignore if this isn't the newest request
        if (current !== this.inflight) {
          return;
        }
        this.setState({ results, pending: false });
      })
      .catch(console.error);
    // keep track that this is now the newest request
    this.inflight = current;
  }, 500);

  render() {
    const { results } = this.state;
    return (
      <Card
        front={
          <div className={styles.search}>
            <input
              className={styles.input}
              placeholder={'Suchen...'}
              onChange={(e) => {
                const term = e.target.value;
                if (term !== '') {
                  this.search(term);
                } else {
                  this.search.cancel();
                  this.setState({ results: null });
                }
              }}
            />
            <div className={classNames(styles.results, { [styles.hasResults]: !!results })}>
              {results && (
                results.length > 0
                  ? results.map(renderResult)
                  : <div className={styles.result}>
                      <dl>
                        <dd>No results</dd>
                      </dl>
                    </div>
              )}
            </div>
          </div>
        }
        disableFlip={true}
      />
    );
  }
}
