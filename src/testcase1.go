package main

import "fmt"

func testcase1() {
	var ht *Hash
	var err error

	ht, err = NewHash(107)
	if err != nil {
		panic(err)
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

	dumper(ht, "Original hash table")

	fmt.Println("Do some find() and remove() operations...")

	testFind(ht, "skulduggery")
	testFind(ht, "lidless")
	testFind(ht, "defenestrated")
	testFind(ht, "peircing")

	testRemove(ht, "garrison")
	testRemove(ht, "infractions")

	testFind(ht, "garrison")
	testFind(ht, "lidless")

	dumper(ht, "Hash table after finds and removes")

	fmt.Println("Next insert should reuse DELETED slots...")

	ht.Insert("undying") // 25 -> 28

	dumper(ht, "Final hash table")
}
