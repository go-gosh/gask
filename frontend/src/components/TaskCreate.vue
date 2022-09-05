<template>
  <el-form :rules="rules" ref="Form" :model="form">
    <el-form-item label="Title" prop="title">
      <el-input v-model="form.title"/>
    </el-form-item>

    <el-form-item label="Root task">
      <el-select v-model="form.parent_id"
                 filterable
                 remote
                 reserve-keyword
                 placeholder="Please enter a keyword"
                 :remote-method="remoteMethod">
        <el-option
            v-for="item in rootTasks"
            :key="item.id"
            :label="item.title"
            :value="item.id"
        />
      </el-select>
    </el-form-item>

    <el-form-item>
      <el-col :span="8">
        <el-form-item label="Point">
          <el-input-number v-model="form.point" :min="1" :max="1000"/>
        </el-form-item>
      </el-col>
      <el-col :span="8">
        <el-form-item label="Check">
          <el-switch v-model="form.is_check"/>
        </el-form-item>
      </el-col>
      <el-col :span="8">
        <el-form-item label="Star">
          <el-rate v-model="form.star" :max="4" show-text :texts="star.texts" :colors="star.colors"/>
        </el-form-item>
      </el-col>
    </el-form-item>
    <el-form-item label="Category">
      <el-input v-model="form.category"/>
    </el-form-item>
    <el-form-item label="Task detail">
      <el-input v-model="form.detail"/>
    </el-form-item>
    <el-form-item label="Task start_at">
      <el-date-picker
          v-model="form.start_at"
          type="datetime"
          value-format="X"
          placeholder="Select date and time"
      />
    </el-form-item>
    <el-form-item label="Task dead_line">
      <el-date-picker
          v-model="form.dead_line"
          type="datetime"
          value-format="X"
          placeholder="Select date and time"
      />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="onSubmit">Create</el-button>
      <el-button>Cancel</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import {dayjs, ElMessage} from "element-plus";
import {taskCreate, taskDetail} from "@/request/api/task";

export default {
  name: "TaskCreate",
  data() {
    return {
      rules: {
        title: [{required: true, message: "Please input Task title", trigger: "blur"},],
      },
      star: {
        colors: {
          1: '#99A9BF',
          2: '#F7BA2A',
          3: '#FF9900',
          4: '#FF0000',
        },
        texts: ['Nothing', 'Urgent', 'Important', 'Important and urgent']
      },
      rootTasks: [],
      form: {
        parent_id: null,
        point: 100,
        is_check: null,
        star: 3,
        category: null,
        title: null,
        detail: null,
        start_at: dayjs().unix().toString(),
        dead_line: null,
      }
    }
  },
  methods: {
    remoteMethod(query) {
      if (query) {
        taskDetail(query).then(res => {
          this.rootTasks = [res.data]
        })
        return
      }
      this.rootTasks = []
    },
    onSubmit() {
      this.$refs.Form.validate((valid) => {
        if (valid) {
          taskCreate(this.form).then(res => {
            ElMessage.success(res.data)
          })
        } else {
          console.log('error submit!')
        }
      })
    }
  },
}
</script>

<style scoped>

</style>