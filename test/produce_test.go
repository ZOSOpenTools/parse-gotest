package main

import (
	"fmt"
	"testing"
)

func TestPass1(t *testing.T) {
	// pass
}

func TestPass2(t *testing.T) {
	// pass
}

func TestFail1(t *testing.T) {
	fmt.Println("fail1 message")
	t.Fatal("Fail1 message")
}

func TestFail2(t *testing.T) {
	fmt.Println("fail2 message")
	t.Fatal("Fail2 message")
}

func TestSkip(t *testing.T) {
	t.Skip("skipping test")
}
