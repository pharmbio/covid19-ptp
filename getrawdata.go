// Workflow written in SciPipe.
// For more information about SciPipe, see: http://scipipe.org
package main

import (
	sp "github.com/scipipe/scipipe"
	//spc "github.com/scipipe/scipipe/components"
	"fmt"
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
	unPackDB.In("xzfile").From(dlExcapeDB.Out("excapexz"))
	unPackDB.SetOut("unxzed", "{i:xzfile|%.xz}")

	dlSupplTables := []*sp.Process{}
	for _, tableId := range []int{1, 2, 3, 4} {
		dlSupplTable := wf.NewProc(fmt.Sprintf("dlsuppl%d", tableId),
			fmt.Sprintf("let i={p:tableId}+2 && curl -o {o:suppl} https://www.biorxiv.org/content/biorxiv/early/2020/03/23/2020.03.22.002386/DC$i/embed/media-$i.xlsx?download=true"))
		dlSupplTable.InParam("tableId").FromInt(tableId)
		dlSupplTable.SetOut("suppl", "raw/gordonetal.suppl0{p:tableId}.xlsx")

		dlSupplTables = append(dlSupplTables, dlSupplTable)
	}

	xlsx2Csv := wf.NewProc("xlsx2csv", "ssconvert --export-type=Gnumeric_stf:stf_csv {i:xlsx} {o:csv}")
	// Connect all xlsx2CSV process to inport
	for _, p := range dlSupplTables {
		xlsx2Csv.In("xlsx").From(p.Out("suppl"))
	}
	xlsx2Csv.SetOut("csv", "{i:xlsx|%.xlsx}.csv")

	csv2Tsv := wf.NewProc("csv2tsv", "cat {i:csv} | sed 's/,/\t/g' > {o:tsv}")
	csv2Tsv.In("csv").From(xlsx2Csv.Out("csv"))
	csv2Tsv.SetOut("tsv", "{i:csv|%.csv}.tsv")

	wf.Run()
}
