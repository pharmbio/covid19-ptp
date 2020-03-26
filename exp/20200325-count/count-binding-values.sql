SELECT  
    t.pref_name,
    count(cs.canonical_smiles)
FROM compound_structures cs,
    assays a,
    activities ac,
    target_dictionary t, 
    component_synonyms csyn,
    target_components tc
WHERE t.target_type = 'SINGLE PROTEIN' AND
    t.organism = 'Homo sapiens' AND
    a.assay_type = 'B' AND
    tc.tid = t.tid AND
    t.tid=a.tid AND
    a.assay_id = ac.assay_id AND
    ac.molregno = cs.molregno AND
    tc.component_id = csyn.component_id AND
    ac.standard_relation = '=' AND
    (
        (ac.standard_type = 'IC50'   AND ac.standard_units = 'nM' ) OR
        (ac.standard_type = 'Ki'     AND ac.standard_units = 'nM' ) OR
        (ac.standard_type = 'EC50'   AND ac.standard_units = 'nM' ) OR
        (ac.standard_type = 'Log Ki' AND ac.standard_units IS NULL) OR
        (ac.standard_type = 'Kd'     AND ac.standard_units = 'nM' ) OR
        (ac.standard_type = 'Log Ki' AND ac.standard_units IS NULL)
    ) AND (
        csyn.component_synonym = 'F2RL1' OR
        csyn.component_synonym = 'HDAC2' OR
        csyn.component_synonym = 'ATP6AP1' OR
        csyn.component_synonym = 'LOX' OR
        csyn.component_synonym = 'ABCC1' OR
        csyn.component_synonym = 'SIGMAR1' OR
        csyn.component_synonym = 'COMT' OR
        csyn.component_synonym = 'PRKACA' OR
        csyn.component_synonym = 'SIGMAR1' OR
        csyn.component_synonym = 'TMEM97' OR
        csyn.component_synonym = 'PTGES2' OR
        csyn.component_synonym = 'BRD2' OR
        csyn.component_synonym = 'BRD4' OR
        csyn.component_synonym = 'SLC6A15' OR
        csyn.component_synonym = 'IMPDH2' OR
        csyn.component_synonym = 'NDUFs' OR
        csyn.component_synonym = 'MARK2' OR
        csyn.component_synonym = 'MARK3' OR
        csyn.component_synonym = 'GLA' OR
        csyn.component_synonym = 'RIPK1' OR
        csyn.component_synonym = 'CSNK2B' OR
        csyn.component_synonym = 'CSNK2A2' OR
        csyn.component_synonym = 'SLC1A3' OR
        csyn.component_synonym = 'DNMT1' OR
        csyn.component_synonym = 'DCTPP1' OR
        csyn.component_synonym = 'TBK1'
    )
GROUP BY (t.pref_name);
