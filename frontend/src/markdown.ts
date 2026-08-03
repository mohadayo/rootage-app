import DOMPurify from 'dompurify'
import { marked } from 'marked'

/**
 * Markdown を HTML に変換する。
 * marked は HTML をそのまま通すため、v-html に渡す前に必ずここでサニタイズする。
 */
export function renderMarkdown(content: string): string {
  return DOMPurify.sanitize(marked.parse(content, { async: false }))
}
