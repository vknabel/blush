package fsmodule_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/google/go-cmp/cmp"
	"github.com/vknabel/lithia/registry"
	"github.com/vknabel/lithia/registry/fsmodule"
)

type (
	testCase struct {
		name  string
		cwd   string
		base  registry.LogicalURI
		files map[string]string
		want  []testWant
	}
	testWant struct {
		uri     registry.LogicalURI
		sources map[registry.LogicalURI]string
	}
)

func TestDiscoderFSModules(t *testing.T) {
	fs := memfs.New()

	tests := []testCase{
		{
			name: "basic test",
			cwd:  "/github.com/vknabel/lithia-example",
			base: "memory:///github.com/vknabel/lithia-example",
			files: map[string]string{
				"/github.com/vknabel/lithia-example/Potfile":          "module potfile",
				"/github.com/vknabel/lithia-example/tools/fmt.lithia": "module tools",
				"/github.com/vknabel/lithia-example/cmd/main.lithia":  "module cmd",
				"/github.com/vknabel/lithia-example/app/root.lithia":  "module app",

				"/github.com/vknabel/lithia-example/app/views/body.lithia": "module views",
			},
			want: []testWant{
				{
					uri: "memory:///github.com/vknabel/lithia-example/app",
					sources: map[registry.LogicalURI]string{
						"memory:///github.com/vknabel/lithia-example/app/root.lithia": "module app",
					},
				},
				{
					uri: "memory:///github.com/vknabel/lithia-example/app/views",
					sources: map[registry.LogicalURI]string{
						"memory:///github.com/vknabel/lithia-example/app/views/body.lithia": "module views",
					},
				},
				{
					uri: "memory:///github.com/vknabel/lithia-example/cmd",
					sources: map[registry.LogicalURI]string{
						"memory:///github.com/vknabel/lithia-example/cmd/main.lithia": "module cmd",
					},
				},
				{
					uri: "memory:///github.com/vknabel/lithia-example/tools",
					sources: map[registry.LogicalURI]string{
						"memory:///github.com/vknabel/lithia-example/tools/fmt.lithia": "module tools",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for path, cont := range tt.files {
				base := filepath.Base(path)
				err := fs.MkdirAll(base, 0644)
				if err != nil {
					t.Error(err)
				}

				file, err := fs.Create(path)
				if err != nil {
					t.Error(err)
				}
				_, err = file.Write([]byte(cont))
				if err != nil {
					t.Error(err)
				}
				err = file.Close()
				if err != nil {
					t.Error(err)
				}
			}

			pkgfs, err := fs.Chroot(tt.cwd)
			if err != nil {
				t.Error(err)
			}

			mods, err := fsmodule.DiscoverModules(tt.base, pkgfs)
			if err != nil {
				t.Error(err)
			}

			if len(mods) != len(tt.want) {
				t.Errorf("want %d modules, got %d (%s)", len(tt.want), len(mods), mods)
			}

			for i, mod := range mods {
				want := tt.want[i]
				t.Run(fmt.Sprintf("mod %d.", i), func(t *testing.T) {
					if mod.URI() != want.uri {
						t.Errorf("got uri %q, want %q", mod.URI(), want.uri)
					}
					srcs, err := mod.Sources()
					if err != nil {
						t.Error(err)
					}
					if len(srcs) != len(want.sources) {
						t.Errorf("app should have %d file, got %d", len(want.sources), len(srcs))
					}

					for j, src := range srcs {
						t.Run(fmt.Sprintf("src %d.", j), func(t *testing.T) {
							wantsrc, ok := want.sources[src.URI()]
							if !ok {
								t.Errorf("unknown source %q", src.URI())
							}

							contents, err := src.Read()
							if err != nil {
								t.Error(err)
							}
							diff := cmp.Diff(string(contents), wantsrc)
							if diff != "" {
								t.Error(diff)
							}
						})
					}
				})
			}
		})
	}
}
