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
			name: "ErrUnconfiguredFlag",
			args: args{
				conf: mojo.Config{
					Root: mojo.CommandConfig{
						Name: "tldr",
					},
				},
				args: []string{"tldr", "-v"},
			},
			want: rets{
				err: fmt.Errorf("mojo: unconfigured flag: -v"),
			},
		},
		{
			name: "ErrInvalidFlag",
			args: args{
				conf: mojo.Config{
					Root: mojo.CommandConfig{
						Name: "tldr",
						Flags: []mojo.FlagConfig{
							{Name: "-v"},
						},
					},
				},
				args: []string{"tldr", "-v"},
			},
			want: rets{
				err: fmt.Errorf("mojo: invalid flag: -v"),
			},
		},
		{
			name: "DisallowDoubleDash",
			args: args{
				conf: mojo.Config{
					DisallowDoubleDash: true,
					Root: mojo.CommandConfig{
						Name: "tldr",
					},
				},
				args: []string{"tldr", "--", "nmap"},
			},
			want: rets{
				err: fmt.Errorf("mojo: invalid flag: --"),
			},
		},
		{
			name: "Argument",
			args: args{
				conf: mojo.Config{
					Root: mojo.CommandConfig{
						Name: "tldr",
					},
				},
				args: []string{"tldr", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.ArgumentObject{Value: "nmap"},
				},
			},
		},
		{
			name: "BoolFlagAndArgument",
			args: args{
				conf: mojo.Config{
					Root: mojo.CommandConfig{
						Name: "tldr",
						Flags: []mojo.FlagConfig{
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
					mojo.CommandObject{Name: "tldr"},
					mojo.FlagObject{Name: "--verbose", Bool: true},
					mojo.ArgumentObject{Value: "nmap"},
				},
			},
		},
		{
			name: "FlagAndArgument",
			args: args{
				conf: mojo.Config{
					Root: mojo.CommandConfig{
						Name: "tldr",
						Flags: []mojo.FlagConfig{
							{Name: "--level"},
						},
					},
				},
				args: []string{"tldr", "--level", "5", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.FlagObject{Name: "--level", Value: "5"},
					mojo.ArgumentObject{Value: "nmap"},
				},
			},
		},
		{
			name: "MultipleFlagsAndArgument",
			args: args{
				conf: mojo.Config{
					AllowMutipleFlags:      true,
					AllowUnconfiguredFlags: true,
					Root: mojo.CommandConfig{
						Name: "tldr",
					},
				},
				args: []string{"tldr", "-vbl", "5", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.FlagObject{Name: "-v", Bool: true, MultipleFlagsStart: true},
					mojo.FlagObject{Name: "-b", Bool: true},
					mojo.FlagObject{Name: "-l", Value: "5", MultipleFlagsEnd: true},
					mojo.ArgumentObject{Value: "nmap"},
				},
			},
		},
		{
			name: "CombinedMultipleFlagValuesAndArgument",
			args: args{
				conf: mojo.Config{
					AllowMutipleFlags:      true,
					AllowUnconfiguredFlags: true,
					Root: mojo.CommandConfig{
						Name: "tldr",
					},
				},
				args: []string{"tldr", "-vl=5", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.FlagObject{Name: "-v", Bool: true, MultipleFlagsStart: true},
					mojo.FlagObject{Name: "-l", Value: "5", MultipleFlagsEnd: true, CombinedFlagValues: true},
					mojo.ArgumentObject{Value: "nmap"},
				},
			},
		},
		{
			name: "SubcommandAndArgument",
			args: args{
				conf: mojo.Config{
					Root: mojo.CommandConfig{
						Name: "tldr",
						Commands: []mojo.CommandConfig{
							{Name: "add"},
						},
					},
				},
				args: []string{"tldr", "add", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.CommandObject{Name: "add"},
					mojo.ArgumentObject{Value: "nmap"},
				},
			},
		},
		{
			name: "SubcommandOrArgument",
			args: args{
				conf: mojo.Config{
					Root: mojo.CommandConfig{
						Name: "tldr",
						Commands: []mojo.CommandConfig{
							{Name: "add"},
						},
					},
				},
				args: []string{"tldr", "nmap"},
			},
			want: rets{
				objs: []mojo.Object{
					mojo.CommandObject{Name: "tldr"},
					mojo.ArgumentObject{Value: "nmap"},
				},
			},
		},
		{
			name: "SubcommandAndFlagAndArgument",
			args: args{
				conf: mojo.Config{
					Root: mojo.CommandConfig{
						Name: "tldr",
						Commands: []mojo.CommandConfig{
							{
								Name: "add",
								Flags: []mojo.FlagConfig{
									{Name: "--level"},
								},
							},
						},
						Flags: []mojo.FlagConfig{
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
					mojo.CommandObject{Name: "tldr"},
					mojo.CommandObject{Name: "add"},
					mojo.FlagObject{Name: "--level", Value: "5"},
					mojo.ArgumentObject{Value: "nmap"},
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
