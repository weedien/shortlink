package command

import (
	"context"
	"shortlink/internal/common/constant"
	"shortlink/internal/common/lock"
	"shortlink/internal/user/domain/group"
)

type CreateGroupCommand struct {
	Username  string
	GroupName string
}

type CreateGroupHandler struct {
	repo   group.Repository
	locker lock.DistributedLock
}

func NewCreateGroupHandler(repo group.Repository) CreateGroupHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return CreateGroupHandler{repo: repo}
}

func (h CreateGroupHandler) Handle(ctx context.Context, cmd CreateGroupCommand) (err error) {
	lockKey := constant.LockGroupCreateKey + cmd.GroupName
	if _, err = h.locker.Acquire(ctx, lockKey, -1); err != nil {
		return err
	}
	defer func() {
		_ = h.locker.Release(ctx, lockKey)
	}()
	// 最大分组数限制
	var size int
	if size, err = h.repo.GetGroupSize(ctx); err != nil {
		return err
	}
	if size >= group.MaxGroupSize {
		return group.ErrGroupSizeExceed
	}
	// 生成唯一分组ID
	retryCount, maxRetries := 0, 10
	var gid string
	for retryCount < maxRetries {
		if gid, err = h.repo.UniqueReturnGid(ctx); err != nil {
			return err
		}
		if gid != "" {
			g := group.NewGroup(gid, cmd.Username, cmd.GroupName, size)
			if err = h.repo.CreateGroup(ctx, g); err != nil {
				return err
			}
			break
		}
		retryCount++
	}
	if gid == "" {
		return group.ErrGenGroupUniqueID
	}
	return nil
}

// CreateGroup 创建分组
// 为其他领域提供服务
func (h CreateGroupHandler) CreateGroup(
	ctx context.Context,
	username, groupName string,
) error {
	return h.Handle(ctx, CreateGroupCommand{
		Username:  username,
		GroupName: groupName,
	})
}
