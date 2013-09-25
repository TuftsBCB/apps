package blast

import (
	"fmt"
	"strings"

	"github.com/TuftsBCB/seq"
)

// ExampleBlast demonstrates a very simple protein BLAST search. Note that
// you'll need to change `dbPath` to your own local BLAST database. The one
// I used in the example is a BLAST database containing all of the protein
// sequences from each strain of yeast from http://www.yeastgenome.org.
func ExampleBlast() {
	dbPath := "/home/andrew/research/repeats/data/blast/amino"
	sequence := seq.Sequence{
		Name: "YAL001C",
		Residues: []seq.Residue(`
	MVLTIYPDELVQIVSDKIASNKGKITLNQLWDISGKYFDLSDKKVKQFVLSCVILKKDIE
	VYCDGAITTKNVTDIIGDANHSYSVGITEDSLWTLLTGYTKKESTIGNSAFELLLEVAKS
	GEKGINTMDLAQVTGQDPRSVTGRIKKINHLLTSSQLIYKGHVVKQLKLKKFSHDGVDSN
	PYINIRDHLATIVEVVKRSKNGIRQIIDLKRELKFDKEKRLSKAFIAAIAWLDEKEYLKK
	VLVVSPKNPAIKIRCVKYVKDIPDSKGSPSFEYDSNSADEDSVSDSKAAFEDEDLVEGLD
	NFNATDLLQNQGLVMEEKEDAVKNEVLLNRFYPLQNQTYDIADKSGLKGISTMDVVNRIT
	GKEFQRAFTKSSEYYLESVDKQKENTGGYRLFRIYDFEGKKKFFRLFTAQNFQKLTNAED
	EISVPKGFDELGKSRTDLKTLNEDNFVALNNTVRFTTDSDGQDIFFWHGELKIPPNSKKT
	PNKNKRKRQVKNSTNASVAGNISNPKRIKLEQHVSTAQEPKSAEDSPSSNGGTVVKGKVV
	NFGGFSARSLRSLQRQRAILKVMNTIGGVAYLREQFYESVSKYMGSTTTLDKKTVRGDVD
	LMVESEKLGARTEPVSGRKIIFLPTVGEDAIQRYILKEKDSKKATFTDVIHDTEIYFFDQ
	TEKNRFHRGKKSVERIRKFQNRQKNAKIKASDDAISKKSTSVNVSDGKIKRRDKKVSAGR
	TTVVVENTKEDKTVYHAGTKDGVQALIRAVVVTKSIKNEIMWDKITKLFPNNSLDNLKKK
	WTARRVRMGHSGWRAYVDKWKKMLVLAIKSEKISLRDVEELDLIKLLDIWTSFDEKEIKR
	PLFLYKNYEENRKKFTLVRDDTLTHSGNDLAMSSMIQREISSLKKTYTRKISASTKDLSK
	SQSDDYIRTVIRSILIESPSTTRNEIEALKNVGNESIDNVIMDMAKEKQIYLHGSKLECT
	DTLPDILENRGNYKDFGVAFQYRCKVNELLEAGNAIVINQEPSDISSWVLIDLISGELLN
	MDVIPMVRNVRPLTYTSRRFEIRTLTPPLIIYANSQTKLNTARKSAVKVPLGKPFSRLWV
	NGSGSIRPNIWKQVVTMVVNEIIFHPGITLSRLQSRCREVLSLHEISEICKWLLERQVLI
	TTDFDGYWVNHNWYSIYEST*
	`),
	}

	blaster := NewBlastp([]seq.Sequence{sequence}, dbPath)
	blaster.SetFlag("evalue", 0.1)

	results, err := Blast(blaster)
	if err != nil {
		fmt.Println(err)
		return
	}

	hit := results.Iterations[0].Hits[0].Def
	fmt.Println(strings.Contains(strings.ToLower(hit), "tfc3"))
	// Output:
	// true
}
