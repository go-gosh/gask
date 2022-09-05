import axios from "axios";
import {ElMessage} from "element-plus";

export const api = axios.create({
    baseURL: "http://localhost:8080/api/v1"
})

api.interceptors.response.use(response => {
    if (response.status === 200) {
        return Promise.resolve(response)
    }
    return Promise.reject(response)
}, error => {
    const {response} = error
    if (response) {
        ElMessage.error({
            message: response.data,
        })
    }
})