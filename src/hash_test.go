package src

import (
	"testing"
)

type FindCase struct {
	str    string
	wanted bool
}

type RemoveCase struct {
	str, wanted string
}

type IndexCase struct {
	str   string
	index int
}

type FindCases []FindCase
type RemoveCases []RemoveCase
type IndexCases []IndexCase

func checkFind(ht *Hash, cases FindCases, t *testing.T) {
	var ok bool
	var err error

	for _, c := range cases {
		if ok, err = ht.Find(c.str); err != nil {
			t.Fatal(err)
		}
		if ok != c.wanted {
			if c.wanted {
				t.Error("should in hash:", c.str)
			} else {
				t.Error("should not in hash:", c.str)
			}
		}
	}
}

func checkRemove(ht *Hash, cases RemoveCases, t *testing.T) {
	var err error

	for _, c := range cases {
		if err = ht.Remove(c.str); err != nil {
			t.Fatal(err)
		}
		if ok, err := ht.Find(c.str); err != nil {
			t.Fatal(err)
		} else if ok {
			t.Fatal("should not have", c.str)
		}
	}
}

func checkIndex(ht *Hash, cases IndexCases, t *testing.T) {

}

func TestHashBasic(t *testing.T) {
	var ht *Hash
	var err error

	ht, err = NewHash(107)
	if err != nil {
		t.Fatal(err)
	}

	ht.Insert("tributes")      // 21
	ht.Insert("skulduggery")   // 22
	ht.Insert("convulse")      // 23
	ht.Insert("frothed")       // 24
	ht.Insert("horrify")       // 25
	ht.Insert("blackmailers")  // 26
	ht.Insert("defenestrated") // 27
	ht.Insert("garrison")      // 23 -> 28
	ht.Insert("lidless")       // 22 -> 29
	// ht.Insert("eye") ;            // 21 -> 30

	ht.Dump()

	finds := FindCases{
		{"skulduggery", true},
		{"lidless", true},
		{"defenestrated", true},
		{"peircing", false},
	}
	checkFind(ht, finds, t)

	removes := RemoveCases{
		{"infractions", ""},
		{"garrison", "garrison"},
	}
	checkRemove(ht, removes, t)

	finds = FindCases{
		{"garrison", false},
		{"lidless", true},
	}
	checkFind(ht, finds, t)

	ht.Insert("undying")

	finds = FindCases{
		{"undying", true},
	}
	checkFind(ht, finds, t)
}
