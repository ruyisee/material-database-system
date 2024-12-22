<!-- 物料管理页面 -->
<template>
  <div class="material-container">
    <!-- 筛选条件 -->
    <el-row>
      <!-- // Material 定义物料模型
        // φ22*5	304	81383220
        // φ45*8	304	81349147
        // φ127*28	16Mn	81386223 -->
      <el-col :span="4">
        <el-input v-model="selectedComponentsData.outerDiameter" placeholder="外径" />
      </el-col>
      <el-col :span="4">
        <el-input v-model="selectedComponentsData.wallThickness" placeholder="壁厚" />
      </el-col>
      <el-col :span="4">
        <el-select v-model="selectedComponentsData.material" placeholder="材质">
          <el-option
            v-for="material in materialList"
            :key="material.code"
            :value="material.code"
            :label="material.code"
          />
        </el-select>
      </el-col>
      <el-col :span="4">
        <el-input v-model="selectedComponentsData.code" placeholder="物料码" />
      </el-col>
      <el-col :span="4">
        <el-button type="primary" @click="doSelectComponents">筛选</el-button>
        <el-button type="primary" @click="doSelectReset">重置</el-button>
      </el-col>
    </el-row>
    <!-- 分页, 每页10条 -->
    <el-table :data="currentPageData" style="width: 100%">
      <el-table-column prop="outerDiameter" label="外径" />
      <el-table-column prop="wallThickness" label="壁厚" />
      <el-table-column prop="material" label="材质" />
      <el-table-column prop="code" label="物料码" />
      <!-- 操作 -->
      <el-table-column label="操作">
        <template #default="scope">
          <el-button type="primary" @click="tableEditComponent(scope.row)">编辑</el-button>
          <el-button type="danger" @click="doDeleteComponent(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页器 -->
    <el-pagination
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="componentsList.length"
      :page-sizes="[10, 20, 30, 50]"
      layout="total, sizes, prev, pager, next"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>

  <el-button type="primary" @click="addComponentDialog = true">添加物料</el-button>
  <el-button type="primary" @click="addMaterialTypeDialog = true">添加材质</el-button>
  <el-button type="primary" @click="importMaterialDialog = true">导入物料</el-button>
  <!-- 添加物料 -->
  <el-dialog v-model="addComponentDialog" title="添加物料">
    <el-input v-model="newComponent.outerDiameter" placeholder="外径" />
    <el-input v-model="newComponent.wallThickness" placeholder="壁厚" />
    <el-select v-model="newComponent.material" placeholder="材质">
      <el-option v-for="material in materialList" :key="material.code" :value="material.code" :label="material.code" />
    </el-select>
    <el-input v-model="newComponent.code" placeholder="物料码" />
    <el-button type="primary" @click="doAddComponent">添加物料</el-button>
  </el-dialog>

  <!-- 编辑物料 -->
  <el-dialog v-model="editComponentDialog" title="编辑物料">
    <el-input v-model="editComponent.outerDiameter" placeholder="外径" />
    <el-input v-model="editComponent.wallThickness" placeholder="壁厚" />
    <el-select v-model="editComponent.material" placeholder="材质">
      <el-option v-for="material in materialList" :key="material.code" :value="material.code" :label="material.code" />
    </el-select>
    <el-input v-model="editComponent.code" placeholder="物料码" />
    <el-button type="primary" @click="doEditComponent">更新</el-button>
  </el-dialog>
  <!-- 添加材质 -->
  <el-dialog v-model="addMaterialTypeDialog" title="添加材质">
    <el-input v-model="newMaterialType.code" placeholder="材质码" />
    <el-input v-model="newMaterialType.remark" placeholder="材质描述" />
    <el-button type="primary" @click="doAddMaterialType">添加材质</el-button>
  </el-dialog>
  <!-- 导入物料 -->
  <el-dialog v-model="importMaterialDialog" title="导入物料">
    <b>导入物料, csv内部格式: 外径,壁厚,材质,物料码. 必须有表头，名字必须对应</b>
    <!-- 选择文件 -->
    <el-upload
      class="upload-demo"
      action="#"
      :auto-upload="false"
      :show-file-list="true"
      :on-change="handleFileChange"
      accept=".csv"
    >
      <template #trigger>
        <el-button type="primary">选择文件</el-button>
      </template>
    </el-upload>
    <el-button type="success" @click="doImportMaterial">导入物料</el-button>
    <div style="white-space: pre-line">{{ importMaterialResult }}</div>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  GoSearchComponents,
  GoAddComponent,
  GoGetAllMaterial,
  GoImportComponents,
  GoUploadFile,
  GoDeleteComponent,
  GoEditComponent
} from '../../wailsjs/go/main/App'

const componentsList = ref([])
const materialList = ref([])

const selectedComponentsData = reactive({
  outerDiameter: '',
  wallThickness: '',
  material: '',
  code: ''
})

const importMaterialDialog = ref(false)
const importMaterialResult = ref('')
const addComponentDialog = ref(false)
const addMaterialTypeDialog = ref(false)
const editComponentDialog = ref(false)
const newComponent = reactive({
  outerDiameter: '',
  wallThickness: '',
  material: '',
  code: ''
})

const editComponent = reactive({
  id: 0,
  outerDiameter: '',
  wallThickness: '',
  material: '',
  code: ''
})

const newMaterialType = reactive({
  code: '',
  remark: ''
})

const doSelectComponents = () => {
  GoSearchComponents(
    selectedComponentsData.outerDiameter,
    selectedComponentsData.wallThickness,
    selectedComponentsData.material,
    selectedComponentsData.code
  )
    .then((selectedComponents) => {
      componentsList.value = []
      selectedComponents.forEach((component) => {
        console.log(component)
        componentsList.value.push(component)
      })
      ElMessage.success('查询成功')
    })
    .catch((err) => {
      console.error('查询失败:', err)
      ElMessage.error('查询失败: ' + err)
    })
}

const doSelectReset = () => {
  selectedComponentsData.outerDiameter = ''
  selectedComponentsData.wallThickness = ''
  selectedComponentsData.material = ''
  selectedComponentsData.code = ''
}

const doDeleteComponent = (component) => {
  console.log('删除物料: ', component)
  //   确认
  ElMessageBox.confirm('确定删除物料: < ' + component.code + ' > 吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(() => {
      GoDeleteComponent(component.code)
        .then(() => {
          ElMessage.success('删除物料成功')
          doSelectComponents()
        })
        .catch((err) => {
          console.error('删除物料失败:', err)
          ElMessage.error('删除物料失败: ' + err)
        })
    })
    .catch(() => {
      ElMessage.info('已取消删除')
    })
}

const tableEditComponent = (component) => {
  editComponent.id = component.id
  editComponent.outerDiameter = component.outerDiameter
  editComponent.wallThickness = component.wallThickness
  editComponent.material = component.material
  editComponent.code = component.code
  editComponentDialog.value = true
}

const doEditComponent = () => {
  console.log('编辑物料: ', editComponent)
  GoEditComponent(
    editComponent.id,
    editComponent.code,
    editComponent.outerDiameter,
    editComponent.wallThickness,
    editComponent.material
  )
    .then(() => {
      ElMessage.success('编辑物料成功')
      doSelectComponents()
    })
    .catch((err) => {
      console.error('编辑物料失败:', err)
      ElMessage.error('编辑物料失败: ' + err)
    })
}

const doAddComponent = () => {
  console.log(newComponent)
  // 调用后端接口
  // 后端接口返回数据
  // 更新materials数据
  GoAddComponent(newComponent.outerDiameter, newComponent.wallThickness, newComponent.material, newComponent.code)
    .then(() => {
      ElMessage.success('添加物料成功')
    })
    .catch((err) => {
      console.error('添加物料失败:', err)
      ElMessage.error('添加物料失败: ' + err)
    })
}

const doAddMaterialType = () => {
  console.log(newMaterialType)
  console.log(newMaterialType.code)
  console.log(newMaterialType.remark)
  GoAddMaterial(newMaterialType.code, newMaterialType.remark)
    .then((result) => {
      // Wails在调用Go函数时，如果返回error，会将error转换为rejected promise
      ElMessage.success('添加材质成功')
      doGetAllMaterial()
    })
    .catch((err) => {
      console.error('添加材质失败:', err)
      ElMessage.error('添加材质失败: ' + err)
    })
}

const doGetAllMaterial = () => {
  GoGetAllMaterial()
    .then((materials) => {
      materialList.value = []
      materials.forEach((material) => {
        console.log(material)
        materialList.value.push(material)
      })
      ElMessage.success('获取材质成功')
    })
    .catch((err) => {
      console.error('获取材质失败:', err)
      ElMessage.error('获取材质失败: ' + err)
    })
}

const selectedFile = ref(null)

// 处理文件选择
const handleFileChange = (uploadFile) => {
  selectedFile.value = uploadFile.raw
}

// 导入物料
const doImportMaterial = async () => {
  if (!selectedFile.value) {
    ElMessage.warning('请先选择文件')
    return
  }

  importMaterialResult.value = ''
  try {
    // 读取文件内容并转换为 Base64
    const reader = new FileReader()
    reader.readAsDataURL(selectedFile.value)

    reader.onload = async () => {
      try {
        // 获取 Base64 字符串（移除 data:application/csv;base64, 前缀）
        const base64Content = reader.result.split(',')[1]
        // 上传文件并获取临时文件路径
        const tempFilePath = await GoUploadFile(base64Content)
        // 导入数据
        const result = await GoImportComponents(tempFilePath)
        ElMessage.success('导入成功')
        importMaterialResult.value = result
        // 刷新物料列表
        doGetAllMaterial()
      } catch (err) {
        console.error('导入物料失败:', err)
        ElMessage.error('导入物料失败: ' + err)
      }
    }
  } catch (err) {
    console.error('读取文件失败:', err)
    ElMessage.error('读取文件失败: ' + err)
  }
}

// 分页相关数据
const currentPage = ref(1)
const pageSize = ref(10)

// 计算当前页的数据
const currentPageData = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return componentsList.value.slice(start, end)
})

// 处理每页条数改变
const handleSizeChange = (val) => {
  pageSize.value = val
  currentPage.value = 1 // 重置到第一页
}

// 处理页码改变
const handleCurrentChange = (val) => {
  currentPage.value = val
}

onMounted(() => {
  doGetAllMaterial()
})
</script>

<style scoped>
.material-container {
  padding: 20px;
}

.el-pagination {
  margin-top: 20px;
  justify-content: center;
}

.upload-demo {
  margin-bottom: 20px;
}
</style>
