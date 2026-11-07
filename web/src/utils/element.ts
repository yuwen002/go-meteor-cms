// Element Plus 服务式 API 的按需入口。
//
// 直接 `from '@/utils/element'` 会命中 es/index.mjs 这个 barrel，把 ElTransfer /
// ElCalendar 等根本没用到的组件一起打进产物（实测多出数百 KB）。模板里的
// <el-*> 组件由 unplugin-vue-components 自动按需解析，服务式调用则统一走这里。
import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/message-box/style/css'

export { ElMessage } from 'element-plus/es/components/message/index'
export { ElMessageBox } from 'element-plus/es/components/message-box/index'

// 纯类型导出不会产生运行时依赖，因此可以安全地从包入口取类型
export type { FormInstance, FormRules, UploadRequestOptions } from 'element-plus'
