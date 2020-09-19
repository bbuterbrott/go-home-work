package hw04_lru_cache //nolint:golint,stylecheck

// List is an interface for double linked lists.
type List interface {
	// Returns length of the list
	Len() int
	// Returns first list element
	Front() *listItem
	// Returns last list element
	Back() *listItem
	// Adds new first element
	PushFront(v interface{}) *listItem
	// Adds new last element
	PushBack(v interface{}) *listItem
	// Deletes element from the list
	Remove(i *listItem)
	// Moves element to the front of the list
	MoveToFront(i *listItem)
	// Clears the list
	Clear()
}

type listItem struct {
	// Value
	Value interface{}
	// Next element
	Next *listItem
	// Previous element
	Prev *listItem
}

type list struct {
	front *listItem
	back  *listItem
	len   int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *listItem {
	return l.front
}

func (l *list) Back() *listItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *listItem {
	if l.front == nil {
		l.front = &listItem{Value: v, Next: nil, Prev: nil}
		l.back = l.front
	} else {
		oldFront := l.front
		l.front = &listItem{Value: v, Next: oldFront, Prev: nil}
		oldFront.Prev = l.front
	}
	l.len++
	return l.front
}

func (l *list) PushBack(v interface{}) *listItem {
	if l.back == nil {
		l.back = &listItem{Value: v, Next: nil, Prev: nil}
		l.front = l.back
	} else {
		oldBack := l.back
		l.back = &listItem{Value: v, Next: nil, Prev: oldBack}
		oldBack.Next = l.back
	}
	l.len++
	return l.back
}

func (l *list) Remove(i *listItem) {
	if l.front == i {
		l.front = i.Next
	}
	if l.back == i {
		l.back = i.Prev
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	if i.Prev != nil {
		i.Prev.Next = i.Next
	}
	l.len--
}

func (l *list) MoveToFront(i *listItem) {
	if l.front == i {
		return
	}
	l.Remove(i)
	l.PushFront(i.Value)
}

func (l *list) Clear() {
	l.front = nil
	l.back = nil
	l.len = 0
}

func NewList() List {
	return &list{}
}
