package coolprojectlogic

type MathematicksInterface interface {
	Sum(a, b int) int
	Multiply(a, b int) int
	Divide(a, b int) int
	Subtract(a, b int) int
}

type Mathematicks struct{}

func NewMathematicks() *Mathematicks {
	return &Mathematicks{}
}

func (M *Mathematicks) Sum(a, b int) int {
	return a + b
}

func (M *Mathematicks) Multiply(a, b int) int {
	return a * b
}

func (M *Mathematicks) Divide(a, b int) int {
	return a / b
}

func (M *Mathematicks) Subtract(a, b int) int {
	return a - b
}
