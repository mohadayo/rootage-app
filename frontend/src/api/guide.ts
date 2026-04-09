import client from './client'

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

export const getPublishedGuides = () =>
  client.get<Guide[]>('/guides')

export const getGuide = (id: string) =>
  client.get<Guide>(`/guides/${id}`)
