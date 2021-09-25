package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func (l list) Len() int {
	return l.len
}

func (l list) Front() *ListItem {
	return l.head
}

func (l list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{v, nil, nil}
	if l.head != nil {
		item.Next = l.head
		l.head.Prev = item
	}
	l.head = item
	if l.tail == nil {
		l.tail = l.head
	}
	l.len++
	return l.head
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{v, nil, nil}
	if l.tail != nil {
		item.Prev = l.tail
		l.tail.Next = item
	}

	l.tail = item
	if l.head == nil {
		l.head = l.tail
	}

	l.len++
	return l.tail
}

func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		if i.Next != nil {
			i.Next.Prev = nil
			l.head = i.Next
		}
		l.len--
		return
	}

	if i.Next == nil {
		if i.Prev != nil {
			i.Prev.Next = nil
			l.tail = i.Prev
		}
		l.len--
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}
	if i.Next == nil {
		l.tail = i.Prev
		i.Prev.Next = nil
		i.Prev = nil
		i.Next = l.head
		l.head = i
		return
	}

	i.Next.Prev = i.Prev
	i.Prev.Next = i.Next

	i.Next = l.head
	i.Prev = nil
	i.Next.Prev = i
	l.head = i
}

func NewList() List {
	return new(list)
}
