/**
 * Mock API Server - 用于前端开发测试
 * 内置管理员账号: admin / admin123
 */
import express from 'express';
import cors from 'cors';

const app = express();
const PORT = 8080;

app.use(cors());
app.use(express.json());

// 内置管理员账号
const ADMIN_USER = {
  id: '00000000-0000-0000-0000-000000000001',
  username: 'admin',
  email: 'admin@novelai.com',
  nickname: '管理员',
  avatar: '',
  role: 'admin',
  created_at: '2024-01-01T00:00:00Z',
};

// Mock JWT Token
const MOCK_TOKEN = 'mock-jwt-token-' + Date.now();

// Mock 项目数据
const mockProjects = [
  {
    id: '10000000-0000-0000-0000-000000000001',
    user_id: '00000000-0000-0000-0000-000000000001',
    title: '玄幻小说：修仙传',
    description: '一个关于修仙的故事',
    genre: '玄幻',
    status: 'writing',
    total_chapters: 10,
    total_words: 50000,
    cover_url: '',
    metadata: {},
    is_deleted: false,
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-15T00:00:00Z',
  },
  {
    id: '10000000-0000-0000-0000-000000000002',
    user_id: '00000000-0000-0000-0000-000000000001',
    title: '都市小说：重生之巅峰',
    description: '重生回到过去，改变命运',
    genre: '都市',
    status: 'writing',
    total_chapters: 6,
    total_words: 30000,
    cover_url: '',
    metadata: {},
    is_deleted: false,
    created_at: '2024-01-10T00:00:00Z',
    updated_at: '2024-01-20T00:00:00Z',
  },
];

// Mock 章节数据
const mockChapters = [
  {
    id: '20000000-0000-0000-0000-000000000001',
    project_id: '10000000-0000-0000-0000-000000000001',
    chapter_number: 1,
    title: '第一章：初入修仙界',
    content: '这是第一章的内容...',
    word_count: 5000,
    status: 'published',
    created_at: '2024-01-01T00:00:00Z',
  },
  {
    id: '20000000-0000-0000-0000-000000000002',
    project_id: '10000000-0000-0000-0000-000000000001',
    chapter_number: 2,
    title: '第二章：筑基之路',
    content: '这是第二章的内容...',
    word_count: 5200,
    status: 'published',
    created_at: '2024-01-02T00:00:00Z',
  },
];

// ==================== 认证接口 ====================

// 登录
app.post('/api/v1/auth/login', (req, res) => {
  const { username, email, password } = req.body;

  // 支持用户名或邮箱登录
  const loginIdentifier = username || email;

  if ((loginIdentifier === 'admin' || loginIdentifier === 'admin@novelai.com') && password === 'admin123') {
    res.json({
      code: 0,
      message: 'success',
      data: {
        token: MOCK_TOKEN,
        user: ADMIN_USER,
      },
    });
  } else {
    res.status(401).json({
      code: 401,
      message: '用户名或密码错误',
    });
  }
});

// 注册
app.post('/api/v1/auth/register', (req, res) => {
  res.json({
    code: 0,
    message: 'success',
    data: {
      token: MOCK_TOKEN,
      user: ADMIN_USER,
    },
  });
});

// 获取当前用户信息
app.get('/api/v1/auth/me', (req, res) => {
  res.json({
    code: 0,
    message: 'success',
    data: ADMIN_USER,
  });
});

// ==================== 项目接口 ====================

// 获取项目列表
app.get('/api/v1/projects', (req, res) => {
  res.json({
    code: 0,
    message: 'success',
    data: {
      items: mockProjects,
      total: mockProjects.length,
      page: 1,
      page_size: 20,
    },
  });
});

// 获取项目详情
app.get('/api/v1/projects/:id', (req, res) => {
  const project = mockProjects.find(p => p.id === req.params.id);
  if (project) {
    res.json({
      code: 0,
      message: 'success',
      data: project,
    });
  } else {
    res.status(404).json({
      code: 404,
      message: '项目不存在',
    });
  }
});

// 创建项目
app.post('/api/v1/projects', (req, res) => {
  const newProject = {
    id: '10000000-0000-0000-0000-' + Date.now(),
    user_id: '00000000-0000-0000-0000-000000000001',
    ...req.body,
    status: req.body.status || 'draft',
    total_chapters: 0,
    total_words: 0,
    cover_url: req.body.cover_url || '',
    metadata: req.body.metadata || {},
    is_deleted: false,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
  };
  mockProjects.push(newProject);
  res.json({
    code: 0,
    message: 'success',
    data: newProject,
  });
});

// 更新项目
app.put('/api/v1/projects/:id', (req, res) => {
  const index = mockProjects.findIndex(p => p.id === req.params.id);
  if (index !== -1) {
    mockProjects[index] = {
      ...mockProjects[index],
      ...req.body,
      updated_at: new Date().toISOString(),
    };
    res.json({
      code: 0,
      message: 'success',
      data: mockProjects[index],
    });
  } else {
    res.status(404).json({
      code: 404,
      message: '项目不存在',
    });
  }
});

// 删除项目
app.delete('/api/v1/projects/:id', (req, res) => {
  const index = mockProjects.findIndex(p => p.id === req.params.id);
  if (index !== -1) {
    mockProjects.splice(index, 1);
    res.json({
      code: 0,
      message: 'success',
    });
  } else {
    res.status(404).json({
      code: 404,
      message: '项目不存在',
    });
  }
});

// ==================== 章节接口 ====================

// 获取章节列表
app.get('/api/v1/projects/:projectId/chapters', (req, res) => {
  const chapters = mockChapters.filter(c => c.project_id === req.params.projectId);
  res.json({
    code: 0,
    message: 'success',
    data: {
      items: chapters,
      total: chapters.length,
    },
  });
});

// 获取章节详情
app.get('/api/v1/chapters/:id', (req, res) => {
  const chapter = mockChapters.find(c => c.id === req.params.id);
  if (chapter) {
    res.json({
      code: 0,
      message: 'success',
      data: chapter,
    });
  } else {
    res.status(404).json({
      code: 404,
      message: '章节不存在',
    });
  }
});

// 创建章节
app.post('/api/v1/projects/:projectId/chapters', (req, res) => {
  const newChapter = {
    id: '20000000-0000-0000-0000-' + Date.now(),
    project_id: req.params.projectId,
    ...req.body,
    status: 'draft',
    word_count: req.body.content?.length || 0,
    created_at: new Date().toISOString(),
  };
  mockChapters.push(newChapter);
  res.json({
    code: 0,
    message: 'success',
    data: newChapter,
  });
});

// ==================== 其他接口 ====================

// 健康检查
app.get('/healthz', (req, res) => {
  res.json({ status: 'ok' });
});

// 404 处理
app.use((req, res) => {
  res.status(404).json({
    code: 404,
    message: 'API not found',
  });
});

// 启动服务器
app.listen(PORT, () => {
  console.log(`
╔════════════════════════════════════════════════════════════╗
║                                                            ║
║   🚀 Mock API Server 已启动                                ║
║                                                            ║
║   地址: http://localhost:${PORT}                            ║
║                                                            ║
║   内置管理员账号:                                           ║
║   用户名: admin                                            ║
║   密码:   admin123                                         ║
║                                                            ║
║   前端地址: http://localhost:3000                          ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
  `);
});
