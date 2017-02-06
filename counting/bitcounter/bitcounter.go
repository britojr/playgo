package bitcounter

import "github.com/willf/bitset"

// BitCounter ..
type BitCounter struct {
	vars   *bitset.BitSet // wich variables are representend
	cardin *[]int         // cardinality of all variables
	vals   []*valToLine   // all the possible values for a variable
}

type valToLine map[int]*bitset.BitSet

// NewBitCounter creates new BitCounter
func NewBitCounter() *BitCounter {
	b := new(BitCounter)
	return b
}

// LoadFromData initializes the BitCounter from a given dataset
func (b *BitCounter) LoadFromData(dataset [][]int, cardinality []int) {
	lin, col := len(dataset), len(dataset[0])
	b.vars = bitset.New(uint(col)).Complement()
	b.vals = make([]*valToLine, col)
	aux := append([]int(nil), cardinality...)
	b.cardin = &aux
	for i, c := range cardinality {
		b.vals[i] = new(valToLine)
		*b.vals[i] = make(map[int]*bitset.BitSet)
		for j := 0; j < c; j++ {
			(*b.vals[i])[j] = bitset.New(uint(lin))
		}
	}
	for i := 0; i < lin; i++ {
		for j := 0; j < col; j++ {
			(*b.vals[j])[dataset[i][j]].Set(uint(i))
		}
	}
}

// Marginalize ..
func (b *BitCounter) Marginalize(vars ...int) (r *BitCounter) {
	r = b.Clone()
	auxvars := bitset.New(r.vars.Len())
	for _, v := range vars {
		auxvars.Set(uint(v))
	}
	r.vars.InPlaceIntersection(auxvars)
	return
}

// SumOut ..
func (b *BitCounter) SumOut(x int) (r *BitCounter) {
	r = b.Clone()
	r.vars.Clear(uint(x))
	return
}

// Clone ..
func (b *BitCounter) Clone() (r *BitCounter) {
	r = new(BitCounter)
	r.cardin = b.cardin
	r.vals = b.vals
	r.vars = b.vars.Clone()
	return
}

// ValueIterator ..
func (b *BitCounter) ValueIterator() (f func() *int) {
	val := make([]int, b.vars.Count())
	f = func() *int {
		if val == nil {
			return nil
		}
		v := b.getCount(val)
		val = b.nextValuation(val)
		return &v
	}
	return
}

// ValueIteratorNonZero ..
func (b *BitCounter) ValueIteratorNonZero() (f func() *int) {
	val := make([]int, b.vars.Count())
	f = func() *int {
		var v int
		for val != nil && v == 0 {
			v = b.getCount(val)
			val = b.nextValuation(val)
		}
		if v != 0 {
			return &v
		}
		return nil
	}
	return
}

func (b *BitCounter) nextValuation(val []int) []int {
	i := 0
	val[i]++
	j, _ := b.vars.NextSet(0)
	for val[i] == b.getCardinality(int(j)) {
		val[i] = 0
		i++
		if i == len(val) {
			return nil
		}
		j, _ = b.vars.NextSet(j + 1)
		val[i]++
	}
	return val
}

func (b *BitCounter) getCardinality(x int) int {
	return len(*b.vals[x])
}

func (b *BitCounter) getCount(val []int) int {
	j, _ := b.vars.NextSet(0)
	aux := (*b.vals[j])[val[0]].Clone()
	for i := 1; i < len(val); i++ {
		j, _ = b.vars.NextSet(j + 1)
		aux.InPlaceIntersection((*b.vals[j])[val[i]])
	}
	return int(aux.Count())
}