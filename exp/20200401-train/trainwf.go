// Workflow written in SciPipe.
// For more information about SciPipe, see: http://scipipe.org
package main

import (
	sp "github.com/scipipe/scipipe"
	spc "github.com/scipipe/scipipe/components"
)

func main() {
	// Create a workflow, using 4 cpu cores
	wf := sp.NewWorkflow("train-models", 18)

	excapeDB := spc.NewFileSource(wf, "excapedb", "../../raw/excapedb.v2.tsv")

	// extractGISA extracts a file with only Gene symbol, id (orig entry), SMILES,
	// and the Activity flag into a .tsv file, for easier subsequent parsing.
	// ATTENTION: The sorting order (Gene, SMILES, Activity) is super important,
	// for the following component, `deduplicateSmiles` to function properly!
	extractGISA := wf.NewProc("extract_gisa", `awk -F "\t" '{ print $9 "\t" $2 "\t" $12 "\t" $4 }' {i:excapedb} | sort -uV -k 1,1 -k 3,3 -k 4,4 > {o:gisa}`)
	extractGISA.SetOut("gisa", "{i:excapedb|%.tsv}.gisa.tsv")
	extractGISA.In("excapedb").From(excapeDB.Out())

	// deduplicateSmiles prints the previous line, unless it has the same values on
	deduplicateSmiles := wf.NewProc("deduplicate_smiles",
		`awk -F "\t" '((( prev1 != $1 ) && ( prev1 != "")) || (( prev3 != $3 ) && ( prev3 != "" ))) && !isconflicting[prev1,prev3] { print prev1 "\t" prev2 "\t" prev3 "\t" prev4 }
				( seen[$1,$3] > 0 ) && ( activity[$1,$3] != $4 ) { isconflicting[$1,$3] = true }
				{ seen[$1,$3]++; activity[$1,$3] = $4; prev1 = $1; prev2 = $2; prev3 = $3; prev4 = $4 }
				END { print }' \
			{i:gisa} > {o:dedup}`)
	deduplicateSmiles.In("gisa").From(extractGISA.Out("gisa"))
	deduplicateSmiles.SetOut("dedup", "{i:gisa|%.tsv}.dedup.tsv")

	// Run the workflow
	wf.Run()
}
