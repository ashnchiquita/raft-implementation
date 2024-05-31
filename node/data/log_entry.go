package data

type Option func(l LogEntry) LogEntry

type LogEntry struct {
	Term    int    `json:"term"`
	Command string `json:"command"`

	// Optional field (maybe zero Value)
	Value string `json:"value"`
}

// CONSTRUCTOR
func NewLogEntry(term int, command string, options ...Option) *LogEntry {
	l := &LogEntry{
		Term:    term,
		Command: command,
	}

	for _, option := range options {
		option(*l)
	}

	return l
}

func WithValue(value string) Option {
	return func(l LogEntry) LogEntry {
		l.Value = value
		return l
	}
}
