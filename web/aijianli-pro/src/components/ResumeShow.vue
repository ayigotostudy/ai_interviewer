<script setup>
import { computed } from 'vue';

// 定义props
const props = defineProps({
  content: {
    type: String,
    default: ''
  }
});

// 处理简历内容，转换为HTML
const processedContent = computed(() => {
  if (!props.content) return '';
  
  let content = props.content;
  
  // 1. 处理标题
  content = content.replace(/^# (.*$)/gm, '<h1 class="resume-title">$1</h1>');
  
  // 2. 处理二级标题
  content = content.replace(/^## (.*$)/gm, '<h2 class="section-title">$1</h2>');
  
  // 3. 处理联系信息行（包含 | 的行）
  content = content.replace(/^([^#\n]*\|[^#\n]*)$/gm, '<div class="contact-info">$1</div>');
  
  // 4. 处理 ::: start ... ::: end 区块 - 改进版本
  // 处理标准格式：::: start\n内容\n::: end
  content = content.replace(/::: start\s*\n([\s\S]*?)\n\s*::: end/g, (match, blockContent) => {
    console.log('找到标准:::区块:', blockContent);
    
    // 清理blockContent，移除多余的空行和空格
    const cleanedContent = blockContent.trim();
    const lines = cleanedContent.split('\n').filter(line => line.trim());
    console.log('处理后的行:', lines);
    
    if (lines.length >= 3) {
      const first = lines[0].replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
      const second = lines[1].replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
      const third = lines[2].replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
      const description = lines.slice(3).join(' ');
      
      const result = `
        <div class="experience-item">
          <div class="experience-header">
            <div class="company-info">
              <h4 class="company-name">${first}</h4>
              <span class="position">${second}</span>
            </div>
            <span class="duration">${third}</span>
          </div>
          ${description ? `<p class="description">${description}</p>` : ''}
        </div>
      `;
      
      console.log('生成的HTML:', result);
      return result;
    }
    console.log(':::区块行数不足，返回原始内容');
    return match;
  });
  
  // 处理紧凑格式：::: start 内容 ::: end
  content = content.replace(/::: start\s+([^:]+?)\s+::: end/g, (match, blockContent) => {
    console.log('找到紧凑:::区块:', blockContent);
    
    const lines = blockContent.split(/\s*\n\s*/).filter(line => line.trim());
    console.log('处理后的行:', lines);
    
    if (lines.length >= 3) {
      const first = lines[0].replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
      const second = lines[1].replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
      const third = lines[2].replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
      const description = lines.slice(3).join(' ');
      
      const result = `
        <div class="experience-item">
          <div class="experience-header">
            <div class="company-info">
              <h4 class="company-name">${first}</h4>
              <span class="position">${second}</span>
            </div>
            <span class="duration">${third}</span>
          </div>
          ${description ? `<p class="description">${description}</p>` : ''}
        </div>
      `;
      
      console.log('生成的HTML:', result);
      return result;
    }
    console.log('紧凑:::区块行数不足，返回原始内容');
    return match;
  });
  
  // 5. 处理粗体文本
  content = content.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
  
  // 6. 处理普通段落
  content = content.replace(/^([^#\n]*)$/gm, (match) => {
    if (match.trim() && !match.includes('|') && !match.includes(':::') && !match.startsWith('<')) {
      return `<p class="section-content">${match}</p>`;
    }
    return match;
  });
  
  // 7. 清理多余的空行
  content = content.replace(/\n\s*\n/g, '\n');
  
  return content;
});

console.log('原始内容:', props.content);
console.log('处理后的内容:', processedContent.value);
</script>

<template>
  <div class="resume-container">
    <div v-html="processedContent"></div>
  </div>
</template>

<style scoped>
.resume-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 60px;
  background: #FFFFFF;
  font-family: '阿里巴巴普惠体', 'PingFang SC', 'Helvetica Neue', 'Hiragino Sans GB', 'Microsoft YaHei', system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
  line-height: 1.6;
  color: #333333;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border-radius: 12px;
}

/* 标题样式 */
:deep(.resume-title) {
  font-size: 2.8rem;
  font-weight: 700;
  color: #333333;
  margin: 0 0 40px 0;
  text-align: center;
  letter-spacing: -0.02em;
  line-height: 1.2;
}

/* 分区标题样式 - 醒目的红色 */
:deep(.section-title) {
  font-size: 1.6rem;
  font-weight: 600;
  color: #EA0202;
  margin: 50px 0 25px 0;
  padding-bottom: 15px;
  border-bottom: 2px solid #EA0202;
  position: relative;
  line-height: 1.3;
}

:deep(.section-title::after) {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 80px;
  height: 2px;
  background: #EA0202;
}

/* 联系信息样式 */
:deep(.contact-info) {
  font-size: 1.1rem;
  color: #666666;
  font-weight: 500;
  text-align: center;
  margin-bottom: 40px;
  line-height: 1.8;
  padding: 20px;
  background: #FAFAFA;
  border-radius: 8px;
}

/* 段落内容样式 */
:deep(.section-content) {
  font-size: 1rem;
  color: #333333;
  margin: 0 0 20px 0;
  line-height: 1.8;
  text-align: justify;
}

/* 经历项目样式 */
:deep(.experience-item) {
  background: #FAFAFA;
  border-radius: 12px;
  padding: 30px;
  border-left: 4px solid #EA0202;
  margin-bottom: 25px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

:deep(.experience-item:hover) {
  transform: translateY(-3px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
  background: #F8F8F8;
}

:deep(.experience-header) {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

:deep(.company-info) {
  flex: 1;
}

:deep(.company-name) {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333333;
  margin: 0 0 12px 0;
  line-height: 1.4;
}

:deep(.position) {
  font-size: 1rem;
  font-weight: 500;
  color: #EA0202;
  background: #FFE6E6;
  padding: 6px 16px;
  border-radius: 20px;
  display: inline-block;
  line-height: 1.4;
}

:deep(.duration) {
  font-size: 0.9rem;
  color: #666666;
  font-weight: 500;
  background: #F0F0F0;
  padding: 8px 16px;
  border-radius: 8px;
  white-space: nowrap;
  line-height: 1.4;
}

:deep(.description) {
  font-size: 0.95rem;
  color: #555555;
  margin: 0;
  line-height: 1.8;
  text-align: justify;
}

/* 粗体文本样式 */
:deep(strong) {
  font-weight: 600;
  color: #333333;
}

/* 列表样式 */
:deep(ul) {
  margin: 20px 0;
  padding-left: 25px;
}

:deep(li) {
  margin-bottom: 12px;
  line-height: 1.8;
  color: #555555;
}

:deep(ul li) {
  list-style-type: disc;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .resume-container {
    padding: 30px 20px;
    margin: 10px;
  }
  
  :deep(.resume-title) {
    font-size: 2.2rem;
    margin-bottom: 30px;
  }
  
  :deep(.section-title) {
    font-size: 1.4rem;
    margin: 40px 0 20px 0;
  }
  
  :deep(.experience-header) {
    flex-direction: column;
    gap: 15px;
  }
  
  :deep(.duration) {
    align-self: flex-start;
  }
  
  :deep(.experience-item) {
    padding: 25px 20px;
  }
}

@media (max-width: 480px) {
  .resume-container {
    padding: 20px 15px;
  }
  
  :deep(.resume-title) {
    font-size: 1.8rem;
  }
  
  :deep(.section-title) {
    font-size: 1.2rem;
  }
}
</style>