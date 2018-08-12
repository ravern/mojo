package mojo_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ravernkoh/mojo"
)

func TestParse(t *testing.T) {
	type args struct {
		conf mojo.Config
		args []string
	}

	type rets struct {
		objs mojo.Objects
		err  error
	}

	tests := []struct {
		name string
		args args
		want rets
	}{
		{
			name: "OneArgument",
			args: args{
				conf: mojo.Config{
					Root: mojo.ConfigCommand{
						Name: "tldr",
					},
				},
				args: []string{"tldr", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got rets
			got.objs, got.err = mojo.Parse(test.args.conf, test.args.args)
			if fmt.Sprintf("%v", got.err) != fmt.Sprintf("%v", test.want.err) {
				t.Errorf("got err %v, want err %v", got.err, test.want.err)
				return
			}
			if !reflect.DeepEqual(got.objs, test.want.objs) {
				t.Errorf("got objs %v, want objs %v", got.objs, test.want.objs)
			}
		})
	}
}
