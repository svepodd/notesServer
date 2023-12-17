package list

type node struct {
	index int64 
	value interface{}
	next  *node
}