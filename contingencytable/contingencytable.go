// Package contingecytable implements Contingecy Table
package contingecytable

// Sparse contingey table
type Sparse struct {
	strideMap map[int]int
	countMap  map[int]int
}

// LoadFromData creates a new contingey table from data
func LoadFromData(data [][]int, cardinality []int) *Sparse {
	sp := Sparse{make(map[int]int), make(map[int]int)}
	stride := 1
	for k, v := range cardinality {
		sp.strideMap[k] = stride
		stride *= v
	}
	// TODO count data
	return &sp
}
