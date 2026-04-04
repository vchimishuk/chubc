// Copyright 2017-2024 Viacheslav Chimishuk <vchimishuk@yandex.ru>
//
// This file is part of chubc.
//
// Chubc is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Chubc is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Chubc. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/vchimishuk/chubby"
	"github.com/vchimishuk/opt"
)

var Commands []Command = []Command{
	NewCreatePlaylistCommand(),
	NewDeletePlaylistCommand(),
	NewEventsCommand(),
	NewKillCommand(),
	NewListCommand(),
	NewNextCommand(),
	NewPauseCommand(),
	NewPingCommand(),
	NewPlayCommand(),
	NewPlaylistsCommand(),
	NewPrevCommand(),
	NewRenamePlaylistCommand(),
	NewSeekCommand(),
	NewStatusCommand(),
	NewStopCommand(),
	NewVolumeCommand(),
}

func command(name string) Command {
	i := slices.IndexFunc(Commands, func(c Command) bool {
		return c.Name() == name
	})
	if i == -1 {
		return nil
	}

	return Commands[i]
}

func prog() string {
	return os.Args[0]
}

func fatal(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s: %s\n", prog(), msg)
	os.Exit(1)
}

func printUsage(opts []*opt.Desc) {
	fmt.Printf("Usage: %s [OPTIONS] COMMAND [ARG]...\n", prog())
	fmt.Printf("Simple Chub non-interactive client.\n")
	fmt.Printf("\n")
	fmt.Printf("Options:\n")
	fmt.Printf("%s", opt.Usage(opts))
	fmt.Printf("\n")
	fmt.Printf("Commands:\n")
	fmt.Printf("  create-playlist  ")
	fmt.Printf("Create new playlist.\n")
	fmt.Printf("  delete-playlist  ")
	fmt.Printf("Delete existing playlist.\n")
	fmt.Printf("  events           ")
	fmt.Printf("Listen for events and print them to stdout.\n")
	fmt.Printf("  help             ")
	fmt.Printf("Show this help.\n")
	fmt.Printf("  kill             ")
	fmt.Printf("Kill Chub server.\n")
	fmt.Printf("  list             ")
	fmt.Printf("List VFS directory contents.\n")
	fmt.Printf("  next             ")
	fmt.Printf("Move playback to the next track in the playlist.\n")
	fmt.Printf("  pause            ")
	fmt.Printf("Toggle pause.\n")
	fmt.Printf("  ping             ")
	fmt.Printf("Ping Chub server.\n")
	fmt.Printf("  play             ")
	fmt.Printf("Start playing track or directory.\n")
	fmt.Printf("  playlists        ")
	fmt.Printf("Print list of existing playlists.\n")
	fmt.Printf("  prev             ")
	fmt.Printf("Move playback to the previous track in the playlist.\n")
	fmt.Printf("  rename-playlist  ")
	fmt.Printf("Rename playlist.\n")
	fmt.Printf("  seek             ")
	fmt.Printf("Seek playback time.\n")
	fmt.Printf("  status           ")
	fmt.Printf("Print Chub player current status.\n")
	fmt.Printf("  stop             ")
	fmt.Printf("Stop playback.\n")
	fmt.Printf("  volume           ")
	fmt.Printf("Set playback volume.\n")
}

func printCommandUsage(cmd Command) {
	opts := ""
	if len(cmd.Options()) != 0 {
		opts = " [OPTIONS]"
	}
	args := ""
	n, _ := cmd.Args()
	if n > 0 {
		args = " [ARG]..."
	}
	fmt.Printf("Usage: %s %s%s%s\n", prog(), cmd.Name(), opts, args)
	if opts != "" {
		fmt.Printf("\n")
		fmt.Printf("Options:\n")
		fmt.Printf("%s", opt.Usage(cmd.Options()))
	}
}

func main() {
	optDescs := []*opt.Desc{
		{"h", "host", opt.ArgString, "HOST",
			"server host name"},
		{"", "help", opt.ArgNone, "",
			"display this help"},
		{"p", "port", opt.ArgInt, "PORT",
			"server port"}}

	opts, args, err := opt.Parse(os.Args[1:], optDescs, true)
	if err != nil {
		fatal("invalid parameters: %s", err)
	}

	help := opts.Has("help")
	if help {
		printUsage(optDescs)
		os.Exit(0)
	}
	if len(args) == 0 || args[0] == "help" {
		printUsage(optDescs)
		os.Exit(0)
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

	cmd := command(args[0])
	if cmd == nil {
		printUsage(optDescs)
		os.Exit(0)
	}

	opts, args, err = opt.Parse(args[1:], cmd.Options(), false)
	mina, maxa := cmd.Args()
	if err != nil || len(args) < mina || len(args) > maxa {
		printCommandUsage(cmd)
		os.Exit(1)
	}

	err = cmd.Exec(c, opts, args)
	if err != nil {
		fatal("%s", err)
	}

	os.Exit(0)
}
