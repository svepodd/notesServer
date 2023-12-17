package mp

import (
	"fmt"
	"Notes/gates/storage"
	"reflect"
	"sort"
	"sync"
)

type Map struct {
	nextIndex int64
	mp        map[int64]interface{}
	mtx       sync.RWMutex
}

func NewMap() *Map {
	fmt.Println("NewMap")
	return &Map{nextIndex: 1, mp: make(map[int64]interface{})}
}

// Len возвращает длину списка
func (m *Map) Len() (l int64) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return int64(len(m.mp))
}

// NextIndex возвращает индекс следующего добавляемого элемента
func (m *Map) NextIndex() (index int64) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	return m.nextIndex
}

// Add добавляет элемент в список и возвращает его индекс
func (m *Map) Add(value interface{}) (id int64, err error) {
	if m.nextIndex == 1 {
		m.mp = make(map[int64]interface{})
	}
	if m.nextIndex > 1 && reflect.TypeOf(m.mp[m.nextIndex-1]).Kind() != reflect.TypeOf(value).Kind() { // проверяем тип данных
		return 0, storage.ErrMismatchType // ошибка - не совпадают типы данных
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.mp[m.nextIndex] = value
	m.nextIndex++
	return m.nextIndex - 1, nil
}

func (m *Map) AddToIndex(value interface{}, index int64) (error) {

	if _, exists := m.mp[index]; exists {
		return fmt.Errorf("index %d already exists", index)
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()
	
	m.mp[index] = value
	if index >= m.nextIndex {
		m.nextIndex = index + 1
	}
	return nil
}

// RemoveByIndex удаляет элемент из списка по индексу
func (m *Map) RemoveByIndex(id int64) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	delete(m.mp, id)
	if id == m.nextIndex-1 {
		m.nextIndex--
	}
}

// RemoveByValue удаляет элемент из списка по значению
func (m *Map) RemoveByValue(value interface{}) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for k, v := range m.mp {
		if v == value {
			delete(m.mp, k)
			if k == m.nextIndex-1 {
				m.nextIndex--
			}
			return
		}
	}
}

// RemoveAllByValue удаляет все элементы из списка по значению
func (m *Map) RemoveAllByValue(value interface{}) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for k, v := range m.mp {
		if v == value {
			delete(m.mp, k)
		}
	}

	keys := make([]int64, 0, len(m.mp))
	for k := range m.mp {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	m.nextIndex = keys[len(keys)-1] + 1

}

// GetByIndex возвращает значение элемента по индексу.
// Если элемента с таким индексом нет, то возвращается 0 и false.
func (m *Map) GetByIndex(id int64) (value interface{}, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	value, ok = m.mp[id]
	if !ok {
		return nil, false
	}
	return
}

// GetByValue возвращает индекс первого найденного элемента по значению.
// Если элемента с таким значением нет, то возвращается 0 и false.
func (m *Map) GetByValue(value interface{}) (id int64, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	for k, v := range m.mp {
		if v == value {
			return k, true
		}
	}
	return 0, false
}

// GetAllByValue возвращает индексы всех найденных элементов по значению
// Если элементов с таким значением нет, то возвращается nil и false.
func (m *Map) GetAllByValue(value interface{}) (ids []int64, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	ids = make([]int64, 0)
	for k, v := range m.mp {
		if v == value {
			ids = append(ids, k)
			return ids, true
		}
	}
	return nil, false
}

// GetAll возвращает все элементы списка
// Если список пуст, то возвращается nil и false.
func (m *Map) GetAll() (values []interface{}, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()

	if len(m.mp) == 0 {
		return nil, false
	}

	for _, v := range m.mp {
		values = append(values, v) 
	}
	return values, true
}

// Clear очищает список
func (m *Map) Clear() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.mp = make(map[int64]interface{})
	m.nextIndex = 1
}

// Print выводит список в консоль
func (m *Map) Print() {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	if len(m.mp) == 0 {
		fmt.Printf("Map is empty!\n")
		return
	}

	fmt.Println(m.mp)
}
