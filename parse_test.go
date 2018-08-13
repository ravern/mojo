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
			name: "Argument",
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
		{
			name: "BoolFlagAndArgument",
			args: args{
				conf: mojo.Config{
					Root: mojo.ConfigCommand{
						Name: "tldr",
						Flags: []mojo.ConfigFlag{
							{
								Name: "--verbose",
								Bool: true,
							},
						},
					},
				},
				args: []string{"tldr", "--verbose", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectFlag{Name: "--verbose", Bool: true},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
		},
		{
			name: "FlagAndArgument",
			args: args{
				conf: mojo.Config{
					Root: mojo.ConfigCommand{
						Name: "tldr",
						Flags: []mojo.ConfigFlag{
							{Name: "--level"},
						},
					},
				},
				args: []string{"tldr", "--level", "5", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectFlag{Name: "--level", Value: "5"},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
		},
		{
			name: "MultipleFlagsAndArgument",
			args: args{
				conf: mojo.Config{
					AllowMutipleFlags:      true,
					AllowUnconfiguredFlags: true,
					Root: mojo.ConfigCommand{
						Name: "tldr",
					},
				},
				args: []string{"tldr", "-vbl", "5", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectFlag{Name: "-v", Bool: true, MultipleFlagsStart: true},
					mojo.ObjectFlag{Name: "-b", Bool: true},
					mojo.ObjectFlag{Name: "-l", Value: "5", MultipleFlagsEnd: true},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
		},
		{
			name: "CombinedMultipleFlagValuesAndArgument",
			args: args{
				conf: mojo.Config{
					AllowMutipleFlags:      true,
					AllowUnconfiguredFlags: true,
					Root: mojo.ConfigCommand{
						Name: "tldr",
					},
				},
				args: []string{"tldr", "-vl=5", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectFlag{Name: "-v", Bool: true, MultipleFlagsStart: true},
					mojo.ObjectFlag{Name: "-l", Value: "5", MultipleFlagsEnd: true, CombinedFlagValues: true},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
		},
		{
			name: "SubcommandAndArgument",
			args: args{
				conf: mojo.Config{
					Root: mojo.ConfigCommand{
						Name: "tldr",
						Commands: []mojo.ConfigCommand{
							{Name: "add"},
						},
					},
				},
				args: []string{"tldr", "add", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectCommand{Name: "add"},
					mojo.ObjectArgument{Value: "nmap"},
				},
			},
		},
		{
			name: "SubcommandOrArgument",
			args: args{
				conf: mojo.Config{
					Root: mojo.ConfigCommand{
						Name: "tldr",
						Commands: []mojo.ConfigCommand{
							{Name: "add"},
						},
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
		{
			name: "SubcommandAndFlagAndArgument",
			args: args{
				conf: mojo.Config{
					Root: mojo.ConfigCommand{
						Name: "tldr",
						Commands: []mojo.ConfigCommand{
							{
								Name: "add",
								Flags: []mojo.ConfigFlag{
									{Name: "--level"},
								},
							},
						},
						Flags: []mojo.ConfigFlag{
							{
								Name: "--level",
								Bool: true,
							},
						},
					},
				},
				args: []string{"tldr", "add", "--level", "5", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.ObjectCommand{Name: "tldr"},
					mojo.ObjectCommand{Name: "add"},
					mojo.ObjectFlag{Name: "--level", Value: "5"},
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
				t.Errorf("want err %v, got err %v", test.want.err, got.err)
				return
			}
			if !reflect.DeepEqual(got.objs, test.want.objs) {
				t.Errorf("want objs %v, got objs %v", test.want.objs, got.objs)
			}
		})
	}
}
