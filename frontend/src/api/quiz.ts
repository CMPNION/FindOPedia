import { api } from './client'

export interface QuizOption {
  label: string
  value: string
}

export interface QuizQuestion {
  id: number
  index: number
  type: 'multiple_choice' | 'true_false'
  question_text: string
  options: QuizOption[]
}

export interface QuizResult {
  status: 'passed' | 'failed'
  score: number
  correct_count: number
  total_count: number
  ownership_claimed: boolean
  owner?: { username: string; claimed_at: string }
}

export interface CollectionItem {
  slug: string
  title: string
  rarity_tier: string
  claimed_at: string
}

export interface Collection {
  username: string
  total: number
  articles: CollectionItem[]
}

export interface LeaderboardEntry {
  Username: string
  Count: number
}

export interface Leaderboards {
  Total: LeaderboardEntry[]
  ByRarity: Record<string, LeaderboardEntry[]>
}

export interface CooldownStatus {
  active: boolean
  next_claim: string
}

export const quizApi = {
  getOrGenerateQuestions: (slug: string, aiProvider?: string, aiApiKey?: string) =>
    api.post<{ questions: QuizQuestion[] }>(`/articles/${slug}/questions`, {
      ai_provider: aiProvider,
      ai_api_key: aiApiKey,
    }),
  submitAttempt: (slug: string, answers: { question_id: number; chosen_answer: string }[]) =>
    api.post<QuizResult>(`/articles/${slug}/attempt`, { answers }),
  getCollection: (username: string) =>
    api.get<Collection>(`/users/${username}/collection`),
  getLeaderboards: () =>
    api.get<Leaderboards>('/leaderboard'),
  getCooldown: () =>
    api.get<CooldownStatus>('/me/cooldown'),
}
