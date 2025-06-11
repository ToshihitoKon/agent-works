package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Executor struct {
	config *Config
}

func NewExecutor(config *Config) *Executor {
	return &Executor{config: config}
}

func (e *Executor) switchContext(contextName string) error {
	context, exists := e.config.Contexts[contextName]
	if !exists {
		return fmt.Errorf("context '%s' not found", contextName)
	}

	e.config.CurrentContext = contextName

	if activateCmd, exists := context.Commands["activate"]; exists {
		if err := e.executeCommand(activateCmd, context.Variables); err != nil {
			return fmt.Errorf("failed to execute activate command: %w", err)
		}
	}

	return e.config.save()
}

func (e *Executor) executeCommand(command string, variables map[string]string) error {
	expandedCommand := e.expandVariables(command, variables)
	
	cmd := exec.Command("sh", "-c", expandedCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (e *Executor) expandVariables(command string, variables map[string]string) string {
	result := command
	for key, value := range variables {
		placeholder := fmt.Sprintf("${%s}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

func (e *Executor) getCurrentContext() *Context {
	if e.config.CurrentContext == "" {
		return nil
	}
	
	if context, exists := e.config.Contexts[e.config.CurrentContext]; exists {
		return &context
	}
	
	return nil
}

func (e *Executor) listContexts() []Context {
	contexts := make([]Context, 0, len(e.config.Contexts))
	for _, context := range e.config.Contexts {
		contexts = append(contexts, context)
	}
	return contexts
}