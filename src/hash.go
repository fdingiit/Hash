package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	LoadFactor = 0.5
	ProbeSlots = 10
)

func hashCode(str string) uint32 {
	var hc uint32

	for _, c := range str {
		hc = hc*uint32(33) + uint32(c)
	}

	return hc
}

/**********************************
 * State
 **********************************/
type State int

const (
	EmptyState    = 0
	OccupiedState = 1
	DeletedState  = 2
)

/**********************************
 * Content
 **********************************/
type content struct {
	state State
	hCode int
	str   *string
}

func (c *content) equal(a *content) bool {
	return *c.str == *a.str && c.state == a.state && c.hCode == a.hCode
}

func (c *content) String() string {
	if c.state == DeletedState {
		return "DELETED"
	}
	if c.str == nil {
		return ""
	}

	return *c.str + " (" + strconv.Itoa(c.hCode) + ")"
}

/**********************************
 * Hashtable
 **********************************/
type hashtable struct {
	used   int
	size   int
	factor float32

	/* if we need to rehash, and do not have a cluster
	 * how many slots we should try to remove at least around a slot
	 */
	rehash_min_cnt int

	slots []*content
}

func newHashtable(size int) (h *hashtable) {
	h = new(hashtable)
	h.used, h.size, h.factor = 0, size, 0.0
	h.rehash_min_cnt = 10
	h.slots = make([]*content, GetPrime(size))

	return
}

func (h *hashtable) getLoadFactor() float32 {
	h.factor = float32(h.used) / float32(h.size)
	return h.factor
}

func (h *hashtable) insert(str *string) (index int, probe int) {
	var hCode int

	if str == nil {
		return -1, 0
	}

	if h.used == h.size {
		return -1, 0
	}

	if i, _ := h.find(*str); i != -1 {
		return -1, 0
	}

	hCode = int(hashCode(*str) % uint32(h.size))
	index = hCode

	for {
		if h.slots[index] == nil || h.slots[index].state == DeletedState {
			if h.slots[index] == nil {
				h.slots[index] = new(content)
			}

			h.slots[index].state = OccupiedState
			h.slots[index].str = str
			h.slots[index].hCode = hCode

			h.used++
			return
		}

		probe++
		index = (index + 1) % h.size
	}
}

func (h *hashtable) insertAndCheck(str *string) (index int, needRehash bool) {
	var probe int

	index, probe = h.insert(str)

	if index == -1 && probe == 0 {
		return -1, false
	}

	if probe > ProbeSlots || h.getLoadFactor() >= LoadFactor {
		return index, true
	}

	return index, false
}

func (h *hashtable) remove(index int) {
	if h.slots == nil {
		return
	}

	if index == -1 {
		return
	}

	h.used--
	h.slots[index].state = DeletedState

	return
}

func (h *hashtable) cleanSlot(index int) {
	if h.slots == nil {
		return
	}

	h.slots[index] = nil
	h.used--

	return
}

func (h *hashtable) find(str string) (index int, probe int) {
	var start int
	var again bool = false

	if h.slots == nil {
		return -1, 0
	}

	start = int(hashCode(str) % uint32(h.size))

	for index = start; index != start || !again; index = (index + 1) % h.size {
		if !again {
			again = true
		}

		if h.slots[index] != nil && h.slots[index].state == OccupiedState && *h.slots[index].str == str {
			return
		}

		probe++
	}

	return -1, probe
}

func (h *hashtable) equal(a *hashtable) bool {
	if h.used != a.used {
		return false
	}
	if h.size != a.size {
		return false
	}

	if (h.slots == nil && a.slots != nil) || (h.slots == nil && a.slots != nil) {
		return false
	}

	for i := 0; i < h.size; i++ {
		if !h.slots[i].equal(a.slots[i]) {
			return false
		}
	}

	return true
}

/**********************************
 * Hash APIs
 **********************************/
type Hash struct {
	/* hashing information */
	rehashing bool

	now, new int

	/* 2 hashtables
	 * 1 for now, 1 for new
	 */
	hts []*hashtable
}

func NewHash(n int) (*Hash, error) {
	var hash *Hash
	var now_size int

	if n < 101 {
		now_size = 101
	} else if n >= 101 && n <= 199999 {
		now_size = n
	} else {
		return nil, errors.New("out_of_range")
	}

	hash = new(Hash)

	hash.now, hash.new = 0, 1
	hash.hts = make([]*hashtable, 2)
	hash.hts[hash.now] = newHashtable(now_size)

	hash.rehashing = false

	return hash, nil
}

func (h *Hash) Equal(a *Hash) bool {
	if h.rehashing || a.rehashing {
		return false
	}

	return h.hts[0].equal(a.hts[0]) && h.hts[1].equal(a.hts[1])
}

func (h *Hash) Copy(a *Hash) {
	// TODO:
}

func (h *Hash) Insert(str string) error {
	var dstr *string = new(string)

	*dstr = str
	return h.insert(dstr)
}

func (h *Hash) Find(str string) (bool, error) {
	var index, arr int
	var err error

	index, arr, err = h.find(str)

	if err != nil {
		return false, err
	}

	if index == -1 || arr == -1 {
		return false, nil
	}

	return true, nil
}

func (h *Hash) Remove(str string) (string, error) {
	var index, arr int
	var err error

	index, arr, err = h.find(str)
	if err != nil {
		return "", err
	}

	if index == -1 || arr == -1 {
		return "", nil
	}

	h.hts[arr].remove(index)

	return *h.hts[arr].slots[index].str, nil
}

func (h *Hash) Dump() {
	for i := range h.hts {
		if h.hts[i] != nil {
			fmt.Fprintf(os.Stdout, "HashTable #%d: size = %d, tableSize = %d\n", i+1, h.hts[i].used, h.hts[i].size)

			for j, s := range h.hts[i].slots {
				if s == nil {
					fmt.Fprintf(os.Stdout, "H%d[%d] =\n", i+1, j)
				} else {
					fmt.Fprintf(os.Stdout, "H%d[%d] = %s\n", i+1, j, s.String())
				}
			}
		}
	}
}

/**********************************
 * Hash implements
 **********************************/
func (h *Hash) find(str string) (int, int, error) {
	var index int
	var err error

	/* continue rehashing if needed */
	if h.rehashing {
		if err = h.rehash(-1); err != nil {
			return -1, -1, err
		}

		/* find in new hashtable */
		index, _ = h.hts[h.new].find(str)
		return index, h.new, nil
	}

	/* find in now hashtable */
	index, _ = h.hts[h.now].find(str)
	return index, h.now, nil
}

func (h *Hash) insert(str *string) error {
	var index int
	var needRehash bool
	var err error

	if h.rehashing {
		if err = h.rehash(-1); err != nil {
			return err
		}

		_, needRehash = h.hts[h.new].insertAndCheck(str)

		/* a very special case that we need to rehash when rehashing*/
		if needRehash {
			return h.specialRehash()
		}
	}

	/* not rehashing, do insert, then check if we need to start a rehash */
	index, needRehash = h.hts[h.now].insertAndCheck(str)
	if needRehash {
		h.rehashing = true
		h.hts[h.new] = newHashtable(4 * h.hts[h.now].used)
		return h.rehash(index)
	}

	return nil
}

/**
 * hashing
 */
func (h *Hash) rehash(index int) error {
	var start, end int
	var cnt int
	var needRehash bool

	if !h.rehashing {
		return errors.New("no in rehashing")
	}

	if h.hts == nil {
		return errors.New("nul h.hts")
	}

	if h.hts[h.now] == nil {
		return errors.New("nul h.hts[h.now]")
	}

	start = index
	for h.hts[h.now].slots[start] != nil || cnt < h.hts[h.now].rehash_min_cnt {
		start = (start - 1) % h.hts[h.now].size
		cnt++
	}

	cnt = 0

	end = index
	for h.hts[h.now].slots[end] != nil || cnt < h.hts[h.now].rehash_min_cnt {
		end = (end + 1) % h.hts[h.now].size
		cnt++
	}

	for start != end {
		var needRehash bool

		if h.hts[h.now].slots[start].state == OccupiedState {
			_, needRehash = h.hts[h.new].insertAndCheck(h.hts[h.now].slots[start].str)
			h.hts[h.now].cleanSlot(start)

			/* a very special case that we need to rehash when rehashing*/
			if needRehash {
				return h.specialRehash()
			}
		}

		start = (start + 1) % h.hts[h.now].size
	}

	/* remove all elements to new hashtable if load factor is less than 0.03 */
	if h.hts[h.now].factor <= 0.03 {
		for i, c := range h.hts[h.now].slots {
			if c != nil && c.state == OccupiedState {
				_, needRehash = h.hts[h.new].insertAndCheck(c.str)
				h.hts[h.now].cleanSlot(i)

				/* a very special case that we need to rehash when rehashing*/
				if needRehash {
					return h.specialRehash()
				}

				h.rehashNormalFinish()
			}
		}

		return nil
	}

	if h.hts[h.now].used == 0 {
		h.rehashNormalFinish()
	}

	return nil
}

func (h *Hash) rehashNormalFinish() {
	h.hts[h.now] = nil
	h.now, h.new = h.new, h.now
	h.rehashing = false
}

func (h *Hash) specialRehash() error {
	var newsize int
	var spec *hashtable

	if !h.rehashing {
		return errors.New("no in rehashing")
	}

	newsize = 4 * (h.hts[0].used + h.hts[1].used)
	if newsize > 199999 {
		return errors.New("out_of_range")
	}

	spec = newHashtable(newsize)

	for _, c := range h.hts[0].slots {
		spec.insert(c.str)
	}

	for _, c := range h.hts[1].slots {
		spec.insert(c.str)
	}

	h.now, h.new = 0, 1
	h.hts[h.now] = spec
	h.hts[h.new] = nil
	h.rehashing = false

	return nil
}

/**********************************
 * Grading functions
 **********************************/
func (h *Hash) isRehashing() bool {
	// TODO:
	return h.rehashing
}

func (h *Hash) tableSize(table int) int {
	// TODO:
	return -1
}

func (h *Hash) size() int {
	// TODO:
	return -1
}

func (h *Hash) at(index, table int) (string, error) {
	// TODO:
	return "", nil
}
