package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	list map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	callback, ok := c.list[cmd.name]
	if !ok {
		return fmt.Errorf("Given command does not exists. Tried to find command via name \"%s\"", cmd.name)
	}
	err := callback(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, callback func(*state, command) error) error {
	if len(name) == 0 {
		return fmt.Errorf("Cannot register command with empty name")
	}
	if callback == nil {
		return fmt.Errorf("Cannot register nil callback with given name \"%s\"", name)
	}
	c.list[name] = callback
	return nil
}
