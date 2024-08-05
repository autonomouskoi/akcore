package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/autonomouskoi/akcore/config"
	"github.com/autonomouskoi/mapset"
)

func TestStringSetUnmarshal(t *testing.T) {
	t.Parallel()
	requireT := require.New(t)

	out := struct {
		Foo []string `yaml:"foo"`
	}{
		Foo: []string{"a", "b", "c"},
	}
	b, err := yaml.Marshal(out)
	requireT.NoError(err, "marshalling list")

	decoded := struct {
		Foo config.StringSet `yaml:"foo"`
	}{}
	requireT.NoError(yaml.Unmarshal(b, &decoded))

	requireT.True(mapset.From("a", "b", "c").
		Equals(decoded.Foo.MapSet),
	)
}
