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
	"strconv"

	"github.com/vchimishuk/chubby"
	"github.com/vchimishuk/opt"
)

type VolumeCommand struct {
}

func NewVolumeCommand() VolumeCommand {
	return VolumeCommand{}
}

func (c VolumeCommand) Name() string {
	return "volume"
}

func (c VolumeCommand) Options() []*opt.Desc {
	return []*opt.Desc{}
}

func (c VolumeCommand) Args() (int, int) {
	return 1, 1
}

func (c VolumeCommand) Exec(ch *chubby.Chubby, opts opt.Options, args []string) error {
	var vols string = args[0]
	var vol int
	var mode chubby.VolumeMode = chubby.VolumeModeAbs
	var err error

	if vols[0] == '-' || vols[0] == '+' {
		mode = chubby.VolumeModeRel
		vol, err = strconv.Atoi(vols)
		if err != nil {
			return err
		}
		if vol < -100 || vol > 100 {
			return errors.New("volume out of range")
		}
	} else {
		vol, err = strconv.Atoi(vols)
		if err != nil {
			return err
		}
		if vol < 0 || vol > 100 {
			return errors.New("volume out of range")
		}
	}

	return ch.Volume(vol, mode)
}
