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

import { useState } from 'react';
import { Row, Col, Card, Button, Form, Stack, Accordion } from 'react-bootstrap';
import { useSearchParams, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { usePageTags, useSkeletonControl } from '@/hooks';
import { Tag, Pagination, QueryGroup, TagsLoader, Icon } from '@/components';
import { formatCount, escapeRemove } from '@/utils';
import { tryNormalLogged } from '@/utils/guard';
import { useQueryTags, following, useQueryTagGroups } from '@/services';
import { loggedUserInfoStore } from '@/stores';
import type * as Type from '@/common/interface';

const sortBtns = ['popular', 'name', 'newest'];

const Tags = () => {
  const [urlSearch, setUrlSearch] = useSearchParams();
  const { t } = useTranslation('translation', { keyPrefix: 'tags' });
  const [searchTag, setSearchTag] = useState('');
  const [viewMode, setViewMode] = useState<'list' | 'group'>('list');
  const { role_id } = loggedUserInfoStore((_) => _.user);

  const page = Number(urlSearch.get('page')) || 1;
  const sort = urlSearch.get('sort') || sortBtns[0];

  const pageSize = 20;
  const {
    data: tags,
    mutate,
    isLoading,
  } = useQueryTags({
    page,
    page_size: pageSize,
    ...(searchTag ? { slug_name: searchTag } : {}),
    ...(sort ? { query_cond: sort } : {}),
  });

  const {
    data: tagGroups,
    mutate: mutateGroups,
    isLoading: groupsLoading,
  } = useQueryTagGroups();

  const { isSkeletonShow } = useSkeletonControl(isLoading || groupsLoading);

  const processTagGroups = (groups: Type.TagGroup[] | undefined) => {
    if (!groups) return [];
    return groups.map((group) => {
      if (!group.group_id || group.group_id === '') {
        return {
          ...group,
          group_name: t('other_group'),
        };
      }
      return group;
    });
  };

  const processedGroups = processTagGroups(tagGroups);

  const handleFollowGroup = (tag: Type.TagInfo) => {
    if (!tryNormalLogged(true)) {
      return;
    }
    following({
      object_id: tag.tag_id,
      is_cancel: tag.is_follower,
    }).then(() => {
      mutateGroups();
    });
  };

  const handleChange = (e) => {
    setSearchTag(e.target.value);
    setUrlSearch((prev) => {
      prev.set('page', '1');
      return prev;
    });
  };

  const handleFollow = (tag) => {
    if (!tryNormalLogged(true)) {
      return;
    }
    following({
      object_id: tag.tag_id,
      is_cancel: tag.is_follower,
    }).then(() => {
      mutate();
    });
  };

  usePageTags({
    title: t('tags', { keyPrefix: 'page_title' }),
  });

  return (
    <Row className="py-4 mb-4">
      <Col xxl={12}>
        <h3 className="mb-4">{t('title')}</h3>
        <div className="d-block d-sm-flex justify-content-between align-items-center flex-wrap">
          <Stack direction="horizontal" gap={3} className="mb-3 mb-sm-0">
            <Form>
              <Form.Group controlId="formBasicEmail">
                <Form.Control
                  value={searchTag}
                  placeholder={t('search_placeholder')}
                  type="search"
                  onChange={handleChange}
                  size="sm"
                />
              </Form.Group>
            </Form>
            {role_id === 2 || role_id === 3 ? (
              <Link
                className="btn btn-outline-primary btn-sm"
                to="/tags/create">
                {t('title', { keyPrefix: 'tag_modal' })}
              </Link>
            ) : null}
            <Button
              variant="outline-secondary"
              size="sm"
              onClick={() => setViewMode(viewMode === 'list' ? 'group' : 'list')}>
              <Icon
                name={viewMode === 'list' ? 'grid-3x3' : 'list'}
                className="me-1"
              />
              {viewMode === 'list' ? t('view_by_group') : t('view_by_list')}
            </Button>
          </Stack>
          {viewMode === 'list' && (
            <QueryGroup
              data={sortBtns}
              currentSort={sort || 'popular'}
              sortKey="sort"
              i18nKeyPrefix="tags.sort_buttons"
            />
          )}
        </div>
      </Col>

      <Col className="mt-4" xxl={12}>
        {viewMode === 'list' ? (
          <>
            <Row>
              {isSkeletonShow ? (
                <TagsLoader />
              ) : (
                tags?.list?.map((tag) => (
                  <Col
                    key={tag.slug_name}
                    xl={3}
                    lg={4}
                    md={4}
                    sm={6}
                    xs={12}
                    className="mb-4">
                    <Card className="h-100">
                      <Card.Body className="d-flex flex-column align-items-start">
                        <Tag className="mb-3" data={tag} showGroup />

                        <div className="small flex-fill text-break text-wrap text-truncate-3 reset-p mb-3">
                          {escapeRemove(tag.excerpt)}
                        </div>
                        <div className="d-flex align-items-center">
                          <Button
                            className={`me-2 ${tag.is_follower ? 'active' : ''}`}
                            variant="outline-primary"
                            size="sm"
                            onClick={() => handleFollow(tag)}>
                            {tag.is_follower
                              ? t('button_following')
                              : t('button_follow')}
                          </Button>
                          <span className="text-secondary small text-nowrap">
                            {formatCount(tag.question_count)} {t('tag_label')}
                          </span>
                        </div>
                      </Card.Body>
                    </Card>
                  </Col>
                ))
              )}
            </Row>
            <div className="d-flex justify-content-center">
              <Pagination
                currentPage={page}
                totalSize={tags?.count || 0}
                pageSize={pageSize}
              />
            </div>
          </>
        ) : (
          <Accordion defaultActiveKey="0" flush>
            {isSkeletonShow ? (
              <TagsLoader />
            ) : (
              processedGroups.map((group, groupIndex) => (
                <Accordion.Item
                  eventKey={String(groupIndex)}
                  key={group.group_id || `group-${groupIndex}`}>
                  <Accordion.Header>
                    <div className="d-flex align-items-center w-100 justify-content-between pe-3">
                      <span>
                        {group.group_name}
                        <span className="ms-2 text-muted small">
                          ({group.tags?.length || 0})
                        </span>
                      </span>
                    </div>
                  </Accordion.Header>
                  <Accordion.Body>
                    <Row>
                      {group.tags?.map((tag) => (
                        <Col
                          key={tag.slug_name}
                          xl={3}
                          lg={4}
                          md={4}
                          sm={6}
                          xs={12}
                          className="mb-4">
                          <Card className="h-100">
                            <Card.Body className="d-flex flex-column align-items-start">
                              <Tag className="mb-3" data={tag} />

                              <div className="small flex-fill text-break text-wrap text-truncate-3 reset-p mb-3">
                                {escapeRemove(tag.excerpt)}
                              </div>
                              <div className="d-flex align-items-center">
                                <Button
                                  className={`me-2 ${tag.is_follower ? 'active' : ''}`}
                                  variant="outline-primary"
                                  size="sm"
                                  onClick={() => handleFollowGroup(tag)}>
                                  {tag.is_follower
                                    ? t('button_following')
                                    : t('button_follow')}
                                </Button>
                                <span className="text-secondary small text-nowrap">
                                  {formatCount(tag.question_count)}{' '}
                                  {t('tag_label')}
                                </span>
                              </div>
                            </Card.Body>
                          </Card>
                        </Col>
                      ))}
                    </Row>
                  </Accordion.Body>
                </Accordion.Item>
              ))
            )}
          </Accordion>
        )}
      </Col>
    </Row>
  );
};

export default Tags;
