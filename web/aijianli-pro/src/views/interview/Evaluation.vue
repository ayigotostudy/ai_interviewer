<template>
  <div class="evaluation-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <button class="back-btn" @click="goBack">
          <span class="icon">←</span>
          返回
        </button>
        <div class="header-title">
          <h1>面试评价分析</h1>
          <p>{{ interviewInfo.candidate }} - {{ interviewInfo.position }}</p>
        </div>
        <div class="header-actions">
          <button class="refresh-btn" @click="loadEvaluation" :disabled="loading">
            <span class="icon">🔄</span>
            {{ loading ? '加载中...' : '刷新' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 主要内容区域 -->
    <div class="page-content">
      <!-- 总体评价卡片 -->
      <div v-if="evaluationData.overallEvaluation" class="evaluation-card overall-score">
        <div class="card-header">
          <h3>总体评价</h3>
          <div class="score-badge" :class="getScoreClass(evaluationData.overallEvaluation?.score)">
            {{ evaluationData.overallEvaluation?.score || 0 }}分
          </div>
        </div>
        <div class="card-content">
          <div class="score-chart" ref="overallScoreChart"></div>
          <div class="score-details">
            <div class="detail-item">
              <span class="label">等级：</span>
              <span class="value">{{ evaluationData.overallEvaluation?.rating || '暂无' }}</span>
            </div>
            <div class="detail-item">
              <span class="label">满分：</span>
              <span class="value">{{ evaluationData.overallEvaluation?.maxScore || 100 }}分</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 胜任力维度卡片 -->
      <div v-if="evaluationData.competencyDimensions" class="evaluation-card competency">
        <div class="card-header">
          <h3>胜任力维度分析</h3>
          <p>多维度能力评估雷达图</p>
        </div>
        <div class="card-content">
          <div class="radar-chart" ref="competencyChart"></div>
        </div>
      </div>

      <!-- 答题内容分析卡片 -->
      <div v-if="evaluationData.answerAnalysis" class="evaluation-card answer-analysis">
        <div class="card-header">
          <h3>答题内容分析</h3>
          <p>关键词提取与岗位匹配度</p>
        </div>
        <div class="card-content">
          <!-- 关键词分析 -->
          <div class="analysis-section">
            <h4>关键词权重分析</h4>
            <div class="wordcloud-chart" ref="keywordChart"></div>
          </div>
          
          <!-- JD匹配度 -->
          <div class="analysis-section">
            <h4>岗位匹配度</h4>
            <div class="jd-match-content">
              <div class="match-chart" ref="jdMatchChart"></div>
              <div class="match-details">
                <div class="match-item">
                  <span class="label">匹配度：</span>
                  <span class="value">{{ evaluationData.answerAnalysis?.jdMatch?.matchPercentage || 0 }}%</span>
                </div>
                <div class="keywords-lists">
                  <div class="keywords-group">
                    <h5>匹配关键词</h5>
                    <div class="keyword-tags">
                      <span 
                        v-for="keyword in evaluationData.answerAnalysis?.jdMatch?.matchedKeywords || []" 
                        :key="keyword"
                        class="keyword-tag matched"
                      >
                        {{ keyword }}
                      </span>
                    </div>
                  </div>
                  <div class="keywords-group">
                    <h5>缺失关键词</h5>
                    <div class="keyword-tags">
                      <span 
                        v-for="keyword in evaluationData.answerAnalysis?.jdMatch?.missingKeywords || []" 
                        :key="keyword"
                        class="keyword-tag missing"
                      >
                        {{ keyword }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-container">
        <div class="loading-spinner"></div>
        <p>正在加载评价数据...</p>
      </div>

      <!-- 无数据状态 -->
      <div v-if="!loading && !evaluationData.overallEvaluation" class="no-data">
        <div class="no-data-icon">📊</div>
        <h3>暂无评价数据</h3>
        <p>该面试尚未生成评价报告，请稍后再试</p>
        <button class="retry-btn" @click="loadEvaluation">重新加载</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import * as echarts from 'echarts'
import { getMeetingDetail, getMeetingEvaluation } from '@/service/meeting'

const router = useRouter()
const route = useRoute()

// 响应式数据
const loading = ref(false)
const evaluationData = ref<any>({})
const interviewInfo = ref({
  candidate: '',
  position: ''
})

// 图表引用
const overallScoreChart = ref<HTMLElement>()
const competencyChart = ref<HTMLElement>()
const keywordChart = ref<HTMLElement>()
const jdMatchChart = ref<HTMLElement>()

// 图表实例
let overallScoreInstance: echarts.ECharts | null = null
let competencyInstance: echarts.ECharts | null = null
let keywordInstance: echarts.ECharts | null = null
let jdMatchInstance: echarts.ECharts | null = null

// 返回上一页
const goBack = () => {
  router.back()
}





// 获取分数等级样式
const getScoreClass = (score: number) => {
  if (!score) return 'default'
  if (score >= 90) return 'excellent'
  if (score >= 80) return 'good'
  if (score >= 70) return 'fair'
  if (score >= 60) return 'pass'
  return 'poor'
}

// 加载面试评价数据
const loadEvaluation = async () => {
  const meetingId = route.params.id as string
  if (!meetingId) return

  loading.value = true
  try {
    // 获取面试详情
    const response = await getMeetingDetail(meetingId)
    if (response.data.code === 1000) {
      const meetingData = response.data.data
      interviewInfo.value = {
        candidate: meetingData.candidate || '候选人',
        position: meetingData.position || '职位'
      }
      
      // 尝试调用真实的面试评价接口
      try {
        const evaluationResponse = await getMeetingEvaluation(meetingId)
        
        // 标准响应结构：{ code: 1000, msg: "success", data: {...} }
        if (evaluationResponse.data && evaluationResponse.data.code === 1000 && evaluationResponse.data.data) {
          // 确保数据是对象而不是字符串
          let dataToSet = evaluationResponse.data.data
          if (typeof dataToSet === 'string') {
            try {
              dataToSet = JSON.parse(dataToSet)

            } catch (parseError) {
              console.error('JSON解析失败:', parseError)
              dataToSet = evaluationResponse.data.data
            }
          }
          
          evaluationData.value = dataToSet

          
          // 等待 DOM 更新后初始化图表
          await nextTick()
          initCharts()
        } else if (evaluationResponse.data && evaluationResponse.data.data) {
          // 如果状态码不是1000但有data字段，尝试使用
          let dataToSet = evaluationResponse.data.data
          if (typeof dataToSet === 'string') {
            try {
              dataToSet = JSON.parse(dataToSet)

            } catch (parseError) {
              console.error('备用方案：JSON解析失败:', parseError)
              dataToSet = evaluationResponse.data.data
            }
          }
          
          if (dataToSet && dataToSet.overallEvaluation) {
            evaluationData.value = dataToSet

            await nextTick()
            initCharts()
          } else {

            await loadMockEvaluationData()
          }
        } else {

          await loadMockEvaluationData()
        }
      } catch (error) {

        // 使用模拟数据作为备选
        await loadMockEvaluationData()
      }
    }
  } catch (error) {
    console.error('加载面试评价失败:', error)
    // 使用模拟数据作为备选
    await loadMockEvaluationData()
  } finally {
    loading.value = false
  }
}

// 加载模拟评价数据（实际项目中应该调用真实接口）
const loadMockEvaluationData = async () => {
  // 模拟 API 延迟
  await new Promise(resolve => setTimeout(resolve, 1000))
  
  evaluationData.value = {
    overallEvaluation: {
      score: 85,
      maxScore: 100,
      rating: "良好",
      chartType: "gauge"
    },
    competencyDimensions: {
      chartType: "radar",
      dimensions: [
        { name: "沟通表达", score: 82, fullMark: 100 },
        { name: "逻辑思维", score: 88, fullMark: 100 },
        { name: "专业知识", score: 90, fullMark: 100 },
        { name: "学习能力", score: 85, fullMark: 100 },
        { name: "抗压性", score: 80, fullMark: 100 },
        { name: "团队合作", score: 78, fullMark: 100 }
      ]
    },
    answerAnalysis: {
      keywordCloud: {
        chartType: "wordcloud",
        keywords: [
          { text: "Spring Boot", value: 32 },
          { text: "微服务", value: 28 },
          { text: "MySQL", value: 25 },
          { text: "分布式系统", value: 22 },
          { text: "问题解决", value: 19 },
          { text: "Java", value: 18 },
          { text: "架构设计", value: 16 },
          { text: "性能优化", value: 14 }
        ]
      },
      jdMatch: {
        chartType: "doughnut",
        matchPercentage: 76,
        matchedKeywords: ["Java", "Spring Boot", "MySQL", "微服务"],
        missingKeywords: ["Redis", "消息队列", "容器化", "Kubernetes"]
      }
    }
  }
  
  // 等待 DOM 更新后初始化图表
  await nextTick()
  console.log('模拟数据加载完成，开始初始化图表')
  initCharts()
}

// 初始化所有图表
const initCharts = () => {
  // 添加延迟确保DOM完全渲染
  setTimeout(() => {
    initOverallScoreChart()
    initCompetencyChart()
    initKeywordChart()
    initJdMatchChart()
  }, 100)
}

// 初始化总体得分图表（仪表盘）
const initOverallScoreChart = () => {
  if (!overallScoreChart.value) return
  
  overallScoreInstance = echarts.init(overallScoreChart.value)
  const score = evaluationData.value.overallEvaluation?.score || 0
  
  const option = {
    series: [{
      type: 'gauge',
      startAngle: 180,
      endAngle: 0,
      min: 0,
      max: 100,
      splitNumber: 10,
      itemStyle: {
        color: '#58D9F9',
        shadowColor: 'rgba(0, 138, 255, 0.45)',
        shadowBlur: 10,
        shadowOffsetX: 2,
        shadowOffsetY: 2
      },
      progress: {
        show: true,
        roundCap: true,
        width: 18
      },
      pointer: {
        icon: 'path://M2090.36389,615.30999 L2090.36389,615.30999 C2091.48372,615.30999 2092.40383,616.194028 2092.44859,617.312956 L2096.90698,728.755929 C2097.05155,732.369577 2094.2393,735.416212 2090.62566,735.56078 C2090.53845,735.564269 2090.45117,735.566014 2090.36389,735.566014 L2090.36389,735.566014 C2086.74736,735.566014 2083.81598,732.634631 2083.81598,729.018082 C2083.81598,728.930965 2083.81772,728.843912 2083.82221,728.756851 L2088.2806,617.312956 C2088.32536,616.194028 2089.24547,615.30999 2090.36389,615.30999 Z',
        length: '75%',
        width: 16,
        offsetCenter: [0, '5%']
      },
      axisLine: {
        roundCap: true,
        lineStyle: {
          width: 18
        }
      },
      axisTick: {
        splitNumber: 2,
        lineStyle: {
          width: 2,
          color: '#999'
        }
      },
      splitLine: {
        length: 12,
        lineStyle: {
          width: 3,
          color: '#999'
        }
      },
      axisLabel: {
        distance: 30,
        color: '#999',
        fontSize: 12
      },
      title: {
        offsetCenter: [0, '30%'],
        fontSize: 20,
        color: '#464646'
      },
      detail: {
        fontSize: 30,
        offsetCenter: [0, '70%'],
        valueAnimation: true,
        formatter: function (value: number) {
          return Math.round(value) + '分'
        },
        color: 'auto'
      },
      data: [{
        value: score,
        name: '综合得分'
      }]
    }]
  }
  
  overallScoreInstance.setOption(option)
}

// 初始化胜任力维度图表（雷达图）
const initCompetencyChart = () => {
  if (!competencyChart.value) return
  
  competencyInstance = echarts.init(competencyChart.value)
  const dimensions = evaluationData.value.competencyDimensions?.dimensions || []
  
  const option = {
    radar: {
      indicator: dimensions.map((dim: any) => ({
        name: dim.name,
        max: dim.fullMark
      })),
      radius: '65%',
      splitNumber: 4,
      axisName: {
        color: '#4285F4',
        fontSize: 12
      },
      splitLine: {
        lineStyle: {
          color: ['#ddd']
        }
      },
      splitArea: {
        show: false
      }
    },
    series: [{
      type: 'radar',
      data: [{
        value: dimensions.map((dim: any) => dim.score),
        name: '能力得分',
        areaStyle: {
          color: 'rgba(66, 133, 244, 0.3)'
        },
        lineStyle: {
          color: '#4285F4',
          width: 2
        },
        itemStyle: {
          color: '#4285F4'
        }
      }]
    }]
  }
  
  competencyInstance.setOption(option)
}

// 初始化关键词图表（使用柱状图替代词云图）
const initKeywordChart = () => {
  if (!keywordChart.value) return
  
  // 确保容器有正确的尺寸
  const container = keywordChart.value
  
  if (container.offsetWidth === 0 || container.offsetHeight === 0) {
    setTimeout(() => initKeywordChart(), 200)
    return
  }
  
  keywordInstance = echarts.init(container)
  const keywords = evaluationData.value.answerAnalysis?.keywordCloud?.keywords || []
  
  if (keywords.length === 0) {
    return
  }
  
  // 确保数据格式正确
  const processedKeywords = keywords.map((item: any) => {
    if (typeof item === 'string') {
      return { name: item, value: 10 }
    }
    return {
      name: item.text || item.name || '未知',
      value: item.value || 10
    }
  })
  
  // 使用柱状图替代词云图
  const option = {
    title: {
      text: '关键词权重分布',
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 'bold'
      }
    },
    tooltip: {
      show: true,
      formatter: function (params: any) {
        return `${params.name}: ${params.value}`
      }
    },
    grid: {
      left: '10%',
      right: '10%',
      top: '15%',
      bottom: '15%'
    },
    xAxis: {
      type: 'category',
      data: processedKeywords.map(item => item.name),
      axisLabel: {
        rotate: 45,
        fontSize: 12
      }
    },
    yAxis: {
      type: 'value',
      name: '权重值'
    },
    series: [{
      type: 'bar',
      data: processedKeywords.map(item => item.value),
      itemStyle: {
        color: function(params: any) {
          const colors = ['#667eea', '#764ba2', '#f093fb', '#f5576c', '#4facfe']
          return colors[params.dataIndex % colors.length]
        },
        borderRadius: [4, 4, 0, 0]
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowColor: 'rgba(0, 0, 0, 0.3)'
        }
      }
    }]
  }
  
  try {
    keywordInstance.setOption(option)
    
    // 强制重绘
    setTimeout(() => {
      keywordInstance?.resize()
    }, 100)
  } catch (error) {
    // 图表设置失败，静默处理
  }
}

// 初始化JD匹配度图表（环形图）
const initJdMatchChart = () => {
  if (!jdMatchChart.value) return
  
  jdMatchInstance = echarts.init(jdMatchChart.value)
  const matchPercentage = evaluationData.value.answerAnalysis?.jdMatch?.matchPercentage || 0
  
  const option = {
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      center: ['50%', '50%'],
      data: [
        {
          value: matchPercentage,
          name: '匹配',
          itemStyle: {
            color: '#52c41a'
          }
        },
        {
          value: 100 - matchPercentage,
          name: '不匹配',
          itemStyle: {
            color: '#f5f5f5'
          }
        }
      ],
      label: {
        show: false
      },
      emphasis: {
        itemStyle: {
          shadowBlur: 10,
          shadowOffsetX: 0,
          shadowColor: 'rgba(0, 0, 0, 0.5)'
        }
      }
    }]
  }
  
  jdMatchInstance.setOption(option)
}

// 监听窗口大小变化，调整图表大小
const handleResize = () => {
  overallScoreInstance?.resize()
  competencyInstance?.resize()
  keywordInstance?.resize()
  jdMatchInstance?.resize()
}

// 强制重绘所有图表
const forceRedraw = () => {
  setTimeout(() => {
    overallScoreInstance?.resize()
    competencyInstance?.resize()
    keywordInstance?.resize()
    jdMatchInstance?.resize()
  }, 500)
}

// 组件挂载
onMounted(() => {
  loadEvaluation()
  window.addEventListener('resize', handleResize)
  
  // 延迟强制重绘
  setTimeout(() => {
    forceRedraw()
  }, 1000)
})

// 组件卸载
onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  overallScoreInstance?.dispose()
  competencyInstance?.dispose()
  keywordInstance?.dispose()
  jdMatchInstance?.dispose()
})
</script>

<style scoped>
.evaluation-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.page-header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
  padding: 1.5rem 2rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  gap: 2rem;
}

.back-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 0.75rem 1.25rem;
  border-radius: 12px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.back-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.header-title h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
  color: #1a202c;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-title p {
  margin: 0.5rem 0 0 0;
  font-size: 1rem;
  color: #718096;
  font-weight: 500;
}

.refresh-btn {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  border: none;
  padding: 0.75rem 1.25rem;
  border-radius: 12px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(16, 185, 129, 0.3);
}

.refresh-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(16, 185, 129, 0.4);
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.page-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
  display: grid;
  gap: 2rem;
  grid-template-columns: 1fr;
}

.evaluation-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.card-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 1.5rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 1.4rem;
  font-weight: 600;
}

.card-header p {
  margin: 0;
  font-size: 0.9rem;
  opacity: 0.9;
}

.score-badge {
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 1.1rem;
  font-weight: 600;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
}

.score-badge.excellent { background: linear-gradient(135deg, #10b981 0%, #059669 100%); }
.score-badge.good { background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%); }
.score-badge.fair { background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%); }
.score-badge.pass { background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%); }
.score-badge.poor { background: linear-gradient(135deg, #6b7280 0%, #4b5563 100%); }
.score-badge.default { background: rgba(255, 255, 255, 0.2); }

.card-content {
  padding: 2rem;
}

.score-chart {
  width: 100%;
  height: 300px;
  margin-bottom: 1.5rem;
}

.score-details {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 1rem;
  background: #f8fafc;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.detail-item .label {
  font-weight: 500;
  color: #64748b;
}

.detail-item .value {
  font-weight: 600;
  color: #1e293b;
}

.radar-chart {
  width: 100%;
  height: 400px;
}

.analysis-section {
  margin-bottom: 2rem;
}

.analysis-section h4 {
  margin: 0 0 1rem 0;
  font-size: 1.2rem;
  font-weight: 600;
  color: #1f2937;
}

.wordcloud-chart {
  width: 100%;
  height: 300px;
  margin-bottom: 1.5rem;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #f9fafb;
  position: relative;
  overflow: hidden;
}

.jd-match-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  align-items: start;
}

.match-chart {
  width: 100%;
  height: 250px;
}

.match-details {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.match-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: #f8fafc;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
}

.match-item .label {
  font-weight: 500;
  color: #64748b;
}

.match-item .value {
  font-size: 1.2rem;
  font-weight: 600;
  color: #10b981;
}

.keywords-lists {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.keywords-group h5 {
  margin: 0 0 0.5rem 0;
  font-size: 1rem;
  font-weight: 600;
  color: #374151;
}

.keyword-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.keyword-tag {
  padding: 0.4rem 0.8rem;
  border-radius: 16px;
  font-size: 0.85rem;
  font-weight: 500;
  border: none;
}

.keyword-tag.matched {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.keyword-tag.missing {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
  border-top: 4px solid #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.no-data {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  text-align: center;
}

.no-data-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

.no-data h3 {
  margin: 0 0 0.5rem 0;
  font-size: 1.5rem;
  color: #374151;
}

.no-data p {
  margin: 0 0 1.5rem 0;
  color: #6b7280;
}

.retry-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}

.retry-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-content {
    padding: 1rem;
    gap: 1rem;
  }
  
  .header-content {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .jd-match-content {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  
  .score-details {
    grid-template-columns: 1fr;
  }
}
</style>
