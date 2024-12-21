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
            :label="material.code + ' - ' + material.remark"
          />
        </el-select>
      </el-col>
      <el-col :span="4">
        <el-input v-model="selectedComponentsData.code" placeholder="物料码" />
      </el-col>
      <el-col :span="4">
        <el-button type="primary" @click="doSelectComponents">筛选</el-button>
      </el-col>
    </el-row>
    <!-- 分页, 每页10条 -->
    <el-table :data="currentPageData" style="width: 100%">
      <el-table-column prop="outerDiameter" label="外径" />
      <el-table-column prop="wallThickness" label="壁厚" />
      <el-table-column prop="material" label="材质" />
      <el-table-column prop="code" label="物料码" />
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
  <!-- 添加物料 -->
  <el-dialog v-model="addComponentDialog" title="添加物料">
    <el-input v-model="newComponent.outerDiameter" placeholder="外径" />
    <el-input v-model="newComponent.wallThickness" placeholder="壁厚" />
    <el-select v-model="newComponent.material" placeholder="材质">
      <el-option
        v-for="material in materialList"
        :key="material.code"
        :value="material.code"
        :label="material.code + ' - ' + material.remark"
      />
    </el-select>
    <el-input v-model="newComponent.code" placeholder="物料码" />
    <el-button type="primary" @click="doAddComponent">添加物料</el-button>
  </el-dialog>
  <!-- 添加材质 -->
  <el-dialog v-model="addMaterialTypeDialog" title="添加材质">
    <el-input v-model="newMaterialType.code" placeholder="材质码" />
    <el-input v-model="newMaterialType.remark" placeholder="材质描述" />
    <el-button type="primary" @click="doAddMaterialType">添加材质</el-button>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { GoSearchComponents, GoAddComponent, GoGetAllMaterial, GoAddMaterial } from '../../wailsjs/go/main/App'

const componentsList = ref([])
const materialList = ref([])

const selectedComponentsData = reactive({
  outerDiameter: '',
  wallThickness: '',
  material: '',
  code: ''
})

const addComponentDialog = ref(false)
const addMaterialTypeDialog = ref(false)
const newComponent = reactive({
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
  ).then((selectedComponents, err) => {
    if (err != null) {
      ElMessage.error('查询失败: ' + err)
    } else {
      componentsList.value = []
      selectedComponents.forEach((component) => {
        console.log(component)
        componentsList.value.push(component)
      })
      ElMessage.success('查询成功')
    }
  })
}

const doAddComponent = () => {
  console.log(newComponent)
  // 调用后端接口
  // 后端接口返回数据
  // 更新materials数据
  GoAddComponent(newComponent.outerDiameter, newComponent.wallThickness, newComponent.material, newComponent.code).then(
    (err) => {
      if (err) {
        ElMessage.error('添加物料失败: ' + err)
      } else {
        ElMessage.success('添加物料成功')
      }
    }
  )
}

const doAddMaterialType = () => {
  console.log(newMaterialType)
  GoAddMaterial(newMaterialType.code, newMaterialType.remark).then((err) => {
    if (err) {
      ElMessage.error('添加材质失败: ' + err)
    } else {
      ElMessage.success('添加材质成功')
      doGetAllMaterial()
    }
  })
}

const doGetAllMaterial = () => {
  GoGetAllMaterial().then((materials, err) => {
    if (err) {
      ElMessage.error('获取材质失败: ' + err)
    } else {
      materialList.value = []
      materials.forEach((material) => {
        console.log(material)
        materialList.value.push(material)
      })
      ElMessage.success('获取材质成功')
    }
  })
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
</style>
