package database

import (
	"github.com/liangweijiang/gedis/interfaces/redis"
	"strings"
)

const flagWrite = 0

const (
	flagReadOnly = 1 << iota
	flagSpecial  // command invoked in Exec
)

var cmdTable = make(map[string]*command)

// ExecFunc is interface for command executor
// args don't include cmd line
type ExecFunc func(db *Database, args [][]byte) redis.Reply

// PreFunc analyses command line when queued command to `multi`
// returns related write keys and read keys
type PreFunc func(args [][]byte) (writeKeys []string, readKeys []string)

// CmdLine is alias for [][]byte, represents a command line
type CmdLine = [][]byte

// UndoFunc returns undo logs for the given command line
// execute from head to tail when undo
type UndoFunc func(db *Database, args [][]byte) []CmdLine

type command struct {
	name     string
	executor ExecFunc
	// prepare returns related keys command
	prepare PreFunc
	// undo generates undo-log before command actually executed, in case the command needs to be rolled back
	undo UndoFunc
	// arity means allowed number of cmdArgs, arity < 0 means len(args) >= -arity.
	// for example: the arity of `get` is 2, `mget` is -2
	arity int
	flags int

	// extra is a pointer to a commandExtra struct that holds additional information about the command.
	// This can include context, configuration details, or other auxiliary data needed for command execution.
	extra *commandExtra
}

type commandExtra struct {
	signs    []string
	firstKey int
	lastKey  int
	keyStep  int
}

func registerCommand(name string, executor ExecFunc, prepare PreFunc, rollback UndoFunc, arity int, flags int) *command {
	name = strings.ToLower(name)
	cmd := &command{
		name:     name,
		executor: executor,
		prepare:  prepare,
		undo:     rollback,
		arity:    arity,
		flags:    flags,
	}
	cmdTable[name] = cmd
	return cmd
}

func isReadOnlyCommand(name string) bool {
	name = strings.ToLower(name)
	cmd := cmdTable[name]
	if cmd == nil {
		return false
	}
	return cmd.flags&flagReadOnly > 0
}

func (c *command) validateArity(cmdArgs [][]byte) bool {
	argLen := len(cmdArgs)
	if c.arity >= 0 {
		return argLen == c.arity
	}
	return argLen >= -c.arity
}

func (c *command) attachCommandExtra(signs []string, firstKey int, lastKey int, keyStep int) {
	c.extra = &commandExtra{
		signs:    signs,
		firstKey: firstKey,
		lastKey:  lastKey,
		keyStep:  keyStep,
	}
}
