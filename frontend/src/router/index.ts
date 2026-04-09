import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { enableRAG } from '@/config'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { guest: true, title: 'ログイン' },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/RegisterView.vue'),
      meta: { guest: true, title: '新規登録' },
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: () => import('@/views/ForgotPasswordView.vue'),
      meta: { guest: true, title: 'パスワードリセット' },
    },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: () => import('@/views/ResetPasswordView.vue'),
      meta: { guest: true, title: 'パスワード再設定' },
    },
    {
      path: '/',
      name: 'home',
      component: () => import('@/views/HomeView.vue'),
      meta: { auth: true, title: 'ホーム' },
    },
    {
      path: '/categories',
      name: 'categories',
      component: () => import('@/views/CategoryListView.vue'),
      meta: { auth: true, title: 'クイズ' },
    },
    {
      path: '/quiz/:sessionId',
      name: 'quiz',
      component: () => import('@/views/QuizSessionView.vue'),
      meta: { auth: true, title: 'クイズ' },
    },
    {
      path: '/quiz/:sessionId/result',
      name: 'quiz-result',
      component: () => import('@/views/QuizResultView.vue'),
      meta: { auth: true, title: 'クイズ結果' },
    },
    {
      path: '/quiz/review',
      name: 'quiz-review',
      component: () => import('@/views/QuizReviewView.vue'),
      meta: { auth: true, title: '復習' },
    },
    ...(enableRAG ? [{
      path: '/rag',
      name: 'rag',
      component: () => import('@/views/RagChatView.vue'),
      meta: { auth: true, title: '社内検索AI' },
    }] : []),
    {
      path: '/guides',
      name: 'guides',
      component: () => import('@/views/GuidesListView.vue'),
      meta: { auth: true, title: '実践ガイド' },
    },
    {
      path: '/guides/:id',
      name: 'guide-detail',
      component: () => import('@/views/GuideDetailView.vue'),
      meta: { auth: true, title: '実践ガイド' },
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/admin/AdminDashboardView.vue'),
      meta: { auth: true, admin: true, title: '管理画面' },
    },
    {
      path: '/admin/questions',
      name: 'admin-questions',
      component: () => import('@/views/admin/AdminQuestionsView.vue'),
      meta: { auth: true, admin: true, title: 'クイズ問題の編集' },
    },
    {
      path: '/admin/categories',
      name: 'admin-categories',
      component: () => import('@/views/admin/AdminCategoriesView.vue'),
      meta: { auth: true, admin: true, title: 'カテゴリの編集' },
    },
    {
      path: '/admin/documents',
      name: 'admin-documents',
      component: () => import('@/views/admin/AdminDocumentsView.vue'),
      meta: { auth: true, admin: true, title: 'RAG文書の登録' },
    },
    {
      path: '/admin/guides',
      name: 'admin-guides',
      component: () => import('@/views/admin/AdminGuidesView.vue'),
      meta: { auth: true, admin: true, title: 'ガイド記事の編集' },
    },
    {
      path: '/admin/user-progress',
      name: 'admin-user-progress',
      component: () => import('@/views/admin/AdminUserProgressView.vue'),
      meta: { auth: true, admin: true, title: '新人の学習状況' },
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  document.title = to.meta.title ? `${to.meta.title} - rootage` : 'rootage - 新人エンジニア育成アプリ'

  if (to.meta.auth && !auth.isLoggedIn) {
    return { name: 'login' }
  }
  if (to.meta.admin && !auth.isAdmin) {
    return { name: 'home' }
  }
  if (to.meta.guest && auth.isLoggedIn) {
    return { name: 'home' }
  }
})

export default router
