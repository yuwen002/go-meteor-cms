<template>
  <el-card>
    <el-form :inline="true" :model="query" class="search-bar">
      <el-form-item label="用户名">
        <el-input v-model="query.username" placeholder="操作人" clearable />
      </el-form-item>
      <el-form-item label="模块">
        <el-input v-model="query.module" placeholder="业务模块" clearable />
      </el-form-item>
      <el-form-item label="结果">
        <el-select v-model="query.status" placeholder="全部" clearable style="width: 120px">
          <el-option label="成功" :value="0" />
          <el-option label="失败" :value="1" />
        </el-select>
      </el-form-item>
      <el-form-item label="时间">
        <el-date-picker
          v-model="timeRange"
          type="daterange"
          value-format="YYYY-MM-DD"
          range-separator="至"
          start-placeholder="开始"
          end-placeholder="结束"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="loadList">查询</el-button>
        <el-button @click="resetQuery">重置</el-button>
      </el-form-item>
    </el-form>

    <el-table :data="list" v-loading="loading" border>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="username" label="操作人" width="120" />
      <el-table-column prop="module" label="模块" width="120" />
      <el-table-column prop="action" label="操作" width="140" />
      <el-table-column prop="method" label="方法" width="80" />
      <el-table-column prop="path" label="路径" min-width="200" />
      <el-table-column prop="ip" label="IP" width="130" />
      <el-table-column label="结果" width="80">
        <template #default="{ row }">
          <el-tag :type="row.status === 0 ? 'success' : 'danger'">{{ row.status === 0 ? '成功' : '失败' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="duration_ms" label="耗时(ms)" width="100" />
      <el-table-column prop="created_at" label="时间" width="180" />
    </el-table>

    <el-pagination
      class="pager"
      background
      layout="total, prev, pager, next"
      :total="total"
      :current-page="query.page"
      :page-size="query.page_size"
      @current-change="(p: number) => { query.page = p; loadList() }"
    />
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { getOperationLogs } from '@/api/log'

const loading = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const timeRange = ref<[string, string]>(['', ''])
const query = reactive({
  page: 1,
  page_size: 10,
  username: '',
  module: '',
  status: undefined as number | undefined,
  start_time: '',
  end_time: '',
})

async function loadList() {
  loading.value = true
  try {
    query.start_time = timeRange.value?.[0] || ''
    query.end_time = timeRange.value?.[1] || ''
    const res = await getOperationLogs(query)
    list.value = res.list || []
    total.value = res.total || 0
  } finally {
    loading.value = false
  }
}

function resetQuery() {
  query.username = ''
  query.module = ''
  query.status = undefined
  timeRange.value = ['', '']
  query.page = 1
  loadList()
}

onMounted(loadList)
</script>

<style scoped>
.search-bar {
  margin-bottom: 12px;
}
.pager {
  margin-top: 12px;
  justify-content: flex-end;
}
</style>
