package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"syscall"
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

	if runCmd, exists := context.Commands["run"]; exists {
		if err := e.executeCommand(runCmd, context.Variables); err != nil {
			return fmt.Errorf("failed to execute run command: %w", err)
		}
	}

	e.config.CurrentContext = contextName
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

func (e *Executor) executeCommandWithOutput(command string, variables map[string]string) (string, error) {
	expandedCommand := e.expandVariables(command, variables)
	
	cmd := exec.Command("sh", "-c", expandedCommand)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	
	// Get exit status code
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
	}
	
	// Build detailed output
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Command: %s\n", expandedCommand))
	result.WriteString(fmt.Sprintf("Exit Code: %d\n", exitCode))
	result.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	stdoutStr := stdout.String()
	stderrStr := stderr.String()
	
	if stdoutStr != "" {
		result.WriteString("STDOUT:\n")
		result.WriteString(stdoutStr)
		if !strings.HasSuffix(stdoutStr, "\n") {
			result.WriteString("\n")
		}
	}
	
	if stderrStr != "" {
		if stdoutStr != "" {
			result.WriteString("\n")
		}
		result.WriteString("STDERR:\n")
		result.WriteString(stderrStr)
		if !strings.HasSuffix(stderrStr, "\n") {
			result.WriteString("\n")
		}
	}
	
	if stdoutStr == "" && stderrStr == "" {
		result.WriteString("(no output)")
	}
	
	return result.String(), err
}

func (e *Executor) executeJobWithOutput(command string, variables map[string]string) (string, int, error) {
	expandedCommand := e.expandVariables(command, variables)
	
	cmd := exec.Command("sh", "-c", expandedCommand)
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	
	// Get exit status code
	exitCode := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
	}
	
	// Build detailed output
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Command: %s\n", expandedCommand))
	result.WriteString(fmt.Sprintf("Exit Code: %d\n", exitCode))
	result.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	
	stdoutStr := stdout.String()
	stderrStr := stderr.String()
	
	if stdoutStr != "" {
		result.WriteString("STDOUT:\n")
		result.WriteString(stdoutStr)
		if !strings.HasSuffix(stdoutStr, "\n") {
			result.WriteString("\n")
		}
	}
	
	if stderrStr != "" {
		if stdoutStr != "" {
			result.WriteString("\n")
		}
		result.WriteString("STDERR:\n")
		result.WriteString(stderrStr)
		if !strings.HasSuffix(stderrStr, "\n") {
			result.WriteString("\n")
		}
	}
	
	if stdoutStr == "" && stderrStr == "" {
		result.WriteString("(no output)")
	}
	
	return result.String(), exitCode, err
}

func (e *Executor) listContexts() []Context {
	var names []string
	for name := range e.config.Contexts {
		names = append(names, name)
	}
	
	sort.Strings(names)
	
	contexts := make([]Context, 0, len(names))
	for _, name := range names {
		contexts = append(contexts, e.config.Contexts[name])
	}
	
	return contexts
}