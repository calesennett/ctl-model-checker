package algorithms

import (
	"fmt"
	"github.com/calesennett/ctl-model-checker/fsm"
	mtx "github.com/skelterjohn/go.matrix"
	"regexp"
)

func Run(sm fsm.StateMachine, comps []string) {
	//E := sm.ToMatrix()
	//elems := make(map[int]float64)
	//h0 := mtx.MakeSparseMatrix(elems, 1, len(sm.States))
	for _, comp := range comps {
		tokens := tokenize(comp)
		fmt.Println(tokens)
		//label := strings.Split(comp, " ")[1]
		//for _, s := range sm.States {
		//	if s.HasLabel(label) {
		//		h0.Set(1, s.ID, 1)
		//	} else {
		//		h0.Set(1, s.ID, 0)
		//	}
		//}
		//if strings.Split(comp, " ")[0] == "EG" {
		//	result := Global(h0, h0, E)
		//	fmt.Println(comp + ":")
		//	fmt.Println(result)
		//} else if strings.Split(comp, " ")[0] == "EX" {
		//	result := Next(h0, E)
		//	fmt.Println(comp + ":")
		//	fmt.Println(result)
		//} else if strings.Split(comp, " ")[0] == "E" {
		//	labelg := strings.Split(comp, " ")[1]
		//	labelf := strings.Split(comp, " ")[3]
		//	g := mtx.MakeSparseMatrix(make(map[int]float64), 1, len(sm.States))
		//	f := mtx.MakeSparseMatrix(make(map[int]float64), 1, len(sm.States))
		//	for _, s := range sm.States {
		//		if s.HasLabel(labelg) {
		//			g.Set(1, s.ID, 1)
		//		}
		//		if s.HasLabel(labelf) {
		//			f.Set(1, s.ID, 1)
		//		}
		//	}
		//	result := Until(f, g, E)
		//	fmt.Println(comp + ":")
		//	fmt.Println(result)
		//}
	}
}

func tokenize(comp string) []string {
	r := regexp.MustCompile("\\(|\\)|or|and|not|EX [a-z]+|EG [a-z]+|E [a-z]+ until [a-z]+")
	matched := r.FindAllString(comp, -1)
	for _, match := range matched {
		fmt.Println(match)
	}
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
