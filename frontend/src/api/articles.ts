import { api } from './client'

export interface Owner {
  username: string
  claimed_at: string
}

export interface Article {
  id: number
  wikipedia_id: number
  title: string
  slug: string
  content: string
  content_length: number
  rarity_tier: 'common' | 'uncommon' | 'rare' | 'epic' | 'legendary'
  summary: string
  owner: Owner | null
}

export const articlesApi = {
  getRandom: () => api.get<Article>('/articles/random'),
  getBySlug: (slug: string) => api.get<Article>(`/articles/${slug}`),
}
