package recycle_bin

import (
	"shortlink/internal/app/recycle_bin/command"
	"shortlink/internal/app/recycle_bin/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	SaveToRecycleBin      command.SaveToRecycleBinHandler
	RemoveFromRecycleBin  command.RemoveFromRecycleBinHandler
	RecoverFromRecycleBin command.RecoverFromRecycleBinHandler
}

type Queries struct {
	PageDisabledLink query.PageRecycleBinHandler
}
