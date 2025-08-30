<script setup>
import { computed } from 'vue';

// 定义props
const props = defineProps({
  content: {
    type: String,
    default: ''
  },
  // 自定义标题样式（黑白为默认）
  titleStyle: { type: Object, default: () => ({}) },
  // 是否启用 h2 左侧蓝色竖线装饰
  h2Bar: { type: Boolean, default: false }
});

// 处理简历内容，转换为HTML（增强版）
const processedContent = computed(() => {
  if (!props.content) return '';

  let content = props.content;

  // 标题
  content = content.replace(/^# (.*$)/gm, '<h1 class="resume-title">$1</h1>');
  content = content.replace(/^## (.*$)/gm, '<h2 class="section-title">$1</h2>');

  // 联系信息：支持半角|与全角｜
  content = content.replace(/^([^#\n]*[\|｜][^#\n]*)$/gm, '<div class="contact-info">$1</div>');

  // ::: 区块（兼容形式：以 ::: 为分隔符的多行）
  content = content.replace(/::: start[\t ]*\n([\s\S]*?)\n[\t ]*::: end/g, (match, inner) => {
    const normalized = inner
      .replace(/\n:::\n/g, '\n') // 把“:::”行当作分隔符去掉
      .replace(/\n:::\s*$/g, '')
      .trim();
    const lines = normalized.split(/\n+/).filter(Boolean);
    if (lines.length >= 3) {
      const first = lines[0].replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
      const second = lines[1].replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
      const third = lines[2].replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');
      const rest = lines.slice(3).join(' ');
      return `
        <div class="experience-item">
          <div class="experience-header">
            <div class="company-info">
              <h4 class="company-name">${first}</h4>
              <span class="position">${second}</span>
            </div>
            <span class="duration">${third}</span>
          </div>
          ${rest ? `<p class="description">${rest}</p>` : ''}
        </div>
      `;
    }
    return match;
  });

  // 行内 code 作为技术标签徽章
  content = content.replace(/`([^`]+?)`/g, '<span class="tag-badge">$1</span>');

  // 列表：将以 "- " 开头的行转 li，并把连续 li 包裹为 ul
  content = content.replace(/^- (.*)$/gm, '<li>$1</li>');
  content = content.replace(/(?:<li>[\s\S]*?<\/li>)(?:(?:\n)?<li>[\s\S]*?<\/li>)+/g, (m) => `<ul>${m}</ul>`);

  // 粗体
  content = content.replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>');

  // 普通段落：非标题、非区块、非列表
  content = content.replace(/^([^#\n<][^\n]*)$/gm, (match) => {
    if (match.trim() && !match.includes(':::') && !match.startsWith('<') && !/^<li>/.test(match)) {
      return `<p class="section-content">${match}</p>`;
    }
    return match;
  });

  // 清理多余空行
  content = content.replace(/\n{2,}/g, '\n');
  return content;
});

// 将样式映射到 CSS 变量，供 h1/h2 使用
const containerStyle = computed(() => {
  const s = props.titleStyle || {}
  const toPx = (v) => (v === 0 || v) ? (typeof v === 'number' ? `${v}px` : String(v)) : '0px'
  return {
    '--title-bg': s.bg || 'transparent',
    '--title-color': s.color || '#000',
    '--title-shadow-x': toPx(s.shadowX ?? 0),
    '--title-shadow-y': toPx(s.shadowY ?? 0),
    '--title-shadow-blur': toPx(s.shadowBlur ?? 0),
    '--title-shadow-color': s.shadowColor || 'rgba(0,0,0,0)'
  }
})

console.log('原始内容:', props.content);
console.log('处理后的内容:', processedContent.value);
</script>

<template>
  <div class="resume-container" :style="containerStyle" :class="{ 'h2-bar': props.h2Bar }">
    <div v-html="processedContent"></div>
  </div>
</template>

<style scoped>
.resume-container {
  max-width: 210mm;
  margin: 0 auto;
  padding: 24pt;
  background: #fff;
  font-family: 'Helvetica Neue', Arial, Helvetica, sans-serif;
  line-height: 1.5;
  color: #000;
}

/* 标题样式 */
:deep(.resume-title) {
  font-size: 18pt;
  font-weight: 700;
  margin: 10pt 0 14pt 0;
  text-align: center;
  line-height: 1.2;
  /* 可自定义 */
  color: var(--title-color) !important;
  background: var(--title-bg) !important;
  box-shadow: var(--title-shadow-x) var(--title-shadow-y) var(--title-shadow-blur) var(--title-shadow-color) !important;
  padding: 4pt 6pt;
}

/* 二级标题（黑白风） */
:deep(.section-title) {
  font-size: 14pt;
  font-weight: 600;
  margin: 16pt 0 8pt 0;
  padding-bottom: 2pt;
  border-bottom: 1pt solid #000;
  /* 可自定义（允许覆盖黑白风） */
  color: var(--title-color) !important;
  background: var(--title-bg) !important;
  box-shadow: var(--title-shadow-x) var(--title-shadow-y) var(--title-shadow-blur) var(--title-shadow-color) !important;
  padding-top: 2pt;
  padding-left: 4pt;
}

/* 联系信息样式 */
:deep(.contact-info) {
  font-size: 11pt;
  color: #333;
  font-weight: 400;
  text-align: center;
  margin-bottom: 12pt;
  line-height: 1.6;
}

/* 段落内容样式 */
:deep(.section-content) {
  font-size: 11pt;
  color: #000;
  margin: 0 0 6pt 0;
  line-height: 1.5;
}

/* 技术标签徽章 */
:deep(.tag-badge) {
  display: inline-block;
  padding: 2px 8px;
  margin: 0 6px 6px 0;
  background: #F1F5F9;
  color: #0F172A;
  border: 1px solid #E2E8F0;
  border-radius: 999px;
  font-size: 0.85rem;
}

/* 经历项目样式（黑白简洁） */
:deep(.experience-item) { background: transparent; border-left: 0; padding: 0; margin-bottom: 10pt; }

:deep(.experience-item:hover) {
  transform: translateY(-3px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.12);
  background: #F8F8F8;
}

:deep(.experience-header) { display:flex; justify-content:space-between; align-items:flex-start; margin-bottom: 6pt; }

:deep(.company-info) {
  flex: 1;
}

:deep(.company-name) { font-size: 12pt; font-weight: 600; color:#000; margin:0 0 4pt 0; line-height:1.3; }

:deep(.position) { font-size: 11pt; font-weight: 600; color:#000; }

:deep(.duration) { font-size: 10pt; color:#000; font-weight:600; white-space:nowrap; }

:deep(.description) { font-size: 11pt; color:#000; margin:0; line-height:1.5; }

/* 粗体文本样式 */
:deep(strong) { font-weight: 700; color:#000; }

/* 列表样式 */
:deep(ul), :deep(ol) { margin: 8pt 0 10pt 0; padding-left: 1em; }
:deep(li) { margin-bottom: 4pt; line-height: 1.5; color:#000; }
:deep(li)::marker { color:#000; }

/* —— 可选：二级标题左侧蓝色竖线（仅在容器带有 .h2-bar 时启用） —— */
.h2-bar :deep(h2) {
  display: flex;
  align-items: center;
  border-bottom: none !important;
  padding-bottom: 0 !important;
  font-size: 16px;
  font-weight: 600;
  margin-top: 1.2em;
  margin-bottom: 0.8em;
  color: #000;
}

.h2-bar :deep(h2::before) {
  content: "";
  display: inline-block;
  width: 4px;
  height: 1.2em;
  background-color: #2b6cb0;
  margin-right: 8px;
  border-radius: 1px;
}

/* 覆盖旧的 section-title 下边框 */
.h2-bar :deep(.section-title) { border-bottom: none !important; padding-bottom: 0 !important; }

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