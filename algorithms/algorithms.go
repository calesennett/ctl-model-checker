package algorithms

import (
	"fmt"
	"github.com/calesennett/ctl-model-checker/fsm"
	"github.com/robertkrimen/otto"
	mtx "github.com/skelterjohn/go.matrix"
	"regexp"
	"strconv"
	"strings"
)

func Run(sm fsm.StateMachine, comps []string) {
	E := sm.ToMatrix()
	elems := make(map[int]float64)
	h0 := mtx.MakeSparseMatrix(elems, 1, len(sm.States))
	for _, comp := range comps {
		tokens := tokenize(comp)
		for i, tok := range tokens {
			expr := regexp.MustCompile("EX [a-z]+|EG [a-z]+|E [a-z]+ until [a-z]+")
			if expr.MatchString(tok) {
				label := strings.Split(tok, " ")[1]
				for _, s := range sm.States {
					if s.HasLabel(label) {
						h0.Set(1, s.ID, 1)
					} else {
						h0.Set(1, s.ID, 0)
					}
				}
				if strings.Split(tok, " ")[0] == "EG" {
					result := Global(h0, h0, E)
					fmt.Println(tok + ":")
					fmt.Println(result)
					res := strconv.Itoa(sm.Satisfies(result))
					tokens[i] = res
				} else if strings.Split(tok, " ")[0] == "EX" {
					result := Next(h0, E)
					fmt.Println(tok + ":")
					fmt.Println(result)
					res := strconv.Itoa(sm.Satisfies(result))
					tokens[i] = res
				} else if strings.Split(tok, " ")[0] == "E" {
					labelg := strings.Split(tok, " ")[1]
					labelf := strings.Split(tok, " ")[3]
					g := mtx.MakeSparseMatrix(make(map[int]float64), 1, len(sm.States))
					f := mtx.MakeSparseMatrix(make(map[int]float64), 1, len(sm.States))
					for _, s := range sm.States {
						if s.HasLabel(labelg) {
							g.Set(1, s.ID, 1)
						}
						if s.HasLabel(labelf) {
							f.Set(1, s.ID, 1)
						}
					}
					result := Until(f, g, E)
					fmt.Println(tok + ":")
					fmt.Println(result)
					res := strconv.Itoa(sm.Satisfies(result))
					tokens[i] = res
				}
			}
		}
		// make string from expression
		expr := strings.Join(tokens, "")
		expr = regexp.MustCompile("and").ReplaceAllString(expr, "&&")
		expr = regexp.MustCompile("or").ReplaceAllString(expr, "||")
		expr = regexp.MustCompile("not").ReplaceAllString(expr, "!")
		vm := otto.New()
		vm.Set("expr", expr)
		vm.Run(`
			var res = eval(expr);
		`)
		value, _ := vm.Get("res")
		res, _ := value.ToBoolean()
		fmt.Printf("\nFSM satisfies %v?\n%v\n\n", comp, res)
	}
}

func tokenize(comp string) []string {
	r := regexp.MustCompile("\\(|\\)|or|and|not|EX [a-z]+|EG [a-z]+|E [a-z]+ until [a-z]+")
	matched := r.FindAllString(comp, -1)
	return matched
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
