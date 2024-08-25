package command

import (
	"context"
	"shortlink/internal/common/constant"
	"shortlink/internal/common/error_no"
	"shortlink/internal/common/lock"
	"shortlink/internal/user/domain/user"
)

type UserRegisterHandler struct {
	repo         user.Repository
	locker       lock.DistributedLock
	groupService GroupService
}

func (h UserRegisterHandler) Handle(ctx context.Context, cmd user.UserRegisterVo) error {
	if exist, err := h.repo.CheckUserExist(ctx, cmd.Username); err != nil {
		return err
	} else if exist {
		return error_no.UserExist
	}
	// 获取分布式锁
	lockKey := constant.LockUserRegisterKey + cmd.Username
	if _, err := h.locker.Acquire(ctx, lockKey, -1); err != nil {
		return error_no.LockAcquireFailed
	}
	defer func() {
		_ = h.locker.Release(ctx, lockKey)
	}()
	// 再次检查用户是否存在
	if exist, err := h.repo.CheckUserExist(ctx, cmd.Username); err != nil {
		return err
	} else if exist {
		return error_no.UserExist
	}
	// 创建用户
	user := user.NewUser(cmd.Username, cmd.Password, cmd.RealName, cmd.Email, cmd.Phone)
	if err := h.repo.CreateUser(ctx, &user); err != nil {
		return err
	}
	// 加入分组
	if err := h.groupService.CreateGroup(ctx, user.Name(), "默认分组"); err != nil {
		return err
	}
	// 加入布隆过滤器
	if err := h.repo.AddUserToBloomFilter(ctx, user.Name()); err != nil {
		return err
	}
	return nil
}
