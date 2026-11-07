// 轻量 Markdown 渲染：零依赖，供后台预览与前台阅读共用。
//
// 安全约束（正文是用户可编辑内容，必须以不可信输入对待）：
//  1. 所有文本先做 HTML 转义，再套用有限的 Markdown 语法，绝不把原始 HTML 放进 v-html；
//  2. 链接与图片地址只允许 http(s)、站内绝对路径与 # 锚点，显式拒绝 // 协议相对地址
//     （否则 ![](//evil.com/t.gif) 会变成外部追踪请求），并阻断 javascript: 等协议；
//  3. 链接一律附加 rel="noopener noreferrer"。
//
// 支持：标题、粗体、斜体、删除线、行内代码、链接、图片、有序/无序/任务列表、
// 引用、代码块（含语言标注）、表格、分割线、段落。不支持原始 HTML（有意收敛攻击面）。

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
  // 协议相对 URL 以 // 开头，会被浏览器解析成外部域名，必须拒绝
  if (url.startsWith('//')) {
    return '#'
  }
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
    .replace(/~~([^~]+)~~/g, '<del>$1</del>')
    .replace(/(^|[^*])\*([^*\n]+)\*/g, '$1<em>$2</em>')
    .replace(/`([^`]+)`/g, '<code>$1</code>')
}

function isTableDivider(s: string): boolean {
  return s.includes('|') && s.includes('-') && /^[\s|:-]+$/.test(s)
}

function tableCells(line: string): string[] {
  return line
    .replace(/^\s*\|/, '')
    .replace(/\|\s*$/, '')
    .split('|')
    .map((cell) => cell.trim())
}

function renderTable(lines: string[]): string {
  const header = tableCells(lines[0])
  const body = lines.slice(2).map(tableCells)
  let out = '<table><thead><tr>'
  for (const cell of header) {
    out += `<th>${inline(escapeHtml(cell))}</th>`
  }
  out += '</tr></thead><tbody>'
  for (const row of body) {
    out += '<tr>'
    for (let i = 0; i < header.length; i++) {
      out += `<td>${inline(escapeHtml(row[i] ?? ''))}</td>`
    }
    out += '</tr>'
  }
  return out + '</tbody></table>'
}

/** 把 Markdown 源码渲染为可安全用于 v-html 的 HTML 字符串 */
export function renderMarkdown(src: string): string {
  const lines = String(src ?? '').replace(/\r\n/g, '\n').split('\n')
  const html: string[] = []
  let code: string[] | null = null
  let codeLang = ''
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
      const cls = codeLang ? ` class="language-${codeLang.replace(/[^\w+-]/g, '')}"` : ''
      html.push(`<pre><code${cls}>${escapeHtml(code.join('\n'))}</code></pre>`)
      code = null
      codeLang = ''
    }
  }

  for (let i = 0; i < lines.length; i++) {
    const raw = lines[i]
    const line = raw.trimEnd()

    const fence = /^```([\w+-]*)\s*$/.exec(line.trim())
    if (fence) {
      flushParagraph()
      flushList()
      if (code) {
        flushCode()
      } else {
        code = []
        codeLang = fence[1]
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

    // 表格：以 | 开头的行 + 下一行是分隔符行
    if (t.startsWith('|') && i + 1 < lines.length && isTableDivider(lines[i + 1].trim())) {
      flushParagraph()
      flushList()
      const table: string[] = []
      while (i < lines.length && lines[i].trim().startsWith('|')) {
        table.push(lines[i].trim())
        i++
      }
      i--
      html.push(renderTable(table))
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

    // 任务列表
    const task = /^[-*+]\s+\[([ xX])\]\s+(.*)$/.exec(t)
    if (task) {
      flushParagraph()
      if (list !== 'ul') {
        flushList()
        html.push('<ul class="task-list">')
        list = 'ul'
      }
      const checked = task[1].toLowerCase() === 'x' ? ' checked' : ''
      html.push(
        `<li class="task-item"><input type="checkbox" disabled${checked} /> ${inline(escapeHtml(task[2]))}</li>`,
      )
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
    .replace(/[#>*_`~|\-]/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
  return text.length > max ? `${text.slice(0, max)}…` : text
}
