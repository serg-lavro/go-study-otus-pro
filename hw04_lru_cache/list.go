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
	length int
	front  *ListItem
	back   *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.length == 0 {
		l.front = item
		l.back = item
	} else {
		item.Next = l.front
		l.front.Prev = item
		l.front = item
	}
	l.length++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}
	if l.length == 0 {
		l.front = item
		l.back = item
	} else {
		item.Next = nil
		l.back.Next = item
		item.Prev = l.back
		l.back = item
	}
	l.length++
	return item
}

func (l *list) Remove(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else { // remove front
		l.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else { // remove tail
		l.back = i.Prev
	}
	i.Next = nil
	i.Prev = nil
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if l.front == i {
		return
	}
	l.Remove(i)
	i.Next = l.front
	i.Prev = nil
	if l.front != nil {
		l.front.Prev = i
	}
	l.front = i
	if l.back == nil {
		l.back = i
	}
	l.length++
}
