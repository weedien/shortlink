package service

import (
	"gorm.io/gorm"
	"shortlink/internal/user/adapter"
	"shortlink/internal/user/app/group"
	"shortlink/internal/user/app/group/command"
	"shortlink/internal/user/app/group/query"
)

func NewGroupApplication(db *gorm.DB, linkService query.ShortLinkService) group.Application {

	repository := adapter.NewGroupRepositoryImpl(db)

	a := group.Application{
		Commands: group.Commands{
			CreateGroup: command.NewCreateGroupHandler(repository),
			UpdateGroup: command.NewUpdateGroupHandler(repository),
			DeleteGroup: command.NewDeleteGroupHandler(repository),
			SortGroup:   command.NewSortGroupHandler(repository),
		},
		Queries: group.Queries{
			ListGroup: query.NewListGroupHandler(repository, linkService),
		},
	}

	return a
}
