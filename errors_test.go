package errors

import (
	"errors"
	"fmt"
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

func TestWrap(t *testing.T) {
	type args struct {
		errs []error
	}

	var Err1 = CustomErr{e: fmt.Errorf("err1")}
	var Err2 = fmt.Errorf("err2")
	var Err3 = fmt.Errorf("err3")

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
			var ec *echain
			if err := ec.Unwrap(); err != nil {
				t.Errorf("Unwrap() error = %v", err)
			}
		})

}
