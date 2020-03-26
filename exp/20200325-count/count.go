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

	// Optionally just plot and exit
	if len(args) > 0 && args[0] == "plot" {
		wf.PlotGraphPDF("workflow.dot")
		return
	}

	// Set up workflow
	excapeDB := spc.NewFileSource(wf, "excapedb", "../../raw/pubchem.chembl.dataset4publication_inchi_smiles.tsv")
	covidData := spc.NewFileSource(wf, "coviddata", "../../raw/coviddata.tsv")

	// Reproduce the following:
	// > We applied a two step filtering strategy to determine the final list of
	// > reported interactors which relied on two different scoring stringency
	// > cutoffs. In the first step, we chose all protein interactions that
	// > possess a MiST score ≥ 0.7, a SAINTexpress BFDR ≤ 0.05 and an average
	// > spectral count ≥ 2.
	filterTargetsStep1 := wf.NewProc("filter-highspec-targets",
		`cat {i:coviddata} \
		| awk -F'\t' '( $4 >= 0.7 && $5 <= 0.05 && $6 >= 2.0 ) { print $3 }' \
		| sort \
		| uniq \
		| sed 's/\"//g' \
		> {o:highspectargets}`)
	filterTargetsStep1.In("coviddata").From(covidData.Out())
	filterTargetsStep1.SetOut("highspectargets", "dat/targets.step1.tsv")

	countPerGene := wf.NewProc("count-per-gene", "cat {i:excapedb} | awk -F'\t' '{ c[$9]++ } END { for (key in c) { print key \"\t\" c[key] } }' > {o:genecounts}")
	countPerGene.SetOut("genecounts", "dat/genecounts.tsv")
	countPerGene.In("excapedb").From(excapeDB.Out())

	// Run the workflow
	wf.Run()
}
