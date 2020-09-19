package hw04_lru_cache //nolint:golint,stylecheck

type List interface {
	// длина списка
	Len() int
	// первый элемент списка
	Front() *listItem
	// последний элемент списка
	Back() *listItem
	// добавить значение в начало
	PushFront(v interface{}) *listItem
	// добавить значение в конец
	PushBack(v interface{}) *listItem
	// удалить элемент
	Remove(i *listItem)
	// переместить элемент в начало
	MoveToFront(i *listItem)
}

type listItem struct {
	// значение
	Value interface{}
	// следующий элемент
	Next *listItem
	// предыдущий элемент
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

func NewList() List {
	return &list{}
}
