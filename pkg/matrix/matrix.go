package matrix

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMatrixColDiff   = errors.New("matrix: each column length on matrix is different")
	ErrIndexOutOfBound = errors.New("matrix: index out of bound")
	ErrNilMatrix       = errors.New("matrix: matrix is nil")
	ErrColRowDiff      = errors.New("matrix: col and row are different")
	ErrDimensionDiff   = errors.New("matrix: dimensions are different")
	ErrNotSquareMatrix = errors.New("matrix: matrix is not square")
	ErrZeroDeterminant = errors.New("matrix: cannot inverse, determinant is zero")
)

// matrix store matrix data and error
type matrix struct {
	row, col int
	values   [][]float64
	err      error
}

//<editor-fold desc="constructor">

// Of construct empty matrix with given row & col
func Of(row, col int) *matrix {
	values := make([][]float64, row)

	for i := range values {
		values[i] = make([]float64, col)
	}

	return &matrix{
		row:    row,
		col:    col,
		values: values,
		err:    nil,
	}
}

// From construct matrix from given slices
// returns matrix with ErrMatrixColDiff whenever the columns are different
func From(values [][]float64) *matrix {
	row := len(values)
	rows := make([][]float64, row)

	col := len(values[0])

	for i := range rows {
		currentCol := len(values[i])

		if currentCol != col {
			return errMatrix(ErrMatrixColDiff)
		}

		rows[i] = values[i]
	}

	return &matrix{
		row:    row,
		col:    col,
		values: values,
		err:    nil,
	}
}

// errMatrix construct empty matrix with given error
func errMatrix(err error) *matrix {
	return &matrix{
		err: err,
	}
}

// Identity will construct square identity matrix with given row
func Identity(row int) *matrix {
	result := Of(row, row)

	for i := 0; i < row; i++ {
		for j := 0; j < row; j++ {
			var value float64

			if i == j {
				value = 1
			}

			result.Set(i, j, value)
		}
	}

	return result
}

//</editor-fold>

//<editor-fold desc="attribute">
// String print matrix form cleanly
func (m *matrix) String() string {
	if m.HasErr() {
		return fmt.Sprintf("{[ matrix has error > %v ]}", m.Err().Error())
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf("{ %vx%v: ", m.Row(), m.Col()))
	b.WriteString("[\n")

	for i := range m.values {
		b.WriteString("  ")
		for j := range m.GetRow(i) {
			b.WriteString(fmt.Sprintf("%10.5f, ", m.Get(i, j)))
		}
		b.WriteString("\n")
	}
	b.WriteString("]}")
	return b.String()
}

// Err return error held by matrix
func (m *matrix) Err() error {
	if m == nil {
		*m = *errMatrix(ErrNilMatrix)
		return ErrNilMatrix
	}
	return m.err
}

// setErr build error on matrix
func (m *matrix) setErr(err error) *matrix {
	if m == nil {
		return errMatrix(ErrNilMatrix)
	}

	m.err = err
	return m
}

// HasErr return true if matrix has error
func (m *matrix) HasErr() bool {
	return m.Err() != nil
}

// Row return matrix's row
func (m *matrix) Row() int {
	if m.HasErr() {
		return 0
	}

	return m.row
}

// Col return matrix's col
func (m *matrix) Col() int {
	if m.HasErr() {
		return 0
	}

	return m.col
}

// Get return value at given index
func (m *matrix) Get(row, col int) float64 {
	if m.HasErr() {
		return 0
	}

	if row > m.Row() || row < 0 || col > m.Col() || col < 0 {
		m.setErr(ErrIndexOutOfBound)
		return 0
	}

	return m.values[row][col]
}

// GetRow return row slice at given row
func (m *matrix) GetRow(row int) []float64 {
	if m.HasErr() {
		return []float64{}
	}

	if row > m.Row() || row < 0 {
		m.setErr(ErrIndexOutOfBound)
		return []float64{}
	}

	return m.values[row][:]
}

// Set build value at given index
func (m *matrix) Set(row, col int, value float64) *matrix {
	if m.HasErr() {
		return m
	}

	if row > m.Row() || row < 0 || col > m.Col() || col < 0 {
		return m.setErr(ErrIndexOutOfBound)
	}

	m.values[row][col] = value

	return m
}

// SetRow build row slice at given row
func (m *matrix) SetRow(row int, values []float64) *matrix {
	if m.HasErr() {
		return m
	}

	if row > m.Row() || row < 0 || len(values) > m.Col() {
		return m.setErr(ErrIndexOutOfBound)
	}

	m.values[row] = values[:]

	return m
}

//</editor-fold>

//<editor-fold desc="operation">

// DotProduct return new matrix as the dot product result
func (m *matrix) DotProduct(other *matrix) *matrix {
	if m.HasErr() {
		return m
	}

	if other.HasErr() {
		return other
	}

	if m.Col() != other.Row() {
		return errMatrix(ErrColRowDiff)
	}

	result := Of(m.Row(), other.Col())

	for i := 0; i < result.Row(); i++ {
		for j := 0; j < result.Col(); j++ {
			var value float64

			for k := 0; k < m.Col(); k++ {
				value += m.Get(i, k) * other.Get(k, j)
			}

			if result.Set(i, j, value).HasErr() {
				return result
			}
		}
	}

	return result
}

// Add return new matrix as the Addition result
func (m *matrix) Add(other *matrix) *matrix {
	if m.HasErr() {
		return m
	}

	if other.HasErr() {
		return other
	}

	if m.Row() != m.Row() && m.Col() != other.Col() {
		return errMatrix(ErrDimensionDiff)
	}

	result := Of(m.Row(), m.Col())

	for i := range result.values {
		for j := range result.GetRow(i) {
			value := m.Get(i, j) + other.Get(i, j)
			if result.Set(i, j, value).HasErr() {
				return result
			}
		}
	}

	return result
}

// Subtract return new matrix as the Subtraction result
func (m *matrix) Subtract(other *matrix) *matrix {
	if m.HasErr() {
		return m
	}

	if other.HasErr() {
		return other
	}

	if m.Row() != m.Row() && m.Col() != other.Col() {
		return errMatrix(ErrDimensionDiff)
	}

	result := Of(m.Row(), m.Col())

	for i := range result.values {
		for j := range result.GetRow(i) {
			value := m.Get(i, j) - other.Get(i, j)
			if result.Set(i, j, value).HasErr() {
				return result
			}
		}
	}

	return result
}

// Transpose return new matrix as the transpose result
func (m *matrix) Transpose() *matrix {
	if m.HasErr() {
		return m
	}

	result := Of(m.Col(), m.Row())

	for i := range m.values {
		for j := range m.GetRow(i) {
			result.Set(j, i, m.Get(i, j))
		}
	}

	return result
}

// Determinant return determinant value from matrix
func (m *matrix) Determinant() float64 {
	if m.HasErr() {
		return 0
	}

	if m.Col() != m.Row() || m.Row() == 0 {
		m.setErr(ErrNotSquareMatrix)
		return 0
	}

	switch m.Row() {
	case 1:
		return m.Get(0, 0)
	case 2:
		return m.determinant2()
	}

	return m.determinant()
}

// DeterminantFromCofactor return determinant value from matrix & its cofactor
func (m *matrix) DeterminantFromCofactor(cofactor *matrix) float64 {
	if m.HasErr() || cofactor.HasErr() {
		return 0
	}

	if m.Col() != m.Row() || m.Row() == 0 || cofactor.Row() == 0 || cofactor.Col() != cofactor.Row() {
		m.setErr(ErrNotSquareMatrix)
		return 0
	}

	var result float64

	for i := 0; i < m.Col(); i++ {
		result += m.Get(0, i) * cofactor.Get(0, i)
	}

	return result
}

// Minor return new matrix as the minor result
func (m *matrix) Minor() *matrix {
	length := m.Row() * m.Col()

	result := Of(m.Row(), m.Col())

	for i := 0; i < length; i++ {
		colMod := i % m.Row()
		rowMod := i / m.Row()

		subSlice := make([][]float64, 0)
		for j := range m.values {
			rows := make([]float64, 0)
			for k := range m.GetRow(j) {
				if colMod != k && j != rowMod {
					rows = append(rows, m.Get(j, k))
				}
			}

			if len(rows) != 0 {
				subSlice = append(subSlice, rows)
			}
		}

		result.Set(rowMod, colMod, From(subSlice).Determinant())
	}

	return result
}

// Cofactor return new matrix as the cofactor result
func (m *matrix) Cofactor() *matrix {
	multiplier := 1.0

	result := Of(m.Row(), m.Col())

	for i := range m.values {
		for j := range m.GetRow(i) {
			result.Set(i, j, m.Get(i, j)*multiplier)
			multiplier *= -1
		}

		if m.Col()%2 == 0 {
			multiplier *= -1
		}
	}

	return result
}

// Inverse return new matrix as the inverse result
func (m *matrix) Inverse() *matrix {
	if m.HasErr() {
		return m
	}

	if m.Col() != m.Row() || m.Row() == 0 {
		return m.setErr(ErrNotSquareMatrix)
	}

	minor := m.Minor()
	cofactor := minor.Cofactor()
	determinant := m.DeterminantFromCofactor(cofactor)
	if determinant == 0 {
		return m.setErr(ErrZeroDeterminant)
	}

	adJoint := cofactor.Transpose()

	return adJoint.inverse(1 / determinant)
}

// Flatten return new matrix as the flatten result
func (m *matrix) Flatten() *matrix {
	if m.HasErr() {
		return m
	}

	var rows []float64
	for i := range m.values {
		rows = append(rows, m.GetRow(i)...)
	}

	return From([][]float64{rows})
}

// IsEqual compare the value of matrices
func (m *matrix) IsEqual(other *matrix) bool {
	if m.HasErr() || other.HasErr() {
		return false
	}

	if m.Row() != other.Row() && m.Col() != other.Col() {
		return false
	}

	for i := range m.values {
		for j := range m.GetRow(i) {
			if other.Get(i, j) != m.Get(i, j) {
				return false
			}
		}
	}

	return true
}

//</editor-fold>

//<editor-fold desc="private method">
// determinant2 return determinant value from matrix 2x2
func (m *matrix) determinant2() float64 {
	ad := m.Get(0, 0) * m.Get(1, 1)
	bc := m.Get(0, 1) * m.Get(1, 0)

	return ad - bc
}

// determinant return determinant value from matrix > 2x2
func (m *matrix) determinant() float64 {
	header := m.GetRow(0)
	body := m.values[1:][:]

	var result float64

	for i := range header {
		subSlice := make([][]float64, 0)
		for j := range body {
			rows := make([]float64, 0)
			for k := range body[j] {
				if i != k {
					// get body that not in a row of header
					rows = append(rows, body[j][k])
				}
			}
			subSlice = append(subSlice, rows)
		}

		det := From(subSlice).Determinant()
		if i%2 != 0 {
			det *= -1
		}

		result += det * header[i]
	}

	return result
}

// inverse return new matrix as the inverse result
func (m *matrix) inverse(determinant float64) *matrix {
	if determinant == 0 {
		return m.setErr(ErrZeroDeterminant)
	}

	result := Of(m.Row(), m.Col())

	for i := range m.values {
		for j := range m.GetRow(i) {
			value := determinant * m.Get(i, j)
			result.Set(i, j, value)
		}
	}

	return result
}

//</editor-fold>
