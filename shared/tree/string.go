/*
  multiple tree implementations
*/
package tree

type any interface{}

type StringTree interface {
	Parent() StringTree
	Root() StringTree

	Path() string
	Value() any

	Len() int

	Get(path string) (any, bool)
	Put(path string, value any)
	Remove(path string) (any, bool)
}

type stringTreeImpl struct {
	parent *stringTreeImpl
	path   string
	value  any

	children map[rune]*stringTreeImpl
}

func NewStringTree() StringTree {
	return &stringTreeImpl{
		nil,
		"",
		nil,
		make(map[rune]*stringTreeImpl),
	}
}

func newStringTreeImpl(path string, value any, parent *stringTreeImpl) *stringTreeImpl {
	return &stringTreeImpl{
		path:     path,
		value:    value,
		parent:   parent,
		children: make(map[rune]*stringTreeImpl),
	}
}

func (tree *stringTreeImpl) Parent() StringTree {
	return tree.parent
}

func (tree *stringTreeImpl) Root() StringTree {
	current := tree
	for current.parent != nil {
		current = current.parent
	}

	return current
}

func (tree *stringTreeImpl) Path() string {
	return tree.path
}

func (tree *stringTreeImpl) key() rune {
	return rune(tree.path[len(tree.path)-1])
}

func (tree *stringTreeImpl) Value() any {
	return tree.value
}

func (tree *stringTreeImpl) Len() int {
	var result int

	for _, child := range tree.children {
		if child.value != nil {
			result += 1
		}
		result += child.Len()
	}

	return result
}

func (tree *stringTreeImpl) Get(path string) (any, bool) {
	current := tree

	for _, character := range path {
		if child, found := current.children[character]; found {
			current = child
		} else {
			break
		}
	}

	return current, current != nil
}

func (tree *stringTreeImpl) Put(path string, value any) {
	current := tree

	for index, character := range path {
		child, found := current.children[character]

		if !found {
			child = newStringTreeImpl(path[0:index], nil, current)
		}

		current = child
	}

	current.value = value
}

func (tree *stringTreeImpl) delete(child *stringTreeImpl) {
	delete(tree.children, child.key())

	if tree.Len() <= 0 && tree.parent != nil {
		tree.parent.delete(tree)
	}
}

func (tree *stringTreeImpl) Remove(path string) (any, bool) {
	current := tree

	for _, character := range path {
		if child, found := current.children[character]; found {
			current = child
		} else {
			return nil, false
		}
	}

	value := current.value
	current.value = nil

	if current.Len() <= 0 {
		current.parent.delete(current)
		current.value = nil
		current.children = nil
	}

	return value, true
}
