package matrix

import (
	"reflect"
	"testing"
)

func TestOf(t *testing.T) {
	row := 5
	col := 4
	m := Of(row, col)

	t.Run("check column", func(t *testing.T) {
		if m.Col() != col {
			t.Errorf("expected: %d, got: %d", col, m.Col())
		}
	})

	t.Run("check error", func(t *testing.T) {
		got := m.HasErr()

		if got != false {

			t.Errorf("expected: %v, got: %v", false, got)
		}
	})

	t.Run("check row", func(t *testing.T) {
		if m.Row() != row {
			t.Errorf("expected: %d, got: %d", row, m.Row())
		}
	})

	t.Run("check value", func(t *testing.T) {
		length := len(m.values)
		if length != row {
			t.Errorf("expected: %d, got: %d", row, length)
		}

		for _, rows := range m.values {
			length := len(rows)
			if length != col {
				t.Errorf("expected: %d, got: %d", col, length)
			}

			values := []float64{0, 0, 0, 0}
			if !reflect.DeepEqual(rows, values) {
				t.Errorf("expected: %v, got: %v", values, rows)
			}
		}
	})
}

func BenchmarkOf(b *testing.B) {
	b.Run("1x1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Of(1, 1)
		}
	})

	b.Run("2x2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Of(2, 2)
		}
	})

	b.Run("3x3", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Of(3, 3)
		}
	})

	b.Run("10x10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Of(10, 10)
		}
	})

	b.Run("25x25", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Of(25, 25)
		}
	})
}

func TestFrom(t *testing.T) {
	t.Run("check valid", func(t *testing.T) {
		row := 2
		col := 5
		values := [][]float64{
			{1, 2, 3, 4, 5},
			{6, 7, 8, 9, 0},
		}

		m := From(values[:][:])

		t.Run("check column", func(t *testing.T) {
			if m.Col() != col {
				t.Errorf("expected: %d, got: %d", col, m.Col())
			}
		})

		t.Run("check error", func(t *testing.T) {
			got := m.HasErr()

			if got != false {

				t.Errorf("expected: %v, got: %v", false, got)
			}
		})

		t.Run("check row", func(t *testing.T) {
			if m.Row() != row {
				t.Errorf("expected: %d, got: %d", row, m.Row())
			}
		})

		t.Run("check value", func(t *testing.T) {
			length := len(m.values)
			if length != row {
				t.Errorf("expected: %d, got: %d", row, length)
			}

			for _, rows := range m.values {
				length := len(rows)
				if length != col {
					t.Errorf("expected: %d, got: %d", col, length)
				}
			}

			expected0 := []float64{1, 2, 3, 4, 5}
			got0 := m.GetRow(0)

			expected1 := []float64{6, 7, 8, 9, 0}
			got1 := m.GetRow(1)

			if !reflect.DeepEqual(got0, expected0) {
				t.Errorf("expected: %v, got: %v", expected0, got0)
			}

			if !reflect.DeepEqual(got1, expected1) {
				t.Errorf("expected: %v, got: %v", expected1, got1)
			}
		})
	})

	t.Run("check not valid", func(t *testing.T) {
		values := [][]float64{
			{1, 2, 3, 4, 5},
			{6, 7, 8, 9},
		}

		m := From(values[:][:])

		t.Run("check column", func(t *testing.T) {
			if m.Col() != 0 {
				t.Errorf("expected: %d, got: %d", 0, m.Col())
			}
		})

		t.Run("check error", func(t *testing.T) {
			got := m.HasErr()

			if got != true {
				t.Errorf("expected: %v, got: %v", true, got)
			}

			if m.Err() != ErrMatrixColDiff {
				t.Errorf("expected: %v, got: %v", ErrMatrixColDiff, m.Err())
			}
		})

		t.Run("check row", func(t *testing.T) {
			if m.Row() != 0 {
				t.Errorf("expected: %d, got: %d", 0, m.Row())
			}
		})

		t.Run("check value", func(t *testing.T) {
			length := len(m.values)
			if length != 0 {
				t.Errorf("expected: %d, got: %d", 0, length)
			}

			for _, rows := range m.values {
				length := len(rows)
				if length != 0 {
					t.Errorf("expected: %d, got: %d", 0, length)
				}
			}
		})
	})
}

func BenchmarkFrom(b *testing.B) {
	b.Run("1x2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			From([][]float64{
				{1, 2},
			})
		}
	})
	b.Run("3x2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			From([][]float64{
				{1, 2},
				{1, 2},
				{1, 2},
			})
		}
	})
	b.Run("5x10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			From([][]float64{
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			})
		}
	})
	b.Run("8x7", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			From([][]float64{
				{1, 2, 3, 4, 5, 6, 7},
				{1, 2, 3, 4, 5, 6, 7},
				{1, 2, 3, 4, 5, 6, 7},
				{1, 2, 3, 4, 5, 6, 7},
				{1, 2, 3, 4, 5, 6, 7},
				{1, 2, 3, 4, 5, 6, 7},
				{1, 2, 3, 4, 5, 6, 7},
				{1, 2, 3, 4, 5, 6, 7},
			})
		}
	})
}

func TestIdentity(t *testing.T) {

}

func BenchmarkIdentity(b *testing.B) {

}
