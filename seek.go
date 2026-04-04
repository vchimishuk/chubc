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
	"errors"

	"github.com/vchimishuk/chubby"
	"github.com/vchimishuk/chubby/time"
	"github.com/vchimishuk/opt"
)

type SeekCommand struct {
}

func NewSeekCommand() SeekCommand {
	return SeekCommand{}
}

func (c SeekCommand) Name() string {
	return "seek"
}

func (c SeekCommand) Options() []*opt.Desc {
	return []*opt.Desc{}
}

func (c SeekCommand) Args() (int, int) {
	return 1, 1
}

func (c SeekCommand) Exec(ch *chubby.Chubby, opts opt.Options, args []string) error {
	var st string
	var mod chubby.SeekMode
	if args[0][0] == '-' {
		st = args[0][1:]
		mod = chubby.SeekModeBackward
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

	return ch.Seek(t, mod)
}
