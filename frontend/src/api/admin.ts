import client from './client'
import type { Category } from './quiz'

export interface Question {
  id: string
  category_id: string
  text: string
  choices: string[]
  correct_index: number
  explanation: string
  created_at: string
}

export interface Document {
  id: string
  title: string
  filename: string
  uploaded_at: string
}

export const getAdminCategories = () =>
  client.get<Category[]>('/admin/categories')

export const createCategory = (data: { name: string; description: string }) =>
  client.post<Category>('/admin/categories', data)

export const updateCategory = (id: string, data: { name: string; description: string }) =>
  client.put(`/admin/categories/${id}`, data)

export const deleteCategory = (id: string) =>
  client.delete(`/admin/categories/${id}`)

export const getAdminQuestions = (categoryId?: string) =>
  client.get<Question[]>('/admin/questions', { params: categoryId ? { category_id: categoryId } : {} })

export const createQuestion = (data: {
  category_id: string
  text: string
  choices: string[]
  correct_index: number
  explanation: string
  difficulty?: string
}) => client.post<Question>('/admin/questions', data)

export const updateQuestion = (id: string, data: {
  text: string
  choices: string[]
  correct_index: number
  explanation: string
  difficulty?: string
}) => client.put(`/admin/questions/${id}`, data)

export const deleteQuestion = (id: string) =>
  client.delete(`/admin/questions/${id}`)

export interface ImportResult {
  imported: number
  errors: { row: number; message: string }[]
}

export const importQuestions = (file: File) => {
  const formData = new FormData()
  formData.append('file', file)
  return client.post<ImportResult>('/admin/questions/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export const createDocumentFromText = (title: string, content: string) =>
  client.post<Document>('/admin/documents/text', { title, content })

export const getDocuments = () =>
  client.get<Document[]>('/admin/documents')

export const uploadDocument = (title: string, file: File) => {
  const formData = new FormData()
  formData.append('title', title)
  formData.append('file', file)
  return client.post<Document>('/admin/documents', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

export const deleteDocument = (id: string) =>
  client.delete(`/admin/documents/${id}`)

export const reindexDocument = (id: string) =>
  client.post(`/admin/documents/${id}/reindex`)

// Guide Categories

export interface GuideCategory {
  id: string
  name: string
  sort_order: number
}

export interface Guide {
  id: string
  category_id: string
  category_name: string
  title: string
  content: string
  is_published: boolean
  sort_order: number
  created_at: string
  updated_at: string
}

export const getAdminGuideCategories = () =>
  client.get<GuideCategory[]>('/admin/guide-categories')

export const createGuideCategory = (data: { name: string; sort_order: number }) =>
  client.post<GuideCategory>('/admin/guide-categories', data)

export const updateGuideCategory = (id: string, data: { name: string; sort_order: number }) =>
  client.put(`/admin/guide-categories/${id}`, data)

export const deleteGuideCategory = (id: string) =>
  client.delete(`/admin/guide-categories/${id}`)

export const getAdminGuides = () =>
  client.get<Guide[]>('/admin/guides')

export const createGuide = (data: { category_id: string; title: string; content: string; is_published: boolean; sort_order: number }) =>
  client.post<Guide>('/admin/guides', data)

export const updateGuide = (id: string, data: { category_id: string; title: string; content: string; is_published: boolean; sort_order: number }) =>
  client.put(`/admin/guides/${id}`, data)

export const deleteGuide = (id: string) =>
  client.delete(`/admin/guides/${id}`)

// Admin Summary

export interface AdminSummary {
  total_users: number
  total_questions: number
  total_documents: number
  total_guides: number
}

export const getAdminSummary = () =>
  client.get<AdminSummary>('/admin/summary')

// User Progress

export interface UserCategoryStatsItem {
  category_id: string
  category_name: string
  total_answers: number
  correct: number
  percentage: number
}

export interface UserProgressItem {
  id: string
  name: string
  email: string
  created_at: string
  last_quiz_at: string | null
  total_correct: number
  total_questions: number
  percentage: number
  category_stats: UserCategoryStatsItem[] | null
}

export const getUserProgress = () =>
  client.get<UserProgressItem[]>('/admin/user-progress')

export const resetPassword = (userId: string, password: string) =>
  client.post(`/admin/users/${userId}/reset-password`, { password })
