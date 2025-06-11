package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

type CLI struct {
	executor *Executor
}

func NewCLI(executor *Executor) *CLI {
	return &CLI{executor: executor}
}

func (c *CLI) Run(args []string) error {
	if len(args) < 2 {
		c.showUsage()
		return nil
	}

	switch args[1] {
	case "list", "ls":
		return c.listContexts()
	case "current":
		return c.showCurrent()
	case "switch", "sw":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s switch <context-name>\n", args[0])
			return fmt.Errorf("context name required")
		}
		return c.switchContext(args[2])
	case "add":
		return c.addContext(args[2:])
	case "remove", "rm":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "Usage: %s remove <context-name>\n", args[0])
			return fmt.Errorf("context name required")
		}
		return c.removeContext(args[2])
	case "tui":
		return c.startTUI()
	case "help", "-h", "--help":
		c.showUsage()
		return nil
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[1])
		c.showUsage()
		return fmt.Errorf("unknown command")
	}
}

func (c *CLI) showUsage() {
	fmt.Printf(`Usage: any-context-switcher <command> [arguments]

Commands:
  list, ls              List all contexts
  current               Show current context
  switch, sw <name>     Switch to context
  add                   Add new context (interactive)
  remove, rm <name>     Remove context
  tui                   Start TUI mode
  help                  Show this help

Examples:
  any-context-switcher list
  any-context-switcher switch development
  any-context-switcher tui
`)
}

func (c *CLI) listContexts() error {
	contexts := c.executor.listContexts()
	
	if len(contexts) == 0 {
		fmt.Println("No contexts configured")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tLABEL\tSTATUS\tDESCRIPTION")
	fmt.Fprintln(w, "----\t-----\t------\t-----------")

	current := c.executor.getCurrentContext()
	for _, context := range contexts {
		marker := " "
		if current != nil && current.Name == context.Name {
			marker = "*"
		}
		fmt.Fprintf(w, "%s%s\t%s\t%s\t%s\n", 
			marker, context.Name, context.Label, context.Status, context.Description)
	}
	
	return w.Flush()
}

func (c *CLI) showCurrent() error {
	current := c.executor.getCurrentContext()
	if current == nil {
		fmt.Println("No context is currently active")
		return nil
	}

	fmt.Printf("Current context: %s\n", current.Name)
	fmt.Printf("Label: %s\n", current.Label)
	fmt.Printf("Status: %s\n", current.Status)
	if current.Description != "" {
		fmt.Printf("Description: %s\n", current.Description)
	}
	
	return nil
}

func (c *CLI) switchContext(name string) error {
	if err := c.executor.switchContext(name); err != nil {
		return err
	}
	
	fmt.Printf("Switched to context: %s\n", name)
	return nil
}

func (c *CLI) addContext(args []string) error {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	name := fs.String("name", "", "Context name (required)")
	label := fs.String("label", "", "Context label (required)")
	description := fs.String("description", "", "Context description")
	status := fs.String("status", "inactive", "Context status")
	
	fs.Parse(args)
	
	if *name == "" || *label == "" {
		fmt.Fprintf(os.Stderr, "Usage: any-context-switcher add -name <name> -label <label> [-description <desc>] [-status <status>]\n")
		return fmt.Errorf("name and label are required")
	}

	context := Context{
		Name:        *name,
		Label:       *label,
		Description: *description,
		Status:      *status,
		Commands:    make(map[string]string),
		Variables:   make(map[string]string),
	}

	c.executor.config.Contexts[*name] = context
	if err := c.executor.config.save(); err != nil {
		return err
	}

	fmt.Printf("Added context: %s\n", *name)
	return nil
}

func (c *CLI) removeContext(name string) error {
	if _, exists := c.executor.config.Contexts[name]; !exists {
		return fmt.Errorf("context '%s' not found", name)
	}

	delete(c.executor.config.Contexts, name)
	
	if c.executor.config.CurrentContext == name {
		c.executor.config.CurrentContext = ""
	}

	if err := c.executor.config.save(); err != nil {
		return err
	}

	fmt.Printf("Removed context: %s\n", name)
	return nil
}

func (c *CLI) startTUI() error {
	tui := NewTUI(c.executor)
	return tui.Run()
}