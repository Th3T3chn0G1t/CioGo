// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright (C) 2022 Emily "TTG" Banerjee <prs.ttg+cionom@pm.me>

package cionom

import (
	"fmt"
	"os"
)

func Assert(Message string, Condition bool) {
	if Condition {
		fmt.Printf("Error: %s\n", Message)
		os.Exit(1)
	}
}

func AssertError(Error error, Context string) {
	Assert(fmt.Sprintf("%s while %s", Error, Context), Error != nil)
}
