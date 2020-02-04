/*
*
* copyright - Bōryoku-tekina soro - 2020
* makiko -  malagasy cryptocurrency
 */

package main

import (
	"fmt"
	"log"

	"github.com/boryoku-tekina/makiko/chain"
	"github.com/boryoku-tekina/makiko/tests"
)

func main() {
	chain.InitChain()
	tests.Donate("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7", 100)
	tests.Donate("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7", 50)

	tests.GetBalanceOf("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7")
	tests.Donate("1PrZapno38xz6g7ZHzwtxb3SM3uKUw8EE6", 70)
	tests.Donate("1PrZapno38xz6g7ZHzwtxb3SM3uKUw8EE6", 20)

	tests.GetBalanceOf("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7")
	tests.GetBalanceOf("1PrZapno38xz6g7ZHzwtxb3SM3uKUw8EE6")

	// log.Panic("look at 1st test : 150, 90")

	// PASSED 150, 90

	// Second test
	// send 70 coin from 1 to 2
	// amount of 1 must be 150 - 70 = 80
	// amount of 2 must be 90 + 70 = 160
	tests.Send("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7", "1PrZapno38xz6g7ZHzwtxb3SM3uKUw8EE6", 70)
	tests.GetBalanceOf("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7")
	tests.GetBalanceOf("1PrZapno38xz6g7ZHzwtxb3SM3uKUw8EE6")
	// log.Panic("look at 2nd test : 80, 160")

	// PASSED 80, 160
	fmt.Println("STEP2 PASSEDXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")

	// THIRD TEST
	// must be 90, 150

	tests.Send("1PrZapno38xz6g7ZHzwtxb3SM3uKUw8EE6", "1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7", 159)

	tests.GetBalanceOf("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7")
	tests.GetBalanceOf("1PrZapno38xz6g7ZHzwtxb3SM3uKUw8EE6")

	// log.Panic("look at 3rd  : 239, 1")

	// PASSED

	// 4th TEST
	// must be 100, 140

	tests.Send("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7", "1PrZapno38xz6g7ZHzwtxb3SM3uKUw8EE6", 9)

	tests.GetBalanceOf("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7")
	tests.GetBalanceOf("1PrZapno38xz6g7ZHzwtxb3SM3uKUw8EE6")

	// log.Panic("look at 4th test : 230, 10")

	// PASSED

	// 4th TEST
	// must be 200, 200

	tests.Donate("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7", 100)
	tests.Donate("1PrZapno38xz6g7ZHzwtxb3SM3uKUw8EE6", 60)

	tests.GetBalanceOf("1KHaWQQ3GHmWN2d417YbtA3L6v65b11Ya7")
	tests.GetBalanceOf("1PrZapno38xz6g7ZHzwtxb3SM3uKUw8EE6")

	log.Panic("look at 5th test : 330, 70")
}
