package list

import (
	"errors"
	"fmt"
	"Notes/gates/storage"
	"reflect"
	"sync"
)

type List struct {
	len       int64
	firstNode *node
	mtx       sync.RWMutex
}

// NewList создает новый список
func NewList() (l *List) {
	fmt.Println("NewList")
	return &List{len: 0, firstNode: nil}
}

// Len возвращает длину списка
func (l *List) Len() (len int64) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	return len
}

// NextIndex возвращает индекс следующего добавляемого элемента
func (l *List) NextIndex() (index int64) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	if l.firstNode == nil{
		return 1
	}

	currentNode := l.firstNode
	for ; currentNode.next != nil; currentNode = currentNode.next {}
	return currentNode.index + 1
}


// Add добавляет элемент в список и возвращает его индекс
func (l *List) Add(value interface{}) (id int64, err error) {
	if l.len > 0 && reflect.TypeOf(l.firstNode.value).Kind() != reflect.TypeOf(value).Kind() {
		return 0, storage.ErrMismatchType
	}

	l.mtx.Lock()
	defer l.mtx.Unlock()

	newNode := &node{value: value}

	if l.firstNode == nil {
		newNode.index = 1
		l.firstNode = newNode
		l.len++
		return 1, nil
	}

	l.len++
	currentNode := l.firstNode
	for ; currentNode.next != nil; currentNode = currentNode.next {
	}
	newNode.index = currentNode.index + 1
	currentNode.next = newNode

	return currentNode.index + 1, nil
}

// AddByIndex добавляет элемент в список по индексу
func (l *List) AddToIndex(value interface{}, index int64) (error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if index > l.len + 1 {
		return errors.New("index out of range")
	}

	newNode := &node{value: value, index: index}

	if l.firstNode == nil {
		l.firstNode = newNode
		l.len++
		return nil
	}

	if index == 0 {
		newNode.next = l.firstNode
		l.firstNode = newNode
	} else {
		currentNode := l.firstNode
		for i := int64(0); i < index-1; i++ {
			currentNode = currentNode.next
		}
		newNode.next = currentNode.next
		currentNode.next = newNode
	}

	l.len++
	return nil

}

// RemoveByIndex удаляет элемент из списка по индексу
func (l *List) RemoveByIndex(id int64) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.firstNode == nil {
		return
	}
	if l.firstNode.index == id {
		l.firstNode = l.firstNode.next
		l.len--
		return
	}
	prevNode := l.firstNode
	for currentNode := l.firstNode.next; currentNode != nil; currentNode = currentNode.next {
		if currentNode.index == id {
			prevNode.next = currentNode.next
			l.len--
			return
		}
		prevNode = currentNode
	}

}

// RemoveByValue удаляет элемент из списка по значению
func (l *List) RemoveByValue(value interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.firstNode == nil {
		return
	}

	if l.firstNode.value == value {
		l.firstNode = l.firstNode.next
		l.len--
		return
	}
	prevNode := l.firstNode
	for currentNode := l.firstNode.next; currentNode != nil; currentNode = currentNode.next {
		if currentNode.value == value {
			prevNode.next = currentNode.next
			l.len--
			return
		}
		prevNode = currentNode
	}
}

// RemoveAllByValue удаляет все элементы из списка по значению
func (l *List) RemoveAllByValue(value interface{}) {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if l.firstNode == nil {
		return
	}

	if l.firstNode.value == value {
		l.firstNode = l.firstNode.next
		l.len--
	}
	prevNode := l.firstNode
	for currentNode := l.firstNode.next; currentNode != nil; currentNode = currentNode.next {
		if currentNode.value == value {
			prevNode.next = currentNode.next
			l.len--
			currentNode = prevNode
		}
		prevNode = currentNode
	}
}

// GetByIndex возвращает значение элемента по индексу.
// Если элемента с таким индексом нет, то возвращается 0 и false.
func (l *List) GetByIndex(id int64) (value interface{}, ok bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	if l.firstNode == nil {
		return 0, false
	}
	if l.firstNode.index == id {
		return l.firstNode.value, true
	}
	for currentNode := l.firstNode.next; currentNode != nil; currentNode = currentNode.next {
		if currentNode.index == id {
			return currentNode.value, true
		}
	}
	return 0, false
}

// GetByValue возвращает индекс первого найденного элемента по значению.
// Если элемента с таким значением нет, то возвращается 0 и false.
func (l *List) GetByValue(value interface{}) (id int64, ok bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	if l.firstNode == nil {
		return 0, false
	}

	if l.firstNode.value == value {
		return l.firstNode.index, true
	}
	for currentNode := l.firstNode.next; currentNode != nil; currentNode = currentNode.next {
		if currentNode.value == value {
			return currentNode.index, true
		}
	}
	return
}

// GetAllByValue возвращает индексы всех найденных элементов по значению
// Если элементов с таким значением нет, то возвращается nil и false.
func (l *List) GetAllByValue(value interface{}) (ids []int64, ok bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	if l.firstNode == nil {
		return nil, false
	}

	if l.firstNode.value == value {
		ids = append(ids, l.firstNode.index)
	}
	for currentNode := l.firstNode.next; currentNode != nil; currentNode = currentNode.next {
		if currentNode.value == value {
			ids = append(ids, currentNode.index)
		}
	}
	if len(ids) == 0 {
		return nil, false
	}
	return ids, true
}

// GetAll возвращает все элементы списка
// Если список пуст, то возвращается nil и false.
func (l *List) GetAll() (values []interface{}, ok bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	if l.firstNode == nil {
		return nil, false
	}
	values = append(values, l.firstNode.value)

	for currentNode := l.firstNode.next; currentNode != nil; currentNode = currentNode.next {
		values = append(values, currentNode.value)
	}
	return values, true
}

// Clear очищает список
func (l *List) Clear() {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	l.firstNode = nil 
	l.len = 0
}

// Print выводит список в консоль
func (l *List) Print() {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	if l.firstNode == nil {
		fmt.Println("Empty")
		return
	}
	n := l.firstNode
	fmt.Print("[")
	for ; n.next != nil; n = n.next {
		fmt.Print(n.value)
		fmt.Print(", ")
	}
	fmt.Printf("%v]\n", n.value)
}
