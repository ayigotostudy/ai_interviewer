<template>
  <div class="resume-tool">
    <div class="toolbar">
      <label class="btn ghost round">
        <span class="icon">📥</span> 导入简历
        <input ref="fileInput" type="file" accept=".md,.txt" @change="onImport" />
      </label>
      <button class="btn" @click="loadTemplate"><span class="icon">📄</span> 简历模板</button>
      <button class="btn" @click="showHelp"><span class="icon">📘</span> 使用教程</button>
      <button class="btn" @click="toggleWatermark"><span class="icon">💧</span> 移除水印</button>
      <span class="spacer" />
      <span class="export-count">累计导出 {{ exportCount }} 次</span>
      <button class="btn" @click="saveToFile"><span class="icon">💾</span> 保存</button>
      <button class="btn primary round" @click="exportPdf"><span class="icon">🖨️</span> 导出 PDF</button>
      <button class="btn" @click="exportPng"><span class="icon">🖼️</span> 导出图片</button>
    </div>

    <div class="workspace">
      <div class="panel">
        <header>Markdown 编辑器</header>
        <textarea v-model="markdown" class="editor" placeholder="# 开始书写你的简历..." />
      </div>

      <div class="panel">
        <header>
          <div class="preview-bar">
            <span>实时预览</span>
            <div class="controls">
              <label>标题背景色 <input type="color" :value="titleStyle.bg" @input="onColor('bg', $event)" /></label>
              <label>文字颜色 <input type="color" :value="titleStyle.color" @input="onColor('color', $event)" /></label>
              <label>X偏移(px) <input type="number" v-model.number="titleStyle.shadowX" class="num" /></label>
              <label>Y偏移(px) <input type="number" v-model.number="titleStyle.shadowY" class="num" /></label>
              <label>模糊(px) <input type="number" v-model.number="titleStyle.shadowBlur" class="num" /></label>
              <label>阴影颜色 <input type="color" :value="titleStyle.shadowColor" @input="onColor('shadowColor', $event)" /></label>
              <label class="switch">
                <input type="checkbox" v-model="h2Bar" />
                <span>二级标题左侧竖线</span>
              </label>
              <button class="btn" @click="resetTitleStyle">重置</button>
            </div>
          </div>
        </header>
        <div class="preview" ref="previewEl">
          <ResumeShow :content="markdown" :titleStyle="titleStyle" :h2Bar="h2Bar" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watchEffect, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import ResumeShow from '@/components/ResumeShow.vue'
import { getResumeDetail } from '@/service/resume'

const fileInput = ref<HTMLInputElement>()
const previewEl = ref<HTMLElement>()

const LS_KEY = 'resume_tool_md'
const LS_EXPORT = 'resume_tool_export_count'

const route = useRoute()
const markdown = ref<string>('')
const titleStyle = ref<{ bg?: string; color?: string; shadowX?: number; shadowY?: number; shadowBlur?: number; shadowColor?: string }>({
  bg: '#ffffff00',
  color: '#000000',
  shadowX: 0,
  shadowY: 0,
  shadowBlur: 0,
  shadowColor: '#00000000'
})
const h2Bar = ref(false)
const exportCount = computed(() => Number(localStorage.getItem(LS_EXPORT) || '0'))

onMounted(async () => {
  const id = (route.query.id as string) || ''
  if (id) {
    try {
      const resp = await getResumeDetail(id)
      const data = resp?.data?.data || resp?.data
      if (data && data.content) {
        markdown.value = data.content
        localStorage.setItem(LS_KEY, markdown.value)
        return
      }
    } catch (e) { /* ignore and fallback */ }
  }
  const saved = localStorage.getItem(LS_KEY)
  if (saved) markdown.value = saved
  else markdown.value = defaultTemplate
})

const onImport = async (e: Event) => {
  const input = e.target as HTMLInputElement
  const f = input.files && input.files[0]
  if (!f) return
  const txt = await f.text()
  markdown.value = txt
  localStorage.setItem(LS_KEY, markdown.value)
}

const loadTemplate = () => {
  markdown.value = defaultTemplate
  localStorage.setItem(LS_KEY, markdown.value)
}

const showHelp = () => {
  alert('使用说明:\n1) 左侧输入 Markdown，右侧实时预览。\n2) 顶部可导入/模板/保存/导出。\n3) 反引号 `内容` 渲染为标签，::: start ... ::: end 构造成经历卡片。')
}

const toggleWatermark = () => {
  document.body.classList.toggle('no-watermark')
}

const saveToFile = () => {
  const blob = new Blob([markdown.value], { type: 'text/markdown' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = 'resume.md'
  a.click()
}

const bumpExport = () => {
  const n = Number(localStorage.getItem(LS_EXPORT) || '0') + 1
  localStorage.setItem(LS_EXPORT, String(n))
}

// 纯前端无额外库：使用浏览器打印为 PDF
const setPrintOnly = (on: boolean) => {
  if (!previewEl.value) return
  if (on) previewEl.value.classList.add('print-only')
  else previewEl.value.classList.remove('print-only')
}

const exportPdf = async () => {
  setPrintOnly(true)
  await nextTick()
  window.print()
  setTimeout(() => setPrintOnly(false), 300)
  bumpExport()
}

// 纯浏览器方案：将预览区域克隆 + 内联样式 -> SVG foreignObject -> Canvas -> PNG
const exportPng = async () => {
  if (!previewEl.value) return
  const node = previewEl.value
  const width = node.offsetWidth
  const height = node.scrollHeight

  const clone = node.cloneNode(true) as HTMLElement
  // 内联常用样式，保证 foreignObject 呈现一致
  const inline = (src: Element, dst: Element) => {
    const computed = window.getComputedStyle(src as Element)
    const cssText = [
      'font-family','font-size','font-weight','font-style','line-height','color','background','background-color','padding','padding-top','padding-right','padding-bottom','padding-left','margin','margin-top','margin-right','margin-bottom','margin-left','border','border-radius','box-shadow','text-align','white-space','letter-spacing','word-break','width','height','display','align-items','justify-content','gap'
    ]
    ;(dst as HTMLElement).style.background = 'white'
    cssText.forEach(k => {
      const v = computed.getPropertyValue(k)
      if (v) (dst as HTMLElement).style.setProperty(k, v)
    })
    Array.from(src.children).forEach((child, i) => {
      inline(child as Element, (dst as Element).children[i])
    })
  }
  inline(node, clone)

  const wrapper = document.createElement('div')
  wrapper.setAttribute('xmlns','http://www.w3.org/1999/xhtml')
  wrapper.style.background = '#fff'
  wrapper.appendChild(clone)

  const serialized = new XMLSerializer().serializeToString(wrapper)
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}"><foreignObject width="100%" height="100%">${serialized}</foreignObject></svg>`
  const blob = new Blob([svg], { type: 'image/svg+xml;charset=utf-8' })
  const url = URL.createObjectURL(blob)

  const img = new Image()
  img.onload = () => {
    const canvas = document.createElement('canvas')
    canvas.width = width
    canvas.height = height
    const ctx = canvas.getContext('2d')!
    ctx.fillStyle = '#ffffff'
    ctx.fillRect(0,0,width,height)
    ctx.drawImage(img, 0, 0)
    URL.revokeObjectURL(url)
    canvas.toBlob((b) => {
      if (!b) return
      const a = document.createElement('a')
      a.href = URL.createObjectURL(b)
      a.download = 'resume.png'
      a.click()
      URL.revokeObjectURL(a.href)
      bumpExport()
    }, 'image/png')
  }
  img.onerror = () => {
    alert('图片导出失败：请尝试使用浏览器打印为 PDF 或更新浏览器。')
    URL.revokeObjectURL(url)
  }
  img.src = url
}

// 同步到本地存储
watchEffect(() => {
  localStorage.setItem(LS_KEY, markdown.value)
})

const defaultTemplate = `# Go高级开发工程师 - 高并发架构专家
(+86)155-0799-7194 ｜ codecv@163.com ｜ 微信号：codecvresume

## 教育背景
::: start
**2015.09 - 2019.07**
:::
icon:fudan **复旦大学** \`985\` \`211\` \`双一流\`
:::
**计算机科学与技术专业 - 学士**
::: end

## 工作经历
::: start
**2021.04 - 至今**
:::
**字节跳动**
:::
**基础架构部**
:::
**Go高级开发工程师**
::: end
- **背景：** 公司业务快速增长，核心服务面临**百万QPS**挑战
- **职责：** 负责高性能RPC框架开发与优化
- **成果：** 性能提升40%，延迟降低35%
`

const resetTitleStyle = () => {
  titleStyle.value = {}
}

// 颜色选择器在拖动时部分浏览器会发出 input 事件带 hex/rgb，统一读取 raw value
const onColor = (key: 'bg'|'color'|'shadowColor', e: Event) => {
  const v = (e.target as HTMLInputElement).value
  titleStyle.value = { ...titleStyle.value, [key]: v }
}
</script>

<style scoped>
.resume-tool{min-height:100vh;background:#F3F4F6}
.toolbar{position:sticky;top:0;z-index:5;background:#fff;border-bottom:1px solid #E5E7EB;padding:.5rem 1rem;display:flex;gap:.5rem;align-items:center}
.toolbar .btn{display:inline-flex;align-items:center;gap:.4rem;background:linear-gradient(#ffffff,#f7f7f7);border:1px solid #E5E7EB;padding:.45rem .9rem;border-radius:10px;cursor:pointer;color:#111;box-shadow:0 1px 0 #ffffff inset,0 1px 2px rgba(0,0,0,.04);transition:all .15s ease}
.toolbar .btn:hover{border-color:#CBD5E1;box-shadow:0 2px 6px rgba(0,0,0,.06);transform:translateY(-1px)}
.toolbar .btn:active{transform:translateY(0);box-shadow:0 1px 3px rgba(0,0,0,.06)}
.toolbar .btn.primary{background:linear-gradient(#111,#000);color:#fff;border-color:#000;box-shadow:0 1px 0 rgba(255,255,255,.08) inset,0 4px 10px rgba(0,0,0,.15)}
.toolbar .btn.primary:hover{box-shadow:0 6px 14px rgba(0,0,0,.2);transform:translateY(-1px)}
.toolbar .btn.round{border-radius:999px}
.toolbar .btn.ghost{background:transparent;border-color:#E5E7EB}
.toolbar .btn.ghost:hover{background:#F8FAFC}
.toolbar .btn .icon{font-size:1rem;line-height:1}
.toolbar .spacer{flex:1}
.toolbar input[type=file]{display:none}

.workspace{max-width:1400px;margin:0 auto;padding:1rem;display:grid;grid-template-columns:1fr 1fr;gap:1rem}
.panel{background:#fff;border:1px solid #E5E7EB;border-radius:12px;display:flex;flex-direction:column;min-height:0}
.panel>header{padding:.75rem 1rem;border-bottom:1px solid #E5E7EB;font-weight:600}
.preview-bar{display:flex;justify-content:space-between;align-items:center;gap:1rem}
.controls{display:flex;flex-wrap:wrap;gap:.5rem;align-items:center}
.controls label{display:flex;gap:.25rem;align-items:center;font-weight:400}
.controls .num{width:80px}
.controls .switch{gap:.4rem;cursor:pointer}
.editor{flex:1;border:none;outline:none;padding:1rem;min-height:600px;font:13px/1.5 ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,monospace}
.preview{padding:1rem;overflow:auto}

/* 仅打印预览区域（隐藏其它） */
@media print {
  body * { visibility: hidden; }
  .preview, .preview * { visibility: visible; }
  .preview { position: absolute; left: 0; top: 0; right: 0; }
}

@media (max-width: 1024px){
  .workspace{grid-template-columns:1fr}
}
</style>

