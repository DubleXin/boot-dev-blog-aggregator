package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/DubleXin/aggregator/internal/config"
	"github.com/DubleXin/aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func checkErr(e error) {
	if e != nil {
		println("Error, ocured:\n\t", e.Error())
		os.Exit(1)
	}
}

func main() {
	config, err := config.ReadConfig()
	checkErr(err)

	s := state{
		cfg: config,
	}
	cmds := commands{
		list: map[string]func(*state, command) error{},
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerResetUsers)
	cmds.register("users", handlerGetUsers)

	db, err := sql.Open("postgres", config.DBUrl)
	checkErr(err)
	s.db = database.New(db)

	cmd, err := readInput()
	checkErr(err)

	err = cmds.run(&s, cmd)
	checkErr(err)

}

func readInput() (command, error) {
	args := os.Args

	if len(args) < 2 {
		return command{}, fmt.Errorf("Failed to read command from args, expected non zero args length, but got:\n\t\t%v", args)
	}

	if len(args) == 2 {
		return command{
			name: args[1],
			args: []string{},
		}, nil
	}
	return command{
		name: args[1],
		args: args[2:],
	}, nil
}
