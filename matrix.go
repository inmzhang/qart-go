package qart_go

import (
	"reflect"
)

func QRMatrix(version int, ecl int, bits string) (matrix [][]int) {
	size := (version-1)*4 + 21
	matrix = makeSquareMatrix(size)
	addFinderAndSeparator(matrix)
	addAlignment(version, matrix)
	addTiming(matrix)
	addDarkAndReserving(version, matrix)
	maskMatrix := deepCopy(matrix)
	placeBits(bits, matrix)
	var numMask int
	numMask, matrix = mask(maskMatrix, matrix)
	addFormatInfo(ecl, numMask, matrix)
	addVersionInfo(version, matrix)
	return
}

func addFinderAndSeparator(matrix [][]int) {
	size := len(matrix)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			val := 0
			switch i {
			case 0, 6:
				if j != 7 {
					val = 1
				}
			case 1, 5:
				if j == 0 || j == 6 {
					val = 1
				}
			case 7:
				val = 0
			default:
				if j != 1 && j != 5 && j != 7 {
					val = 1
				}
			}
			matrix[i][j], matrix[i][size-j-1], matrix[size-i-1][j] = val, val, val
		}
	}
}

func addAlignment(ver int, matrix [][]int) {
	if ver < 2 {
		return
	}

	coords := AlignLocation[ver-2]
	for _, row := range coords {
		for _, col := range coords {
			if matrix[row][col] == -1 {
				addAnAlignment(row, col, matrix)
			}
		}
	}
}

func addAnAlignment(row, col int, matrix [][]int) {
	for i := row - 2; i <= row+2; i++ {
		for j := col - 2; j <= col+2; j++ {
			if i == row-2 || i == row+2 || j == col-2 || j == col+2 {
				matrix[i][j] = 1
			} else {
				matrix[i][j] = 0
			}
		}
	}
	matrix[row][col] = 1
}

func addTiming(matrix [][]int) {
	for i := 8; i < len(matrix)-8; i++ {
		var val int
		if i%2 == 0 {
			val = 1
		} else {
			val = 0
		}
		matrix[i][6] = val
		matrix[6][i] = val
	}
}

func addDarkAndReserving(ver int, matrix [][]int) {
	size := len(matrix)
	for j := 0; j < 8; j++ {
		matrix[8][j], matrix[8][size-j-1], matrix[j][8], matrix[size-j-1][8] = 0, 0, 0, 0
	}
	matrix[8][8] = 0
	matrix[8][6], matrix[6][8], matrix[size-8][8] = 1, 1, 1

	if ver > 6 {
		for i := 0; i < 6; i++ {
			for _, j := range []int{-9, -10, -11} {
				matrix[i][size+j], matrix[size+j][i] = 0, 0
			}
		}
	}
}

func placeBits(bits string, m [][]int) {
	up := true
	bitIdx := 0
	for a := len(m) - 1; a > 0; a -= 2 {
		val := a
		if val <= 6 {
			val--
		}
		var rangeI []int
		if up {
			for i := len(m) - 1; i >= 0; i-- {
				rangeI = append(rangeI, i)
			}
		} else {
			for i := 0; i < len(m); i++ {
				rangeI = append(rangeI, i)
			}
		}
		for _, i := range rangeI {
			for _, j := range []int{val, val - 1} {
				if m[i][j] == -1 {
					m[i][j] = int(bits[bitIdx] - '0')
					bitIdx++
				}
			}
		}
		up = !up
	}
}

func mask(mm [][]int, matrix [][]int) (int, [][]int) {
	mps := getMaskPatterns(mm)
	scores := make([]int, len(mps))
	for i, mp := range mps {
		for j, row := range mp {
			for k, val := range row {
				mps[i][j][k] = val ^ matrix[j][k]
			}
		}
		scores[i] = computeScore(mp)
	}
	best := 0
	for i, score := range scores {
		if score < scores[best] {
			best = i
		}
	}
	return best, mps[best]
}

func getMaskPatterns(mm [][]int) [][][]int {
	mm[len(mm)-8][8] = -1
	for i := range mm {
		for j := range mm {
			if mm[i][j] != -1 {
				mm[i][j] = 0
			}
		}
	}
	mps := make([][][]int, 8)
	for i := 0; i < 8; i++ {
		mp := make([][]int, len(mm))
		for j := range mm {
			mp[j] = make([]int, len(mm[j]))
			copy(mp[j], mm[j])
		}
		for row := range mp {
			for column := range mp {
				if mp[row][column] == -1 && formula(i, row, column) {
					mp[row][column] = 1
				} else {
					mp[row][column] = 0
				}
			}
		}
		mps[i] = mp
	}
	return mps
}

func computeScore(matrix [][]int) (score int) {
	// evaluation1: count the number of consecutive 1s or 0s in each row and column
	ev1 := func(ma [][]int) (sc int) {
		for _, mi := range ma {
			j := 0
			for j < len(mi)-4 {
				n := 4
				for isIn(mi[j:j+n+1], [][]int{makeOnes(n + 1), make([]int, n+1)}) {
					n++
					if j+n+1 > len(mi) {
						break
					}
				}
				if n > 4 {
					sc += n - 2
					j += n
				} else {
					j++
				}
			}
		}
		return
	}
	transposed := transpose(matrix)
	score += ev1(matrix) + ev1(transposed)

	// evaluation2: count the number of 2x2 squares of the same color
	size := len(matrix)
	for i := 0; i < size-1; i++ {
		for j := 0; j < size-1; j++ {
			ele := matrix[i][j]
			if ele == matrix[i+1][j] && ele == matrix[i][j+1] && ele == matrix[i+1][j+1] {
				score += 3
			}
		}
	}

	// evaluation3: count the number of special patterns
	ev3 := func(ma [][]int) (sc int) {
		for _, mi := range ma {
			j := 0
			for j < len(mi)-10 {
				rowSlice := mi[j : j+11]
				if isIn(rowSlice, [][]int{{1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0}}) {
					sc += 40
					j += 7
				} else if isIn(rowSlice, [][]int{{0, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1}}) {
					sc += 40
					j += 4
				} else {
					j++
				}
			}
		}
		return
	}
	score += ev3(matrix) + ev3(transposed)

	// evaluation4: count the percentage of dark cells and apply penalty
	numDark := 0
	for _, row := range matrix {
		for _, cell := range row {
			numDark += cell
		}
	}
	percent := float64(numDark) / float64(size*size) * 100
	s := int((50-percent)/5) * 5
	if s >= 0 {
		score += 2 * s
	} else {
		score -= 2 * s
	}

	return
}

func addFormatInfo(ecl int, numMask int, matrix [][]int) {
	var formats []int
	for _, format := range FormatInfo[ecl][numMask] {
		formats = append(formats, int(format-'0'))
	}
	numRows := len(matrix)
	numCols := len(matrix[0])
	for j := 0; j < 6; j++ {
		matrix[8][j] = formats[j]
		matrix[numRows-1-j][8] = formats[j]
		matrix[8][numCols-1-j] = formats[len(formats)-1-j]
		matrix[j][8] = formats[len(formats)-1-j]
	}
	matrix[8][7] = formats[6]
	matrix[numRows-7][8] = formats[6]
	matrix[8][8] = formats[7]
	matrix[8][numCols-8] = formats[7]
	matrix[7][8] = formats[8]
	matrix[8][numCols-7] = formats[8]
}

func addVersionInfo(version int, matrix [][]int) {
	if version < 7 {
		return
	}
	var versionInfo []int
	for _, v := range VersionInfo[version-7] {
		versionInfo = append(versionInfo, int(v-'0'))
	}
	numRows := len(matrix)
	numCols := len(matrix[0])
	var count int
	for j := 5; j > -1; j-- {
		for _, i := range []int{-9, -10, -11} {
			matrix[numRows+i][j] = versionInfo[count]
			matrix[j][numCols+i] = versionInfo[count]
			count++
		}
	}
}

func isIn(s []int, lst [][]int) bool {
	for _, l := range lst {
		if reflect.DeepEqual(s, l) {
			return true
		}
	}
	return false
}

func makeOnes(n int) (ones []int) {
	for i := 0; i < n; i++ {
		ones = append(ones, 1)
	}
	return
}

func transpose(m [][]int) (transposed [][]int) {
	n := len(m)
	transposed = make([][]int, n)
	for i := 0; i < n; i++ {
		transposed[i] = make([]int, n)
		for j := 0; j < n; j++ {
			transposed[i][j] = m[j][i]
		}
	}
	return transposed
}

func formula(i, row, column int) bool {
	switch i {
	case 0:
		return (row+column)%2 == 0
	case 1:
		return row%2 == 0
	case 2:
		return column%3 == 0
	case 3:
		return (row+column)%3 == 0
	case 4:
		return (row/2+column/3)%2 == 0
	case 5:
		return (row*column)%2+(row*column)%3 == 0
	case 6:
		return ((row*column)%2+(row*column)%3)%2 == 0
	case 7:
		return ((row+column)%2+(row*column)%3)%2 == 0
	default:
		return false
	}
}

func makeSquareMatrix(size int) [][]int {
	matrix := make([][]int, size)
	for i := 0; i < size; i++ {
		matrix[i] = make([]int, size)
		for j := 0; j < size; j++ {
			matrix[i][j] = -1
		}
	}
	return matrix
}

func deepCopy(src [][]int) (dst [][]int) {
	dst = make([][]int, len(src))
	for i := range src {
		dst[i] = make([]int, len(src[i]))
		copy(dst[i], src[i])
	}
	return
}
