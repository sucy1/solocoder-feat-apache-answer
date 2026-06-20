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

package controller

import (
	"github.com/apache/answer/internal/base/handler"
	"github.com/apache/answer/internal/base/middleware"
	"github.com/apache/answer/internal/base/pager"
	"github.com/apache/answer/internal/base/reason"
	"github.com/apache/answer/internal/schema"
	"github.com/apache/answer/internal/service/permission"
	"github.com/apache/answer/internal/service/rank"
	"github.com/apache/answer/internal/service/tag_common"
	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/errors"
)

// TagGroupController tag group controller
type TagGroupController struct {
	tagGroupService *tag_common.TagGroupService
	rankService     *rank.RankService
}

// NewTagGroupController new controller
func NewTagGroupController(
	tagGroupService *tag_common.TagGroupService,
	rankService *rank.RankService,
) *TagGroupController {
	return &TagGroupController{tagGroupService: tagGroupService, rankService: rankService}
}

// AddTagGroup add tag group
// @Summary add tag group
// @Description add tag group
// @Security ApiKeyAuth
// @Tags Tag
// @Accept json
// @Produce json
// @Param data body schema.AddTagGroupReq true "tag group"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/tag/group [post]
func (tc *TagGroupController) AddTagGroup(ctx *gin.Context) {
	req := &schema.AddTagGroupReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	isAdminModerator := middleware.GetUserIsAdminModerator(ctx)
	if !isAdminModerator {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := tc.rankService.CheckOperationPermission(ctx, req.UserID, permission.TagAdd, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = tc.tagGroupService.AddTagGroup(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// RemoveTagGroup remove tag group
// @Summary remove tag group
// @Description remove tag group
// @Security ApiKeyAuth
// @Tags Tag
// @Accept json
// @Produce json
// @Param data body schema.RemoveTagGroupReq true "tag group"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/tag/group [delete]
func (tc *TagGroupController) RemoveTagGroup(ctx *gin.Context) {
	req := &schema.RemoveTagGroupReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	isAdminModerator := middleware.GetUserIsAdminModerator(ctx)
	if !isAdminModerator {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := tc.rankService.CheckOperationPermission(ctx, req.UserID, permission.TagDelete, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = tc.tagGroupService.RemoveTagGroup(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// UpdateTagGroup update tag group
// @Summary update tag group
// @Description update tag group
// @Security ApiKeyAuth
// @Tags Tag
// @Accept json
// @Produce json
// @Param data body schema.UpdateTagGroupReq true "tag group"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/tag/group [put]
func (tc *TagGroupController) UpdateTagGroup(ctx *gin.Context) {
	req := &schema.UpdateTagGroupReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	isAdminModerator := middleware.GetUserIsAdminModerator(ctx)
	if !isAdminModerator {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := tc.rankService.CheckOperationPermission(ctx, req.UserID, permission.TagEdit, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = tc.tagGroupService.UpdateTagGroup(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// GetTagGroup get tag group
// @Summary get tag group
// @Description get tag group
// @Tags Tag
// @Accept json
// @Produce json
// @Param id query string false "tag group id"
// @Param name query string false "tag group name"
// @Success 200 {object} handler.RespBody{data=schema.GetTagGroupResp}
// @Router /answer/api/v1/tag/group [get]
func (tc *TagGroupController) GetTagGroup(ctx *gin.Context) {
	req := &schema.GetTagGroupReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := tc.tagGroupService.GetTagGroup(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// GetTagGroupList get tag group list
// @Summary get tag group list
// @Description get tag group list
// @Tags Tag
// @Produce json
// @Success 200 {object} handler.RespBody{data=[]schema.GetTagGroupResp}
// @Router /answer/api/v1/tag/group/list [get]
func (tc *TagGroupController) GetTagGroupList(ctx *gin.Context) {
	resp, err := tc.tagGroupService.GetTagGroupList(ctx)
	handler.HandleResponse(ctx, err, resp)
}

// GetTagGroupPage get tag group page
// @Summary get tag group page
// @Description get tag group page
// @Tags Tag
// @Produce json
// @Param page query int false "page size"
// @Param page_size query int false "page size"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.GetTagGroupResp}}
// @Router /answer/api/v1/tag/group/page [get]
func (tc *TagGroupController) GetTagGroupPage(ctx *gin.Context) {
	req := &schema.GetTagGroupWithPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := tc.tagGroupService.GetTagGroupPage(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if pager.ValPageOutOfRange(resp.Count, req.Page, req.PageSize) {
		handler.HandleResponse(ctx, errors.NotFound(reason.RequestFormatError), nil)
		return
	}
	handler.HandleResponse(ctx, err, resp)
}

// GetTagGroupWithTags get tag group with tags
// @Summary get tag group with tags
// @Description get tag group with tags
// @Tags Tag
// @Produce json
// @Success 200 {object} handler.RespBody{data=[]schema.TagGroupWithTagsResp}
// @Router /answer/api/v1/tag/group/with-tags [get]
func (tc *TagGroupController) GetTagGroupWithTags(ctx *gin.Context) {
	resp, err := tc.tagGroupService.GetTagGroupWithTags(ctx)
	handler.HandleResponse(ctx, err, resp)
}

// UpdateTagGroupForTag update tag group for tag
// @Summary update tag group for tag
// @Description update tag group for tag
// @Security ApiKeyAuth
// @Tags Tag
// @Accept json
// @Produce json
// @Param data body schema.UpdateTagGroupReqForTag true "tag group for tag"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/tag/group/tag [put]
func (tc *TagGroupController) UpdateTagGroupForTag(ctx *gin.Context) {
	req := &schema.UpdateTagGroupReqForTag{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	isAdminModerator := middleware.GetUserIsAdminModerator(ctx)
	if !isAdminModerator {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := tc.rankService.CheckOperationPermission(ctx, req.UserID, permission.TagEdit, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = tc.tagGroupService.UpdateTagGroupForTag(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}
