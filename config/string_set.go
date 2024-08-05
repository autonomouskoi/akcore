package config

import (
	"gopkg.in/yaml.v3"

	"github.com/autonomouskoi/mapset"
)

type StringSet struct {
	mapset.MapSet[string]
}

func (s *StringSet) UnmarshalYAML(value *yaml.Node) error {
	ss := []string{}
	if err := value.Decode(&ss); err != nil {
		return err
	}
	s.MapSet = mapset.From(ss...)
	return nil
}
