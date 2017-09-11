package main

import "fmt"

func testcase3() {
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
	ht.Insert("eye")           // 21 -> 30, should trigger inc. rehash

	dumper(ht, "Hash table after rehash triggered")

	fmt.Println("Do some insertions to make long linear probe sequence in the second table.")

	ht.Insert("wildcat")        // 18 (new table)
	ht.Insert("weightlessness") // 69 (new table)
	ht.Insert("sorceress")      // 33 (new table)
	ht.Insert("enchantress")    // 33 (new table) really.
	ht.Insert("witchery")       // 67 -> 68 (new table)
	ht.Insert("heliosphere")    // 67 -> 72 (new table)
	ht.Insert("obstruct")       // 67 -> 73 (new table)

	dumper(ht, "Hash table insertions.")

	fmt.Println("One more insertion in slot 67 should make us give up on rehashing.")

	ht.Insert("peripatetic") // 67 -> 77 (new table)

	dumper(ht, "Hash table giving up on rehashing.")
}
