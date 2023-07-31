package coolprojectlogic

import "testing"

func TestAddition(t *testing.T) {
	m := NewMathematicks()
	if m.Sum(1, 1) != 2 {
		t.Error("1 + 1 should equal 2")
	}
}

func TestSubtraction(t *testing.T) {
	m := NewMathematicks()
	if m.Subtract(1, 1) != 0 {
		t.Error("1 - 1 should equal 0")
	}
}

func TestMultiplication(t *testing.T) {
	m := NewMathematicks()
	if m.Multiply(1, 1) != 1 {
		t.Error("1 * 1 should equal 1")
	}
}

func TestDivision(t *testing.T) {
	m := NewMathematicks()
	if m.Divide(1, 1) != 1 {
		t.Error("1 / 1 should equal 1")
	}
}
