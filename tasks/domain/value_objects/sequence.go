package value_objects

import (
	"io/ioutil"
	"log"
	"strconv"
)

type Sequence int

func (seq *Sequence) Next() {
	*seq++
}

func (seq *Sequence) Current() int {
	return int(*seq)
}

func SaveSequence(path string, seq *Sequence) error {
	bytes := []byte(strconv.Itoa(seq.Current()))
	return ioutil.WriteFile(path, bytes, 0644)
}

func NewSequence(index int) *Sequence {
	sequence := Sequence(index)
	return &sequence
}

func NewSequenceFromFS(path string) *Sequence {
	payload, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	index, err2 := strconv.ParseInt(string(payload), 10, 0)
	if err2 != nil {
		log.Fatalf("invalid sequence fount in file: %v.\n\tOriginal Error: %v\n", path, err2)
	}
	return NewSequence(int(index))
}
