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
	"sort"

	"github.com/vchimishuk/chubby"
	"github.com/vchimishuk/opt"
)

type PlaylistsCommand struct {
}

func NewPlaylistsCommand() PlaylistsCommand {
	return PlaylistsCommand{}
}

func (c PlaylistsCommand) Name() string {
	return "playlists"
}

func (c PlaylistsCommand) Options() []*opt.Desc {
	return []*opt.Desc{}
}

func (c PlaylistsCommand) Args() (int, int) {
	return 0, 0
}

func (c PlaylistsCommand) Exec(ch *chubby.Chubby, opts opt.Options, args []string) error {
	plists, err := ch.Playlists()
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
