package puzzle

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrNotSolved = errors.New("puzzle not solved")
)

type PuzzleSolver struct {
	puzzle     *Puzzle
	iterations uint64
	isSolved   bool
	solution   [SolutionSize]byte
}

func NewPuzzleSolver(puzzle *Puzzle) *PuzzleSolver {
	return &PuzzleSolver{
		puzzle: puzzle,
	}
}

func (p *PuzzleSolver) Solve() bool {
	if p.isSolved {
		return true
	}

	solution := make([]byte, SolutionSize)
	binary.LittleEndian.PutUint64(solution, p.iterations)

	if p.puzzle.IsValidSolution([SolutionSize]byte(solution)) {
		fmt.Println(p.iterations)
		p.isSolved = true
		p.solution = [SolutionSize]byte(solution)

		return true
	}

	p.iterations++
	return false
}

func (p *PuzzleSolver) GetSolution() ([SolutionSize]byte, error) {
	if !p.isSolved {
		return [SolutionSize]byte{}, ErrNotSolved
	}

	return p.solution, nil
}
