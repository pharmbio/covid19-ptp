// Workflow written in SciPipe.
// For more information about SciPipe, see: http://scipipe.org
package main

import (
	"fmt"
	sp "github.com/scipipe/scipipe"
	spc "github.com/scipipe/scipipe/components"
)

func main() {
	// ------------------------------------------------------------------------
	// Initiate workflow, with 18 cores
	// ------------------------------------------------------------------------
	wf := sp.NewWorkflow("get-rawdata", 4)

	// ------------------------------------------------------------------------
	// Download ExCAPE DB
	// ------------------------------------------------------------------------
	dbFileName := "pubchem.chembl.dataset4publication_inchi_smiles_v2.tsv.xz"
	dlExcapeDB := wf.NewProc("download_excapedb", fmt.Sprintf("wget https://zenodo.org/record/2543724/files/%s -O {o:excapexz}", dbFileName))
	dlExcapeDB.SetOut("excapexz", "raw/excapedb.v2.tsv.xz")

	unPackDB := wf.NewProc("unPackDB", "xzcat {i:xzfile} > {o:unxzed}")
	unPackDB.In("xzfile").From(dlExcapeDB.Out("excapexz"))
	unPackDB.SetOut("unxzed", "{i:xzfile|%.xz}")

	// ------------------------------------------------------------------------
	// Download Supplemental tables from Gordon et. al [1]
	// [1] https://doi.org/10.1101/2020.03.22.002386
	// ------------------------------------------------------------------------
	dlSupplTables := []*sp.Process{}
	for _, tableId := range []int{1, 2, 3, 4} {
		dlSupplTable := wf.NewProc(fmt.Sprintf("dlsuppl%d", tableId),
			fmt.Sprintf("let i={p:tableId}+2 && curl -o {o:suppl} https://www.biorxiv.org/content/biorxiv/early/2020/03/23/2020.03.22.002386/DC$i/embed/media-$i.xlsx?download=true"))
		dlSupplTable.InParam("tableId").FromInt(tableId)
		dlSupplTable.SetOut("suppl", "raw/gordonetal.suppl0{p:tableId}.xlsx")
		dlSupplTables = append(dlSupplTables, dlSupplTable)
	}

	// ------------------------------------------------------------------------
	// Convert file formats: xlsx -> csv
	// ------------------------------------------------------------------------
	xlsx2Csv := wf.NewProc("xlsx2csv", "ssconvert --export-type=Gnumeric_stf:stf_csv {i:xlsx} {o:csv}")
	// Connect all xlsx2CSV process to inport
	for _, p := range dlSupplTables {
		xlsx2Csv.In("xlsx").From(p.Out("suppl"))
	}
	xlsx2Csv.SetOut("csv", "{i:xlsx|%.xlsx}.csv")

	// ------------------------------------------------------------------------
	// Convert file formats: csv -> tsv
	// ------------------------------------------------------------------------
	csv2Tsv := wf.NewProc("csv2tsv", "csvtool -t COMMA -u TAB cat {i:csv} | tr ',' '.' > {o:tsv}")
	csv2Tsv.In("csv").From(xlsx2Csv.Out("csv"))
	csv2Tsv.SetOut("tsv", "{i:csv|%.csv}.tsv")

	// Extract targets from table 3 and 4
	supplTable3 := spc.NewFileSource(wf, "suppltbl3", "raw/gordonetal.suppl03.tsv")
	supplTable4 := spc.NewFileSource(wf, "suppltbl4", "raw/gordonetal.suppl04.tsv")
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
	extractTbl34.SetOut("genes", "raw/gordonetal.targets-tbl3-4.tsv")

	// ------------------------------------------------------------------------
	// Run workflow
	// ------------------------------------------------------------------------
	wf.Run()
}
