package streams

// IteratorType controls where to start in a stream.
type IteratorType string

const (
	IteratorAtSequence    IteratorType = "AT_SEQUENCE_NUMBER"
	IteratorAfterSequence IteratorType = "AFTER_SEQUENCE_NUMBER"
	IteratorTrimHorizon   IteratorType = "TRIM_HORIZON"
	IteratorLatest        IteratorType = "LATEST"
)
