package errors

import (
	"errors"
	"fmt"
	"io/fs"
	"testing"
)

type CustomErr struct {
	e error
}

func (e CustomErr) Error() string {
	return e.e.Error()
}
func (e CustomErr) Unwrap() error {
	return errors.Unwrap(e.e)
}

func TestEqualityAndNils(t *testing.T) {
	err1 := New("1")
	err2 := New("2")
	erra := Wrap(err1, err2)
	errb := Wrap(err1, err2)
	if erra == errb {
		t.Errorf("errors are identical, but shoulnd't be")
	}

	if Wrap() != nil {
		t.Errorf("Empty wrap wasn't nil")
	}
	if Wrap(nil) != nil {
		t.Errorf("Single nil wrap wasn't nil")
	}
	if Wrap(nil, nil, nil) != nil {
		t.Errorf("Multiple nil wrap wasn't nil")
	}
	var e error
	if Wrap(e) != nil {
		t.Errorf("Single explicit nil err wasn't wrapped as nil")
	}
	if Wrap(e, e) != nil {
		t.Errorf("Multiple explicit nil err wasn't wrapped as nil")
	}
	if Wrap(e, nil) != nil {
		t.Errorf("Single explicit nil err wasn't wrapped with explicit nils as nil")
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		errs []error
	}

	var Err1 = CustomErr{e: fmt.Errorf("err1")}
	var Err2 = fmt.Errorf("err2")
	var Err3 = errors.New("err3")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "several",
			args:    args{errs: []error{Err1, Err2, Err3}},
			wantErr: true,
		},
		{
			name:    "one",
			args:    args{errs: []error{Err1}},
			wantErr: true,
		},
		{
			name:    "two",
			args:    args{errs: []error{Err2, Err3}},
			wantErr: true,
		},
		{
			name:    "none",
			args:    args{errs: nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Wrap(tt.args.errs...); (err != nil) == tt.wantErr {
				for _, e := range tt.args.errs {
					if !errors.Is(err, e) {
						t.Errorf("err %q was not wrapped in %q", e, err)
					}
					if errors.Is(err, Err1) {
						var foo = CustomErr{e: fmt.Errorf("err1")}
						if !errors.As(err, &foo) {
							t.Errorf("Err1 was not recognized as a CustomErr")
						}
					}
				}
			} else {
				t.Errorf("Wrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_echain_Error(t *testing.T) {
	type fields struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "single",
			fields: fields{err: Wrap(errors.New("a"))},
			want:   "a",
		},
		{
			name:   "single specific echain",
			fields: fields{err: Wrap(&echain{err: errors.New("b")})},
			want:   "b",
		},
		{
			name:   "none",
			fields: fields{err: Wrap(nil)},
			want:   "",
		},
		{
			name:   "two",
			fields: fields{err: Wrap(errors.New("a"), errors.New("b"))},
			want:   "a: b",
		},
		{
			name:   "five",
			fields: fields{err: Wrap(errors.New("a"), errors.New("b"), errors.New("c"), errors.New("d"), errors.New("e"))},
			want:   "a: b: c: d: e",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ec := tt.fields.err
			if ec == nil {
				return
			}
			if got := ec.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_echain_Unwrap(t *testing.T) {
	type fields struct {
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "regular error",
			fields:  fields{err: errors.New("regular")},
			wantErr: true,
		},
		{
			name:    "nil",
			fields:  fields{err: nil},
			wantErr: false,
		},
		{
			name:    "wrapped error",
			fields:  fields{err: Wrap(errors.New("a"), errors.New("b"))},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ec := &echain{err: tt.fields.err}
			if err := ec.Unwrap(); (err != nil) != tt.wantErr {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// also check for true nil
	t.Run("true nil",
		func(t *testing.T) {
			var ec = &echain{}
			if err := ec.Unwrap(); err != nil {
				t.Errorf("Unwrap() error = %v", err)
			}
		})

}

// local implementation of stlib's errors.New
func localNew(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func Test_stdlib_interaction(t *testing.T) {
	stderr := errors.New("stderr")
	lstderr := localNew("almoststderr")
	myerr := &echain{err: errors.New("myerr")}
	wrapped := Wrap(stderr, myerr, lstderr)
	if !errors.Is(wrapped, stderr) {
		t.Errorf("couldn't find %s in chain", stderr)
	}
	if !errors.Is(wrapped, lstderr) {
		t.Errorf("couldn't find %s in chain", lstderr)
	}
	if !errors.Is(wrapped, myerr) {
		t.Errorf("couldn't find %s in chain", myerr)
	}
	if !errors.As(wrapped, new(*echain)) {
		t.Errorf("couldn't find type %T in chain", new(*echain))
	}
	if !errors.As(wrapped, new(*errorString)) {
		t.Errorf("couldn't find type %T in chain", new(*errorString))
	}

	if !Is(wrapped, stderr) {
		t.Errorf("couldn't find %s in chain", stderr)
	}
	if !Is(wrapped, lstderr) {
		t.Errorf("couldn't find %s in chain", lstderr)
	}
	if !Is(wrapped, myerr) {
		t.Errorf("couldn't find %s in chain", myerr)
	}
	if !As(wrapped, new(*echain)) {
		t.Errorf("couldn't find type %T in chain", new(*echain))
	}
	if !As(wrapped, new(*errorString)) {
		t.Errorf("couldn't find type %T in chain", new(*errorString))
	}

}

type errorUncomparable struct {
	f []string
}

func (errorUncomparable) Error() string {
	return "uncomparable error"
}

func (errorUncomparable) Is(target error) bool {
	_, ok := target.(errorUncomparable)
	return ok
}

type errorT struct{ s string }

func (e errorT) Error() string { return fmt.Sprintf("errorT(%s)", e.s) }

type poser struct {
	msg string
	f   func(error) bool
}

var poserPathErr = &fs.PathError{Op: "poser"}

func (p *poser) Error() string     { return p.msg }
func (p *poser) Is(err error) bool { return p.f(err) }
func (p *poser) As(err interface{}) bool {
	switch x := err.(type) {
	case **poser:
		*x = p
	case *errorT:
		*x = errorT{"poser"}
	case **fs.PathError:
		*x = poserPathErr
	default:
		return false
	}
	return true
}
func TestIs(t *testing.T) {
	err1 := errors.New("1")
	erra := Wrap(errors.New(""), err1)
	errb := Wrap(errors.New(""), erra)

	err3 := errors.New("3")

	poser := &poser{"either 1 or 3", func(err error) bool {
		return err == err1 || err == err3
	}}

	testCases := []struct {
		err    error
		target error
		match  bool
	}{
		{nil, nil, true},
		{err1, nil, false},
		{err1, err1, true},
		{erra, err1, true},
		{errb, err1, true},
		{err1, err3, false},
		{erra, err3, false},
		{errb, err3, false},
		{poser, err1, true},
		{poser, err3, true},
		{poser, erra, false},
		{poser, errb, false},
		{errorUncomparable{}, errorUncomparable{}, true},
		{errorUncomparable{}, &errorUncomparable{}, false},
		{&errorUncomparable{}, errorUncomparable{}, true},
		{&errorUncomparable{}, &errorUncomparable{}, false},
		{errorUncomparable{}, err1, false},
		{&errorUncomparable{}, err1, false},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			if got := errors.Is(tc.err, tc.target); got != tc.match {
				t.Errorf("Is(%v, %v) = %v, want %v", tc.err, tc.target, got, tc.match)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	err1 := errors.New("1")
	err2 := errors.New("2")
	erra := Wrap(nil, err1)

	testCases := []struct {
		err  error
		want error
	}{
		{nil, nil},
		{Wrap(nil), nil},
		{err1, nil},
		{erra, err1},
		{Wrap(erra), erra},
		{Wrap(erra, err2), erra},
		{Wrap(erra, err2), err2},
	}
	for i, tc := range testCases {
		if got := errors.Unwrap(tc.err); !errors.Is(got, tc.want) {
			t.Errorf("%d: Unwrap(%v/%p) = %v/%p, want %v/%p", i, tc.err, tc.err, got, got, tc.want, tc.want)
		}
	}
}
