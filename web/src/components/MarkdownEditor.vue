<template>
  <div class="md-editor">
    <div class="md-toolbar">
      <el-button size="small" @click="prefixLine('## ')">H2</el-button>
      <el-button size="small" @click="prefixLine('### ')">H3</el-button>
      <el-button size="small" @click="wrap('**', '**', '加粗文字')">B</el-button>
      <el-button size="small" @click="wrap('*', '*', '斜体文字')">I</el-button>
      <el-button size="small" @click="wrap('`', '`', '代码')">代码</el-button>
      <el-button size="small" @click="wrap('[', '](https://)', '链接文字')">链接</el-button>
      <el-button size="small" @click="prefixLine('- ')">列表</el-button>
      <el-button size="small" @click="prefixLine('> ')">引用</el-button>
      <el-button size="small" @click="insertText('\n```\n代码\n```\n')">代码块</el-button>
      <el-button size="small" @click="insertText('\n---\n')">分割线</el-button>
      <el-button size="small" type="primary" plain :loading="uploading" @click="pickImage">
        插入图片
      </el-button>
      <span class="md-spacer" />
      <el-button-group>
        <el-button size="small" :type="mode === 'edit' ? 'primary' : undefined" @click="mode = 'edit'">
          编辑
        </el-button>
        <el-button size="small" :type="mode === 'preview' ? 'primary' : undefined" @click="mode = 'preview'">
          预览
        </el-button>
      </el-button-group>
    </div>

    <textarea
      v-show="mode === 'edit'"
      ref="textareaRef"
      class="md-textarea"
      :value="modelValue"
      :rows="rows"
      placeholder="使用 Markdown 写作，右上角可切换预览…"
      @input="onInput"
    />
    <div v-show="mode === 'preview'" class="md-preview" v-html="rendered" />
    <input ref="fileRef" type="file" accept="image/*" hidden @change="onImagePicked" />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { renderMarkdown } from '@/utils/markdown'
import { uploadFile } from '@/api/upload'

const props = withDefaults(defineProps<{ modelValue: string; rows?: number }>(), { rows: 18 })
const emit = defineEmits<{ 'update:modelValue': [value: string] }>()

const textareaRef = ref<HTMLTextAreaElement | null>(null)
const fileRef = ref<HTMLInputElement | null>(null)
const mode = ref<'edit' | 'preview'>('edit')
const uploading = ref(false)
const rendered = computed(() => renderMarkdown(props.modelValue))

function onInput(e: Event) {
  emit('update:modelValue', (e.target as HTMLTextAreaElement).value)
}

function cursor(): { start: number; end: number } {
  const el = textareaRef.value
  const len = props.modelValue.length
  return { start: el?.selectionStart ?? len, end: el?.selectionEnd ?? len }
}

function restore(pos: number, end = pos) {
  nextTick(() => {
    const el = textareaRef.value
    if (!el) return
    el.focus()
    el.setSelectionRange(pos, end)
  })
}

function wrap(before: string, after: string, placeholder: string) {
  const { start, end } = cursor()
  const value = props.modelValue
  const selected = value.slice(start, end) || placeholder
  emit('update:modelValue', value.slice(0, start) + before + selected + after + value.slice(end))
  restore(start + before.length, start + before.length + selected.length)
}

function insertText(text: string) {
  const { start, end } = cursor()
  const value = props.modelValue
  emit('update:modelValue', value.slice(0, start) + text + value.slice(end))
  restore(start + text.length)
}

function prefixLine(prefix: string) {
  const { start } = cursor()
  const value = props.modelValue
  const lineStart = value.lastIndexOf('\n', Math.max(0, start - 1)) + 1
  emit('update:modelValue', value.slice(0, lineStart) + prefix + value.slice(lineStart))
  restore(start + prefix.length)
}

function pickImage() {
  fileRef.value?.click()
}

async function onImagePicked(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  input.value = ''
  if (!file) return
  uploading.value = true
  try {
    const res = await uploadFile(file)
    insertText(`\n![${file.name}](${res.url})\n`)
    ElMessage.success('图片已插入')
  } catch {
    ElMessage.error('图片上传失败')
  } finally {
    uploading.value = false
  }
}
</script>

<style scoped>
.md-editor {
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
  overflow: hidden;
}

.md-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
  padding: 8px;
  background: var(--el-fill-color-light);
  border-bottom: 1px solid var(--el-border-color);
}

.md-spacer {
  flex: 1;
}

.md-textarea {
  display: block;
  width: 100%;
  box-sizing: border-box;
  padding: 14px 16px;
  border: 0;
  outline: none;
  resize: vertical;
  font: 14px/1.7 ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  color: var(--el-text-color-primary);
  background: var(--el-bg-color);
}

.md-preview {
  min-height: 320px;
  max-height: 720px;
  overflow: auto;
  padding: 14px 20px;
  line-height: 1.8;
  color: var(--el-text-color-primary);
}

/* v-html 注入的节点不带 scoped 属性，必须用 :deep 才能命中 */
.md-preview :deep(h1),
.md-preview :deep(h2),
.md-preview :deep(h3) {
  margin: 20px 0 10px;
  font-weight: 600;
  line-height: 1.4;
}
.md-preview :deep(h1) { font-size: 26px; }
.md-preview :deep(h2) { font-size: 22px; }
.md-preview :deep(h3) { font-size: 18px; }
.md-preview :deep(p) { margin: 10px 0; }
.md-preview :deep(ul),
.md-preview :deep(ol) { padding-left: 24px; margin: 10px 0; }
.md-preview :deep(blockquote) {
  margin: 12px 0;
  padding: 4px 14px;
  border-left: 4px solid var(--el-color-primary-light-5);
  color: var(--el-text-color-secondary);
  background: var(--el-fill-color-lighter);
}
.md-preview :deep(pre) {
  padding: 12px 14px;
  overflow: auto;
  border-radius: 4px;
  background: var(--el-fill-color-dark);
}
.md-preview :deep(code) {
  padding: 2px 5px;
  border-radius: 3px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  background: var(--el-fill-color);
}
.md-preview :deep(pre code) { padding: 0; background: transparent; }
.md-preview :deep(img) { max-width: 100%; border-radius: 4px; }
.md-preview :deep(a) { color: var(--el-color-primary); }
.md-preview :deep(hr) { margin: 20px 0; border: 0; border-top: 1px solid var(--el-border-color); }
</style>
