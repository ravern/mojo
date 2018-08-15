package mojo_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ravernkoh/mojo"
)

func TestObjects_Flags(t *testing.T) {
	type args struct {
		objs mojo.Objects
		name string
	}

	type rets struct {
		objs []mojo.FlagObject
	}

	tests := []struct {
		name string
		args args
		want rets
	}{
		{
			name: "None",
			args: args{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.ArgumentObject{Value: "nmap"},
				},
				name: "-v",
			},
			want: rets{},
		},
		{
			name: "Many",
			args: args{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.FlagObject{Name: "-v", Value: "a"},
					mojo.FlagObject{Name: "-v", Value: "b"},
				},
				name: "-v",
			},
			want: rets{
				objs: []mojo.FlagObject{
					{Name: "-v", Value: "a"},
					{Name: "-v", Value: "b"},
				},
			},
		},
		{
			name: "One",
			args: args{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.FlagObject{Name: "--verbose", Bool: true},
				},
				name: "--verbose",
			},
			want: rets{
				objs: []mojo.FlagObject{
					{Name: "--verbose", Bool: true},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got rets
			got.objs = test.args.objs.Flags(test.args.name)
			if !reflect.DeepEqual(got.objs, test.want.objs) {
				t.Errorf("want objs %v, got objs %v", test.want.objs, got.objs)
			}
		})
	}
}

func TestObjects_Flag(t *testing.T) {
	type args struct {
		objs mojo.Objects
		name string
	}

	type rets struct {
		obj mojo.FlagObject
		err error
	}

	tests := []struct {
		name string
		args args
		want rets
	}{
		{
			name: "ErrFlagNotFound",
			args: args{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.ArgumentObject{Value: "nmap"},
				},
				name: "-v",
			},
			want: rets{
				err: fmt.Errorf("mojo: flag not found: -v"),
			},
		},
		{
			name: "ErrTooManyFlags",
			args: args{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.FlagObject{Name: "-v", Bool: true},
					mojo.FlagObject{Name: "-v", Bool: true},
				},
				name: "-v",
			},
			want: rets{
				err: fmt.Errorf("mojo: too many flags: -v"),
			},
		},
		{
			name: "BoolFlagAndArgument",
			args: args{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.FlagObject{Name: "--verbose", Bool: true},
				},
				name: "--verbose",
			},
			want: rets{
				obj: mojo.FlagObject{Name: "--verbose", Bool: true},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got rets
			got.obj, got.err = test.args.objs.Flag(test.args.name)
			if fmt.Sprintf("%v", got.err) != fmt.Sprintf("%v", test.want.err) {
				t.Errorf("want err %v, got err %v", test.want.err, got.err)
				return
			}
			if !reflect.DeepEqual(got.obj, test.want.obj) {
				t.Errorf("want obj %v, got obj %v", test.want.obj, got.obj)
			}
		})
	}
}

func TestObjects_Argument(t *testing.T) {
	type args struct {
		objs mojo.Objects
		i    int
	}

	type rets struct {
		obj mojo.ArgumentObject
		err error
	}

	tests := []struct {
		name string
		args args
		want rets
	}{
		{
			name: "ErrArgumentNotFound",
			args: args{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.ArgumentObject{Value: "nmap"},
				},
				i: 1,
			},
			want: rets{
				err: fmt.Errorf("mojo: argument not found: 1"),
			},
		},
		{
			name: "BoolFlagAndArgument",
			args: args{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.FlagObject{Name: "--verbose", Bool: true},
					mojo.ArgumentObject{Value: "netstat"},
					mojo.ArgumentObject{Value: "nmap"},
				},
				i: 1,
			},
			want: rets{
				obj: mojo.ArgumentObject{Value: "nmap"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got rets
			got.obj, got.err = test.args.objs.Argument(test.args.i)
			if fmt.Sprintf("%v", got.err) != fmt.Sprintf("%v", test.want.err) {
				t.Errorf("want err %v, got err %v", test.want.err, got.err)
				return
			}
			if !reflect.DeepEqual(got.obj, test.want.obj) {
				t.Errorf("want obj %v, got obj %v", test.want.obj, got.obj)
			}
		})
	}
}
