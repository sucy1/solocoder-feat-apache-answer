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

import { FC } from 'react';
import { Badge } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router-dom';

import { FormatTime, Icon, Empty, Pagination } from '@/components';
import { useGetUserAchievementList } from '@/services';

import './index.scss';

interface Props {
  visible: boolean;
  userId?: string;
}

const ACHIEVEMENT_TYPE_CONFIG: Record<
  string,
  { icon: string; color: string; reputation: number }
> = {
  register: { icon: 'person-plus-fill', color: 'primary', reputation: 100 },
  first_question: { icon: 'question-circle-fill', color: 'info', reputation: 50 },
  first_answer: { icon: 'chat-square-text-fill', color: 'success', reputation: 30 },
  answer_accepted: { icon: 'check-circle-fill', color: 'warning', reputation: 100 },
  consecutive_login_7: { icon: 'calendar-check-fill', color: 'danger', reputation: 200 },
  question_upvote: { icon: 'hand-thumbs-up-fill', color: 'info', reputation: 5 },
  answer_upvote: { icon: 'hand-thumbs-up-fill', color: 'success', reputation: 5 },
  badge: { icon: 'medal-fill', color: 'warning', reputation: 0 },
};

const Index: FC<Props> = ({ visible, userId }) => {
  const { t } = useTranslation('translation', { keyPrefix: 'achievement' });
  const [searchParams] = useSearchParams();
  const page = Number(searchParams.get('page') || 1);
  const pageSize = 20;

  const { data, isLoading } = useGetUserAchievementList(
    visible ? userId : null,
    page,
    pageSize,
  );

  const achievements = data?.list || [];
  const count = data?.count || 0;

  const getAchievementConfig = (type: string) => {
    return ACHIEVEMENT_TYPE_CONFIG[type] || {
      icon: 'star-fill',
      color: 'secondary',
      reputation: 0,
    };
  };

  if (!visible) {
    return null;
  }

  const reputationAchievements = achievements.filter(
    (item) => item.achievement_type !== 'badge',
  );
  const badgeAchievements = achievements.filter(
    (item) => item.achievement_type === 'badge',
  );

  return (
    <div className="achievements-container">
      {!isLoading && achievements.length <= 0 && <Empty />}

      {reputationAchievements.length > 0 && (
        <div className="mb-5">
          <h5 className="mb-4">{t('reputation_history')}</h5>
          <div className="timeline">
            {reputationAchievements.map((item) => {
              const config = getAchievementConfig(item.achievement_type);
              return (
                <div className="timeline-item" key={item.id}>
                  <div className="timeline-marker">
                    <Icon
                      name={config.icon}
                      size="20px"
                      className={`text-${config.color}`}
                    />
                  </div>
                  <div className="timeline-content">
                    <div className="d-flex justify-content-between align-items-start mb-1">
                      <div>
                        <Badge bg={config.color} className="me-2">
                          {t(`type.${item.achievement_type}`, item.achievement_type)}
                        </Badge>
                        <span className="text-success fw-bold">+{item.reputation}</span>
                      </div>
                      <div className="small text-secondary">
                        <FormatTime time={item.created_at} />
                      </div>
                    </div>
                    <div className="text-break">{item.description}</div>
                    {item.source && (
                      <div className="small text-secondary mt-1">
                        {t('source')}: {item.source}
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {badgeAchievements.length > 0 && (
        <div className="mb-5">
          <h5 className="mb-4">{t('badge_history')}</h5>
          <div className="timeline">
            {badgeAchievements.map((item) => {
              const config = getAchievementConfig(item.achievement_type);
              return (
                <div className="timeline-item" key={item.id}>
                  <div className="timeline-marker">
                    <Icon
                      name={config.icon}
                      size="20px"
                      className={`text-${config.color}`}
                    />
                  </div>
                  <div className="timeline-content">
                    <div className="d-flex justify-content-between align-items-start mb-1">
                      <div>
                        <Badge bg={config.color} className="me-2">
                          {t('type.badge')}
                        </Badge>
                        <strong>{item.description}</strong>
                      </div>
                      <div className="small text-secondary">
                        <FormatTime time={item.created_at} />
                      </div>
                    </div>
                    {item.source && (
                      <div className="small text-secondary mt-1">
                        {t('source')}: {item.source}
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {count > 0 && (
        <div className="d-flex justify-content-center py-4">
          <Pagination
            pageSize={pageSize}
            totalSize={count}
            currentPage={page}
          />
        </div>
      )}
    </div>
  );
};

export default Index;
