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
	case "init":
		return c.initConfig()
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
  init                  Initialize configuration with example contexts
  list, ls              List all contexts
  current               Show current context
  switch, sw <name>     Switch to context
  add                   Add new context (interactive)
  remove, rm <name>     Remove context
  tui                   Start TUI mode
  help                  Show this help

Examples:
  any-context-switcher init
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

func (c *CLI) initConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Configuration file already exists at: %s\n", configPath)
		fmt.Print("Do you want to overwrite it? (y/N): ")
		
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("Initialization cancelled.")
			return nil
		}
	}

	exampleConfig := &Config{
		CurrentContext: "",
		Theme: ColorTheme{
			Title:       "205",
			Selected:    "199", 
			Border:      "168",
			OutputTitle: "212",
		},
		Contexts: map[string]Context{
			"docker": {
				Name:        "docker",
				Label:       "Docker Services",
				Description: "Start/stop Docker containers",
				Status:      "inactive",
				Commands: map[string]string{
					"activate": "docker-compose up -d && echo 'Docker services started'",
				},
				Variables: map[string]string{
					"COMPOSE_FILE": "docker-compose.yml",
					"PROJECT_NAME": "myapp",
				},
			},
			"vpn": {
				Name:        "vpn",
				Label:       "VPN Connection",
				Description: "Connect to company VPN",
				Status:      "inactive",
				Commands: map[string]string{
					"activate": "echo 'Connecting to VPN: ${VPN_SERVER}' && ping -c 1 ${VPN_SERVER}",
				},
				Variables: map[string]string{
					"VPN_SERVER": "vpn.company.com",
					"VPN_CONFIG": "~/.config/vpn/company.conf",
				},
			},
			"database": {
				Name:        "database",
				Label:       "Database Tunnel",
				Description: "SSH tunnel to database server",
				Status:      "inactive",
				Commands: map[string]string{
					"activate": "echo 'Setting up SSH tunnel to ${DB_HOST}:${DB_PORT}' && nc -z ${DB_HOST} ${DB_PORT}",
				},
				Variables: map[string]string{
					"DB_HOST": "database.company.com",
					"DB_PORT": "5432",
					"LOCAL_PORT": "5433",
				},
			},
			"monitoring": {
				Name:        "monitoring",
				Label:       "System Monitoring",
				Description: "Enable system monitoring tools",
				Status:      "active",
				Commands: map[string]string{
					"activate": "echo 'Monitoring enabled: CPU, Memory, Disk' && ps aux | grep -E '(htop|top|iostat)' | head -3",
				},
				Variables: map[string]string{
					"MONITOR_INTERVAL": "5",
					"LOG_PATH": "/var/log/monitoring",
				},
			},
			"proxy": {
				Name:        "proxy",
				Label:       "HTTP Proxy",
				Description: "Route traffic through proxy server",
				Status:      "inactive",
				Commands: map[string]string{
					"activate": "export http_proxy=${PROXY_URL} && export https_proxy=${PROXY_URL} && echo 'Proxy configured: ${PROXY_URL}'",
				},
				Variables: map[string]string{
					"PROXY_URL": "http://proxy.company.com:8080",
					"NO_PROXY": "localhost,127.0.0.1,.local",
				},
			},
		},
	}

	if err := exampleConfig.save(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Configuration initialized at: %s\n", configPath)
	fmt.Println("Example tool contexts created:")
	fmt.Println("  - docker: Docker Services")
	fmt.Println("  - vpn: VPN Connection")
	fmt.Println("  - database: Database Tunnel")
	fmt.Println("  - monitoring: System Monitoring (active)")
	fmt.Println("  - proxy: HTTP Proxy")
	fmt.Println("\nRun 'any-context-switcher list' to see all contexts.")
	fmt.Println("Run 'any-context-switcher tui' to use the interactive interface.")

	return nil
}