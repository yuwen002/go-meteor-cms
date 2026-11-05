// 轻量 Markdown 渲染：零依赖，供文章预览与前台阅读页共用。
//
// 安全约束（正文是用户可编辑内容，必须以不可信输入对待）：
//  1. 所有文本先做 HTML 转义，再套用有限的 Markdown 语法，绝不把原始 HTML 放进 v-html；
//  2. 链接与图片地址只允许 http(s)、站内绝对路径与 # 锚点，阻断 javascript: 等协议；
//  3. 链接一律附加 rel="noopener noreferrer"。
//
// 支持：标题、粗体、斜体、行内代码、链接、图片、无序/有序列表、引用、
// 代码块、分割线、段落。不支持表格与原始 HTML（有意收敛攻击面）。

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')
}

function safeUrl(raw: string): string {
  const url = raw.trim()
  return /^(https?:\/\/|\/|#)/i.test(url) ? url : '#'
}

function inline(escaped: string): string {
  return escaped
    .replace(/!\[([^\]]*)\]\(([^)\s]+)\)/g, (_m: string, alt: string, url: string) => {
      return `<img src="${safeUrl(url)}" alt="${alt}" loading="lazy" />`
    })
    .replace(/\[([^\]]+)\]\(([^)\s]+)\)/g, (_m: string, label: string, url: string) => {
      return `<a href="${safeUrl(url)}" target="_blank" rel="noopener noreferrer">${label}</a>`
    })
    .replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
    .replace(/(^|[^*])\*([^*\n]+)\*/g, '$1<em>$2</em>')
    .replace(/`([^`]+)`/g, '<code>$1</code>')
}

/** 把 Markdown 源码渲染为可安全用于 v-html 的 HTML 字符串 */
export function renderMarkdown(src: string): string {
  const lines = String(src ?? '').replace(/\r\n/g, '\n').split('\n')
  const html: string[] = []
  let code: string[] | null = null
  let list: 'ul' | 'ol' | null = null
  let paragraph: string[] = []

  const flushParagraph = () => {
    if (paragraph.length > 0) {
      html.push(`<p>${inline(escapeHtml(paragraph.join(' ')))}</p>`)
      paragraph = []
    }
  }
  const flushList = () => {
    if (list) {
      html.push(`</${list}>`)
      list = null
    }
  }
  const flushCode = () => {
    if (code) {
      html.push(`<pre><code>${escapeHtml(code.join('\n'))}</code></pre>`)
      code = null
    }
  }

  for (const raw of lines) {
    const line = raw.trimEnd()
    if (/^```/.test(line.trim())) {
      flushParagraph()
      flushList()
      if (code) {
        flushCode()
      } else {
        code = []
      }
      continue
    }
    if (code) {
      code.push(raw)
      continue
    }

    const t = line.trim()
    if (t === '') {
      flushParagraph()
      flushList()
      continue
    }

    const heading = /^(#{1,6})\s+(.*)$/.exec(t)
    if (heading) {
      flushParagraph()
      flushList()
      const level = heading[1].length
      html.push(`<h${level}>${inline(escapeHtml(heading[2]))}</h${level}>`)
      continue
    }
    if (/^(-{3,}|\*{3,})$/.test(t)) {
      flushParagraph()
      flushList()
      html.push('<hr />')
      continue
    }
    if (t.startsWith('>')) {
      flushParagraph()
      flushList()
      html.push(`<blockquote>${inline(escapeHtml(t.slice(1).trim()))}</blockquote>`)
      continue
    }

    const ul = /^[-*+]\s+(.*)$/.exec(t)
    const ol = /^\d+[.)]\s+(.*)$/.exec(t)
    if (ul || ol) {
      flushParagraph()
      const type: 'ul' | 'ol' = ul ? 'ul' : 'ol'
      if (list !== type) {
        flushList()
        html.push(`<${type}>`)
        list = type
      }
      html.push(`<li>${inline(escapeHtml((ul ?? ol)![1]))}</li>`)
      continue
    }

    flushList()
    paragraph.push(t)
  }

  flushParagraph()
  flushList()
  flushCode()
  return html.join('\n')
}

/** 从 Markdown 源码生成纯文本摘要（列表页/SEO 兜底用） */
export function plainSummary(src: string, max = 120): string {
  const text = String(src ?? '')
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/!\[([^\]]*)\]\([^)]*\)/g, '$1')
    .replace(/\[([^\]]+)\]\([^)]*\)/g, '$1')
    .replace(/[#>*_`~-]/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
  return text.length > max ? `${text.slice(0, max)}…` : text
}
