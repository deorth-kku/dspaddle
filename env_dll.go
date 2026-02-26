//go:build dll

package main

const IsDLL = true

func init() {
	go main()
}
