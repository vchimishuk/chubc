// Copyright 2017 Viacheslav Chimishuk <vchimishuk@yandex.ru>
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
	"errors"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/vchimishuk/chubby"
	"github.com/vchimishuk/chubby/time"
	"github.com/vchimishuk/opt"
)

func fatal(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], msg)
	os.Exit(1)
}

func printUsage(opts []*opt.Desc) {
	fmt.Printf("Usage: %s [OPTIONS] COMMAND [ARG]...\n", os.Args[0])
	fmt.Printf("Simple Chub noninteractive client.\n")
	fmt.Printf("\n")
	fmt.Printf("Options:\n")
	fmt.Printf("%s", opt.Usage(opts))
	fmt.Printf("\n")
	fmt.Printf("Commands:\n")
	fmt.Printf("  create-playlist NAME\n")
	fmt.Printf("    Create playlist with the name specified by NAME parameter.\n")
	fmt.Printf("\n")
	fmt.Printf("  delete-playlist NAME       delete playlist\n")
	fmt.Printf("    Delete existing playlist with the name specified by NAME parameter.\n")
	fmt.Printf("\n")
	fmt.Printf("  help\n")
	fmt.Printf("    Show this help.\n")
	fmt.Printf("\n")
	fmt.Printf("  kill\n")
	fmt.Printf("    Kill Chub server.\n")
	fmt.Printf("\n")
	fmt.Printf("  list PATH\n")
	fmt.Printf("    List VFS directory's contents.\n")
	fmt.Printf("\n")
	fmt.Printf("  next\n")
	fmt.Printf("    Move playback to the next track in the current playlist.\n")
	fmt.Printf("\n")
	fmt.Printf("  pause\n")
	fmt.Printf("    Toggle pause: pause if currently playin or resume pause if any.\n")
	fmt.Printf("\n")
	fmt.Printf("  ping\n")
	fmt.Printf("    Ping server. Ping does nothing just verifies that server can be connected\n")
	fmt.Printf("    and accepts requests successfuly.\n")
	fmt.Printf("\n")
	fmt.Printf("  play PATH\n")
	fmt.Printf("    Start playing track or directory specified by VFS PATH parameter.\n")
	fmt.Printf("\n")
	fmt.Printf("  playlists\n")
	fmt.Printf("    Print list of existing playlists.\n")
	fmt.Printf("\n")
	fmt.Printf("  prev\n")
	fmt.Printf("    Move playback to the previous track in the current playlist.\n")
	fmt.Printf("\n")
	fmt.Printf("  rename-playlist FROM TO\n")
	fmt.Printf("    Rename playlist specified by FROM parameter to new name TO.\n")
	fmt.Printf("\n")
	fmt.Printf("  seek [-|+]TIME\n")
	fmt.Printf("    Seek playback time. TIME parameter must be provided using common\n")
	fmt.Printf("    time format: [[HH:]MM:]SS. In this case TIME is used as an absolute\n")
	fmt.Printf("    track time offset. If - prefix is specified TIME parameter is treated\n")
	fmt.Printf("    as an interval for rewind. + prefix enables fast-forward mode instead.\n")
	fmt.Printf("\n")
	fmt.Printf("  status\n")
	fmt.Printf("    Print player's current status like current track, time, etc.\n")
	fmt.Printf("\n")
	fmt.Printf("  stop\n")
	fmt.Printf("    Stop playback.\n")
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

	help := opts.Bool("help")
	if help {
		printUsage(optDescs)
		os.Exit(0)
	}

	err = checkArgs(args, 1)
	if err != nil || args[0] == "help" {
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
	c := &chubby.CmdClient{}

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
	case "seek":
		err = cmdSeek(c, args[1:])
	case "status":
		err = cmdStatus(c, args[1:])
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

func cmdList(c *chubby.CmdClient, args []string) error {
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

func cmdPlaylists(c *chubby.CmdClient, args []string) error {
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

func cmdRenamePlaylist(c *chubby.CmdClient, args []string) error {
	err := checkArgs(args, 2)
	if err != nil {
		return err
	}

	return c.RenamePlaylist(args[0], args[1])
}

func cmdSeek(c *chubby.CmdClient, args []string) error {
	err := checkArgs(args, 1)
	if err != nil {
		return err
	}

	var st string
	var mod chubby.SeekMode
	if args[0][0] == '-' {
		st = args[0][1:]
		mod = chubby.SeekModeRewind
	} else if args[0][0] == '+' {
		st = args[0][1:]
		mod = chubby.SeekModeForward
	} else {
		st = args[0]
		mod = chubby.SeekModeAbs
	}

	t, err := time.Parse(st)
	if err != nil {
		return errors.New("invalid time format")
	}

	return c.Seek(t, mod)
}

func cmdStatus(c *chubby.CmdClient, args []string) error {
	err := checkArgs(args, 0)
	if err != nil {
		return err
	}

	s, err := c.Status()
	if err != nil {
		return err
	}

	if s.State == chubby.StateStopped {
		fmt.Printf("State: stopped\n")
	} else {
		if s.State == chubby.StatePlaying {
			fmt.Printf("State: playing\n")
		} else {
			fmt.Printf("State: paused\n")
		}
		fmt.Printf("Playlist name: %s\n", s.Playlist.Name)
		fmt.Printf("Playlist position: %d\n", s.PlaylistPos+1)
		fmt.Printf("Track path: %s\n", s.Track.Path)
		fmt.Printf("Track length: %s\n", s.TrackLen)
		fmt.Printf("Track position: %s\n", s.TrackPos)
	}

	return nil
}
