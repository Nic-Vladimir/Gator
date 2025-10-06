package main

import (
	"errors"
	"fmt"
)

type command struct {
	Name string
	Args []string
}

type commandInfo struct {
	handler     func(*state, command) error
	description string
}

type commands struct {
	registeredCommands map[string]commandInfo
	order              []string
}

func (c *commands) register(name string, desc string, f func(*state, command) error) {
	if c.registeredCommands == nil {
		c.registeredCommands = make(map[string]commandInfo)
	}
	c.registeredCommands[name] = commandInfo{handler: f, description: desc}
	c.order = append(c.order, name)
}

func (c *commands) list() []string {
	return c.order
}

func (c *commands) Run(s *state, cmd command) error {
	info, ok := c.registeredCommands[cmd.Name]
	if !ok {
		return errors.New("Unknown command")
	}
	return info.handler(s, cmd)
}

func handlerHelp(s *state, cmd command) error {
	fmt.Println("Available commands:")
	for _, name := range s.registeredCommands.list() {
		info := s.registeredCommands.registeredCommands[name]
		fmt.Printf(" %-12s: %s\n", name, info.description)
	}
	return nil
}
