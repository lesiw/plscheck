// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The maprange command applies the lesiw.io/plscheck/maprange
// analysis to the specified packages of Go source code.
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"lesiw.io/plscheck/maprange"
)

func main() { singlechecker.Main(maprange.Analyzer) }
