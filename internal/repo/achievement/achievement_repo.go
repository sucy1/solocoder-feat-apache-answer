/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package achievement

import (
	"context"

	"github.com/apache/answer/internal/base/data"
	"github.com/apache/answer/internal/base/pager"
	"github.com/apache/answer/internal/base/reason"
	"github.com/apache/answer/internal/entity"
	"github.com/segmentfault/pacman/errors"
	"xorm.io/xorm"
)

// AchievementRepo achievement repository interface
type AchievementRepo interface {
	AddAchievement(ctx context.Context, achievement *entity.UserAchievement) error
	GetUserAchievementList(ctx context.Context, userID string, page, pageSize int) (achievements []*entity.UserAchievement, total int64, err error)
	GetUserBadgeCount(ctx context.Context, userID string) (int64, error)
	GetUserTotalReputation(ctx context.Context, userID string) (int, error)
	HasUserBadge(ctx context.Context, userID, badgeID string) (bool, error)
	AddReputation(ctx context.Context, userID string, source string, reputation int, description string) error
	AwardBadge(ctx context.Context, userID, badgeID string, source string) error
	CheckConsecutiveLoginDays(ctx context.Context, userID string) (int, error)
	UpdateConsecutiveLoginDays(ctx context.Context, userID string, days int) error
}

// achievementRepo achievement repository
type achievementRepo struct {
	data *data.Data
}

// NewAchievementRepo new achievement repository
func NewAchievementRepo(data *data.Data) AchievementRepo {
	return &achievementRepo{
		data: data,
	}
}

// AddAchievement add user achievement
func (ar *achievementRepo) AddAchievement(ctx context.Context, achievement *entity.UserAchievement) error {
	_, err := ar.data.DB.Context(ctx).Insert(achievement)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

// GetUserAchievementList get user achievement list
func (ar *achievementRepo) GetUserAchievementList(ctx context.Context, userID string, page, pageSize int) (
	achievements []*entity.UserAchievement, total int64, err error) {
	achievements = make([]*entity.UserAchievement, 0)
	session := ar.data.DB.Context(ctx).Where("user_id = ?", userID).OrderBy("created_at DESC")
	if page == 0 || pageSize == 0 {
		err = session.Find(&achievements)
	} else {
		total, err = pager.Help(page, pageSize, &achievements, &entity.UserAchievement{}, session)
	}
	if err != nil {
		return nil, 0, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetUserBadgeCount get user badge count
func (ar *achievementRepo) GetUserBadgeCount(ctx context.Context, userID string) (int64, error) {
	count, err := ar.data.DB.Context(ctx).
		Where("user_id = ? AND achievement_type = ?", userID, entity.AchievementTypeBadge).
		Count(&entity.UserAchievement{})
	if err != nil {
		return 0, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return count, nil
}

// GetUserTotalReputation get user total reputation
func (ar *achievementRepo) GetUserTotalReputation(ctx context.Context, userID string) (int, error) {
	session := ar.data.DB.Context(ctx).
		Select("COALESCE(SUM(reputation), 0) as total_reputation").
		Where("user_id = ? AND achievement_type = ?", userID, entity.AchievementTypeReputation)

	var total int
	_, err := session.Table(&entity.UserAchievement{}).Get(&total)
	if err != nil {
		return 0, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return total, nil
}

// HasUserBadge check if user has badge
func (ar *achievementRepo) HasUserBadge(ctx context.Context, userID, badgeID string) (bool, error) {
	exist, err := ar.data.DB.Context(ctx).
		Where("user_id = ? AND achievement_type = ? AND achievement_id = ?",
			userID, entity.AchievementTypeBadge, badgeID).
		Exist(&entity.UserAchievement{})
	if err != nil {
		return false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return exist, nil
}

// AddReputation add reputation to user
func (ar *achievementRepo) AddReputation(ctx context.Context, userID string, source string, reputation int, description string) error {
	achievement := &entity.UserAchievement{
		UserID:          userID,
		AchievementType: entity.AchievementTypeReputation,
		Source:          source,
		Reputation:      reputation,
		Description:     description,
	}

	_, err := ar.data.DB.Transaction(func(session *xorm.Session) (any, error) {
		session = session.Context(ctx)

		_, err := session.Insert(achievement)
		if err != nil {
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		_, err = session.Where("id = ?", userID).Incr("reputation", reputation).Update(&entity.User{})
		if err != nil {
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		return nil, nil
	})
	return err
}

// AwardBadge award badge to user
func (ar *achievementRepo) AwardBadge(ctx context.Context, userID, badgeID string, source string) error {
	_, err := ar.data.DB.Transaction(func(session *xorm.Session) (any, error) {
		session = session.Context(ctx)

		badgeCount, err := session.
			Where("user_id = ? AND achievement_type = ?", userID, entity.AchievementTypeBadge).
			Count(&entity.UserAchievement{})
		if err != nil {
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		if badgeCount >= entity.MaxUserBadges {
			return nil, errors.InternalServer("error.achievement.badge_count_exceed").WithStack()
		}

		achievement := &entity.UserAchievement{
			UserID:          userID,
			AchievementType: entity.AchievementTypeBadge,
			AchievementID:   badgeID,
			Source:          source,
		}

		_, err = session.Insert(achievement)
		if err != nil {
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		return nil, nil
	})
	return err
}

// CheckConsecutiveLoginDays check user consecutive login days
func (ar *achievementRepo) CheckConsecutiveLoginDays(ctx context.Context, userID string) (int, error) {
	user := &entity.User{}
	exist, err := ar.data.DB.Context(ctx).Where("id = ?", userID).Cols("consecutive_login_days").Get(user)
	if err != nil {
		return 0, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if !exist {
		return 0, nil
	}
	return user.ConsecutiveLoginDays, nil
}

// UpdateConsecutiveLoginDays update user consecutive login days
func (ar *achievementRepo) UpdateConsecutiveLoginDays(ctx context.Context, userID string, days int) error {
	_, err := ar.data.DB.Context(ctx).
		Where("id = ?", userID).
		Cols("consecutive_login_days").
		Update(&entity.User{ConsecutiveLoginDays: days})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}
