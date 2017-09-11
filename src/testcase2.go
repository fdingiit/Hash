package main

import "fmt"

func testcase2() {
	var ht *Hash
	var err error

	ht, err = NewHash(107)
	if err != nil {
		panic(err)
	}

	ht.Insert("undulation")  //  9
	ht.Insert("impertinent") // 10
	ht.Insert("maladies")    // 10 -> 11
	ht.Insert("dominions")   // 12

	ht.Insert("waspish")    // 52
	ht.Insert("wildebeest") // 52 -> 53
	ht.Insert("reaction")   // 52 -> 54

	ht.Insert("pawns")       // 43
	ht.Insert("vacuously")   // 43 -> 44
	ht.Insert("firth")       // 43 -> 45
	ht.Insert("manipulator") // 43 -> 46
	ht.Insert("dreariness")  // 43 -> 47

	ht.Insert("insanity")     // 105
	ht.Insert("enthronement") // 105 -> 106
	ht.Insert("falsifiers")   // 105 -> 0
	ht.Insert("ignominious")  // 105 -> 1
	ht.Insert("mummified")    // 105 -> 2

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

	fmt.Println("Inserting \"eye\" should trigger rehash")
	ht.Insert("eye")

	dumper(ht, "Hash table after rehash triggered")

	fmt.Println("Search for \"manipulator\" should move cluster in slots 43-47.")

	ht.Find("manipulator")

	dumper(ht, "Hash table after cluster 43-47 moved.")

	fmt.Println("Do some finds, inserts and removes")

	testFind(ht, "zip")
	testFind(ht, "spaceflight")
	testFind(ht, "frothed")

	ht.Insert("wildcat")
	ht.Insert("weightlessness")
	ht.Insert("sorceress")
	ht.Insert("enchantress")

	dumper(ht, "Hash table after more insertions.")

	fmt.Println("A find on \"ignominious\" + 1 more operation should cause the tables to consolidate down to one table.")

	testFind(ht, "ignominious")
	testFind(ht, "reaction")

	dumper(ht, "Hash table after wrap up.")
}
