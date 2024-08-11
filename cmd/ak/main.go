package main

import (
	"github.com/autonomouskoi/akcore/exe"
	_ "github.com/autonomouskoi/trackstar"
	_ "github.com/autonomouskoi/trackstar/overlay"
	_ "github.com/autonomouskoi/trackstar/rekordboxdb"
	_ "github.com/autonomouskoi/trackstar/stagelinq"
)

func main() {
	exe.Main()
}
