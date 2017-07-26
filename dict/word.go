package dict

type PartOfSpeech string

const (
	Noun      PartOfSpeech = "Noun"
	Adjective PartOfSpeech = "Adjective"
	Verb      PartOfSpeech = "Verb"
)

type Word struct {
	Base         string
	PartOfSpeech PartOfSpeech
	Forms        []string
	Definitions  []string
}
