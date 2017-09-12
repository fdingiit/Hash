package main

import (
	"fmt"
	"math/rand"
	"time"
)

func testcase4() {
	var ht *Hash
	var err error
	var htCnt, sCnt int
	var standard map[string]struct{}

	standard = make(map[string]struct{})

	ht, err = NewHash(883)
	if err != nil {
		panic(err)
	}

	for i := 100; i < 1000; i++ {
		err = ht.Insert(words[i])
		if err != nil {
			panic(err)
		}

		standard[words[i]] = struct{}{}
	}

	for i := 400; i < 500; i++ {
		ht.Remove(words[i])
		delete(standard, words[i])
	}

	for i := 30000; i < 40000; i++ {
		ht.Insert(words[i])
		standard[words[i]] = struct{}{}
	}

	for i := 28000; i < 32000; i++ {
		ht.Remove(words[i])
		delete(standard, words[i])
	}

	rand.Seed(time.Now().Unix())

	for i := 0; i < 10000; i++ {
		r := rand.Int() % numWords

		if ok, err := ht.Find(words[r]); err != nil {
			panic(err)
		} else if ok {
			htCnt++
		}

		if _, ok := standard[words[r]]; ok {
			sCnt++
		}
	}

	if htCnt == sCnt {
		fmt.Println("Passed random find() tests:", "Tcount =", htCnt, ",", "Scount =", sCnt)
	} else {
		fmt.Println("***Failed random find() tests:", "Tcount =", htCnt, ",", "Scount =", sCnt)
	}

	sanityCheck(ht, standard)
}

func sanityCheck(ht *Hash, standard map[string]struct{}) {
	var all []string
	var T map[string]struct{}
	var inSnotT, inTnotS int

	if ht.Size() == len(standard) {
		fmt.Println("Sets sizes are both")
	} else {
		fmt.Println("Sets sizes are different:", "T size =", ht.Size(), ",", "S.size() =", len(standard))
	}

	all = ht.GetAll()

	T = make(map[string]struct{})
	for _, s := range all {
		T[s] = struct{}{}
	}

	for k := range T {
		if _, ok := standard[k]; !ok {
			inTnotS++
		}
	}

	for k := range standard {
		if _, ok := T[k]; !ok {
			inSnotT++
		}
	}

	if inSnotT == 0 && inTnotS == 0 {
		fmt.Println("Passed set equality test")
	} else {
		fmt.Println("***Failed set equality test")
		fmt.Println("   ", inSnotT, "words in set S but not in hash table T.")
		fmt.Println("   ", inTnotS, "words in hash table T but not in set S.")
	}
}
