package event

import (
	"strings"
	"sync"
)

const Sep = "/"

type Dispatch struct {
	mu       sync.Mutex
	handlers []func(event string) bool
	trie     map[string]*Dispatch
}

func (d *Dispatch) Event(pattern string, handler func(event string) bool) {
	if pattern == "" {
		d.event(nil, handler)
		return
	}
	d.event(strings.Split(pattern, Sep), handler)
}

func (d *Dispatch) Happen(event string) bool {
	if event == "" {
		return d.happen(nil)
	}
	return d.happen(strings.Split(event, Sep))
}

func (d *Dispatch) event(pattern []string, handler func(event string) bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(pattern) == 0 {
		d.handlers = append(d.handlers, handler)
		return
	}

	if d.trie == nil {
		d.trie = make(map[string]*Dispatch)
	}
	if d.trie[pattern[0]] == nil {
		d.trie[pattern[0]] = &Dispatch{}
	}
	d.trie[pattern[0]].event(pattern[1:], handler)
}

func (d *Dispatch) happen(event []string) bool {
	d.mu.Lock()
	handlers := d.handlers
	d.mu.Unlock()

	for _, handler := range handlers {
		if handler(strings.Join(event, Sep)) {
			return true
		}
	}

	if len(event) == 0 {
		return false
	}

	d.mu.Lock()
	next := d.trie[event[0]]
	d.mu.Unlock()

	if next != nil {
		return next.happen(event[1:])
	}
	return false
}
