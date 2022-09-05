import {api} from "@/request/api/api";


export const taskPaginate = (params) => {
    return api.get("/task", {params: params})
}

export const taskCreate = (data) => {
    return api.post("/task", data)
}

export const taskUpdate = (id, data) => {
    return api.put("/task" + id, data)
}

export const taskDetail = (id) => {
    return api.get("/task/" + id)
}

export const taskDelete = (id) => {
    return api.delete("/task/" + id)
}
