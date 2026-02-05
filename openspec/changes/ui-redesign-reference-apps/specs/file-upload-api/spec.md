# Specification: File Upload API

## ADDED Requirements

### Requirement: Image Upload
系统 SHALL 提供API端点上传图片文件，支持项目封面、角色头像、地点图片。

#### Scenario: Upload image
- **WHEN** 客户端发送 POST /api/v1/upload 带图片文件（multipart/form-data）
- **THEN** 系统验证文件类型和大小，保存文件，返回文件URL和ID

#### Scenario: Validate file type
- **WHEN** 用户上传非图片文件
- **THEN** 系统返回400错误，说明仅支持JPG、PNG、WebP格式

#### Scenario: Validate file size
- **WHEN** 用户上传超过5MB的图片
- **THEN** 系统返回400错误，说明文件过大

#### Scenario: Generate unique filename
- **WHEN** 系统保存上传的图片
- **THEN** 系统生成唯一文件名（UUID + 原始扩展名），避免文件名冲突

### Requirement: File Storage
系统 SHALL 将上传的文件存储到本地文件系统或云存储。

#### Scenario: Store to local filesystem
- **WHEN** 系统配置为本地存储
- **THEN** 系统将文件保存到指定目录（如 /uploads/images/）

#### Scenario: Store to cloud storage
- **WHEN** 系统配置为云存储（S3、OSS等）
- **THEN** 系统将文件上传到云存储，返回云存储URL

#### Scenario: Organize by date
- **WHEN** 系统保存文件
- **THEN** 系统按日期组织目录结构（如 /uploads/2024/02/04/）

### Requirement: File Metadata
系统 SHALL 在数据库中记录文件元数据。

#### Scenario: Store file metadata
- **WHEN** 系统保存上传的文件
- **THEN** 系统在数据库中创建文件记录（ID、用户ID、文件名、大小、类型、URL、上传时间）

#### Scenario: Associate with entity
- **WHEN** 用户上传项目封面
- **THEN** 系统将文件ID关联到项目的cover_image_url字段

#### Scenario: Track file usage
- **WHEN** 系统保存文件元数据
- **THEN** 系统记录文件的用途（project_cover、character_avatar、location_image）

### Requirement: File Retrieval
系统 SHALL 提供API端点获取和下载文件。

#### Scenario: Get file by ID
- **WHEN** 客户端发送 GET /api/v1/files/:id
- **THEN** 系统返回文件元数据和访问URL

#### Scenario: Download file
- **WHEN** 客户端访问文件URL
- **THEN** 系统返回文件内容（Content-Type正确设置）

#### Scenario: File not found
- **WHEN** 客户端请求不存在的文件ID
- **THEN** 系统返回404错误

### Requirement: File Deletion
系统 SHALL 提供API端点删除文件。

#### Scenario: Delete file
- **WHEN** 客户端发送 DELETE /api/v1/files/:id
- **THEN** 系统删除文件记录和物理文件，返回204

#### Scenario: Check file ownership
- **WHEN** 用户尝试删除文件
- **THEN** 系统验证文件属于该用户，否则返回403错误

#### Scenario: Cascade deletion
- **WHEN** 用户删除项目
- **THEN** 系统自动删除关联的封面图片文件

#### Scenario: Soft delete
- **WHEN** 系统删除文件
- **THEN** 系统标记文件为已删除，延迟30天后物理删除

### Requirement: Image Processing
系统 SHALL 对上传的图片进行处理，生成缩略图和优化尺寸。

#### Scenario: Generate thumbnail
- **WHEN** 用户上传图片
- **THEN** 系统生成200x200的缩略图，用于列表显示

#### Scenario: Resize large images
- **WHEN** 用户上传超过2000x2000的图片
- **THEN** 系统自动缩放到最大2000x2000，保持宽高比

#### Scenario: Optimize file size
- **WHEN** 系统保存图片
- **THEN** 系统压缩图片质量到80%，减少文件大小

### Requirement: Security
系统 SHALL 确保文件上传的安全性。

#### Scenario: Authenticate upload
- **WHEN** 用户上传文件
- **THEN** 系统验证JWT令牌，确认用户身份

#### Scenario: Scan for malware
- **WHEN** 用户上传文件
- **THEN** 系统扫描文件内容，拒绝可疑文件

#### Scenario: Validate image content
- **WHEN** 用户上传图片
- **THEN** 系统验证文件确实是有效的图片格式，拒绝伪装的可执行文件

#### Scenario: Rate limiting
- **WHEN** 用户频繁上传文件
- **THEN** 系统限制上传频率（每分钟最多10个文件）

### Requirement: Storage Quota
系统 SHALL 限制用户的存储配额。

#### Scenario: Check quota
- **WHEN** 用户上传文件
- **THEN** 系统检查用户已使用的存储空间，超过配额则返回403错误

#### Scenario: Display quota usage
- **WHEN** 客户端请求用户存储信息
- **THEN** 系统返回已使用空间和总配额

#### Scenario: Upgrade quota
- **WHEN** 用户升级账户
- **THEN** 系统增加用户的存储配额
