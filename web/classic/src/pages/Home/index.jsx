/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React, { useContext, useEffect, useState } from 'react';
import {
  Button,
  Typography,
  Input,
  ScrollList,
  ScrollItem,
} from '@douyinfe/semi-ui';
import { API, showError, copy, showSuccess } from '../../helpers';
import { useIsMobile } from '../../hooks/common/useIsMobile';
import { API_ENDPOINTS } from '../../constants/common.constant';
import { StatusContext } from '../../context/Status';
import { useActualTheme } from '../../context/Theme';
import { marked } from 'marked';
import { useTranslation } from 'react-i18next';
import {
  IconGithubLogo,
  IconPlay,
  IconFile,
  IconCopy,
} from '@douyinfe/semi-icons';
import { Link } from 'react-router-dom';
import NoticeModal from '../../components/layout/NoticeModal';
import {
  Moonshot,
  OpenAI,
  XAI,
  Zhipu,
  Volcengine,
  Cohere,
  Claude,
  Gemini,
  Suno,
  Minimax,
  Wenxin,
  Spark,
  Qingyan,
  DeepSeek,
  Qwen,
  Midjourney,
  Grok,
  AzureAI,
  Hunyuan,
  Xinference,
} from '@lobehub/icons';

const { Text } = Typography;

const providerItems = [
  { name: 'Moonshot', render: () => <Moonshot size={40} /> },
  { name: 'OpenAI', render: () => <OpenAI size={40} /> },
  { name: 'xAI', render: () => <XAI size={40} /> },
  { name: 'Zhipu', render: () => <Zhipu.Color size={40} /> },
  { name: 'Volcengine', render: () => <Volcengine.Color size={40} /> },
  { name: 'Cohere', render: () => <Cohere.Color size={40} /> },
  { name: 'Claude', render: () => <Claude.Color size={40} /> },
  { name: 'Gemini', render: () => <Gemini.Color size={40} /> },
  { name: 'Suno', render: () => <Suno size={40} /> },
  { name: 'Minimax', render: () => <Minimax.Color size={40} /> },
  { name: 'Wenxin', render: () => <Wenxin.Color size={40} /> },
  { name: 'Spark', render: () => <Spark.Color size={40} /> },
  { name: 'Qingyan', render: () => <Qingyan.Color size={40} /> },
  { name: 'DeepSeek', render: () => <DeepSeek.Color size={40} /> },
  { name: 'Qwen', render: () => <Qwen.Color size={40} /> },
  { name: 'Midjourney', render: () => <Midjourney size={40} /> },
  { name: 'Grok', render: () => <Grok size={40} /> },
  { name: 'Azure AI', render: () => <AzureAI.Color size={40} /> },
  { name: 'Hunyuan', render: () => <Hunyuan.Color size={40} /> },
  { name: 'Xinference', render: () => <Xinference.Color size={40} /> },
  { name: '30+', isCount: true },
];

const Home = () => {
  const { t, i18n } = useTranslation();
  const [statusState] = useContext(StatusContext);
  const actualTheme = useActualTheme();
  const [homePageContentLoaded, setHomePageContentLoaded] = useState(false);
  const [homePageContent, setHomePageContent] = useState('');
  const [noticeVisible, setNoticeVisible] = useState(false);
  const isMobile = useIsMobile();
  const isDemoSiteMode = statusState?.status?.demo_site_enabled || false;
  const docsLink = statusState?.status?.docs_link || '';
  const serverAddress =
    statusState?.status?.server_address || `${window.location.origin}`;
  const endpointItems = API_ENDPOINTS.map((e) => ({ value: e }));
  const [endpointIndex, setEndpointIndex] = useState(0);
  const isChinese = i18n.language.startsWith('zh');

  const displayHomePageContent = async () => {
    setHomePageContent(localStorage.getItem('home_page_content') || '');
    const res = await API.get('/api/home_page_content');
    const { success, message, data } = res.data;
    if (success) {
      let content = data;
      if (!data.startsWith('https://')) {
        content = marked.parse(data);
      }
      setHomePageContent(content);
      localStorage.setItem('home_page_content', content);

      // 如果内容是 URL，则发送主题模式
      if (data.startsWith('https://')) {
        const iframe = document.querySelector('iframe');
        if (iframe) {
          iframe.onload = () => {
            iframe.contentWindow.postMessage({ themeMode: actualTheme }, '*');
            iframe.contentWindow.postMessage({ lang: i18n.language }, '*');
          };
        }
      }
    } else {
      showError(message);
      setHomePageContent('加载首页内容失败...');
    }
    setHomePageContentLoaded(true);
  };

  const handleCopyBaseURL = async () => {
    const ok = await copy(serverAddress);
    if (ok) {
      showSuccess(t('已复制到剪切板'));
    }
  };

  useEffect(() => {
    const checkNoticeAndShow = async () => {
      const lastCloseDate = localStorage.getItem('notice_close_date');
      const today = new Date().toDateString();
      if (lastCloseDate !== today) {
        try {
          const res = await API.get('/api/notice');
          const { success, data } = res.data;
          if (success && data && data.trim() !== '') {
            setNoticeVisible(true);
          }
        } catch (error) {
          console.error('获取公告失败:', error);
        }
      }
    };

    checkNoticeAndShow();
  }, []);

  useEffect(() => {
    displayHomePageContent().then();
  }, []);

  useEffect(() => {
    const timer = setInterval(() => {
      setEndpointIndex((prev) => (prev + 1) % endpointItems.length);
    }, 3000);
    return () => clearInterval(timer);
  }, [endpointItems.length]);

  return (
    <div className='w-full overflow-x-hidden'>
      <NoticeModal
        visible={noticeVisible}
        onClose={() => setNoticeVisible(false)}
        isMobile={isMobile}
      />
      {homePageContentLoaded && homePageContent === '' ? (
        <div className='home-default-page w-full overflow-hidden'>
          <section className='home-hero-stage'>
            <div className='home-ambient-layer' aria-hidden='true'>
              <div className='home-grid-plane' />
              <div className='home-data-ribbon home-data-ribbon-a' />
              <div className='home-data-ribbon home-data-ribbon-b' />
              <div className='home-scanline' />
            </div>

            <div className='home-hero-inner'>
              <div className='home-hero-main'>
                <h1
                  className={`home-hero-title ${isChinese ? 'home-hero-title-cn' : ''}`}
                >
                  <span>{t('统一的')}</span>
                  <span className='home-title-accent'>
                    {t('大模型接口网关')}
                  </span>
                </h1>
                <p className='home-hero-subtitle'>
                  {t('多模型统一接入，只需将基址替换为：')}
                </p>

                <div className='home-endpoint-console'>
                  <Input
                    readOnly
                    value={serverAddress}
                    className='home-endpoint-input'
                    size={isMobile ? 'default' : 'large'}
                    suffix={
                      <div className='home-endpoint-suffix'>
                        <div className='home-endpoint-picker'>
                          <ScrollList
                            bodyHeight={32}
                            style={{ border: 'unset', boxShadow: 'unset' }}
                          >
                            <ScrollItem
                              mode='wheel'
                              cycled={true}
                              list={endpointItems}
                              selectedIndex={endpointIndex}
                              onSelect={({ index }) => setEndpointIndex(index)}
                            />
                          </ScrollList>
                        </div>
                        <Button
                          type='primary'
                          onClick={handleCopyBaseURL}
                          icon={<IconCopy />}
                          className='home-copy-button'
                        />
                      </div>
                    }
                  />
                </div>

                <div className='home-actions'>
                  <Link to='/console' className='home-action-link'>
                    <Button
                      theme='solid'
                      type='primary'
                      size={isMobile ? 'default' : 'large'}
                      className='home-primary-action'
                      icon={<IconPlay />}
                    >
                      {t('获取密钥')}
                    </Button>
                  </Link>
                  {isDemoSiteMode && statusState?.status?.version ? (
                    <Button
                      size={isMobile ? 'default' : 'large'}
                      className='home-secondary-action'
                      icon={<IconGithubLogo />}
                      onClick={() =>
                        window.open(
                          'https://github.com/QuantumNous/new-api',
                          '_blank',
                        )
                      }
                    >
                      {statusState.status.version}
                    </Button>
                  ) : (
                    docsLink && (
                      <Button
                        size={isMobile ? 'default' : 'large'}
                        className='home-secondary-action'
                        icon={<IconFile />}
                        onClick={() => window.open(docsLink, '_blank')}
                      >
                        {t('文档')}
                      </Button>
                    )
                  )}
                </div>
              </div>

              <div className='home-provider-strip'>
                <Text type='tertiary' className='home-provider-label'>
                  {t('支持众多的大模型供应商')}
                </Text>
                <div className='home-provider-grid'>
                  {providerItems.map((item) => (
                    <div
                      key={item.name}
                      className={`home-provider-logo ${item.isCount ? 'home-provider-count' : ''}`}
                      title={item.name}
                    >
                      {item.isCount ? (
                        <Typography.Text className='home-provider-count-text'>
                          30+
                        </Typography.Text>
                      ) : (
                        item.render()
                      )}
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </section>
        </div>
      ) : (
        <div className='overflow-x-hidden w-full'>
          {homePageContent.startsWith('https://') ? (
            <iframe
              src={homePageContent}
              className='w-full h-screen border-none'
            />
          ) : (
            <div
              className='mt-[60px]'
              dangerouslySetInnerHTML={{ __html: homePageContent }}
            />
          )}
        </div>
      )}
    </div>
  );
};

export default Home;
