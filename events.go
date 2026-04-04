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

type EventsCommand struct {
}

func NewEventsCommand() EventsCommand {
	return EventsCommand{}
}

func (c EventsCommand) Name() string {
	return "events"
}

func (c EventsCommand) Options() []*opt.Desc {
	return []*opt.Desc{}
}

func (c EventsCommand) Args() (int, int) {
	return 0, 0
}

func (c EventsCommand) Exec(ch *chubby.Chubby, opts opt.Options, args []string) error {
	chn, err := ch.Events(true)
	if err != nil {
		return err
	}

	for {
		e := <-chn
		if e == nil {
			return nil
		}

		fmt.Printf("%s %s\n", e.Event(), e.Serialize())
	}
}
