package history

const maxHistorySize = 100

type History struct {
	inputs map[int]string
	cursor int
}

func NewHistory() *History {
	return &History{
		map[int]string{},
		0,
	}
}

func (h *History) Reset() *History {
	h.inputs = map[int]string{}
	h.cursor = 0

	return h
}

func (h *History) Add(input string) *History {
	if len(h.inputs) >= maxHistorySize {
		h.pruneOldest()
	}
	h.cursor = len(h.inputs)
	h.inputs[h.cursor] = input

	return h
}

func (h *History) pruneOldest() {
	if len(h.inputs) == 0 {
		return
	}
	lowestKey := 0
	newInputs := make(map[int]string)
	offset := lowestKey + 1
	for k, v := range h.inputs {
		if k != lowestKey {
			newInputs[k-offset] = v
		}
	}
	h.inputs = newInputs
}

func (h *History) GetAll() map[int]string {
	return h.inputs
}

func (h *History) GetCursor() int {
	return h.cursor
}

func (h *History) GetPrevious() *string {
	if input, ok := h.inputs[h.cursor]; ok {
		h.cursor--
		return &input
	}

	return nil
}

func (h *History) GetNext() *string {
	if input, ok := h.inputs[h.cursor+1]; ok {
		h.cursor++
		return &input
	}

	return nil
}
