package next

import (
	"github.com/hiro4bbh/sticker/sticker-util/common"
)

// Command is the interface for an abstract LabelNext command.
type Command interface {
	Parse(args []string) ([]string, error)
	Run() error
}

// TestCommand is the interface for an abstract TestLabelNextCommand.
type TestCommand interface {
	Command
}

var newTestCommand = func(opts common.Options) TestCommand {
	panic("unsupported @testNext")
}

// NewTestCommand returns an new TestCommand.
func NewTestCommand(opts common.Options) TestCommand {
	return newTestCommand(opts)
}

// TrainCommand is the interface for an abstract TrainLabelNextCommand.
type TrainCommand interface {
	Command
}

var newTrainCommand = func(opts common.Options) TrainCommand {
	panic("unsupported @trainNext")
}

// NewTrainCommand returns an new TrainCommand.
func NewTrainCommand(opts common.Options) TrainCommand {
	return newTrainCommand(opts)
}
