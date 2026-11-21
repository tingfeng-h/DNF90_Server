package mysql

import (
	"context"

	"longheng.io/server/internal/modules/dnf/repository"
)

// readOnlyCachedCharacterRepository 在同一个只读 REPEATABLE READ 事务内复用
// withinBuffOwnedSnapshot 已经加载过的 CharacterRecord，避免 apply 回调立刻再发一次
// 重复的 SELECT dnf_characters / stats / flags。
type readOnlyCachedCharacterRepository struct {
	repository.CharacterRepository
	key    string
	record repository.CharacterRecord
	found  bool
}

func (r *readOnlyCachedCharacterRepository) Load(ctx context.Context, key string) (repository.CharacterRecord, bool, error) {
	if key == r.key {
		return repository.CloneCharacter(r.record), r.found, nil
	}
	return r.CharacterRepository.Load(ctx, key)
}

// readOnlyCachedSettingsRepository 在同一个只读 REPEATABLE READ 事务内复用
// withinBuffOwnedSnapshot 已经预读的 SettingsRecord，避免 apply 回调立刻重复查表。
type readOnlyCachedSettingsRepository struct {
	repository.SettingsRepository
	key    string
	record repository.SettingsRecord
	found  bool
}

func (r *readOnlyCachedSettingsRepository) Load(ctx context.Context, key string) (repository.SettingsRecord, bool, error) {
	if key == r.key {
		return r.record, r.found, nil
	}
	return r.SettingsRepository.Load(ctx, key)
}
