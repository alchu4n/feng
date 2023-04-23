<script lang="ts" setup>
import { reactive } from 'vue'
import { NButton, NCard, NTag, NDescriptions, NDescriptionsItem, NLi, NUl, NList, NListItem } from 'naive-ui'
import { appvuln } from "../wailsjs/go/models";
import { AppScan } from '../wailsjs/go/main/App'


const data = reactive({
    vulns: [] as appvuln.Info[],
})

const scan = () => {
    AppScan().then(res => {
        data.vulns = res;
        console.log(res);
    })
}

</script>

<template>
    <main>
        <div style="text-align: center;">
            <n-button size="small" type="info" @click="scan">扫描</n-button>
        </div>
        <div v-if="data.vulns">
            <n-card v-for="item in data.vulns" :bordered="false">
                <n-descriptions label-placement="left" :column=1 size="small" bordered>
                    <n-descriptions-item label="App">
                        {{ item.path }}
                    </n-descriptions-item>
                    <n-descriptions-item label="执行文件">
                        {{ item.executable_path }}
                    </n-descriptions-item>
                    <n-descriptions-item label="是否能注入动态库">
                        <n-tag :bordered="false" type="warning">{{ item.injectable }}</n-tag>
                    </n-descriptions-item>
                    <n-descriptions-item label="可能存在注入的方式和路径">
                        <div v-for="dylib in item.dylibs">
                            Type:
                            <n-tag size="small" :bordered="false" round type="success">
                                {{ dylib.type }}
                            </n-tag>
                            <br>
                            Path: <n-tag size="small" :bordered="false" round type="info">
                                {{ dylib.path }}
                            </n-tag>
                        </div>
                    </n-descriptions-item>
                </n-descriptions>

            </n-card>
        </div>
    </main>
</template>

<style scoped></style>
