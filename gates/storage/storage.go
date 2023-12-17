package storage

import "errors"

// Storage - интерфейс, представляющий обобщенное хранилище данных.
type Storage interface {

	// Len возвращает количество элементов в хранилище.
	Len() int64

	// NextIndex возвращает индекс следующего добавляемого элемента.
	NextIndex() int64

	// Add добавляет элемент в хранилище и возвращает его уникальный идентификатор и возможную ошибку.
	// value может быть любого типа.
	// Если тип value отличается от типов уже присутствующих в хранилище элементов,
	// возвращается ошибка ErrMismatchType. Если хранилище пусто, тип данных value становится допустимым типом для хранилища,
	// и ошибка не возвращается.
	Add(value interface{}) (int64, error)

	// AddToIndex добавляет элемент в хранилище по указанному индексу.
	AddToIndex(value interface{}, index int64) (error)

	// RemoveByIndex удаляет элемент с указанным индексом из хранилища.
	// Если элемента с таким индексом нет, функция не делает ничего.
	RemoveByIndex(id int64)

	// RemoveByValue удаляет первый найденный элемент с указанным значением из хранилища.
	// Если элемента с таким значением нет, функция не делает ничего.
	RemoveByValue(value interface{})

	// RemoveAllByValue удаляет все элементы с указанным значением из хранилища.
	// Если элементов с таким значением нет, функция не делает ничего.
	RemoveAllByValue(value interface{})

	// GetByIndex возвращает значение элемента с указанным индексом.
	// Если элемента с таким индексом нет, возвращается nil и false.
	GetByIndex(id int64) (interface{}, bool)

	// GetByValue возвращает индекс первого найденного элемента с указанным значением.
	// Если элемента с таким значением нет, возвращается 0 и false.
	GetByValue(value interface{}) (int64, bool)

	// GetAllByValue возвращает индексы всех найденных элементов с указанным значением.
	// Если элементов с таким значением нет, возвращается nil и false.
	GetAllByValue(value interface{}) ([]int64, bool)

	// GetAll возвращает все элементы хранилища.
	// Если хранилище пусто, возвращается nil и false.
	GetAll() ([]interface{}, bool)

	// Clear удаляет все элементы из хранилища.
	Clear()

	// Print выводит содержимое хранилища в консоль.
	Print()
}

// ErrMismatchType ошибка, возвращаемая методом Add, если тип добавляемого элемента
// не соответствует типу уже присутствующих в хранилище элементов.
var ErrMismatchType = errors.New("mismatched type: the type of the provided value does not match the type of items already in the storage")

