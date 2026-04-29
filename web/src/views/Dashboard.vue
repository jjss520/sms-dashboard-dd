<template>
  <div class="min-h-screen bg-gray-50 p-4 md:p-8">
    <div class="max-w-6xl mx-auto">
      <div class="flex items-center justify-between mb-6 gap-4">
        <h1 class="text-xl md:text-3xl font-bold text-gray-800 truncate">短信管理面板</h1>
        <div class="flex items-center gap-2 flex-shrink-0">
            <input v-model="searchKeyword" @keyup.enter="handleSearch" type="text" placeholder="搜索内容、发件人、手机号、机型..." class="px-3 py-1.5 border rounded text-sm w-64 outline-none focus:ring-2 focus:ring-blue-500">
            <button @click="handleSearch" class="px-3 py-1.5 md:px-4 md:py-2 bg-green-500 text-white rounded hover:bg-green-600 transition text-sm whitespace-nowrap">搜索</button>
            <button v-if="isSearching" @click="clearSearch" class="px-3 py-1.5 md:px-4 md:py-2 bg-gray-500 text-white rounded hover:bg-gray-600 transition text-sm whitespace-nowrap">清空</button>
            <button @click="fetchGroupedSMS" class="px-3 py-1.5 md:px-4 md:py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition text-sm whitespace-nowrap">刷新</button>
            <button v-if="selectedIds.length > 0" @click="batchDelete" class="px-3 py-1.5 md:px-4 md:py-2 bg-red-500 text-white rounded hover:bg-red-600 transition text-sm whitespace-nowrap">
              删除选中({{ selectedIds.length }})
            </button>
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
      
      <!-- 搜索结果 -->
      <div v-if="isSearching && searchResults.length > 0" class="mb-8">
        <h2 class="text-lg font-bold text-gray-800 mb-4">🔍 搜索结果 (共 {{ searchTotal }} 条)</h2>
        <div class="bg-white shadow rounded-lg overflow-hidden border border-gray-200">
          <table class="min-w-full divide-y divide-gray-200 table-fixed">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-12">
                  <input type="checkbox" @change="toggleSelectAll('search', $event)" :checked="isAllSelected('search')" class="w-4 h-4">
                </th>
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
              <tr v-for="sms in searchResults" :key="sms.id" class="hover:bg-gray-50 transition">
                <td class="px-6 py-4 whitespace-nowrap">
                  <input type="checkbox" :value="sms.id" v-model="selectedIds" class="w-4 h-4">
                </td>
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
        <!-- 分页 -->
        <div v-if="searchTotal > searchPageSize" class="mt-4 flex justify-center gap-2">
          <button @click="searchPage--" :disabled="searchPage <= 1" class="px-3 py-1 bg-gray-200 rounded disabled:opacity-50">上一页</button>
          <span class="px-3 py-1">第 {{ searchPage }} / {{ Math.ceil(searchTotal / searchPageSize) }} 页</span>
          <button @click="searchPage++" :disabled="searchPage >= Math.ceil(searchTotal / searchPageSize)" class="px-3 py-1 bg-gray-200 rounded disabled:opacity-50">下一页</button>
        </div>
      </div>

      <!-- 搜索无结果 -->
      <div v-if="isSearching && searchResults.length === 0" class="bg-white shadow rounded-lg border border-gray-200 p-12 text-center text-gray-500 mb-8">
        未找到匹配的短信
      </div>

      <!-- 最新记录 -->
      <div v-if="latestSMS.length > 0" class="mb-8">
        <h2 class="text-lg font-bold text-gray-800 mb-4">🔥 最新记录</h2>
        <div class="bg-white shadow rounded-lg overflow-hidden border border-gray-200">
          <table class="min-w-full divide-y divide-gray-200 table-fixed">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-12">
                  <input type="checkbox" @change="toggleSelectAll('latest', $event)" :checked="isAllSelected('latest')" class="w-4 h-4">
                </th>
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
                <td class="px-6 py-4 whitespace-nowrap">
                  <input type="checkbox" :value="sms.id" v-model="selectedIds" class="w-4 h-4">
                </td>
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
        <!-- 最新记录 - 显示更多/显示全部 -->
        <div class="mt-4 text-center flex justify-center gap-2">
          <button v-if="!latestShowAll" @click="loadMoreLatest" :disabled="latestLoading" class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition disabled:opacity-50">
            {{ latestLoading ? '加载中...' : '显示更多' }}
          </button>
          <button v-if="!latestShowAll" @click="showAllLatest" :disabled="latestLoading" class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition disabled:opacity-50">
            显示全部
          </button>
          <span v-if="latestShowAll" class="text-gray-500 text-sm">已显示全部</span>
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
            <input type="checkbox" @change="toggleSelectAll(group.device, $event)" :checked="isAllSelected(group.device)" class="w-4 h-4">
          </div>
          
          <!-- 分组表格 -->
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 table-fixed">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider w-12"></th>
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
                  <td class="px-6 py-4 whitespace-nowrap">
                    <input type="checkbox" :value="sms.id" v-model="selectedIds" class="w-4 h-4">
                  </td>
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

          <!-- 显示更多/显示全部按钮 -->
          <div v-if="group.hasMore || !group.showAll" class="mt-4 text-center flex justify-center gap-2">
            <button 
              v-if="!group.showAll"
              @click="loadMore(group.device)" 
              :disabled="group.loading"
              class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition disabled:opacity-50"
            >
              {{ group.loading ? '加载中...' : '显示更多' }}
            </button>
            <button 
              v-if="!group.showAll"
              @click="showAllGroup(group.device)" 
              :disabled="group.loading"
              class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600 transition disabled:opacity-50"
            >
              显示全部
            </button>
          </div>
          <div v-else class="mt-4 text-center">
            <span class="text-gray-500 text-sm">已显示全部 ({{ group.total }} 条)</span>
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
  showAll?: boolean
}

const latestSMS = ref<SMS[]>([])
const groupedSMS = ref<GroupedSMS[]>([])
const selectedIds = ref<number[]>([])
const router = useRouter()

// 搜索相关
const searchKeyword = ref('')
const isSearching = ref(false)
const searchResults = ref<SMS[]>([])
const searchTotal = ref(0)
const searchPage = ref(1)
const searchPageSize = 20

// 最新记录相关
const latestOffset = ref(10)
const latestLoading = ref(false)
const latestShowAll = ref(false)

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
      offset: 10,
      showAll: false
    }))
    selectedIds.value = [] // 清空选中
    // 重置最新记录状态
    latestOffset.value = 10
    latestShowAll.value = false
    // 如果在搜索状态,退出搜索
    if (isSearching.value) {
      isSearching.value = false
      searchKeyword.value = ''
    }
  } catch (e) {
    console.error(e)
    if (axios.isAxiosError(e) && e.response?.status === 401) {
        logout()
    }
  }
}

// 最新记录 - 显示更多
const loadMoreLatest = async () => {
  if (latestLoading.value || latestShowAll.value) return
  
  latestLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await axios.get('/api/sms/list', {
      params: {
        page: Math.floor(latestOffset.value / 10) + 1,
        pageSize: 10
      },
      headers: { Authorization: `Bearer ${token}` }
    })
    
    latestSMS.value = [...latestSMS.value, ...res.data.data]
    latestOffset.value += 10
  } catch (e) {
    console.error(e)
    alert('加载失败')
  } finally {
    latestLoading.value = false
  }
}

// 最新记录 - 显示全部
const showAllLatest = async () => {
  if (latestShowAll.value) return
  
  latestLoading.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await axios.get('/api/sms/list', {
      params: {
        page: 1,
        pageSize: 1000 // 获取足够多的数据
      },
      headers: { Authorization: `Bearer ${token}` }
    })
    
    latestSMS.value = res.data.data
    latestShowAll.value = true
  } catch (e) {
    console.error(e)
    alert('加载失败')
  } finally {
    latestLoading.value = false
  }
}

const handleSearch = async () => {
  if (!searchKeyword.value.trim()) return
  
  isSearching.value = true
  searchPage.value = 1
  await executeSearch()
}

const executeSearch = async () => {
  try {
    const token = localStorage.getItem('token')
    const res = await axios.get('/api/sms/search', {
      params: {
        keyword: searchKeyword.value,
        page: searchPage.value,
        pageSize: searchPageSize
      },
      headers: { Authorization: `Bearer ${token}` }
    })
    searchResults.value = res.data.data || []
    searchTotal.value = res.data.pagination?.total || 0
    selectedIds.value = [] // 清空选中
  } catch (e) {
    console.error(e)
    alert('搜索失败')
  }
}

const clearSearch = () => {
  searchKeyword.value = ''
  isSearching.value = false
  searchResults.value = []
  searchTotal.value = 0
  searchPage.value = 1
}

// 监听搜索页码变化
import { watch } from 'vue'
watch(searchPage, () => {
  if (isSearching.value) {
    executeSearch()
  }
})

const loadMore = async (device: string) => {
  const group = groupedSMS.value.find(g => g.device === device)
  if (!group || group.loading || group.showAll) return

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

// 分组 - 显示全部
const showAllGroup = async (device: string) => {
  const group = groupedSMS.value.find(g => g.device === device)
  if (!group || group.showAll) return

  group.loading = true
  try {
    const token = localStorage.getItem('token')
    const res = await axios.get('/api/sms/load-all', {
      params: { device },
      headers: { Authorization: `Bearer ${token}` }
    })
    
    group.smsList = res.data.smsList
    group.total = res.data.total
    group.hasMore = false
    group.showAll = true
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

const toggleSelectAll = (type: string, event: Event) => {
  const target = event.target as HTMLInputElement
  const checked = target.checked
  
  if (type === 'latest') {
    if (checked) {
      latestSMS.value.forEach(sms => {
        if (!selectedIds.value.includes(sms.id)) {
          selectedIds.value.push(sms.id)
        }
      })
    } else {
      const latestIds = latestSMS.value.map(sms => sms.id)
      selectedIds.value = selectedIds.value.filter(id => !latestIds.includes(id))
    }
  } else if (type === 'search') {
    // 搜索结果全选
    if (checked) {
      searchResults.value.forEach(sms => {
        if (!selectedIds.value.includes(sms.id)) {
          selectedIds.value.push(sms.id)
        }
      })
    } else {
      const searchIds = searchResults.value.map(sms => sms.id)
      selectedIds.value = selectedIds.value.filter(id => !searchIds.includes(id))
    }
  } else {
    // 分组全选
    const group = groupedSMS.value.find(g => g.device === type)
    if (group) {
      if (checked) {
        group.smsList.forEach(sms => {
          if (!selectedIds.value.includes(sms.id)) {
            selectedIds.value.push(sms.id)
          }
        })
      } else {
        const groupIds = group.smsList.map(sms => sms.id)
        selectedIds.value = selectedIds.value.filter(id => !groupIds.includes(id))
      }
    }
  }
}

const isAllSelected = (type: string): boolean => {
  if (type === 'latest') {
    return latestSMS.value.length > 0 && latestSMS.value.every(sms => selectedIds.value.includes(sms.id))
  } else if (type === 'search') {
    return searchResults.value.length > 0 && searchResults.value.every(sms => selectedIds.value.includes(sms.id))
  } else {
    const group = groupedSMS.value.find(g => g.device === type)
    if (group && group.smsList.length > 0) {
      return group.smsList.every(sms => selectedIds.value.includes(sms.id))
    }
    return false
  }
}

const batchDelete = async () => {
  if (selectedIds.value.length === 0) return
  
  if (!confirm(`确定要删除选中的 ${selectedIds.value.length} 条短信吗？`)) return
  
  try {
    const token = localStorage.getItem('token')
    await axios.post('/api/sms/batch-delete', {
      ids: selectedIds.value
    }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    alert(`成功删除 ${selectedIds.value.length} 条短信`)
    selectedIds.value = []
    fetchGroupedSMS()
  } catch (e) {
    console.error(e)
    alert('批量删除失败')
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
