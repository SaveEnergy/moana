package store

import (
	"errors"
	"io"
	"testing"
)

type fakeExecResult struct {
	n   int64
	err error
}

func (fakeExecResult) LastInsertId() (int64, error) { return 0, nil }

func (f fakeExecResult) RowsAffected() (int64, error) { return f.n, f.err }

func TestExecExactlyOneRow(t *testing.T) {
	t.Parallel()
	errZero := errors.New("zero")
	if err := execExactlyOneRow(fakeExecResult{n: 1}, errZero); err != nil {
		t.Fatalf("one row: %v", err)
	}
	if err := execExactlyOneRow(fakeExecResult{n: 0}, errZero); !errors.Is(err, errZero) {
		t.Fatalf("zero rows: got %v want %v", err, errZero)
	}
	if err := execExactlyOneRow(fakeExecResult{n: 2}, errZero); err == nil {
		t.Fatal("two rows: got nil, want error")
	}
	if err := execExactlyOneRow(fakeExecResult{n: 0, err: io.ErrUnexpectedEOF}, errZero); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatalf("rows affected error: got %v", err)
	}
}
