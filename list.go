// Copyright 2026 Viacheslav Chimishuk <vchimishuk@yandex.ru>
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
	"path"
	"strconv"
	"strings"

	"github.com/vchimishuk/chubby"
	"github.com/vchimishuk/opt"
)

type ListCommand struct {
}

func NewListCommand() ListCommand {
	return ListCommand{}
}

func (c ListCommand) Name() string {
	return "list"
}

func (c ListCommand) Options() []*opt.Desc {
	return []*opt.Desc{
		{"f", "", opt.ArgString, "FORMAT", "list item format"},
	}
}

func (c ListCommand) Args() (int, int) {
	return 1, 1
}

func (c ListCommand) Exec(ch *chubby.Chubby, opts opt.Options, args []string) error {
	f := opts.StringOr("f", "%f%/")

	entries, err := ch.List(args[0])
	if err != nil {
		return err
	}
	for _, e := range entries {
		vars := map[string]string{}

		if e.IsDir() {
			p, f := path.Split(e.Dir().Path)
			vars["dir"] = "true"
			vars["path"] = p
			vars["file"] = f
		} else {
			p, f := path.Split(e.Track().Path)
			vars["dir"] = "false"
			vars["path"] = p
			vars["file"] = f
			vars["artist"] = e.Track().Artist
			vars["album"] = e.Track().Album
			vars["year"] = strconv.Itoa(e.Track().Year)
			vars["title"] = e.Track().Title
			vars["number"] = strconv.Itoa(e.Track().Number)
			vars["length"] = e.Track().Length.String()
		}

		fmt.Print(format(f, vars))
		fmt.Println()
	}

	return nil
}

func format(fmt string, vars map[string]string) string {
	var spec bool = false
	var s strings.Builder

	for _, r := range fmt {
		if !spec && r == '%' {
			spec = true
		} else if spec {
			switch r {
			case '%':
				s.WriteRune(r)
			case '/':
				if vars["dir"] == "true" {
					s.WriteRune('/')
				}
			case 'A':
				s.WriteString(vars["album"])
			case 'a':
				s.WriteString(vars["artist"])
			case 'f':
				s.WriteString(vars["file"])
			case 'l':
				s.WriteString(vars["length"])
			case 'n':
				s.WriteString(vars["number"])
			case 'p':
				s.WriteString(vars["path"])
			case 't':
				s.WriteString(vars["title"])
			case 'y':
				s.WriteString(vars["year"])
			default:
				s.WriteRune('%')
				s.WriteRune(r)
			}
			spec = false
		} else {
			s.WriteRune(r)
		}
	}

	return s.String()
}
