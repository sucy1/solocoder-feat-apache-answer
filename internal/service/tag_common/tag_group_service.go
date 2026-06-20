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

package tag_common

import (
	"context"
	"strconv"

	"github.com/apache/answer/internal/base/pager"
	"github.com/apache/answer/internal/base/reason"
	"github.com/apache/answer/internal/entity"
	"github.com/apache/answer/internal/schema"
	"github.com/apache/answer/internal/service/unique"
	"github.com/jinzhu/copier"
	"github.com/segmentfault/pacman/errors"
)

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

type TagGroupService struct {
	tagGroupRepo      TagGroupRepo
	tagCommonRepo     TagCommonRepo
	tagRepo           TagRepo
	uniqueIDRepo      unique.UniqueIDRepo
	tagCommonService  *TagCommonService
}

func NewTagGroupService(
	tagGroupRepo TagGroupRepo,
	tagCommonRepo TagCommonRepo,
	tagRepo TagRepo,
	uniqueIDRepo unique.UniqueIDRepo,
	tagCommonService *TagCommonService,
) *TagGroupService {
	return &TagGroupService{
		tagGroupRepo:     tagGroupRepo,
		tagCommonRepo:    tagCommonRepo,
		tagRepo:          tagRepo,
		uniqueIDRepo:     uniqueIDRepo,
		tagCommonService: tagCommonService,
	}
}

func (ts *TagGroupService) AddTagGroup(ctx context.Context, req *schema.AddTagGroupReq) error {
	_, exist, err := ts.tagGroupRepo.GetTagGroupByName(ctx, req.Name)
	if err != nil {
		return err
	}
	if exist {
		return errors.BadRequest(reason.TagAlreadyExist)
	}

	_, exist, err = ts.tagGroupRepo.GetTagGroupBySlugName(ctx, req.SlugName)
	if err != nil {
		return err
	}
	if exist {
		return errors.BadRequest(reason.TagAlreadyExist)
	}

	tagGroup := &entity.TagGroup{}
	_ = copier.Copy(tagGroup, req)
	tagGroup.Status = entity.TagGroupStatusAvailable
	tagGroup.ID, err = ts.uniqueIDRepo.GenUniqueIDStr(ctx, tagGroup.TableName())
	if err != nil {
		return err
	}

	return ts.tagGroupRepo.AddTagGroup(ctx, tagGroup)
}

func (ts *TagGroupService) RemoveTagGroup(ctx context.Context, req *schema.RemoveTagGroupReq) error {
	_, exist, err := ts.tagGroupRepo.GetTagGroupByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.TagNotFound)
	}

	return ts.tagGroupRepo.RemoveTagGroup(ctx, req.ID)
}

func (ts *TagGroupService) UpdateTagGroup(ctx context.Context, req *schema.UpdateTagGroupReq) error {
	tagGroup, exist, err := ts.tagGroupRepo.GetTagGroupByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.TagNotFound)
	}

	cols := make([]string, 0)
	if len(req.Name) > 0 {
		oldTagGroup, exist, err := ts.tagGroupRepo.GetTagGroupByName(ctx, req.Name)
		if err != nil {
			return err
		}
		if exist && oldTagGroup.ID != req.ID {
			return errors.BadRequest(reason.TagAlreadyExist)
		}
		tagGroup.Name = req.Name
		cols = append(cols, "name")
	}

	if len(req.SlugName) > 0 {
		oldTagGroup, exist, err := ts.tagGroupRepo.GetTagGroupBySlugName(ctx, req.SlugName)
		if err != nil {
			return err
		}
		if exist && oldTagGroup.ID != req.ID {
			return errors.BadRequest(reason.TagAlreadyExist)
		}
		tagGroup.SlugName = req.SlugName
		cols = append(cols, "slug_name")
	}

	if req.SortOrder >= 0 {
		tagGroup.SortOrder = req.SortOrder
		cols = append(cols, "sort_order")
	}

	if len(cols) == 0 {
		return nil
	}

	return ts.tagGroupRepo.UpdateTagGroup(ctx, tagGroup, cols)
}

func (ts *TagGroupService) GetTagGroup(ctx context.Context, req *schema.GetTagGroupReq) (*schema.GetTagGroupResp, error) {
	var (
		tagGroup *entity.TagGroup
		exist    bool
		err      error
	)

	if len(req.ID) > 0 {
		tagGroup, exist, err = ts.tagGroupRepo.GetTagGroupByID(ctx, req.ID)
	} else if len(req.Name) > 0 {
		tagGroup, exist, err = ts.tagGroupRepo.GetTagGroupByName(ctx, req.Name)
	} else {
		return nil, errors.BadRequest(reason.RequestFormatError)
	}

	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.NotFound(reason.TagNotFound)
	}

	resp := &schema.GetTagGroupResp{}
	_ = copier.Copy(resp, tagGroup)
	resp.CreatedAt = tagGroup.CreatedAt.Unix()
	resp.UpdatedAt = tagGroup.UpdatedAt.Unix()
	return resp, nil
}

func (ts *TagGroupService) GetTagGroupList(ctx context.Context) ([]*schema.GetTagGroupResp, error) {
	tagGroups, err := ts.tagGroupRepo.GetTagGroupList(ctx, nil)
	if err != nil {
		return nil, err
	}

	resp := make([]*schema.GetTagGroupResp, 0, len(tagGroups))
	for _, tagGroup := range tagGroups {
		item := &schema.GetTagGroupResp{}
		_ = copier.Copy(item, tagGroup)
		item.CreatedAt = tagGroup.CreatedAt.Unix()
		item.UpdatedAt = tagGroup.UpdatedAt.Unix()
		resp = append(resp, item)
	}
	return resp, nil
}

func (ts *TagGroupService) GetTagGroupPage(ctx context.Context, req *schema.GetTagGroupWithPageReq) (*pager.PageModel, error) {
	page := req.Page
	pageSize := req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	tagGroups, total, err := ts.tagGroupRepo.GetTagGroupPage(ctx, page, pageSize, nil)
	if err != nil {
		return nil, err
	}

	resp := make([]*schema.GetTagGroupResp, 0, len(tagGroups))
	for _, tagGroup := range tagGroups {
		item := &schema.GetTagGroupResp{}
		_ = copier.Copy(item, tagGroup)
		item.CreatedAt = tagGroup.CreatedAt.Unix()
		item.UpdatedAt = tagGroup.UpdatedAt.Unix()
		resp = append(resp, item)
	}
	return pager.NewPageModel(total, resp), nil
}

func (ts *TagGroupService) UpdateTagGroupForTag(ctx context.Context, req *schema.UpdateTagGroupReqForTag) error {
	tag, exist, err := ts.tagCommonRepo.GetTagByID(ctx, req.TagID, false)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.TagNotFound)
	}

	if req.GroupID > 0 {
		groupIDStr := strconv.FormatInt(req.GroupID, 10)
		_, exist, err := ts.tagGroupRepo.GetTagGroupByID(ctx, groupIDStr)
		if err != nil {
			return err
		}
		if !exist {
			return errors.BadRequest(reason.TagNotFound)
		}
	}

	tag.GroupID = req.GroupID
	return ts.tagRepo.UpdateTag(ctx, tag)
}

func (ts *TagGroupService) GetTagGroupWithTags(ctx context.Context) ([]*schema.TagGroupWithTagsResp, error) {
	tagGroups, err := ts.tagGroupRepo.GetTagGroupList(ctx, nil)
	if err != nil {
		return nil, err
	}

	allTags, err := ts.tagCommonRepo.GetTagListByName(ctx, "", false, false)
	if err != nil {
		return nil, err
	}

	groupTagMap := make(map[int64][]*entity.Tag)
	for _, tag := range allTags {
		if tag.MainTagID > 0 {
			continue
		}
		groupTagMap[tag.GroupID] = append(groupTagMap[tag.GroupID], tag)
	}

	resp := make([]*schema.TagGroupWithTagsResp, 0, len(tagGroups)+1)
	for _, tagGroup := range tagGroups {
		groupID := tagGroup.ID
		groupIDInt, _ := strconv.ParseInt(groupID, 10, 64)
		tags := groupTagMap[groupIDInt]

		groupResp := &schema.GetTagGroupResp{}
		_ = copier.Copy(groupResp, tagGroup)
		groupResp.CreatedAt = tagGroup.CreatedAt.Unix()
		groupResp.UpdatedAt = tagGroup.UpdatedAt.Unix()

		tagRespList := make([]*schema.GetTagBasicResp, 0, len(tags))
		for _, tag := range tags {
			tagItem := &schema.GetTagBasicResp{}
			_ = copier.Copy(tagItem, tag)
			tagItem.TagID = tag.ID
			tagRespList = append(tagRespList, tagItem)
		}

		resp = append(resp, &schema.TagGroupWithTagsResp{
			Group: groupResp,
			Tags:  tagRespList,
		})

		delete(groupTagMap, groupIDInt)
	}

	otherTags := groupTagMap[0]
	if len(otherTags) > 0 {
		otherGroup := &schema.GetTagGroupResp{
			ID:       "0",
			Name:     "其他",
			SlugName: "other",
		}

		tagRespList := make([]*schema.GetTagBasicResp, 0, len(otherTags))
		for _, tag := range otherTags {
			tagItem := &schema.GetTagBasicResp{}
			_ = copier.Copy(tagItem, tag)
			tagItem.TagID = tag.ID
			tagRespList = append(tagRespList, tagItem)
		}

		resp = append(resp, &schema.TagGroupWithTagsResp{
			Group: otherGroup,
			Tags:  tagRespList,
		})
	}

	return resp, nil
}
