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
	"github.com/apache/answer/internal/base/reason"
	"github.com/apache/answer/internal/schema"
	"github.com/apache/answer/internal/service/achievement"
	usercommon "github.com/apache/answer/internal/service/user_common"
	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/errors"
)

type AchievementController struct {
	achievementService *achievement.AchievementService
	userCommon         *usercommon.UserCommon
}

func NewAchievementController(
	achievementService *achievement.AchievementService,
	userCommon *usercommon.UserCommon,
) *AchievementController {
	return &AchievementController{
		achievementService: achievementService,
		userCommon:         userCommon,
	}
}

// GetUserAchievementSummary get user achievement summary
// @Summary get user achievement summary
// @Description get user achievement summary by username
// @Tags api-achievement
// @Accept json
// @Produce json
// @Param username query string true "username"
// @Success 200 {object} handler.RespBody{data=schema.GetUserAchievementSummaryResp}
// @Router /answer/api/v1/achievement/summary [get]
func (ac *AchievementController) GetUserAchievementSummary(ctx *gin.Context) {
	req := &schema.GetUserAchievementReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := ac.achievementService.GetUserAchievementSummary(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// GetUserAchievementList get user achievement list
// @Summary get user achievement list
// @Description get user achievement list by username
// @Tags api-achievement
// @Accept json
// @Produce json
// @Param username query string true "username"
// @Param page query int false "page"
// @Param page_size query int false "page size"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.GetUserAchievementResp}}
// @Router /answer/api/v1/achievement/list [get]
func (ac *AchievementController) GetUserAchievementList(ctx *gin.Context) {
	req := &schema.GetUserAchievementReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := ac.achievementService.GetUserAchievementList(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// GetMyAchievementSummary get my achievement summary
// @Summary get my achievement summary
// @Description get current login user achievement summary
// @Tags api-achievement
// @Accept json
// @Produce json
// @Success 200 {object} handler.RespBody{data=schema.GetUserAchievementSummaryResp}
// @Router /answer/api/v1/achievement/me/summary [get]
func (ac *AchievementController) GetMyAchievementSummary(ctx *gin.Context) {
	userID := middleware.GetLoginUserIDFromContext(ctx)
	if len(userID) == 0 {
		handler.HandleResponse(ctx, errors.Unauthorized(reason.UnauthorizedError), nil)
		return
	}

	userInfo, exist, err := ac.userCommon.GetUserBasicInfoByID(ctx, userID)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !exist {
		handler.HandleResponse(ctx, errors.BadRequest(reason.UserNotFound), nil)
		return
	}

	req := &schema.GetUserAchievementReq{
		Username: userInfo.Username,
	}

	resp, err := ac.achievementService.GetUserAchievementSummary(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}
