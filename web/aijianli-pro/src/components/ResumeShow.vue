<script setup>
import { computed } from 'vue';
import markdownit from 'markdown-it';
function parseText(input) {
  const trimmed = input.replace(/^:::\s*start\s*|\s*:::\s*end\s*$/gi, '');
  const sections = trimmed.split(/\s*:::\s*/).filter(section => section.trim() !== '');
  const result = sections.map(section => {
    return section.split('\n').filter(line => line.trim() !== '');
  });
  
  return result;
}

function FlexboxBlock(md) {
  // 识别 :::start ... :::end 区块
  function flexbox_container(state, startLine, endLine, silent) {
    const start = state.bMarks[startLine] + state.tShift[startLine];
    const max = state.eMarks[startLine];
    const marker = ':::';
    const markerLen = marker.length;

    // 以 :::start 开始
    if (state.src.slice(start, start + markerLen) !== marker) return false;
    const params = state.src.slice(start + markerLen, max).trim();
    if (!/^start\s*$/.test(params)) return false;

    // 查找结尾 :::end
    let nextLine = startLine + 1;
    let endMarkerLine = -1;
    while (nextLine < endLine) {
      const lineStart = state.bMarks[nextLine] + state.tShift[nextLine];
      const lineMax = state.eMarks[nextLine];
      if (
        state.src.slice(lineStart, lineStart + markerLen) === marker &&
        /^end\s*$/.test(state.src.slice(lineStart + markerLen, lineMax).trim())
      ) {
        endMarkerLine = nextLine;
        break;
      }
      nextLine++;
    }
    if (endMarkerLine === -1) return false;

    if (silent) return true;

    // 捕获整个 :::start 到 :::end 之间的原始文本
    const contentLines = state.getLines(startLine, endMarkerLine + 1, state.blkIndent, false);
    const token = state.push('flexbox_block', '', 0);
    token.content = contentLines;
    token.map = [startLine, endMarkerLine + 1];

    state.line = endMarkerLine + 1; // 跳过已处理行
    return true;
  }


  md.block.ruler.before('fence', 'flexbox_block', flexbox_container);

  // 渲染规则
  md.renderer.rules.flexbox_block = (tokens, idx) => {
    const content = tokens[idx].content;
    // 解析区块内容
    const groups = parseText(content);

    let html = `<div class="flexbox">\n`;
    for (const group of groups) {
      html += `  <div>\n`;
      for (const line of group) {
        html += `    <span>${md.renderInline(line)}</span>\n`;
      }
      html += `  </div>\n`;
    }
    html += `</div>\n`;
    return html;
  };
}
const md = markdownit({ html:true, linkify: true, typographer: true });
md.use(FlexboxBlock);

const input = `# 王xx
女 ｜ 28.000000岁 ｜ 药剂师 ｜ 本科 ｜ 138-1234-5678 ｜ wangjinlan@example.com
## 自我评价
具备扎实的药学背景知识，有较强的研发能力与实验技能，善于团队协作，积极面对挑战，能够快速适应并投入新工作环境。

## 工作经历

::: start
**广州药业**

:::
**药剂师**

:::
**2019年8月至今**

::: end
负责药品的研发和药效分析，制定并实施药物管理流程；

协助药品注册及相关审核工作，协调跨部门合作项目。

## 项目经历

::: start
**新型抗生素研发**

:::
**项目负责人**

:::
**领导一支由5人组成的团队，进行新型抗生素的研发**

:::
**2020年6月至2021年12月**

::: end
项目成功获得了地方研发经费支持，并在国家核心期刊发表相关研究论文。

## 专业技能
熟练掌握药品研发流程，熟悉GMP规范；

擅长使用HPLC等分析仪器；

具备良好的文献检索与数据分析能力。

## 教育背景与资质
::: start
**广东药科大学** | 药学 | 学士学位

:::
**2015年9月至2019年7月**

::: end
药师执业资格证，GMP培训合格证`;

const html = computed(() => md.render(input));
</script>

<template>
  <div class="markdown-body" v-html="html"></div>
</template>

<style>
.flexbox {
  display: flex;
  gap: 16px;
  margin-bottom: 8px;
}
.flexbox > div {
  flex: 1 1 0;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 4px 8px;
  background: #f3f3f3;
  border-radius: 4px;
}
.markdown-body h1,h2,h3,h4 {
  margin: 1em 0 0.5em;
}
.markdown-body p {
  margin: .5em 0;
}
</style>