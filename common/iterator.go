package common

type RuneIterator struct {
	Items []rune
	Index int
}

func (t RuneIterator) Get() rune {
	char := t.Items[t.Index]
	return char
}
func (t *RuneIterator) Eat() rune {
	char := t.Items[t.Index]
	t.Index++
	return char
}

func (t RuneIterator) HasNext() bool {
	return t.Index+1 < len(t.Items)
}

func (t RuneIterator) IsOutOfBounds() bool {
	return t.Index >= len(t.Items)
}

func (t RuneIterator) GetNext() rune {
	return t.Items[t.Index+1]
}
