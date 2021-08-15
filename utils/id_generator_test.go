/*
 * File : id_generator_test.go
 * CreateDate : 2019-12-15 21:38:03
 * */

package utils

import (
	"fmt"
	"testing"
)

func TestNextId(t *testing.T) {
	t.Log(NextId())
	fmt.Println(NextId())
}

func BenchmarkNextId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Println(NextId())
	}
}

/* vim: set tabstop=4 set shiftwidth=4 */
