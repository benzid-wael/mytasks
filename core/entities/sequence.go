package entities

type Sequence int

func (seq *Sequence) Next() {
	*seq++
}

func (seq *Sequence) Current() int {
	return int(*seq)
}
