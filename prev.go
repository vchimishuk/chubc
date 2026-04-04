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
	"github.com/vchimishuk/chubby"
	"github.com/vchimishuk/opt"
)

type PrevCommand struct {
}

func NewPrevCommand() PrevCommand {
	return PrevCommand{}
}

func (c PrevCommand) Name() string {
	return "prev"
}

func (c PrevCommand) Options() []*opt.Desc {
	return []*opt.Desc{}
}

func (c PrevCommand) Args() (int, int) {
	return 0, 0
}

func (c PrevCommand) Exec(ch *chubby.Chubby, opts opt.Options, args []string) error {
	return ch.Prev()
}
