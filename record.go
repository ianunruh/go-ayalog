package ayalog

type Level uint8

const (
	// ErrorLevel designates very serious errors.
	ErrorLevel Level = iota + 1

	// WarnLevel designates hazardous situations.
	WarnLevel

	// InfoLevel designates useful information.
	InfoLevel

	// DebugLevel designates lower priority information.
	DebugLevel

	// TraceLevel designates very low priority, often extremely verbose, information.
	TraceLevel
)

type Record struct {
	Target string
	Module string

	Level Level

	// File is the name of the file where the log was emitted.
	File string

	// Line is the line number of the file where the log was emitted.
	Line uint32

	Message string
}
