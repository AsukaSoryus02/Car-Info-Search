# CarRag - 车辆信息管理系统

一个功能完善的车辆信息管理系统，用于记录、查询和管理车辆相关信息。系统采用前后端分离架构，提供友好的用户界面和高效的数据处理能力。

## 功能特点

- **车辆信息管理**：添加、编辑、删除和查看车辆信息
- **高级搜索**：按品牌、车型、燃油类型等多条件筛选车辆
- **数据可视化**：直观展示车辆数据统计和分析
- **响应式设计**：适配桌面和移动设备的界面
- **高性能**：采用缓存机制提升数据访问速度
- **数据安全**：文件备份和原子写入保证数据安全

## 技术栈

### 前端
- **框架**：Vue.js 3 (Composition API)
- **UI组件**：Ant Design Vue
- **路由**：Vue Router
- **HTTP客户端**：Axios
- **构建工具**：Vite

### 后端
- **语言**：Go
- **Web框架**：Gin
- **缓存**：Redis
- **存储**：JSON文件（带备份机制）
- **日志**：自定义多级日志系统

## 系统架构

```
+----------------+      +----------------+      +----------------+
|                |      |                |      |                |
|  前端应用      | <--> |  后端API服务   | <--> |  数据存储层    |
|  (Vue.js)      |      |  (Go/Gin)      |      |  (JSON/Redis)  |
|                |      |                |      |                |
+----------------+      +----------------+      +----------------+
```

## 项目结构

```
├── src/                  # 前端源代码
│   ├── api/              # API接口
│   │   └── carApi.js     # 车辆API接口
│   ├── components/       # 公共组件
│   │   └── AppHeader.vue # 应用头部组件
│   ├── data/             # 前端数据文件
│   │   └── carModels.json # 车型数据
│   ├── router/           # 路由配置
│   │   └── index.js      # 路由定义
│   ├── utils/            # 工具函数
│   │   └── formatter.js  # 格式化工具
│   ├── views/            # 前端视图组件
│   │   ├── CarForm.vue   # 车辆信息表单
│   │   ├── CarList.vue   # 车辆信息列表
│   │   └── Home.vue      # 首页
│   ├── App.vue           # 主应用组件
│   └── main.js           # 前端入口文件
├── server/               # 后端源代码
│   ├── config/           # 配置文件
│   │   └── config.go     # 应用配置
│   ├── controllers/      # 控制器
│   │   └── car_controller.go # 车辆控制器
│   ├── data/             # 后端数据存储目录
│   │   └── cars.json     # 车辆信息数据文件
│   ├── middleware/       # 中间件
│   │   └── logger.go     # 日志中间件
│   ├── models/           # 数据模型
│   │   └── car.go        # 车辆模型定义
│   ├── repositories/     # 数据访问层
│   │   └── file_repository.go # 文件存储实现
│   ├── utils/            # 工具类
│   │   ├── helpers.go    # 辅助函数
│   │   ├── logger.go     # 日志工具
│   │   ├── redis.go      # Redis缓存
│   │   └── storage.go    # 文件存储
│   ├── main.go           # 后端主程序
│   ├── go.mod            # Go模块定义
│   └── go.sum            # Go依赖校验文件
├── index.html            # HTML入口文件
├── package.json          # 前端依赖配置
├── package-lock.json     # 依赖版本锁定文件
└── vite.config.js        # Vite配置文件
```

## 安装和运行

### 前提条件

- Node.js 16+
- Go 1.18+
- Redis 6+ (可选，用于缓存)

### 前端

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

### 后端

```bash
# 进入后端目录
cd server

# 下载依赖
go mod download

# 运行后端服务
go run main.go

# 构建可执行文件
go build -o carrag-server
```

## API文档

### 车辆信息API

| 方法   | 路径                  | 描述                 | 参数                         |
|--------|----------------------|---------------------|------------------------------|
| GET    | /api/cars            | 获取所有车辆信息      | 无                           |
| GET    | /api/cars/:id        | 获取指定ID的车辆信息  | id: 车辆ID                   |
| POST   | /api/cars            | 创建新的车辆信息      | 请求体: 车辆信息JSON          |
| PUT    | /api/cars/:id        | 更新指定ID的车辆信息  | id: 车辆ID, 请求体: 更新数据  |
| DELETE | /api/cars/:id        | 删除指定ID的车辆信息  | id: 车辆ID                   |
| GET    | /api/cars/brand/:brand | 获取指定品牌的车辆   | brand: 车辆品牌              |

### 车辆信息数据结构

```json
{
  "id": "a1b2c3d4-5e6f78",       // 车辆唯一标识符
  "brand": "丰田",              // 品牌
  "model": "卡罗拉",            // 车型
  "fuelConsumption": 6.2,       // 油耗(L/100km)
  "fuelType": "汽油",           // 燃油类型
  "mileage": 15000,             // 行驶里程(km)
  "annualMileage": 12000,       // 年均行驶里程(km)
  "storageEnvironment": "车库",  // 存放环境
  "usageScenario": ["通勤", "家用"], // 使用场景
  "remarks": "车况良好",         // 备注
  "createdAt": "2023-01-01T12:00:00Z", // 创建时间
  "updatedAt": "2023-02-01T12:00:00Z"  // 更新时间
}
```

## 数据存储机制

系统采用分层存储架构：

1. **Redis缓存层**：用于存储热点数据，提高访问速度
2. **文件存储层**：使用JSON文件持久化存储数据
   - 采用原子写入机制确保数据完整性
   - 自动创建数据备份，防止意外损坏
   - 使用文件锁机制防止并发写入冲突

当Redis不可用时，系统会自动降级为仅使用文件存储，确保系统可用性。

## 日志系统

系统实现了多级日志记录机制：

- **DEBUG**：详细的调试信息，用于开发和排错
- **INFO**：常规操作信息，记录系统正常运行状态
- **WARNING**：警告信息，表示可能的问题但不影响系统运行
- **ERROR**：错误信息，表示系统遇到问题但可以继续运行
- **FATAL**：致命错误，导致系统无法继续运行
