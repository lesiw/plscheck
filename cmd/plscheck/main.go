package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"lesiw.io/plscheck/deprecated"
	"lesiw.io/plscheck/embeddirective"
	"lesiw.io/plscheck/fillreturns"
	"lesiw.io/plscheck/infertypeargs"
	"lesiw.io/plscheck/maprange"
	"lesiw.io/plscheck/modernize"
	"lesiw.io/plscheck/nonewvars"
	"lesiw.io/plscheck/noresultvalues"
	"lesiw.io/plscheck/recursiveiter"
	"lesiw.io/plscheck/simplifycompositelit"
	"lesiw.io/plscheck/simplifyrange"
	"lesiw.io/plscheck/simplifyslice"
	"lesiw.io/plscheck/unusedfunc"
	"lesiw.io/plscheck/unusedparams"
	"lesiw.io/plscheck/unusedvariable"
	"lesiw.io/plscheck/yield"
)

func main() {
	multichecker.Main(
		deprecated.Analyzer,
		embeddirective.Analyzer,
		fillreturns.Analyzer,
		infertypeargs.Analyzer,
		maprange.Analyzer,
		modernize.Analyzer,
		nonewvars.Analyzer,
		noresultvalues.Analyzer,
		recursiveiter.Analyzer,
		simplifycompositelit.Analyzer,
		simplifyrange.Analyzer,
		simplifyslice.Analyzer,
		unusedfunc.Analyzer,
		unusedparams.Analyzer,
		unusedvariable.Analyzer,
		yield.Analyzer,
	)
}
