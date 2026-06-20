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
import { Card, Tooltip, OverlayTrigger } from 'react-bootstrap';
import { useTranslation } from 'react-i18next';
import classnames from 'classnames';

import Icon from '../Icon';
import FormatTime from '../FormatTime';
import * as Type from '@/common/interface';

import './index.scss';

interface IProps {
  data: Type.AchievementBadge;
  showEarnedTime?: boolean;
}

const Index: FC<IProps> = ({ data, showEarnedTime = false }) => {
  const { t } = useTranslation('translation', { keyPrefix: 'achievement' });

  const renderTooltip = () => (
    <Tooltip id={`badge-tooltip-${data.badge_id}`}>
      <div className="text-start">
        <strong>{data.badge_name}</strong>
        <div className="small">{data.badge_description}</div>
        {data.earned && data.earned_at && showEarnedTime && (
          <div className="small mt-1">
            {t('earned_at')}: <FormatTime time={data.earned_at} />
          </div>
        )}
      </div>
    </Tooltip>
  );

  return (
    <OverlayTrigger placement="top" overlay={renderTooltip}>
      <Card
        className={classnames(
          'text-center achievement-badge-card h-100',
          !data.earned && 'badge-not-earned',
        )}>
        <Card.Body className="p-3">
          {data.icon.startsWith('http') ? (
            <img
              src={data.icon}
              width={64}
              height={64}
              alt={data.badge_name}
              className="mb-2"
            />
          ) : (
            <Icon
              name={data.icon}
              size="64px"
              className={classnames(
                'lh-1 mb-2',
                data.earned && 'earned-badge',
              )}
            />
          )}
          <h6 className="mb-0 small text-truncate">{data.badge_name}</h6>
          {data.earned && data.earned_at && showEarnedTime && (
            <div className="small text-secondary mt-1">
              <FormatTime time={data.earned_at} />
            </div>
          )}
        </Card.Body>
      </Card>
    </OverlayTrigger>
  );
};

export default Index;
