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
	"strconv"

	"github.com/apache/answer/internal/base/constant"
	"github.com/apache/answer/internal/base/data"
	"github.com/apache/answer/internal/base/pager"
	"github.com/apache/answer/internal/base/reason"
	"github.com/apache/answer/internal/entity"
	"github.com/apache/answer/internal/repo/achievement"
	"github.com/apache/answer/internal/schema"
	"github.com/apache/answer/internal/service/badge"
	"github.com/apache/answer/internal/service/config"
	"github.com/apache/answer/internal/service/eventqueue"
	usercommon "github.com/apache/answer/internal/service/user_common"
	"github.com/apache/answer/pkg/uid"
	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
)

const (
	ReputationRegister         = 100
	ReputationFirstQuestion    = 50
	ReputationAnswerAccepted   = 100
	ReputationConsecutiveLogin = 200
	ReputationFirstAnswer      = 30
	ReputationUpvote           = 5

	ConsecutiveLoginDaysRequired = 7

	BadgeHandlerRegister         = "UserRegister"
	BadgeHandlerFirstQuestion    = "FirstQuestion"
	BadgeHandlerAnswerAccepted   = "AnswerAccepted"
	BadgeHandlerConsecutiveLogin = "ConsecutiveLogin"
	BadgeHandlerFirstAnswer      = "FirstAnswer"
)

type AchievementService struct {
	data                *data.Data
	achievementRepo     achievement.AchievementRepo
	badgeService        *badge.BadgeService
	badgeAwardService   *badge.BadgeAwardService
	userCommon          *usercommon.UserCommon
	eventQueueService   eventqueue.Service
	configService       *config.ConfigService
}

func NewAchievementService(
	data *data.Data,
	achievementRepo achievement.AchievementRepo,
	badgeService *badge.BadgeService,
	badgeAwardService *badge.BadgeAwardService,
	userCommon *usercommon.UserCommon,
	eventQueueService eventqueue.Service,
	configService *config.ConfigService,
) *AchievementService {
	as := &AchievementService{
		data:              data,
		achievementRepo:   achievementRepo,
		badgeService:      badgeService,
		badgeAwardService: badgeAwardService,
		userCommon:        userCommon,
		eventQueueService: eventQueueService,
		configService:     configService,
	}
	eventQueueService.RegisterHandler(as.Handler)
	return as
}

func (as *AchievementService) Handler(ctx context.Context, msg *schema.EventMsg) error {
	switch msg.EventType {
	case constant.EventUserUpdate:
		if msg.GetExtra("action") == "register" {
			if err := as.HandleUserRegister(ctx, msg.UserID); err != nil {
				log.Errorf("handle user register achievement error: %v", err)
			}
		}
	case constant.EventQuestionCreate:
		if err := as.HandleFirstQuestion(ctx, msg.UserID); err != nil {
			log.Errorf("handle first question achievement error: %v", err)
		}
	case constant.EventQuestionAccept:
		if err := as.HandleAnswerAccepted(ctx, msg.AnswerUserID); err != nil {
			log.Errorf("handle answer accepted achievement error: %v", err)
		}
	case constant.EventAnswerCreate:
		if err := as.HandleFirstAnswer(ctx, msg.UserID); err != nil {
			log.Errorf("handle first answer achievement error: %v", err)
		}
	case constant.EventQuestionVote:
		if voteAmount := msg.GetExtra("vote_up_amount"); voteAmount == "1" {
			if err := as.HandleQuestionUpvote(ctx, msg.QuestionUserID); err != nil {
				log.Errorf("handle question upvote achievement error: %v", err)
			}
		}
	case constant.EventAnswerVote:
		if voteAmount := msg.GetExtra("vote_up_amount"); voteAmount == "1" {
			if err := as.HandleAnswerUpvote(ctx, msg.AnswerUserID); err != nil {
				log.Errorf("handle answer upvote achievement error: %v", err)
			}
		}
	}
	return nil
}

func (as *AchievementService) GetUserAchievementSummary(
	ctx context.Context,
	req *schema.GetUserAchievementReq,
) (*schema.GetUserAchievementSummaryResp, error) {
	userInfo, exist, err := as.userCommon.GetUserBasicInfoByUserName(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.BadRequest(reason.UserNotFound)
	}
	req.UserID = userInfo.ID

	totalReputation, err := as.achievementRepo.GetUserTotalReputation(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	badgeCount, err := as.achievementRepo.GetUserBadgeCount(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	badgeReq := &schema.GetUserBadgeAwardListReq{
		Username: req.Username,
		UserID:   req.UserID,
		Limit:    10,
	}
	badges, _, err := as.badgeAwardService.GetUserRecentBadgeAwardList(ctx.(*gin.Context), badgeReq)
	if err != nil {
		return nil, err
	}

	achievements, _, err := as.achievementRepo.GetUserAchievementList(ctx, req.UserID, 1, 10)
	if err != nil {
		return nil, err
	}

	recentActivities := as.formatAchievementList(ctx.(*gin.Context), achievements)

	return &schema.GetUserAchievementSummaryResp{
		TotalReputation:  totalReputation,
		BadgeCount:       int(badgeCount),
		MaxBadges:        as.getMaxBadges(ctx),
		Badges:           badges,
		RecentActivities: recentActivities,
	}, nil
}

func (as *AchievementService) GetUserAchievementList(
	ctx context.Context,
	req *schema.GetUserAchievementReq,
) (*pager.PageModel, error) {
	userInfo, exist, err := as.userCommon.GetUserBasicInfoByUserName(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.BadRequest(reason.UserNotFound)
	}
	req.UserID = userInfo.ID

	page, pageSize := pager.ValPageAndPageSize(req.Page, req.PageSize)
	achievements, total, err := as.achievementRepo.GetUserAchievementList(ctx, req.UserID, page, pageSize)
	if err != nil {
		return nil, err
	}

	resp := as.formatAchievementList(ctx.(*gin.Context), achievements)
	return pager.NewPageModel(total, resp), nil
}

func (as *AchievementService) AddReputation(
	ctx context.Context,
	event *schema.AchievementReputationEvent,
) error {
	if event.Reputation <= 0 {
		return nil
	}
	return as.achievementRepo.AddReputation(ctx, event.UserID, event.Source, event.Reputation, event.Description)
}

func (as *AchievementService) getMaxBadges(ctx context.Context) int {
	val, err := as.configService.GetStringValue(ctx, constant.ConfigKeyAchievementMaxBadges)
	if err != nil || val == "" {
		return entity.DefaultMaxUserBadges
	}
	max, err := strconv.Atoi(val)
	if err != nil || max <= 0 {
		return entity.DefaultMaxUserBadges
	}
	return max
}

func (as *AchievementService) AwardBadge(
	ctx context.Context,
	event *schema.AchievementBadgeEvent,
) error {
	hasBadge, err := as.achievementRepo.HasUserBadge(ctx, event.UserID, event.BadgeID)
	if err != nil {
		return err
	}
	if hasBadge {
		return nil
	}

	badgeCount, err := as.achievementRepo.GetUserBadgeCount(ctx, event.UserID)
	if err != nil {
		return err
	}
	if int(badgeCount) >= as.getMaxBadges(ctx) {
		return errors.InternalServer("error.achievement.badge_count_exceed")
	}

	return as.badgeAwardService.Award(ctx, event.BadgeID, event.UserID, entity.BadgeEmptyAwardKey)
}

func (as *AchievementService) HandleUserRegister(ctx context.Context, userID string) error {
	if len(userID) == 0 {
		return nil
	}

	reputationEvent := &schema.AchievementReputationEvent{
		UserID:      userID,
		Source:      entity.AchievementSourceRegister,
		Reputation:  ReputationRegister,
		Description: "User registration reward",
	}
	if err := as.AddReputation(ctx, reputationEvent); err != nil {
		log.Errorf("add register reputation error: %v", err)
	}

	badges := as.getBadgesByHandler(ctx, BadgeHandlerRegister)
	for _, b := range badges {
		badgeEvent := &schema.AchievementBadgeEvent{
			BadgeID: b.ID,
			UserID:  userID,
		}
		if err := as.AwardBadge(ctx, badgeEvent); err != nil {
			log.Errorf("award register badge error: %v", err)
		}
	}
	return nil
}

func (as *AchievementService) HandleFirstQuestion(ctx context.Context, userID string) error {
	if len(userID) == 0 {
		return nil
	}

	achievements, _, err := as.achievementRepo.GetUserAchievementList(ctx, userID, 0, 0)
	if err != nil {
		return err
	}

	hasFirstQuestion := false
	for _, a := range achievements {
		if a.Source == entity.AchievementSourceFirstQuestion {
			hasFirstQuestion = true
			break
		}
	}
	if hasFirstQuestion {
		return nil
	}

	reputationEvent := &schema.AchievementReputationEvent{
		UserID:      userID,
		Source:      entity.AchievementSourceFirstQuestion,
		Reputation:  ReputationFirstQuestion,
		Description: "First question reward",
	}
	if err := as.AddReputation(ctx, reputationEvent); err != nil {
		log.Errorf("add first question reputation error: %v", err)
	}

	badges := as.getBadgesByHandler(ctx, BadgeHandlerFirstQuestion)
	for _, b := range badges {
		badgeEvent := &schema.AchievementBadgeEvent{
			BadgeID: b.ID,
			UserID:  userID,
		}
		if err := as.AwardBadge(ctx, badgeEvent); err != nil {
			log.Errorf("award first question badge error: %v", err)
		}
	}
	return nil
}

func (as *AchievementService) HandleAnswerAccepted(ctx context.Context, userID string) error {
	if len(userID) == 0 {
		return nil
	}

	reputationEvent := &schema.AchievementReputationEvent{
		UserID:      userID,
		Source:      entity.AchievementSourceAnswerAccepted,
		Reputation:  ReputationAnswerAccepted,
		Description: "Answer accepted reward",
	}
	if err := as.AddReputation(ctx, reputationEvent); err != nil {
		log.Errorf("add answer accepted reputation error: %v", err)
	}

	badges := as.getBadgesByHandler(ctx, BadgeHandlerAnswerAccepted)
	for _, b := range badges {
		badgeEvent := &schema.AchievementBadgeEvent{
			BadgeID: b.ID,
			UserID:  userID,
		}
		if err := as.AwardBadge(ctx, badgeEvent); err != nil {
			log.Errorf("award answer accepted badge error: %v", err)
		}
	}
	return nil
}

func (as *AchievementService) HandleConsecutiveLogin(ctx context.Context, userID string) error {
	if len(userID) == 0 {
		return nil
	}

	days, err := as.achievementRepo.CheckConsecutiveLoginDays(ctx, userID)
	if err != nil {
		return err
	}

	newDays := days + 1
	if err := as.achievementRepo.UpdateConsecutiveLoginDays(ctx, userID, newDays); err != nil {
		return err
	}

	if newDays >= ConsecutiveLoginDaysRequired && days < ConsecutiveLoginDaysRequired {
		reputationEvent := &schema.AchievementReputationEvent{
			UserID:      userID,
			Source:      entity.AchievementSourceConsecutiveLogin,
			Reputation:  ReputationConsecutiveLogin,
			Description: "7 days consecutive login reward",
		}
		if err := as.AddReputation(ctx, reputationEvent); err != nil {
			log.Errorf("add consecutive login reputation error: %v", err)
		}

		badges := as.getBadgesByHandler(ctx, BadgeHandlerConsecutiveLogin)
		for _, b := range badges {
			badgeEvent := &schema.AchievementBadgeEvent{
				BadgeID: b.ID,
				UserID:  userID,
			}
			if err := as.AwardBadge(ctx, badgeEvent); err != nil {
				log.Errorf("award consecutive login badge error: %v", err)
			}
		}
	}
	return nil
}

func (as *AchievementService) HandleFirstAnswer(ctx context.Context, userID string) error {
	if len(userID) == 0 {
		return nil
	}

	achievements, _, err := as.achievementRepo.GetUserAchievementList(ctx, userID, 0, 0)
	if err != nil {
		return err
	}

	hasFirstAnswer := false
	for _, a := range achievements {
		if a.Source == entity.AchievementSourceFirstAnswer {
			hasFirstAnswer = true
			break
		}
	}
	if hasFirstAnswer {
		return nil
	}

	reputationEvent := &schema.AchievementReputationEvent{
		UserID:      userID,
		Source:      entity.AchievementSourceFirstAnswer,
		Reputation:  ReputationFirstAnswer,
		Description: "First answer reward",
	}
	if err := as.AddReputation(ctx, reputationEvent); err != nil {
		log.Errorf("add first answer reputation error: %v", err)
	}

	badges := as.getBadgesByHandler(ctx, BadgeHandlerFirstAnswer)
	for _, b := range badges {
		badgeEvent := &schema.AchievementBadgeEvent{
			BadgeID: b.ID,
			UserID:  userID,
		}
		if err := as.AwardBadge(ctx, badgeEvent); err != nil {
			log.Errorf("award first answer badge error: %v", err)
		}
	}
	return nil
}

func (as *AchievementService) HandleQuestionUpvote(ctx context.Context, userID string) error {
	if len(userID) == 0 {
		return nil
	}

	reputationEvent := &schema.AchievementReputationEvent{
		UserID:      userID,
		Source:      entity.AchievementSourceQuestionUpvote,
		Reputation:  ReputationUpvote,
		Description: "Question upvote reward",
	}
	return as.AddReputation(ctx, reputationEvent)
}

func (as *AchievementService) HandleAnswerUpvote(ctx context.Context, userID string) error {
	if len(userID) == 0 {
		return nil
	}

	reputationEvent := &schema.AchievementReputationEvent{
		UserID:      userID,
		Source:      entity.AchievementSourceAnswerUpvote,
		Reputation:  ReputationUpvote,
		Description: "Answer upvote reward",
	}
	return as.AddReputation(ctx, reputationEvent)
}

func (as *AchievementService) formatAchievementList(
	ctx *gin.Context,
	achievements []*entity.UserAchievement,
) []*schema.GetUserAchievementResp {
	resp := make([]*schema.GetUserAchievementResp, 0, len(achievements))
	for _, a := range achievements {
		item := &schema.GetUserAchievementResp{
			ID:          uid.EnShortID(a.ID),
			Type:        a.AchievementType,
			Source:      a.Source,
			Reputation:  a.Reputation,
			Description: a.Description,
			CreatedAt:   a.CreatedAt.Unix(),
		}

		if a.AchievementType == entity.AchievementTypeBadge && len(a.AchievementID) > 0 {
			badgeInfo, err := as.badgeService.GetBadgeInfo(ctx, a.AchievementID, a.UserID)
			if err == nil {
				item.BadgeInfo = &schema.GetUserBadgeAwardListResp{
					ID:          badgeInfo.ID,
					Name:        badgeInfo.Name,
					Icon:        badgeInfo.Icon,
					EarnedCount: badgeInfo.EarnedCount,
					Level:       badgeInfo.Level,
				}
			}
		}
		resp = append(resp, item)
	}
	return resp
}

func (as *AchievementService) getBadgesByHandler(ctx context.Context, handler string) []*entity.Badge {
	badges := make([]*entity.Badge, 0)
	err := as.data.DB.Context(ctx).Where("handler = ?", handler).Find(&badges)
	if err != nil {
		log.Errorf("error getting badge by handler %s: %v", handler, err)
		return nil
	}
	return badges
}
