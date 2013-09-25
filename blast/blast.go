package blast

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
)

// Blast executes the search query described by blaster. Search results are
// returned from Blast's XML output format mode.
func Blast(blaster Blaster) (*BlastResults, error) {
	args := blaster.CmdArgs()
	args = append(args, "-outfmt", "5")

	cmd := exec.Command(blaster.Executable(), args...)
	cmd.Stdin = blaster.Stdin()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	decoder := xml.NewDecoder(stdout)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("%s\n\nstderr: %s", err, onStderr(stderr))
	}

	results := new(BlastResults)
	if err := decoder.Decode(&results); err != nil {
		return nil, fmt.Errorf("%s\n\nstderr: %s", err, onStderr(stderr))
	}
	return results, nil
}

func onStderr(r io.Reader) string {
	all, _ := ioutil.ReadAll(r)
	return string(all)
}

// BlastResults is the top-level struct for representing XML output of the
// BLAST family of programs. Subsequent XML elements are represented with
// other `Blast*` types.
//
// The types are meant to be comprehensive with respect to NCBI's DTD found
// here: http://www.ncbi.nlm.nih.gov/dtd/NCBI_BlastOutput.dtd.
// Note that the meat is really here:
// http://www.ncbi.nlm.nih.gov/dtd/NCBI_BlastOutput.mod.dtd.
type BlastResults struct {
	XMLName    xml.Name         `xml:"BlastOutput"`
	Program    string           `xml:"BlastOutput_program"`
	Version    string           `xml:"BlastOutput_version"`
	Reference  string           `xml:"BlastOutput_reference"`
	DB         string           `xml:"BlastOutput_db"`
	QueryID    string           `xml:"BlastOutput_query-ID"`
	QueryDef   string           `xml:"BlastOutput_query-def"`
	QueryLen   int              `xml:"BlastOutput_query-len"`
	QuerySeq   string           `xml:"BlastOutput_query-seq"`
	Params     BlastParams      `xml:"BlastOutput_param>Parameters"`
	Iterations []BlastIteration `xml:"BlastOutput_iterations>Iteration"`
}

type BlastParams struct {
	XMLName     xml.Name `xml:"Parameters"`
	Matrix      string   `xml:"Parameters_matrix"`
	Expect      float64  `xml:"Parameters_exect"`
	Include     float64  `xml:"Parameters_include"`
	ScMatch     int      `xml:"Parameters_sc-match"`
	ScMismatch  int      `xml:"Parameters_sc-mismatch"`
	GapOpen     int      `xml:"Parameters_gap-open"`
	GapExtend   int      `xml:"Parameters_gap-extend"`
	Filter      string   `xml:"Parameters_filter"`
	Pattern     string   `xml:"Parameters_pattern"`
	EntrezQuery string   `xml:"Parameters_entrez-query"`
}

type BlastIteration struct {
	XMLName  xml.Name        `xml:"Iteration"`
	Num      int             `xml:"Iteration_iter-num"`
	QueryID  string          `xml:"Iteration_query-ID"`
	QueryDef string          `xml:"Iteration_query-def"`
	QueryLen int             `xml:"Iteration_query-len"`
	Hits     []BlastHit      `xml:"Iteration_hits>Hit"`
	Stats    BlastStatistics `xml:"Iteration_stat>Statistics"`
	Message  string          `xml:"Iteration_message"`
}

type BlastStatistics struct {
	XMLName      xml.Name `xml:"Statistics"`
	NumSequences int      `xml:"Statistics_db-num"`
	Length       int      `xml:"Statistics_db-len"`
	HSPLength    int      `xml:"Statistics_hsp-len"`
	EffSpace     float64  `xml:"Statistics_eff-space"`
	Kappa        float64  `xml:"Statistics_kappa"`
	Lambda       float64  `xml:"Statistics_lambda"`
	Entropy      float64  `xml:"Statistics_entropy"`
}

type BlastHit struct {
	XMLName   xml.Name   `xml:"Hit"`
	Num       int        `xml:"Hit_num"`
	Id        string     `xml:"Hit_id"`
	Def       string     `xml:"Hit_def"`
	Accession string     `xml:"Hit_accession"`
	Length    int        `xml:"Hit_len"`
	Hsps      []BlastHSP `xml:"Hit_hsps>Hsp"`
}

type BlastHSP struct {
	XMLName     xml.Name `xml:"Hsp"`
	Num         int      `xml:"Hsp_num"`
	BitScore    float64  `xml:"Hsp_bit-score"`
	Score       float64  `xml:"Hsp_score"`
	EValue      float64  `xml:"Hsp_evalue"`
	QueryFrom   int      `xml:"Hsp_query-from"`
	QueryTo     int      `xml:"Hsp_query-to"`
	HitFrom     int      `xml:"Hsp_hit-from"`
	HitTo       int      `xml:"Hsp_hit-to"`
	PatternFrom int      `xml:"Hsp_pattern-from"`
	PatternTo   int      `xml:"Hsp_pattern-to"`
	QueryFrame  int      `xml:"Hsp_query-frame"`
	HitFrame    int      `xml:"Hsp_hit-frame"`
	Identity    int      `xml:"Hsp_identity"`
	Positive    int      `xml:"Hsp_positive"`
	Gaps        int      `xml:"Hsp_gaps"`
	AlignLength int      `xml:"Hsp_align-len"`
	Density     int      `xml:"Hsp_density"`
	AlignQuery  string   `xml:"Hsp_qseq"`
	AlignHit    string   `xml:"Hsp_hseq"`
	AlignMiddle string   `xml:"Hsp_midline"`
}
