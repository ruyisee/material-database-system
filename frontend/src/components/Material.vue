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
        <el-input v-model="selectedMaterialsData.outerDiameter" placeholder="外径" />
      </el-col>
      <el-col :span="4">
        <el-input v-model="selectedMaterialsData.wallThickness" placeholder="壁厚" />
      </el-col>
      <el-col :span="4">
        <el-input v-model="selectedMaterialsData.material" placeholder="材质" />
      </el-col>
      <el-col :span="4">
        <el-input v-model="selectedMaterialsData.code" placeholder="物料码" />
      </el-col>
      <el-col :span="4">
        <el-button type="primary" @click="doSelectMaterials">筛选</el-button>
      </el-col>
    </el-row>
    <el-table :data="materials" style="width: 100%">
      <el-table-column prop="outerDiameter" label="外径" />
      <el-table-column prop="wallThickness" label="壁厚" />
      <el-table-column prop="material" label="材质" />
      <el-table-column prop="code" label="物料码" />
    </el-table>
  </div>

  <el-button type="primary" @click="addMaterialDialog = true">添加物料</el-button>
  <!-- 添加物料 -->
  <el-dialog v-model="addMaterialDialog" title="添加物料">
    <el-input v-model="newMaterial.outerDiameter" placeholder="外径" />
    <el-input v-model="newMaterial.wallThickness" placeholder="壁厚" />
    <el-input v-model="newMaterial.material" placeholder="材质" />
    <el-input v-model="newMaterial.code" placeholder="物料码" />
    <el-button type="primary" @click="doAddMaterial">添加</el-button>
  </el-dialog>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { GoSearchMaterials, GoAddMaterial } from '../../wailsjs/go/main/App'

const materials = ref([])

const selectedMaterialsData = reactive({
  outerDiameter: '',
  wallThickness: '',
  material: '',
  code: ''
})

const addMaterialDialog = ref(false)
const newMaterial = reactive({
  outerDiameter: '',
  wallThickness: '',
  material: '',
  code: ''
})

const doSelectMaterials = () => {
  GoSearchMaterials(
    selectedMaterialsData.outerDiameter,
    selectedMaterialsData.wallThickness,
    selectedMaterialsData.material
  ).then((selectedMaterials, err) => {
    materials.value = []
    if (err != null) {
      ElMessage.error('查询失败' + err)
    } else {
      selectedMaterials.forEach((material) => {
        console.log(material)
        materials.value.push(material)
      })
      ElMessage.success('查询成功')
    }
  })
}

const doAddMaterial = () => {
  console.log(newMaterial)
  // 调用后端接口
  // 后端接口返回数据
  // 更新materials数据
  GoAddMaterial(newMaterial.outerDiameter, newMaterial.wallThickness, newMaterial.material, newMaterial.code).then(
    (err) => {
      if (err) {
        ElMessage.error('添加失败' + err)
      } else {
        ElMessage.success('添加成功')
      }
    }
  )
}
</script>

<style scoped>
.material-container {
  padding: 20px;
}
</style>
