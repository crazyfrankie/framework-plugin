package logging

type Level int

const (
	LevelDebug Level = -1
	LevelInfo  Level = 0
	LevelWarn  Level = 1
	LevelError Level = 2
)
