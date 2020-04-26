package cronos

//OptionsList - list of options
type OptionsList []Options

// Options - dependency injection options
type Options struct {
	key      key
	task     func(objectKey key, digo *Apollo)
	priority int
}

func (ol OptionsList) Len() int {
	return len(ol)
}

func (ol OptionsList) Less(i, j int) bool {
	return ol[i].priority < ol[j].priority
}

func (ol OptionsList) Swap(i, j int) {
	ol[i], ol[j] = ol[j], ol[i]
}
