package link

import (
	"shortlink/internal/app/link/command"
	"shortlink/internal/app/link/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateLink      command.CreateLinkHandler
	CreateLinkBatch command.CreateLinkBatchHandler
	UpdateLink      command.UpdateLinkHandler
}

type Queries struct {
	PageLink       query.PageLinkHandler
	ListGroupCount query.ListGroupCountHandler
	GetOriginalUrl query.GetOriginalUrlHandler
}
