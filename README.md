Random String Generator
=======================

This package helps to generates strings based on the given regular expressions and limit


Usage
=====

```go
package main

import (
	"fmt"

	"github.com/befika/strGen"
)

func main() {
	// generate a single string
	str, err := reggen.Generate("[-+]?[0-9]{1,16}[.][0-9]{1,6}", 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(str)

	// create a reusable generator
	g, err := reggen.NewGenerator("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{8}")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 5; i++ {
		// 10 is the maximum number of times star, range or plus should repeat
		// i.e. [0-9]+ will generate at most 10 characters if this is set to 10
		fmt.Println(g.Generate(10))
	}
}
```
```go
//Sample out put 1

pattern:= "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{8}"
/*
	0e3ec3cf-fc85-e4e5-0e61f52b
	db200804-d124-29a6-5554e9d8
	927480e8-f2ed-8ab6-6ef94518
	1fdece34-7eac-e85e-de8609b7
	137c9b74-42ba-4552-1c04909e
*/

//Sample out put 2

pattern:= "[-+]?[0-9]{1,16}[.][0-9]{1,6}"
/*
	+3962810.775
	-87.637623
	+412841024.72902
	80326766.193
*/


//Sample out put 3

pattern:= "(1[0-2]|0[1-9])(:[0-5][0-9]){2} (A|P)M"
/*
	10:30:55 AM
	11:08:12 AM
	11:55:58 AM
	10:58:38 PM
	08:01:47 AM
*/
