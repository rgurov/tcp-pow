package puzzle

import (
	"crypto/rand"
	"crypto/sha256"
)

const (
	SolutionSize = 8
	PuzzleSize   = sha256.Size
)

type Puzzle struct {
	initialHash [PuzzleSize]byte
	complexity  int
}

func NewPuzzle(initialHash [PuzzleSize]byte, complexity int) *Puzzle {

	return &Puzzle{
		initialHash: initialHash,
		complexity:  complexity,
	}
}

func NewRandomPuzzle(complexity int) *Puzzle {
	randomBytes := make([]byte, PuzzleSize)
	_, _ = rand.Read(randomBytes)

	hash := sha256.New()
	_, _ = hash.Write(randomBytes)

	return &Puzzle{
		initialHash: [PuzzleSize]byte(randomBytes),
		complexity:  complexity,
	}
}

func (p *Puzzle) GetInitialHash() [PuzzleSize]byte {
	return p.initialHash
}

func (p *Puzzle) IsValidSolution(solution [SolutionSize]byte) bool {
	hash := sha256.New()
	hash.Write(p.initialHash[:])
	hash.Write(solution[:])
	sum := hash.Sum(nil)

	leadingZeros := 0
	for i := range sum {
		if (sum[i] >> 4) == 0 {
			leadingZeros++
		} else {
			break
		}

		if sum[i]&0x0f == 0 {
			leadingZeros++
		} else {
			break
		}
	}

	return leadingZeros >= p.complexity
}
