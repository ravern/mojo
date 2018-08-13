package mojo_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ravernkoh/mojo"
)

func TestObjects_Assemble(t *testing.T) {
	type args struct {
		objs mojo.Objects
	}

	type rets struct {
		args []string
		err  error
	}

	tests := []struct {
		name string
		args args
		want rets
	}{
		{
			name: "Argument",
			args: args{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
			want: rets{
				args: []string{"tldr", "nmap"},
			},
		},
		{
			name: "BoolFlagAndArgument",
			args: args{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectFlag{Name: "--verbose", Bool: true},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
			want: rets{
				args: []string{"tldr", "--verbose", "nmap"},
			},
		},
		{
			name: "FlagAndArgument",
			args: args{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectFlag{Name: "--level", Value: "5"},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
			want: rets{
				args: []string{"tldr", "--level", "5", "nmap"},
			},
		},
		{
			name: "MultipleFlagsAndArgument",
			args: args{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectFlag{Name: "-v", Bool: true, MultipleFlagsStart: true},
					mojo.ObjectFlag{Name: "-b", Bool: true},
					mojo.ObjectFlag{Name: "-l", Value: "5", MultipleFlagsEnd: true},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
			want: rets{
				args: []string{"tldr", "-vbl", "5", "nmap"},
			},
		},
		{
			name: "CombinedMultipleFlagValuesAndArgument",
			args: args{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectFlag{Name: "-v", Bool: true, MultipleFlagsStart: true},
					mojo.ObjectFlag{Name: "-l", Value: "5", MultipleFlagsEnd: true, CombinedFlagValues: true},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
			want: rets{
				args: []string{"tldr", "-vl=5", "nmap"},
			},
		},
		{
			name: "SubcommandAndArgument",
			args: args{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectCommand{Name: "add"},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
			want: rets{
				args: []string{"tldr", "add", "nmap"},
			},
		},
		{
			name: "SubcommandOrArgument",
			args: args{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
			want: rets{
				args: []string{"tldr", "nmap"},
			},
		},
		{
			name: "SubcommandAndFlagAndArgument",
			args: args{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectCommand{Name: "add"},
					mojo.ObjectFlag{Name: "--level", Value: "5"},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
			want: rets{
				args: []string{"tldr", "add", "--level", "5", "nmap"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got rets
			got.args, got.err = test.args.objs.Assemble()
			if fmt.Sprintf("%v", got.err) != fmt.Sprintf("%v", test.want.err) {
				t.Errorf("want err %v, got err %v", test.want.err, got.err)
				return
			}
			if !reflect.DeepEqual(got.args, test.want.args) {
				t.Errorf("want objs %v, got objs %v", test.want.args, got.args)
			}
		})
	}
}
