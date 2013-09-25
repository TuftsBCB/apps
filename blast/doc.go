/*
Package blast provides functions and types to help with running any of the
BLAST suite of programs. Namely, this package defines an interface `Blaster`
whereby values of types that implement it can execute a BLAST search using
the `Blast` function in this package.

Finally, the results of a BLAST search are captured as XML data and loaded
into the `BlastResults` structure automatically.
*/
package blast
