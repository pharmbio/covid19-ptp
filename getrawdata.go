// Workflow written in SciPipe.
// For more information about SciPipe, see: http://scipipe.org
package main

import (
	"fmt"
	sp "github.com/scipipe/scipipe"
	//spc "github.com/scipipe/scipipe/components"
)

func main() {
	// --------------------------------
	// Create a pipeline runner
	// --------------------------------
	wf := sp.NewWorkflow("get-rawdata", 18)

	dbFileName := "pubchem.chembl.dataset4publication_inchi_smiles.tsv.xz"
	dlExcapeDB := wf.NewProc("dlDB", fmt.Sprintf("wget https://zenodo.org/record/173258/files/%s -O {o:excapexz}", dbFileName))
	dlExcapeDB.SetOut("excapexz", "raw/"+dbFileName)

	unPackDB := wf.NewProc("unPackDB", "xzcat {i:xzfile} > {o:unxzed}")
	unPackDB.SetOut("unxzed", "{i:xzfile|%.xz}")

	// --------------------------------
	// Connect workflow dependency network
	// --------------------------------
	unPackDB.In("xzfile").From(dlExcapeDB.Out("excapexz"))

	// --------------------------------
	// Run the pipeline!
	// --------------------------------
	wf.Run()
}
