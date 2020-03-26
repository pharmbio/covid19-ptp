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

	filterHighSpecTargets := wf.NewProc("filter-highspec-targets", "cat {i:coviddata} | awk -F'\t' '( $4 >= 0.99 ) { print $3 }' | sort | uniq | sed 's/\"//g' > {o:highspectargets}")
	filterHighSpecTargets.In("coviddata").From(covidData.Out())
	filterHighSpecTargets.SetOut("highspectargets", "dat/highspectargets.tsv")

	countPerGene := wf.NewProc("count-per-gene", "cat {i:excapedb} | awk -F'\t' '{ c[$9]++ } END { for (key in c) { print key \"\t\" c[key] } }' > {o:genecounts}")
	countPerGene.SetOut("genecounts", "dat/genecounts.tsv")
	countPerGene.In("excapedb").From(excapeDB.Out())

	// Run the workflow
	wf.Run()
}
