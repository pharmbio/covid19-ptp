// Workflow written in SciPipe.
// For more information about SciPipe, see: http://scipipe.org
package main

import (
	sp "github.com/scipipe/scipipe"
	spc "github.com/scipipe/scipipe/components"
	"os"
)

func main() {
	args := os.Args[1:]

	// Create a workflow, using 4 cpu cores
	wf := sp.NewWorkflow("count", 4)

	// Set up workflow
	excapeDB := spc.NewFileSource(wf, "excapedb", "../../raw/pubchem.chembl.dataset4publication_inchi_smiles.tsv")
	supplTbl1 := spc.NewFileSource(wf, "suppltbl1", "../../raw/gordonetal.suppl01.tsv")

	// Reproduce the following:
	// > We applied a two step filtering strategy to determine the final list of
	// > reported interactors which relied on two different scoring stringency
	// > cutoffs. In the first step, we chose all protein interactions that
	// > possess a MiST score ≥ 0.7, a SAINTexpress BFDR ≤ 0.05 and an average
	// > spectral count ≥ 2.
	filterTargetsStep1 := wf.NewProc("filter-highspec-targets",
		`cat {i:suppltbl1} \
		| awk -F'\t' '( $4 >= 0.7 && $5 <= 0.05 && $6 >= 2.0 ) { print $3 }' \
		| sort \
		| uniq \
		| sed 's/\"//g' \
		> {o:highspectargets}`)
	filterTargetsStep1.In("suppltbl1").From(supplTbl1.Out())
	filterTargetsStep1.SetOut("highspectargets", "dat/targets.step1.tsv")

	countPerGene := wf.NewProc("count-per-gene", "cat {i:excapedb} | awk -F'\t' '{ c[$9]++ } END { for (key in c) { print key \"\t\" c[key] } }' > {o:genecounts}")
	countPerGene.SetOut("genecounts", "dat/genecounts.tsv")
	countPerGene.In("excapedb").From(excapeDB.Out())

	// Extract targets from table 3 and 4
	supplTable3 := spc.NewFileSource(wf, "suppltbl3", "../../raw/gordonetal.suppl03.tsv")
	supplTable4 := spc.NewFileSource(wf, "suppltbl4", "../../raw/gordonetal.suppl04.tsv")
	extractTbl34 := wf.NewProc("extract-targets-tbl3-4",
		`cat {i:tbl3} {i:tbl4} \
		 | awk -F'\t' '
			( $2 ~ /\/[A-Z]/ ) { a=$2; b=$2; sub(/\/.*/, "", a); sub(/.*\//, "", b); print a; print b }         # For lines with multiple full gene names (like ABC1/DEF4), print both on separate lines
			( $2 ~ /\/[0-9]/ ) { a=$2; b=$2; sub(/\/[0-9]/, "", a); sub(/[0-9]\//, "", b); print a; print b  }  # For lines with genes with multiple numbers (like ABC1/2), print both on separate lines
			( $2 !~ /\// ) { print $2 }                                                                         # For the rest (lines without a slash in col 2), just print normally' \
		 | awk '( $1 ~ /^[A-Z0-9]{3,10}$/ ) # Remove everything that doesnt look like a gene name' \
		 | sort \
		 | uniq \
		 > {o:genes}`)
	extractTbl34.In("tbl3").From(supplTable3.Out())
	extractTbl34.In("tbl4").From(supplTable4.Out())
	extractTbl34.SetOut("genes", "dat/targets.tbl3-4.tsv")

	// Optionally just plot and exit
	if len(args) > 0 && args[0] == "plot" {
		wf.PlotGraphPDF("workflow.dot")
		return
	}

	// Run the workflow
	wf.Run()
}
