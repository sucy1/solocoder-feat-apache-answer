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

package tag

import (
	"context"

	"github.com/apache/answer/internal/base/data"
	"github.com/apache/answer/internal/base/reason"
	"github.com/apache/answer/internal/entity"
	"github.com/segmentfault/pacman/errors"
	"xorm.io/builder"
)

// TagGroupRepo tag group repository
type TagGroupRepo interface {
	AddTagGroup(ctx context.Context, tagGroup *entity.TagGroup) error
	RemoveTagGroup(ctx context.Context, id string) error
	UpdateTagGroup(ctx context.Context, tagGroup *entity.TagGroup, cols []string) error
	GetTagGroupByID(ctx context.Context, id string) (tagGroup *entity.TagGroup, exist bool, err error)
	GetTagGroupByName(ctx context.Context, name string) (tagGroup *entity.TagGroup, exist bool, err error)
	GetTagGroupBySlugName(ctx context.Context, slugName string) (tagGroup *entity.TagGroup, exist bool, err error)
	GetTagGroupList(ctx context.Context, tagGroup *entity.TagGroup) (tagGroups []*entity.TagGroup, err error)
	GetTagGroupPage(ctx context.Context, page, pageSize int, tagGroup *entity.TagGroup) (tagGroups []*entity.TagGroup, total int64, err error)
}

// tagGroupRepo tag group repository
type tagGroupRepo struct {
	data *data.Data
}

// NewTagGroupRepo new repository
func NewTagGroupRepo(data *data.Data) TagGroupRepo {
	return &tagGroupRepo{
		data: data,
	}
}

// AddTagGroup add tag group
func (tr *tagGroupRepo) AddTagGroup(ctx context.Context, tagGroup *entity.TagGroup) (err error) {
	_, err = tr.data.DB.Context(ctx).Insert(tagGroup)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// RemoveTagGroup delete tag group
func (tr *tagGroupRepo) RemoveTagGroup(ctx context.Context, id string) (err error) {
	session := tr.data.DB.Context(ctx).Where(builder.Eq{"id": id})
	_, err = session.Update(&entity.TagGroup{Status: entity.TagGroupStatusDeleted})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// UpdateTagGroup update tag group
func (tr *tagGroupRepo) UpdateTagGroup(ctx context.Context, tagGroup *entity.TagGroup, cols []string) (err error) {
	session := tr.data.DB.Context(ctx).Where(builder.Eq{"id": tagGroup.ID})
	if len(cols) > 0 {
		session.Cols(cols...)
	}
	_, err = session.Update(tagGroup)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetTagGroupByID get tag group by id
func (tr *tagGroupRepo) GetTagGroupByID(ctx context.Context, id string) (tagGroup *entity.TagGroup, exist bool, err error) {
	tagGroup = &entity.TagGroup{}
	exist, err = tr.data.DB.Context(ctx).ID(id).Get(tagGroup)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetTagGroupByName get tag group by name
func (tr *tagGroupRepo) GetTagGroupByName(ctx context.Context, name string) (tagGroup *entity.TagGroup, exist bool, err error) {
	tagGroup = &entity.TagGroup{}
	session := tr.data.DB.Context(ctx).Where(builder.Eq{"name": name, "status": entity.TagGroupStatusAvailable})
	exist, err = session.Get(tagGroup)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetTagGroupBySlugName get tag group by slug name
func (tr *tagGroupRepo) GetTagGroupBySlugName(ctx context.Context, slugName string) (tagGroup *entity.TagGroup, exist bool, err error) {
	tagGroup = &entity.TagGroup{}
	session := tr.data.DB.Context(ctx).Where(builder.Eq{"slug_name": slugName, "status": entity.TagGroupStatusAvailable})
	exist, err = session.Get(tagGroup)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetTagGroupList get tag group list all
func (tr *tagGroupRepo) GetTagGroupList(ctx context.Context, tagGroup *entity.TagGroup) (tagGroups []*entity.TagGroup, err error) {
	tagGroups = make([]*entity.TagGroup, 0)
	session := tr.data.DB.Context(ctx).Where(builder.Eq{"status": entity.TagGroupStatusAvailable})
	session.OrderBy("sort_order ASC, id DESC")
	err = session.Find(&tagGroups, tagGroup)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetTagGroupPage get tag group page
func (tr *tagGroupRepo) GetTagGroupPage(ctx context.Context, page, pageSize int, tagGroup *entity.TagGroup) (tagGroups []*entity.TagGroup, total int64, err error) {
	tagGroups = make([]*entity.TagGroup, 0)
	session := tr.data.DB.Context(ctx).Where(builder.Eq{"status": entity.TagGroupStatusAvailable})
	if tagGroup != nil {
		session.And(tagGroup)
	}
	total, err = session.Count(&entity.TagGroup{})
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		return
	}
	session.OrderBy("sort_order ASC, id DESC")
	err = session.Limit(pageSize, (page-1)*pageSize).Find(&tagGroups)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}
