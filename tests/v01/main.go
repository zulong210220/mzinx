/*
 * File : main.go
 * CreateDate : 2019-12-07 15:53:28
 * */

package main

import "mzinx/znet"

func main() {
	s := znet.NewServer("[zinx v01]")
	s.Serve()
}

/* vim: set tabstop=4 set shiftwidth=4 */
