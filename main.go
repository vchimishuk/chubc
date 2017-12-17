// Copyright 2017 Viacheslav Chimishuk <vchimishuk@yandex.ru>
//
// This file is part of chubc.
//
// Chub is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Chub is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Chub. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/vchimishuk/chubby"
	"github.com/vchimishuk/opt"
)

func fatal(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}

func printUsage(opts []*opt.Desc) {
	fmt.Printf("Usage: %s [OPTIONS] COMMAND\n", os.Args[0])
	fmt.Printf("Simple Chub noninteractive client.\n")
	fmt.Printf("\n")
	fmt.Printf("Commands:\n")
	fmt.Printf("  create-playlist  create playlist\n")
	fmt.Printf("  delete-playlist  delete playlist\n")
	fmt.Printf("  help             show this help\n")
	fmt.Printf("  kill             kill server\n")
	fmt.Printf("  list             list directory contents\n")
	fmt.Printf("  next             play next track\n")
	fmt.Printf("  pause            toggle pause state\n")
	fmt.Printf("  ping             ping server\n")
	fmt.Printf("  play             play path\n")
	fmt.Printf("  playlists        list playlists\n")
	fmt.Printf("  prev             play previous track\n")
	fmt.Printf("  rename-playlist  rename playlist\n")
	fmt.Printf("  stop             stop playback\n")
	fmt.Printf("\n")
	fmt.Printf("Options:\n")
	fmt.Printf("%s", opt.Usage(opts))
}

func main() {
	optDescs := []*opt.Desc{
		{"h", "host", opt.ArgString, "HOST",
			"server host name"},
		{"", "help", opt.ArgNone, "",
			"display this help"},
		{"p", "port", opt.ArgInt, "PORT",
			"server port"}}

	opts, args, err := opt.Parse(os.Args[1:], optDescs)
	if err != nil {
		fatal("invalid parameters: %s", err)
	}

	if opts.Bool("help") {
		printUsage(optDescs)
		os.Exit(0)
	} else if args[0] == "help" {
		err = checkArgs(args, 1)
		if err != nil {
			fatal("%s", err)
		}
		printUsage(optDescs)
		os.Exit(0)
	}

	if len(args) == 0 {
		fatal("missing command parameter")
	}

	defaultHost, ok := os.LookupEnv("CHUBC_HOST")
	if !ok {
		defaultHost = "localhost"
	}
	host := opts.StringOr("host", defaultHost)
	defaultPortStr, ok := os.LookupEnv("CHUBC_PORT")
	if !ok {
		defaultPortStr = "5115"
	}
	defaultPort, err := strconv.Atoi(defaultPortStr)
	if err != nil {
		fatal("invalid port number: %s", defaultPortStr)
	}
	port := opts.IntOr("port", defaultPort)
	c := &chubby.Chubby{}

	err = c.Connect(host, port)
	if err != nil {
		fatal("unnable to connect to remote host: %s", err)
	}
	defer c.Close()

	switch args[0] {
	case "delete-playlist":
		err = oneArgCmd(c.DeletePlaylist, args[1:])
	case "create-playlist":
		err = oneArgCmd(c.CreatePlaylist, args[1:])
	case "kill":
		err = noArgsCmd(c.Kill, args[1:])
	case "list":
		err = cmdList(c, args[1:])
	case "next":
		err = noArgsCmd(c.Next, args[1:])
	case "pause":
		err = noArgsCmd(c.Pause, args[1:])
	case "ping":
		err = noArgsCmd(c.Ping, args[1:])
	case "play":
		err = oneArgCmd(c.Play, args[1:])
	case "playlists":
		err = cmdPlaylists(c, args[1:])
	case "prev":
		err = noArgsCmd(c.Prev, args[1:])
	case "rename-playlist":
		err = cmdRenamePlaylist(c, args[1:])
	case "stop":
		err = noArgsCmd(c.Stop, args[1:])
	default:
		err = fmt.Errorf("'%s' is not a valid command", args[0])
	}
	if err != nil {
		fatal("%s", err)
	}
}

func checkArgs(args []string, expected int) error {
	if len(args) < expected {
		return errors.New("not enough arguments")
	} else if len(args) > expected {
		return errors.New("too many arguments")
	} else {
		return nil
	}
}

func noArgsCmd(cmd func() error, args []string) error {
	err := checkArgs(args, 0)
	if err != nil {
		return err
	}

	return cmd()
}

func oneArgCmd(cmd func(string) error, args []string) error {
	err := checkArgs(args, 1)
	if err != nil {
		return err
	}

	return cmd(args[0])
}

func cmdList(c *chubby.Chubby, args []string) error {
	err := checkArgs(args, 1)
	if err != nil {
		return err
	}

	entries, err := c.List(args[0])
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			_, name := path.Split(e.Dir().Path)
			fmt.Printf("%s/\n", name)
		} else {
			_, name := path.Split(e.Track().Path)
			fmt.Printf("%s\n", name)
		}
	}

	return nil
}

func cmdPlaylists(c *chubby.Chubby, args []string) error {
	err := checkArgs(args, 0)
	if err != nil {
		return err
	}

	plists, err := c.Playlists()
	if err != nil {
		return err
	}

	sort.Slice(plists, func(i, j int) bool {
		return plists[i].Name < plists[j].Name
	})
	for _, pl := range plists {
		fmt.Printf("%s\n", pl.Name)
	}

	return nil
}

func cmdRenamePlaylist(c *chubby.Chubby, args []string) error {
	err := checkArgs(args, 2)
	if err != nil {
		return err
	}

	return c.RenamePlaylist(args[0], args[1])
}
