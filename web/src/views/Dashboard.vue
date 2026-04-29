<template>
  <div class="min-h-screen bg-gray-50 p-4 md:p-8">
    <div class="max-w-6xl mx-auto">
      <div class="flex items-center justify-between mb-6 gap-4">
        <h1 class="text-xl md:text-3xl font-bold text-gray-800 truncate">短信管理面板</h1>
        <div class="flex items-center gap-2 flex-shrink-0">
            <button @click="fetchGroupedSMS" class="px-3 py-1.5 md:px-4 md:py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition text-sm whitespace-nowrap">刷新</button>
            <button @click="showChangePassword = true" class="px-3 py-1.5 md:px-4 md:py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition text-sm whitespace-nowrap">修改密码</button>
            <button @click="logout" class="px-3 py-1.5 md:px-4 md:py-2 bg-red-500 text-white rounded hover:bg-red-600 transition text-sm whitespace-nowrap">退出登录</button>
        </div>
      </div>

      <!-- Change Password Modal -->
      <div v-if="showChangePassword" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
        <div class="bg-white rounded-lg shadow-xl w-full max-w-md p-6">
          <h2 class="text-xl font-bold mb-4">修改密码</h2>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">旧密码</label>
              <input v-model="passwordForm.oldPassword" type="password" class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" placeholder="请输入旧密码">
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">新密码</label>
              <input v-model="passwordForm.newPassword" type="password" class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" placeholder="至少6个字符">
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">确认新密码</label>
              <input v-model="passwordForm.confirmPassword" type="password" class="w-full px-3 py-2 border rounded focus:ring-2 focus:ring-blue-500 outline-none" placeholder="请再次输入新密码">
            </div>
          </div>
          <div class="mt-6 flex justify-end gap-3">
            <button @click="closePasswordModal" class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded transition">取消</button>
            <button @click="changePassword" :disabled="isChanging" class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition disabled:opacity-50">
              {{ isChanging ? '保存中...' : '保存' }}
            </button>
          </div>
        </div>
      </div>
      
      <!-- 最新记录 -->
      <div v-if="latestSMS.length > 0" class="mb-8">
        <h2 class="text-lg font-bold text-gray-800 mb-4">🔥 最新记录</h2>
        <div class="bg-white shadow rounded-lg overflow-hidden border border-gray-200">
          <table class="min-w-full divide-y divide-gray-200 table-fixed">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-16">ID</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-24">机型</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-48">时间</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-32">发件人</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">内容</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-40">卡槽/手机号</th>
                <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider w-24">操作</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="sms in latestSMS" :key="sms.id" class="hover:bg-gray-50 transition">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ sms.id }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{{ sms.device || '-' }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ sms.sendTime }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{{ sms.sender || '-' }}</td>
                <td class="px-6 py-4 text-sm text-gray-900 break-words leading-relaxed">
                  {{ sms.content }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{{ sms.phone || '-' }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <button @click="deleteSMS(sms.id)" class="text-red-600 hover:text-red-900 transition">删除</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 分隔线 -->
      <div v-if="groupedSMS.length > 0" class="border-t-2 border-gray-300 my-8"></div>

      <!-- 分组列表 -->
      <div v-if="groupedSMS.length > 0" class="space-y-8">
        <div v-for="group in groupedSMS" :key="group.device" class="bg-white shadow rounded-lg border border-gray-200 p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-bold text-gray-800">
              📱 {{ group.device }} 
              <span class="text-sm font-normal text-gray-500">(共 {{ group.total }} 条)</span>
            </h3>
          </div>
          
          <!-- 分组表格 -->
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 table-fixed">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-16">ID</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-48">时间</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-32">发件人</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">内容</th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-40">卡槽/手机号</th>
                  <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider w-24">操作</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="sms in group.smsList" :key="sms.id" class="hover:bg-gray-50 transition">
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ sms.id }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ sms.sendTime }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{{ sms.sender || '-' }}</td>
                  <td class="px-6 py-4 text-sm text-gray-900 break-words leading-relaxed">
                    {{ sms.content }}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{{ sms.phone || '-' }}</td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button @click="deleteSMS(sms.id)" class="text-red-600 hover:text-red-900 transition">删除</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- 显示更多按钮 -->
          <div v-if="group.hasMore" class="mt-4 text-center">
            <button 
              @click="loadMore(group.device)" 
              :disabled="group.loading"
              class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition disabled:opacity-50"
            >
              {{ group.loading ? '加载中...' : '显示更多' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="latestSMS.length === 0 && groupedSMS.length === 0" class="bg-white shadow rounded-lg border border-gray-200 p-12 text-center text-gray-500">
        暂无短信
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRouter } from 'vue-router'

interface SMS {
  id: number
  sendTime: string
  content: string
  sender?: string
  phone?: string
  device?: string
}

interface GroupedSMS {
  device: string
  total: number
  hasMore: boolean
  smsList: SMS[]
  loading?: boolean
  offset?: number
}

const latestSMS = ref<SMS[]>([])
const groupedSMS = ref<GroupedSMS[]>([])
const router = useRouter()

// Change Password State
const showChangePassword = ref(false)
const isChanging = ref(false)
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const closePasswordModal = () => {
  showChangePassword.value = false
  passwordForm.value = {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  }
}

const changePassword = async () => {
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    alert('两次输入的新密码不一致')
    return
  }
  if (passwordForm.value.newPassword.length < 6) {
    alert('新密码长度不能少于6个字符')
    return
  }

  isChanging.value = true
  try {
    const token = localStorage.getItem('token')
    await axios.post('/api/change-password', {
      oldPassword: passwordForm.value.oldPassword,
      newPassword: passwordForm.value.newPassword
    }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    alert('密码修改成功！请重新登录。')
    logout()
  } catch (e) {
    console.error(e)
    const msg = axios.isAxiosError(e) ? e.response?.data?.error : '修改密码失败'
    alert(msg || '修改密码失败')
  } finally {
    isChanging.value = false
  }
}

const fetchGroupedSMS = async () => {
  try {
    const token = localStorage.getItem('token')
    if (!token) {
        router.push('/login')
        return
    }
    const res = await axios.get('/api/sms/grouped', {
      headers: { Authorization: `Bearer ${token}` }
    })
    latestSMS.value = res.data.latest || []
    groupedSMS.value = (res.data.groups || []).map((group: any) => ({
      ...group,
      loading: false,
      offset: 10
    }))
  } catch (e) {
    console.error(e)
    if (axios.isAxiosError(e) && e.response?.status === 401) {
        logout()
    }
  }
}

const loadMore = async (device: string) => {
  const group = groupedSMS.value.find(g => g.device === device)
  if (!group || group.loading) return

  group.loading = true
  try {
    const token = localStorage.getItem('token')
    const res = await axios.get('/api/sms/load-more', {
      params: {
        device: device,
        offset: group.offset || 10,
        limit: 10
      },
      headers: { Authorization: `Bearer ${token}` }
    })
    
    // 追加新数据
    group.smsList = [...group.smsList, ...res.data.smsList]
    group.hasMore = res.data.hasMore
    group.offset = (group.offset || 10) + 10
  } catch (e) {
    console.error(e)
    alert('加载失败')
  } finally {
    group.loading = false
  }
}

const deleteSMS = async (id: number) => {
  if (!confirm('确定要删除这条短信吗？')) return
  
  try {
    const token = localStorage.getItem('token')
    await axios.delete(`/api/sms/${id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    fetchGroupedSMS()
  } catch (e) {
    console.error(e)
    alert('删除短信失败')
  }
}

const logout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}

onMounted(() => {
  fetchGroupedSMS()
})
</script>
