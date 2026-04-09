import client from './client'

export interface Category {
  id: string
  name: string
  description: string
}

export interface QuizQuestion {
  id: string
  text: string
  choices: string[]
}

export interface QuizStartResponse {
  session_id: string
  questions: QuizQuestion[]
}

export interface AnswerResponse {
  is_correct: boolean
  correct_index: number
  explanation: string
}

export interface AnswerDetail {
  question_id: string
  question_text: string
  choices: string[]
  selected_index: number
  correct_index: number
  is_correct: boolean
  explanation: string
}

export interface QuizResultResponse {
  session_id: string
  score: number
  total: number
  percentage: number
  answers: AnswerDetail[]
}

export interface QuizHistoryItem {
  session_id: string
  category_id: string
  category_name: string
  score: number
  total: number
  percentage: number
  started_at: string
  finished_at: string | null
}

export interface ReviewItem {
  question_id: string
  question_text: string
  choices: string[]
  correct_index: number
  explanation: string
  category_name: string
}

export interface CategoryStatsItem {
  category_id: string
  category_name: string
  total_answers: number
  correct: number
  percentage: number
}

export interface QuizStatsResponse {
  total_sessions: number
  category_stats: CategoryStatsItem[]
  weak_category: string
}

export const getCategories = () =>
  client.get<Category[]>('/categories')

export const startQuiz = (categoryId: string, difficulty?: string) =>
  client.post<QuizStartResponse>('/quiz/start', { category_id: categoryId, difficulty: difficulty || '' })

export interface DailyStatsItem {
  date: string
  total: number
  correct: number
  percentage: number
}

export interface AchievementItem {
  category_id: string
  category_name: string
  difficulty: string
  best_score: number
  total: number
  percentage: number
  cleared: boolean
}

export interface UserDashboardResponse {
  total_sessions: number
  total_questions: number
  total_correct: number
  overall_percentage: number
  category_stats: CategoryStatsItem[]
  weak_category: string
  daily_stats: DailyStatsItem[]
  current_streak: number
  achievements: AchievementItem[]
}

export const getStats = () =>
  client.get<QuizStatsResponse>('/quiz/stats')

export const getUserDashboard = () =>
  client.get<UserDashboardResponse>('/users/me/stats')

export const submitAnswer = (sessionId: string, questionId: string, selectedIndex: number) =>
  client.post<AnswerResponse>('/quiz/answer', {
    session_id: sessionId,
    question_id: questionId,
    selected_index: selectedIndex,
  })

export const finishQuiz = (sessionId: string) =>
  client.post<QuizResultResponse>('/quiz/finish', { session_id: sessionId })

export const getReview = () =>
  client.get<ReviewItem[]>('/quiz/review')
