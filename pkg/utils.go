// Copyright (c) 2017 Dongsu Park <dpark@posteo.net>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package pkg

import (
	"strconv"
	"time"
)

func StringToInt(instr string) int {
	intVal, err := strconv.Atoi(instr)
	if err != nil {
		return 0
	}
	return intVal
}

func ComposeTimeDuration(hourNum, minNum, secNum, msecNum int) time.Duration {
	return (time.Duration(hourNum) * time.Hour) +
		(time.Duration(minNum) * time.Minute) +
		(time.Duration(secNum) * time.Second) +
		(time.Duration(msecNum) * time.Millisecond)
}

func DurationToClockNums(dur time.Duration) (int, int, int, int) {
	eh := int(dur.Hours())
	em := int(dur.Minutes())
	es := int(dur.Seconds())
	en := int(dur.Nanoseconds())

	return eh, (em - (eh * 60)), (es - (em * 60)), (en / 1000 / 1000 % 1000)
}
