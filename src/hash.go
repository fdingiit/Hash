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
	return c.state == a.state && c.hCode == a.hCode && *c.str == *a.str
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

	slots []*content
}

func newHashtable(size int) (h *hashtable) {
	h = new(hashtable)
	h.used, h.size, h.factor = 0, GetPrime(size), 0.0
	h.slots = make([]*content, h.size)

	return
}

func (h *hashtable) getLoadFactor() float32 {
	/* lazy calculation */
	h.factor = float32(h.used) / float32(h.size)
	return h.factor
}

func (h *hashtable) calcIndex(str string) int {
	return int(hashCode(str)) % h.size
}

func (h *hashtable) insert(str *string) (index int, probe int) {
	var hCode int

	if h.used == h.size {
		return -1, ProbeSlots
	}

	if str == nil {
		return -1, 0
	}

	if i, _ := h.find(*str); i != -1 {
		return -1, 0
	}

	hCode = h.calcIndex(*str)
	index = hCode

	for {
		probe++

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

		index = (index + 1) % h.size
	}
}

func (h *hashtable) insertAndCheck(str *string) (index int, needRehash bool) {
	var probe int

	index, probe = h.insert(str)
	needRehash = probe >= ProbeSlots || h.getLoadFactor() >= LoadFactor

	return
}

func (h *hashtable) find(str string) (index int, probe int) {
	var start int

	if h.slots == nil {
		return -1, 0
	}

	start = h.calcIndex(str)

	for index = start; probe < ProbeSlots; index = (index + 1) % h.size {
		probe++

		if h.slots[index] != nil && h.slots[index].state == OccupiedState && *h.slots[index].str == str {
			return
		}
	}

	return -1, probe
}

func (h *hashtable) findAndCheck(str string) (index int, needRehash bool) {
	var probe int

	index, probe = h.find(str)
	needRehash = probe >= ProbeSlots || h.getLoadFactor() >= LoadFactor

	return
}

func (h *hashtable) remove(str string) (index int, probe int) {
	if h.slots == nil {
		return -1, 0
	}

	index, probe = h.find(str)

	if index != -1 {
		h.slots[index].state = DeletedState
		h.used--
	}

	return
}

func (h *hashtable) removeAndCheck(str string) (index int, needRehash bool) {
	var probe int

	index, probe = h.remove(str)
	needRehash = probe >= ProbeSlots || h.getLoadFactor() >= LoadFactor

	return
}

func (h *hashtable) removeByIndex(index int) int {
	if h.slots == nil {
		return -1
	}

	if h.slots[index] == nil {
		return -1
	}

	if h.slots[index].state == OccupiedState {
		h.slots[index].state = DeletedState
		h.used--
		return index
	}

	return -1
}

func (h *hashtable) cleanSlot(index int) {
	if h.slots == nil {
		return
	}

	h.slots[index] = nil
	h.used--

	return
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
	return h.insert(str)
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

	index = h.hts[arr].removeByIndex(index)

	if index != -1 {
		return *h.hts[arr].slots[index].str, nil
	}

	return "", nil
}

func (h *Hash) Size() int {
	var size int

	size += h.hts[h.now].used

	if h.rehashing {
		size += h.hts[h.new].used
	}

	return size
}

func (h *Hash) GetAll() []string {
	var all []string
	all = []string{}

	for _, s := range h.hts[h.now].slots {
		if s != nil && s.state == OccupiedState {
			all = append(all, *s.str)
		}
	}

	if h.rehashing {
		for _, s := range h.hts[h.new].slots {
			if s != nil && s.state == OccupiedState {
				all = append(all, *s.str)
			}
		}
	}

	return all
}

func (h *Hash) Dump() {
	var cnt int = 0

	for i := range h.hts {
		if h.hts[i] != nil {
			cnt++
			fmt.Fprintf(os.Stdout, "HashTable #%d: size = %d, tableSize = %d\n", cnt, h.hts[i].used, h.hts[i].size)

			for j, s := range h.hts[i].slots {
				if s == nil {
					fmt.Fprintf(os.Stdout, "H%d[%d] =\n", cnt, j)
				} else {
					fmt.Fprintf(os.Stdout, "H%d[%d] = %s\n", cnt, j, s.String())
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
		var rehashIndex int

		/* find in new hashtable first */
		index, _ = h.hts[h.new].find(str)

		if index != -1 {
			return index, h.new, nil
		}

		/* find in now hashtable */
		index, _ = h.hts[h.now].find(str)

		/**
		 * a special case handler that in order to trigger a rehash
		 * we *thought* we found the target string
		 */
		if index == -1 {
			rehashIndex = h.hts[h.now].calcIndex(str)
		} else {
			rehashIndex = index
		}

		if err = h.rehash(rehashIndex); err != nil {
			return -1, -1, err
		}

		return index, h.now, nil
	}

	/* find in now hashtable */
	index, _ = h.hts[h.now].find(str)
	return index, h.now, nil
}

func (h *Hash) insert(str string) error {
	var index int
	var needRehash bool
	var dStr string

	dStr = str

	if h.rehashing {
		var oldIndex int

		oldIndex = h.hts[h.now].calcIndex(str)
		_, needRehash = h.hts[h.new].insertAndCheck(&dStr)

		/* a very special case that we need to rehash when rehashing*/
		if needRehash {
			return h.specialRehash()
		}

		if err := h.rehash(oldIndex); err != nil {
			return err
		}

		/* a very special case that we need to rehash when rehashing*/
		if needRehash {
			return h.specialRehash()
		}

		return nil
	}

	/* not rehashing, do insert, then check if we need to start a rehash */
	index, needRehash = h.hts[h.now].insertAndCheck(&dStr)
	if needRehash {
		h.normalRehashStart()
		return h.rehash(index)
	}

	return nil
}

/**
 * hashing
 */
func (h *Hash) rehash(index int) error {
	var start, end int

	if !h.rehashing {
		return errors.New("not in rehashing")
	}

	if h.hts == nil {
		return errors.New("nul h.hts")
	}

	if h.hts[h.now] == nil {
		return errors.New("nul h.hts[h.now]")
	}

	if h.hts[h.now].slots[index] == nil {
		return nil
	}

	start = index
	for h.hts[h.now].slots[start] != nil {
		start = (start + h.hts[h.now].size - 1) % h.hts[h.now].size
	}
	start = (start + 1) % h.hts[h.now].size

	end = index
	for h.hts[h.now].slots[end] != nil {
		end = (end + 1) % h.hts[h.now].size
	}

	for i := start; i != end && h.hts[h.now].slots[i] != nil; i = (i + 1) % h.hts[h.now].size {
		if h.hts[h.now].slots[i].state == OccupiedState {
			var needRehash bool

			_, needRehash = h.hts[h.new].insertAndCheck(h.hts[h.now].slots[i].str)

			/* a very special case that we need to rehash when rehashing*/
			if needRehash {
				return h.specialRehash()
			}
		}

		h.hts[h.now].cleanSlot(i)
	}

	/* remove all elements to new hashtable if load factor is less than 0.03 */
	if h.hts[h.now].getLoadFactor() <= 0.03 {
		var needRehash bool

		for i := 0; i < h.hts[h.now].size && h.hts[h.now].slots[i] != nil; i++ {
			if h.hts[h.now].slots[i].state == OccupiedState {
				_, needRehash = h.hts[h.new].insertAndCheck(h.hts[h.now].slots[i].str)

				/* a very special case that we need to rehash when rehashing*/
				if needRehash {
					return h.specialRehash()
				}
			}

			h.hts[h.now].cleanSlot(i)
		}
	}

	if h.hts[h.now].used == 0 {
		h.normalRehashFinish()
	}

	return nil
}

func (h *Hash) normalRehashStart() {
	h.rehashing = true
	h.hts[h.new] = newHashtable(4 * h.hts[h.now].used)
}

func (h *Hash) normalRehashFinish() {
	h.hts[h.now] = nil
	h.now, h.new = h.new, h.now
	h.rehashing = false
}

func (h *Hash) specialRehash() error {
	var newSize int
	var spec *hashtable

	if !h.rehashing {
		return errors.New("no in rehashing")
	}

	newSize = 4 * (h.hts[0].used + h.hts[1].used)
	if newSize > 199999 {
		return errors.New("out_of_range")
	}

	spec = newHashtable(newSize)

	for _, c := range h.hts[0].slots {
		if c != nil && c.state == OccupiedState {
			spec.insert(c.str)
		}
	}

	for _, c := range h.hts[1].slots {
		if c != nil && c.state == OccupiedState {
			spec.insert(c.str)
		}
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
