package main

type Elements struct {
	items []Element
}

type ElementAction func(int, Element)
type ElementQuery func(Element) bool

func (this *Elements) Add(e Element) {
	this.items = append(this.items, e)
}

func (this *Elements) Delete(e Element) {
	i := this.IndexOf(e)
	if i > -1 {
		this.items = append(this.items[:i], this.items[i+1:]...)
	}
}

func (this *Elements) IndexOf(e Element) int {
	for i, el := range this.items {
		if el == e {
			return i
		}
	}
	debugf("!Could not find element %v", e)
	return -1
}

func (this *Elements) Each(action ElementAction) {
	for i, e := range this.items {
		action(i, e)
	}
}

func (this *Elements) Any(query ElementQuery) bool {
	for _, e := range this.items {
		if query(e) {
			return true
		}
	}
	return false
}

func (this *Elements) Count() int {
	return len(this.items)
}
