import client from './client'

export interface RAGSource {
  document_id: string
  document_title: string
  chunk_index: number
  content: string
}

export interface RAGResponse {
  answer: string
  sources: RAGSource[]
}

export interface ChatMessage {
  role: string
  content: string
  sources?: RAGSource[]
  created_at: string
}

export interface ChatHistoryItem {
  session_id: string
  messages: ChatMessage[]
  created_at: string
}

export interface HistoryMessage {
  role: string
  content: string
}

export const askQuestion = (question: string, history?: HistoryMessage[]) =>
  client.post<RAGResponse>('/rag/ask', { question, history })
