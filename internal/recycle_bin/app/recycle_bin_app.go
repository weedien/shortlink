package app

import (
	"shortlink/internal/recycle_bin/app/command"
	"shortlink/internal/recycle_bin/app/query"
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
