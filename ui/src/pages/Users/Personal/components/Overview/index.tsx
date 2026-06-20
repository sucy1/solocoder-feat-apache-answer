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

import { FC, memo } from 'react';
import { useTranslation } from 'react-i18next';
import { Row, Col, Card } from 'react-bootstrap';

// import * as Type from '@/common/interface';
import { CardBadge, AchievementBadge, Icon } from '@/components';
import {
  useGetRecentAwardBadges,
  useGetUserAchievementSummary,
} from '@/services';
import TopList from '../TopList';

interface Props {
  username: string;
  visible: boolean;
  introduction: string;
  data;
  userId?: string;
}
const Index: FC<Props> = ({ visible, introduction, data, username, userId }) => {
  const { t } = useTranslation('translation', { keyPrefix: 'personal' });
  const { data: recentBadges } = useGetRecentAwardBadges(
    visible ? username : null,
  );
  const { data: achievementSummary } = useGetUserAchievementSummary(
    visible ? userId : null,
  );
  if (!visible) {
    return null;
  }
  return (
    <div>
      <h5 className="mb-3">{t('about_me')}</h5>
      {introduction ? (
        <div
          className="mb-5 text-break fmt"
          dangerouslySetInnerHTML={{ __html: introduction }}
        />
      ) : (
        <div className="mb-5">{t('about_me_empty')}</div>
      )}

      <Row className="mb-4">
        <Col sm={12} md={6} className="mb-4">
          <h5 className="mb-3">{t('top_answers')}</h5>
          {data?.answer?.length > 0 ? (
            <TopList data={data?.answer} type="answer" />
          ) : (
            <div className="mb-5">{t('content_empty')}</div>
          )}
        </Col>
        <Col sm={12} md={6}>
          <h5 className="mb-3">{t('top_questions')}</h5>
          {data?.question?.length > 0 ? (
            <TopList data={data?.question} type="question" />
          ) : (
            <div className="mb-5">{t('content_empty')}</div>
          )}
        </Col>
      </Row>

      <div className="mb-4">
        <h5 className="mb-3">{t('recent_badges')}</h5>
        {Number(recentBadges?.count) > 0 ? (
          <Row>
            {recentBadges?.list?.map((item) => {
              return (
                <Col sm={6} md={4} lg={3} key={item.id} className="mb-4">
                  <CardBadge
                    data={item}
                    urlSearchParams={`username=${username}`}
                    badgePillType="count"
                  />
                </Col>
              );
            })}
          </Row>
        ) : (
          <div className="mb-5">{t('content_empty')}</div>
        )}
      </div>

      {achievementSummary && (
        <>
          <h5 className="mb-3">{t('achievement.overview')}</h5>
          <Row className="mb-4">
            <Col sm={4} className="mb-3">
              <Card className="text-center h-100">
                <Card.Body>
                  <Icon
                    name="trophy-fill"
                    size="32px"
                    className="text-warning mb-2"
                  />
                  <div className="h3 mb-1">
                    {achievementSummary.total_reputation}
                  </div>
                  <div className="small text-secondary">
                    {t('achievement.total_reputation')}
                  </div>
                </Card.Body>
              </Card>
            </Col>
            <Col sm={4} className="mb-3">
              <Card className="text-center h-100">
                <Card.Body>
                  <Icon
                    name="medal-fill"
                    size="32px"
                    className="text-success mb-2"
                  />
                  <div className="h3 mb-1">
                    {achievementSummary.badges_count}/{achievementSummary.max_badges || 50}
                  </div>
                  <div className="small text-secondary">
                    {t('achievement.badges_earned')}
                  </div>
                </Card.Body>
              </Card>
            </Col>
            <Col sm={4} className="mb-3">
              <Card className="text-center h-100">
                <Card.Body>
                  <Icon
                    name="calendar-check-fill"
                    size="32px"
                    className="text-primary mb-2"
                  />
                  <div className="h3 mb-1">
                    {achievementSummary.consecutive_login_days}
                  </div>
                  <div className="small text-secondary">
                    {t('achievement.consecutive_login_days')}
                  </div>
                </Card.Body>
              </Card>
            </Col>
          </Row>

          {achievementSummary.badges?.length > 0 && (
            <div className="mb-4">
              <h5 className="mb-3">{t('achievement.badges')}</h5>
              <Row>
                {achievementSummary.badges.map((badge) => (
                  <Col
                    sm={4}
                    md={3}
                    lg={2}
                    key={badge.badge_id}
                    className="mb-3">
                    <AchievementBadge data={badge} showEarnedTime />
                  </Col>
                ))}
              </Row>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default memo(Index);
