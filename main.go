package main

import (
	"flag"
	"fmt"
	"os"
	
	"gwcoffey/otis/commands/echo"
	"gwcoffey/otis/commands/wordcount"
)

type Command struct {
	Name string
	Alias string
	Action func([]string)
}

var commands []Command

func setupFlags() {
	flag.Bool("help", false, "show usage for a command")
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n\n", os.Args[0])
		fmt.Println("Commands:")
		for _, cmd := range commands {
			fmt.Printf("  %s (%s)\n", cmd.Name, cmd.Alias)
		}
		fmt.Println("\nFlags:")
		
		flag.PrintDefaults()
	}	
	flag.Parse()
}

func loadCommand(name string, alias string, action func([]string)) {
	commands = append(commands, Command{Name: name, Alias: alias, Action: action})
}

func lookupCommand(name string) func([]string) {
	for _, cmd := range commands {
		if cmd.Name == name || (cmd.Alias != "" && cmd.Alias == name) {
			return cmd.Action
		}
	}
	
	return nil
}

func main() {
	setupFlags()
	
	loadCommand("echo", "e", echo.Echo)
	loadCommand("wordcount", "wc", wordcount.WordCount)

	commandFn := lookupCommand(flag.Arg(0))
	if commandFn != nil {
		commandFn(os.Args[2:])
	} else {
		flag.Usage()
		os.Exit(1)
	}
}