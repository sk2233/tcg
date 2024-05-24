/*
@author: sk
@date: 2023/2/5
*/
package main

import (
	"fmt"
	"testing"
)

func Test_CheckNum(t *testing.T) {
	c1, c2 := GetSetoKaibaCards()
	fmt.Println(len(c1), len(c2))
}
