<template>
  <div>
    <div style="text-align: right">
      <el-button @click="showCreate = true" type="primary">Create</el-button>
    </div>

    <el-dialog v-model="showCreate" title="Create New Task">
      <TaskCreate/>
    </el-dialog>

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
    <el-pagination v-model:current-page="data.page" v-model:page-size="data.page_size"
                   :page-sizes="[10, 30, 50, 100]" :total="data.total" background
                   layout="total, sizes, prev, pager, next, jumper"
                   small @size-change="request(0)"
                   @current-change="request(0)"
    />
  </div>
</template>

<script>
import {dayjs} from "element-plus";
import TaskCreate from "@/components/TaskCreate";
import {taskPaginate} from "@/request/api/task";

export default {
  name: "TaskList",
  components: {TaskCreate},
  created() {
    this.request(0)
  },
  methods: {
    dateFormatter(row, col) {
      let date = row[col.property]
      if (date === undefined) return ""
      let unix = parseInt(date);
      if (unix === 0) return ""
      return dayjs.unix(unix).format('YYYY-MM-DD HH:mm:ss')
    },
    request(parentId) {
      taskPaginate({
        page: this.data.page,
        page_size: this.data.page_size,
        parent_id: parentId,
      }).then(res => {
        this.data = res.data
      })
    }
  },
  data() {
    return {
      showCreate: false,
      data: {
        page: 1,
        page_size: 10,
        total: 0,
        data: [],
      },
    }
  }
}
</script>

<style scoped>
.el-pagination {
  justify-content: center;
}
</style>