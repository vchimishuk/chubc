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

	"github.com/vchimishuk/chubby"
	"github.com/vchimishuk/opt"
)

type StatusCommand struct {
}

func NewStatusCommand() StatusCommand {
	return StatusCommand{}
}

func (c StatusCommand) Name() string {
	return "status"
}

func (c StatusCommand) Options() []*opt.Desc {
	return []*opt.Desc{}
}

func (c StatusCommand) Args() (int, int) {
	return 0, 0
}

func (c StatusCommand) Exec(ch *chubby.Chubby, opts opt.Options, args []string) error {
	s, err := ch.Status()
	if err != nil {
		return err
	}

	fmt.Printf("State: %s\n", s.State)
	fmt.Printf("Volume: %d\n", s.Volume)
	if s.State != chubby.StateStopped {
		fmt.Printf("Playlist name: %s\n", s.Playlist.Name)
		fmt.Printf("Playlist position: %d\n", s.PlaylistPos+1)
		fmt.Printf("Playlist length: %d\n", s.Playlist.Length)
		fmt.Printf("Playlist duration: %s\n", s.Playlist.Duration)
		fmt.Printf("Track path: %s\n", s.Track.Path)
		fmt.Printf("Track duration: %s\n", s.Track.Length)
		fmt.Printf("Track position: %s\n", s.TrackPos)
		fmt.Printf("Track artist: %s\n", s.Track.Artist)
		fmt.Printf("Track album: %s\n", s.Track.Album)
		fmt.Printf("Track title: %s\n", s.Track.Title)
		fmt.Printf("Track year: %d\n", s.Track.Year)
		fmt.Printf("Track number: %d\n", s.Track.Number)
	}

	return nil
}
