package common

import (
	"strings"
	"sync"
)

type DomainTree struct {
	node
	sep string
	sync.RWMutex
}
type node struct {
	value  interface{}
	branch map[string]*node
}

func NewDomainTree(sep string) *DomainTree {
	return &DomainTree{
		node: node{
			value:  nil,
			branch: map[string]*node{},
		},
		sep:     sep,
		RWMutex: sync.RWMutex{},
	}
}

func (t *DomainTree) Store(k string, v interface{}) {
	t.Lock()
	t.store(strings.Split(strings.TrimSuffix(k, t.sep), t.sep), v)
	t.Unlock()
}
func (t *DomainTree) UnsafeStore(k string, v interface{}) {
	t.store(strings.Split(strings.TrimSuffix(k, t.sep), t.sep), v)
}
func (n *node) store(ks []string, v interface{}) {
	l := len(ks)
	switch l {
	case 0:
		return
	case 1:
		k := ks[l-1]

		if k == "*" || k == "**" {
			n.value = v
		}

		b, ok := n.branch[k]
		if ok {
			b.value = v
			return
		}

		n.branch[k] = &node{
			value:  v,
			branch: map[string]*node{},
		}
	default:
		k := ks[l-1]

		b, ok := n.branch[k]
		if !ok {
			b = &node{
				value:  nil,
				branch: map[string]*node{},
			}
			n.branch[k] = b
		}

		b.store(ks[:l-1], v)
	}
}

func (t *DomainTree) Load(k string) interface{} {
	t.RLock()
	v := t.load(strings.Split(strings.TrimSuffix(k, t.sep), t.sep))
	t.RUnlock()
	return v
}
func (t *DomainTree) UnsafeLoad(k string) interface{} {
	return t.load(strings.Split(strings.TrimSuffix(k, t.sep), t.sep))
}
func (n *node) load(ks []string) interface{} {
	l := len(ks)
	switch l {
	case 0:
		return nil
	case 1:
		b, ok := n.branch[ks[l-1]]
		if ok {
			return b.value
		}

		b, ok = n.branch["*"]
		if ok {
			return b.value
		}

		b, ok = n.branch["**"]
		if ok {
			return b.value
		}

		return nil
	default:
		b, ok := n.branch[ks[l-1]]
		if ok {
			s := b.load(ks[:l-1])
			if s != nil {
				return s
			}
		}

		b, ok = n.branch["*"]
		if ok {
			s := b.load(ks[:l-1])
			if s != nil {
				return s
			}
		}

		b, ok = n.branch["**"]
		if ok {
			return b.value
		}

		return nil
	}
}
