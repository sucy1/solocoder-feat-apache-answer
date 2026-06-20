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

package message

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/apache/answer/internal/base/data"
	"github.com/apache/answer/internal/base/pager"
	"github.com/apache/answer/internal/base/reason"
	"github.com/apache/answer/internal/entity"
	"github.com/apache/answer/internal/service/unique"
	"github.com/segmentfault/pacman/errors"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type MessageRepo interface {
	AddMessage(ctx context.Context, message *entity.Message) error
	GetMessageByID(ctx context.Context, id string) (message *entity.Message, exist bool, err error)
	GetMessageList(ctx context.Context, userID, conversationID string, page, pageSize int) (messages []*entity.Message, total int64, err error)
	GetConversationList(ctx context.Context, userID string, page, pageSize int) (conversations []*entity.Message, total int64, err error)
	GetUnreadCountByConversation(ctx context.Context, userID string) (map[string]int64, error)
	GetUnreadMessageCount(ctx context.Context, userID string) (total, system, private int64, err error)
	UpdateMessageStatus(ctx context.Context, id string, status int) error
	MarkConversationAsRead(ctx context.Context, userID, conversationID string) error
	DeleteMessage(ctx context.Context, id string) error
	CleanOldMessages(ctx context.Context, days int) error
}

type MessageBlockRepo interface {
	BlockUser(ctx context.Context, block *entity.MessageBlock) error
	UnblockUser(ctx context.Context, userID, blockedUserID string) error
	IsBlocked(ctx context.Context, userID, blockedUserID string) (bool, error)
	GetBlockedUserList(ctx context.Context, userID string) ([]*entity.MessageBlock, error)
}

type messageRepo struct {
	data         *data.Data
	uniqueIDRepo unique.UniqueIDRepo
}

func NewMessageRepo(data *data.Data, uniqueIDRepo unique.UniqueIDRepo) MessageRepo {
	return &messageRepo{
		data:         data,
		uniqueIDRepo: uniqueIDRepo,
	}
}

func NewMessageBlockRepo(data *data.Data) MessageBlockRepo {
	return &messageRepo{
		data: data,
	}
}

func GetConversationID(userID1, userID2 string) string {
	ids := []string{userID1, userID2}
	sort.Strings(ids)
	return strings.Join(ids, "_")
}

func (mr *messageRepo) AddMessage(ctx context.Context, message *entity.Message) (err error) {
	message.ID, err = mr.uniqueIDRepo.GenUniqueIDStr(ctx, message.TableName())
	if err != nil {
		return err
	}
	_, err = mr.data.DB.Context(ctx).Insert(message)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (mr *messageRepo) GetMessageByID(ctx context.Context, id string) (message *entity.Message, exist bool, err error) {
	message = &entity.Message{}
	exist, err = mr.data.DB.Context(ctx).Where("id = ?", id).Get(message)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		return nil, false, err
	}
	return message, exist, nil
}

func (mr *messageRepo) GetMessageList(ctx context.Context, userID, conversationID string, page, pageSize int) (messages []*entity.Message, total int64, err error) {
	messages = make([]*entity.Message, 0)

	session := mr.data.DB.Context(ctx)
	session.Where(builder.Or(
		builder.Eq{"from_user_id": userID},
		builder.Eq{"to_user_id": userID},
	))
	session.Where(builder.Eq{"conversation_id": conversationID})
	session.Where(builder.Neq{"status": entity.MessageStatusDeleted})
	session.Desc("created_at")

	cond := &entity.Message{}
	total, err = pager.Help(page, pageSize, &messages, cond, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (mr *messageRepo) GetConversationList(ctx context.Context, userID string, page, pageSize int) (conversations []*entity.Message, total int64, err error) {
	conversations = make([]*entity.Message, 0)

	subQuery := builder.Select("MAX(id) as max_id").
		From("message").
		Where(builder.Or(
			builder.Eq{"from_user_id": userID},
			builder.Eq{"to_user_id": userID},
		)).
		Where(builder.Neq{"status": entity.MessageStatusDeleted}).
		GroupBy("conversation_id")

	session := mr.data.DB.Context(ctx)
	session.Join("INNER", subQuery, "message.id = sub.max_id")
	session.Where(builder.Or(
		builder.Eq{"from_user_id": userID},
		builder.Eq{"to_user_id": userID},
	))
	session.Where(builder.Neq{"status": entity.MessageStatusDeleted})
	session.Desc("message.created_at")

	cond := &entity.Message{}
	total, err = pager.Help(page, pageSize, &conversations, cond, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (mr *messageRepo) GetUnreadCountByConversation(ctx context.Context, userID string) (map[string]int64, error) {
	type UnreadCount struct {
		ConversationID string `xorm:"conversation_id"`
		Count          int64  `xorm:"count"`
	}

	result := make([]UnreadCount, 0)
	err := mr.data.DB.Context(ctx).
		Select("conversation_id, COUNT(*) as count").
		Table("message").
		Where(builder.Eq{"to_user_id": userID}).
		Where(builder.Eq{"status": entity.MessageStatusUnread}).
		Where(builder.Neq{"status": entity.MessageStatusDeleted}).
		GroupBy("conversation_id").
		Find(&result)
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	countMap := make(map[string]int64)
	for _, item := range result {
		countMap[item.ConversationID] = item.Count
	}
	return countMap, nil
}

func (mr *messageRepo) GetUnreadMessageCount(ctx context.Context, userID string) (total, system, private int64, err error) {
	session := mr.data.DB.Context(ctx)
	session.Where(builder.Eq{"to_user_id": userID})
	session.Where(builder.Eq{"status": entity.MessageStatusUnread})
	session.Where(builder.Neq{"status": entity.MessageStatusDeleted})

	total, err = session.Count(&entity.Message{})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		return 0, 0, 0, err
	}

	system, err = mr.data.DB.Context(ctx).
		Where(builder.Eq{"to_user_id": userID}).
		Where(builder.Eq{"status": entity.MessageStatusUnread}).
		Where(builder.Eq{"is_from_system": true}).
		Where(builder.Neq{"status": entity.MessageStatusDeleted}).
		Count(&entity.Message{})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		return 0, 0, 0, err
	}

	private, err = mr.data.DB.Context(ctx).
		Where(builder.Eq{"to_user_id": userID}).
		Where(builder.Eq{"status": entity.MessageStatusUnread}).
		Where(builder.Eq{"is_from_system": false}).
		Where(builder.Neq{"status": entity.MessageStatusDeleted}).
		Count(&entity.Message{})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		return 0, 0, 0, err
	}

	return total, system, private, nil
}

func (mr *messageRepo) UpdateMessageStatus(ctx context.Context, id string, status int) (err error) {
	_, err = mr.data.DB.Context(ctx).
		Where("id = ?", id).
		Cols("status").
		Update(&entity.Message{Status: status})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (mr *messageRepo) MarkConversationAsRead(ctx context.Context, userID, conversationID string) (err error) {
	_, err = mr.data.DB.Context(ctx).
		Where("to_user_id = ?", userID).
		And("conversation_id = ?", conversationID).
		And("status = ?", entity.MessageStatusUnread).
		Cols("status").
		Update(&entity.Message{Status: entity.MessageStatusRead})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (mr *messageRepo) DeleteMessage(ctx context.Context, id string) (err error) {
	_, err = mr.data.DB.Context(ctx).
		Where("id = ?", id).
		Cols("status").
		Update(&entity.Message{Status: entity.MessageStatusDeleted})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (mr *messageRepo) CleanOldMessages(ctx context.Context, days int) (err error) {
	cutoffTime := time.Now().AddDate(0, 0, -days)
	_, err = mr.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		_, err = session.Context(ctx).
			Where("status = ?", entity.MessageStatusDeleted).
			And("updated_at < ?", cutoffTime).
			Delete(&entity.Message{})
		if err != nil {
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return
}

func (mr *messageRepo) BlockUser(ctx context.Context, block *entity.MessageBlock) (err error) {
	block.ID, err = mr.uniqueIDRepo.GenUniqueIDStr(ctx, block.TableName())
	if err != nil {
		return err
	}
	_, err = mr.data.DB.Context(ctx).Insert(block)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (mr *messageRepo) UnblockUser(ctx context.Context, userID, blockedUserID string) (err error) {
	_, err = mr.data.DB.Context(ctx).
		Where("user_id = ?", userID).
		And("blocked_user_id = ?", blockedUserID).
		Delete(&entity.MessageBlock{})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (mr *messageRepo) IsBlocked(ctx context.Context, userID, blockedUserID string) (bool, error) {
	exist, err := mr.data.DB.Context(ctx).
		Where("user_id = ?", userID).
		And("blocked_user_id = ?", blockedUserID).
		Exist(&entity.MessageBlock{})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		return false, err
	}
	return exist, nil
}

func (mr *messageRepo) GetBlockedUserList(ctx context.Context, userID string) ([]*entity.MessageBlock, error) {
	blocks := make([]*entity.MessageBlock, 0)
	err := mr.data.DB.Context(ctx).
		Where("user_id = ?", userID).
		Find(&blocks)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		return nil, err
	}
	return blocks, nil
}
