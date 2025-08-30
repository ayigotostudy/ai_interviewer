# 面试评价展示页面

## 功能概述

新创建的面试评价页面 (`/interview/evaluation/:id`) 提供了完整的面试评价数据可视化展示，包括：

### 1. 总体评价
- **仪表盘图表**：使用 ECharts 仪表盘展示综合得分
- **分数等级**：根据得分自动显示等级（优秀、良好、一般、及格、不及格）
- **详细信息**：显示满分值和当前得分

### 2. 胜任力维度分析
- **雷达图**：展示候选人在各个能力维度的得分情况
- **多维度评估**：包括沟通表达、逻辑思维、专业知识、学习能力、抗压性、团队合作等

### 3. 答题内容分析
- **关键词词云**：提取面试回答中的关键技术和概念
- **岗位匹配度**：环形图展示与职位描述的匹配程度
- **关键词分类**：区分匹配和缺失的关键词

## 技术实现

### 依赖库
- **ECharts 5.x**：用于图表渲染
- **Vue 3 Composition API**：响应式数据管理
- **TypeScript**：类型安全

### 图表类型
1. **仪表盘 (Gauge)**：总体得分展示
2. **雷达图 (Radar)**：能力维度分析
3. **词云图 (WordCloud)**：关键词提取
4. **环形图 (Pie)**：岗位匹配度

## 使用方法

### 1. 从面试详情页进入
- 在面试详情页面，点击"查看评价"按钮
- 仅当面试状态为"已完成"时显示

### 2. 从面试房间进入
- 在面试进行中，点击聊天区域右上角的"详细评价"按钮
- 可随时查看当前评价状态

### 3. 直接访问
- 访问路径：`/interview/evaluation/{面试ID}`

## 数据接口

### 主要接口
- `GET /meeting/remark?id={面试ID}`：获取面试评价数据

### 数据结构
```typescript
interface EvaluationData {
  overallEvaluation: {
    score: number;
    maxScore: number;
    rating: string;
    chartType: string;
  };
  competencyDimensions: {
    chartType: string;
    dimensions: Array<{
      name: string;
      score: number;
      fullMark: number;
    }>;
  };
  answerAnalysis: {
    keywordCloud: {
      chartType: string;
      keywords: Array<{
        text: string;
        value: number;
      }>;
    };
    jdMatch: {
      chartType: string;
      matchPercentage: number;
      matchedKeywords: string[];
      missingKeywords: string[];
    };
  };
}
```

## 页面特性

### 响应式设计
- 支持桌面端和移动端
- 图表自动适应容器大小
- 窗口大小变化时自动调整

### 加载状态
- 显示加载动画
- 错误处理和重试机制
- 无数据时的友好提示

### 交互功能
- 图表悬停效果
- 数据刷新按钮
- 返回导航

## 样式设计

### 视觉风格
- 现代化渐变背景
- 毛玻璃效果卡片
- 统一的色彩主题

### 布局结构
- 卡片式布局
- 清晰的信息层次
- 合理的间距和留白

## 扩展性

### 新增图表类型
- 在 `initCharts()` 方法中添加新的图表初始化函数
- 在模板中添加对应的图表容器
- 更新数据结构和接口调用

### 新增评价维度
- 在数据结构中添加新的评价字段
- 创建对应的图表组件
- 更新样式和布局

## 注意事项

1. **ECharts 依赖**：确保已安装 `echarts` 包
2. **数据格式**：接口返回的数据需要符合预定义的结构
3. **图表性能**：大量数据时考虑图表的分页或虚拟滚动
4. **浏览器兼容**：ECharts 5.x 需要现代浏览器支持

## 故障排除

### 常见问题
1. **图表不显示**：检查容器元素是否正确渲染
2. **数据加载失败**：查看网络请求和控制台错误
3. **样式异常**：确认 CSS 类名和样式规则

### 调试建议
- 使用浏览器开发者工具检查网络请求
- 查看控制台错误信息
- 验证数据结构是否符合预期
