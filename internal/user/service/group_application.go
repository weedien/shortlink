package service

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"shortlink/internal/common/lock"
	"shortlink/internal/user/adapter"
	"shortlink/internal/user/app/group"
	"shortlink/internal/user/app/group/command"
	"shortlink/internal/user/app/group/query"
)

func NewGroupApplication(db *gorm.DB, rdb *redis.Client, linkService query.ShortLinkService) group.Application {

	repository := adapter.NewGroupRepositoryImpl(db, rdb)
	locker := lock.NewRedisLock(rdb)

	a := group.Application{
		Commands: group.Commands{
			CreateGroup: command.NewCreateGroupHandler(repository, locker),
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
