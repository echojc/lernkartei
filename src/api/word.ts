const baseUrl = 'http://localhost:9000';

type PartOfSpeech = 'Noun' | 'Adjective' | 'Verb';

export function partOfSpeech(pos: PartOfSpeech): string {
  switch (pos) {
    case 'Noun': return 'S';
    case 'Adjective': return 'A';
    case 'Verb': return 'V';
    default:
     const n: never = pos;
     return n;
  }
};

export interface Result {
  Base: string;
  Definitions: string[];
  Forms: string[];
  PartOfSpeech: PartOfSpeech;
}

export function lookup(word: string): Promise<Result[]> {
  return fetch(`${baseUrl}/?word=${word}`).then(_ => _.json());
}
