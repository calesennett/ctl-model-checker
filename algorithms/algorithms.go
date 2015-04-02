package algorithms

import (
	"fmt"
	"github.com/calesennett/ctl-model-checker/fsm"
	mtx "github.com/skelterjohn/go.matrix"
	"regexp"
	//"strconv"
)

func Run(sm fsm.StateMachine, comps []string) {
	E := sm.ToMatrix()
	for _, comp := range comps {
		var operands []interface{}
		var operators []interface{}
		tokens := tokenize(comp)
		for _, tok := range tokens {
			operatorExpr := regexp.MustCompile("EX|EG|until|and|not|or")
			operandExpr := regexp.MustCompile("[a-z]+")
			rightParen := regexp.MustCompile("\\)")
			if operatorExpr.MatchString(tok) {
				operators = push(tok, operators)
			} else if operandExpr.MatchString(tok) {
				operands = push(vectorize(sm, tok), operands)
			} else if rightParen.MatchString(tok) {
				var op interface{}
				op, operators = pop(operators)
				if op == "EG" {
					var h0 interface{}
					h0, operands = pop(operands)
					if h0, ok := h0.(*mtx.SparseMatrix); ok {
						result := Global(h0, h0, E)
						operands = push(result, operands)
					}
				} else if op == "EX" {
					var h0 interface{}
					h0, operands = pop(operands)
					if h0, ok := h0.(*mtx.SparseMatrix); ok {
						result := Next(h0, E)
						operands = push(result, operands)
					}
				} else if op == "until" {
					var f, g interface{}
					f, operands = pop(operands)
					g, operands = pop(operands)
					if f, ok := f.(*mtx.SparseMatrix); ok {
						if g, ok := g.(*mtx.SparseMatrix); ok {
							result := Until(f, g, E)
							operands = push(result, operands)
						}
					}
				} else if op == "or" {
					var f, g interface{}
					f, operands = pop(operands)
					g, operands = pop(operands)
					if f, ok := f.(*mtx.SparseMatrix); ok {
						if g, ok := g.(*mtx.SparseMatrix); ok {
							result := or(f, g)
							operands = push(result, operands)
						}
					}
				} else if op == "and" {
					var f, g interface{}
					f, operands = pop(operands)
					g, operands = pop(operands)
					if f, ok := f.(*mtx.SparseMatrix); ok {
						if g, ok := g.(*mtx.SparseMatrix); ok {
							result := and(f, g)
							operands = push(result, operands)
						}
					}
				} else if op == "not" {
					var f interface{}
					f, operands = pop(operands)
					if f, ok := f.(*mtx.SparseMatrix); ok {
						result := not(f)
						operands = push(result, operands)
					}
				}
			}
		}
		final := operands[0]
		if final, ok := final.(*mtx.SparseMatrix); ok {
			fmt.Println(comp)
			fmt.Printf("Resulting vector: \n%v\n", final)
			fmt.Printf("FSM satisfies?\n%v\n\n", sm.Satisfies(final))
		}
	}
}

func tokenize(comp string) []string {
	r := regexp.MustCompile("\\(|\\)|EG|EX|[a-z]+|or|and|not|until")
	matched := r.FindAllString(comp, -1)
	return matched
}

func vectorize(sm fsm.StateMachine, label string) *mtx.SparseMatrix {
	elems := make(map[int]float64)
	h0 := mtx.MakeSparseMatrix(elems, 1, len(sm.States))
	if label == "true" {
		for _, s := range sm.States {
			h0.Set(1, s.ID, 1)
		}
	} else if label == "false" {
		for _, s := range sm.States {
			h0.Set(1, s.ID, 0)
		}
	} else {
		for _, s := range sm.States {
			if s.HasLabel(label) {
				h0.Set(1, s.ID, 1)
			} else {
				h0.Set(1, s.ID, 0)
			}
		}
	}
	return h0

}

func push(elem interface{}, slice []interface{}) []interface{} {
	return append(slice, elem)
}

func pop(slice []interface{}) (interface{}, []interface{}) {
	return slice[len(slice)-1], slice[:len(slice)-1]
}

func Global(h0 *mtx.SparseMatrix, hn *mtx.SparseMatrix, E *mtx.SparseMatrix) *mtx.SparseMatrix {
	hNext, _ := hn.TimesSparse(E)
	hNext = and(hNext, h0)
	if mtx.Equals(hNext, hn) {
		return hNext
	} else {
		return Global(h0, hNext, E)
	}
}

func and(m1 *mtx.SparseMatrix, m2 *mtx.SparseMatrix) *mtx.SparseMatrix {
	m1Arr := m1.DenseMatrix().Array()
	m2Arr := m2.DenseMatrix().Array()
	m := mtx.MakeSparseMatrix(make(map[int]float64), 1, len(m1Arr))
	for i, _ := range m1Arr {
		if m1Arr[i] > 0 && m2Arr[i] > 0 {
			m.Set(1, i, 1)
		}
	}
	return m
}

func or(m1 *mtx.SparseMatrix, m2 *mtx.SparseMatrix) *mtx.SparseMatrix {
	m1Arr := m1.DenseMatrix().Array()
	m2Arr := m2.DenseMatrix().Array()
	m := mtx.MakeSparseMatrix(make(map[int]float64), 1, len(m1Arr))
	for i, _ := range m1Arr {
		if m1Arr[i] > 0 || m2Arr[i] > 0 {
			m.Set(1, i, 1)
		}
	}
	return m
}

func not(m *mtx.SparseMatrix) *mtx.SparseMatrix {
	mArr := m.DenseMatrix().Array()
	for i, _ := range mArr {
		if mArr[i] == 0 {
			m.Set(1, i, 1)
		} else if mArr[i] == 1 {
			m.Set(1, i, 0)
		}
	}
	return m
}

func Until(hn *mtx.SparseMatrix, g *mtx.SparseMatrix, E *mtx.SparseMatrix) *mtx.SparseMatrix {
	hNext, _ := hn.TimesSparse(E)
	hNext = and(hNext, g)
	hNext = or(hn, hNext)
	if mtx.Equals(hNext, hn) {
		return hNext
	} else {
		return Until(hNext, g, E)
	}
}

func Next(h0 *mtx.SparseMatrix, E *mtx.SparseMatrix) *mtx.SparseMatrix {
	result, _ := h0.TimesSparse(E)
	resArr := result.DenseMatrix().Array()
	m := mtx.MakeSparseMatrix(make(map[int]float64), 1, len(resArr))
	for i, val := range resArr {
		if val > 0 {
			m.Set(1, i, 1)
		}
	}
	return m
}
