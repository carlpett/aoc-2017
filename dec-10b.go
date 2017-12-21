package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	dataSize  = 256
	blockSize = 16
	rounds    = 64
)

func main() {
	t := time.Now()
	data := make([]byte, dataSize)
	for idx := range data {
		data[idx] = byte(idx)
	}

	lens, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	magic := []byte{17, 31, 73, 47, 23}
	lens = append(lens, magic...)
	hasher := KnotHasher{data: data}

	for r := 0; r < rounds; r++ {
		hasher.hash(lens)
	}
	denseHash := make([]byte, dataSize/blockSize)
	for block := 0; block < dataSize/blockSize; block++ {
		denseHash[block] = xor(hasher.data[block*blockSize : (block+1)*blockSize])
	}
	fmt.Printf("B: %x\n", denseHash)
	fmt.Println(time.Since(t))
}

func xor(bs []byte) byte {
	var sum byte = 0
	for _, b := range bs {
		sum ^= b
	}
	return sum
}

type circularSlice struct {
	s   []byte
	pos int
	l   int
}

func newCircularSlice(s []byte, pos, l int) (c circularSlice) {
	return circularSlice{
		s:   s,
		pos: pos,
		l:   l,
	}
}

func (p circularSlice) mapIndex(i int) int { return (p.pos + i) % len(p.s) }
func (p circularSlice) Len() int           { return p.l }
func (p circularSlice) Less(i, j int) bool {
	return p.s[p.mapIndex(i)] < p.s[p.mapIndex(j)]
}
func (p circularSlice) Swap(i, j int) {
	p.s[p.mapIndex(i)], p.s[p.mapIndex(j)] = p.s[p.mapIndex(j)], p.s[p.mapIndex(i)]
}

func reverse(cs circularSlice) {
	for i := cs.Len()/2 - 1; i >= 0; i-- {
		opp := cs.Len() - 1 - i
		cs.Swap(i, opp)
	}
}

type KnotHasher struct {
	pos  int
	skip int
	data []byte
}

func (kh *KnotHasher) hash(bs []byte) {
	for _, b := range bs {
		reverse(newCircularSlice(kh.data, kh.pos, int(b)))
		kh.pos = (kh.pos + int(b) + kh.skip) % len(kh.data)
		kh.skip++
	}
}
