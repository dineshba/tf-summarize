// command_executor.go
package parser

// CommandExecutor defines an interface for executing commands.
type CommandExecutor interface {
	CombinedOutput(name string, args ...string) ([]byte, error)
}
