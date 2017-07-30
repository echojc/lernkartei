import * as classNames from 'classnames';
import { debounce } from 'lodash';
import * as React from 'react';

import { Card } from './Card';

import * as api from 'api/word';

import * as styles from './newCard.less';

interface Props {
  add: (front: string, back: string[]) => void;
}

interface State {
  input: string;
  pending: boolean;
  results: api.Result[] | null;
}

const renderResult = (add: (r: api.Result) => void) => (r: api.Result): React.ReactElement<{}> => {
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
    <section
      key={r.Base + r.PartOfSpeech}
      className={styles.result}
      onClick={() => add(r)}
    >
      <dl>
        <dt>{r.Base} <span className={styles.extended}>{pos} {forms}</span></dt>
        <dd>{defs}</dd>
      </dl>
    </section>
  );
};

export class NewCard extends React.Component<Props, State> {
  inflight: Promise<void> | null;

  state: State = {
    input: '',
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
        this.inflight = null;
        this.setState({ results, pending: false });
      })
      .catch(console.error);
    // keep track that this is now the newest request
    this.inflight = current;
  }, 500);

  add = (r: api.Result) => {
    const forms = r.PartOfSpeech === 'Noun'
      ? r.Forms
      : [r.Base].concat(r.Forms);

    this.props.add(r.Definitions[0] || '(unknown)', forms);
    this.setState({ input: '', results: null, pending: false });
  }

  render() {
    const { results, input } = this.state;
    return (
      <Card
        front={
          <div className={styles.search}>
            <input
              className={styles.input}
              placeholder={'Suchen...'}
              value={input}
              onChange={(e) => {
                const term = e.target.value;
                if (term !== '') {
                  this.search(term);
                  this.setState({ input: term });
                } else {
                  this.inflight = null;
                  this.search.cancel();
                  this.setState({ input: '', results: null, pending: false });
                }
              }}
            />
            <div className={classNames(styles.results, { [styles.hasResults]: !!results })}>
              {results && (
                results.length > 0
                  ? results.map(renderResult(this.add))
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
