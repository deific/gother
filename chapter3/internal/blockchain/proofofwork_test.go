package blockchain

import (
	"fmt"
	"testing"
)

func TestBlock_GetTarget(t *testing.T) {
	b := &Block{}
	fmt.Printf("target:%x\n", b.GetTarget(0))
	fmt.Printf("target:%x\n", b.GetTarget(1))
	fmt.Printf("target:%x\n", b.GetTarget(2))
	fmt.Printf("target:%x\n", b.GetTarget(10))
	fmt.Printf("target:%x\n", b.GetTarget(12))
	fmt.Printf("target:%x\n", b.GetTarget(20))
}
