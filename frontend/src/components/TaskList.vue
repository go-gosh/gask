<template>
  <div>
    <el-table :data="data.data" stripe border table-layout="auto">
      <el-table-column prop="id" label="ID"/>
      <el-table-column prop="point" label="Point"/>
      <el-table-column prop="is_check" label="Checked">
        <template #default="scope">
          <el-switch v-model="scope.row.is_check" disabled/>
        </template>
      </el-table-column>
      <el-table-column prop="star" label="Star">
        <template #default="scope">
          <el-rate v-model="scope.row.star" :max="4" disabled/>
        </template>
      </el-table-column>
      <el-table-column prop="category" label="Category">
        <template #default="scope">
          <el-tag>
            {{ scope.row.category }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="Title"/>
      <el-table-column prop="detail" label="Detail"/>
      <el-table-column prop="start_at" label="Start Time" :formatter="dateFormatter"/>
      <el-table-column prop="deadline" label="Deadline" :formatter="dateFormatter"/>
      <el-table-column prop="created_at" label="Create Time" :formatter="dateFormatter"/>
      <el-table-column prop="updated_at" label="Update Time" :formatter="dateFormatter"/>
    </el-table>
  </div>
</template>

<script>
import axios from "axios";
import {dayjs} from "element-plus";

export default {
  name: "TaskList",
  created() {
    this.request(0, 1)
  },
  methods: {
    dateFormatter(row, col) {
      let date = row[col.property]
      if (date === undefined) return ""
      let unix = parseInt(date);
      if (unix === 0) return ""
      return dayjs.unix(unix).format('YYYY-MM-DD HH:mm:ss')
    },
    load(row, treeNode, resolve) {
      axios.get("http://localhost:8080/api/v1/task?parent_id=" + row.id + "&page=" + 1)
          .then(res => {
            resolve(res.data.data)
          }).catch(function (error) { // 请求失败处理
        console.log(error);
      })
    },
    request(parentId, page) {
      axios.get("http://localhost:8080/api/v1/task?parent_id=" + parentId + "&page=" + page)
          .then(res => {
            this.data = res.data
          }).catch(function (error) { // 请求失败处理
        console.log(error);
      })
    }
  },
  data() {
    return {
      data: {},
    }
  }
}
</script>

<style scoped>

</style>