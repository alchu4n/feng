<script lang="ts" setup>
import { reactive } from 'vue'
import { NButton, NCard, NTag, NDescriptions, NDescriptionsItem, NLi, NUl, NList, NListItem } from 'naive-ui'
import { macho } from "../wailsjs/go/models";
import { SelectMacho, ParseMacho } from '../wailsjs/go/main/App'

const data = reactive({
  machoPath: "",
  machos: [] as macho.MachOInfo[],
})

const select = () => {
  SelectMacho().then(res => {
    data.machoPath = res;
    data.machoPath ? parse() : data.machos = [];
  })
}

const parse = () => {
  ParseMacho(data.machoPath).then(res => {
    if (res !== null) {
      data.machos = res;
    }
  })
}
</script>

<template>
  <main>
    <div style="text-align: center;">
      <span style="margin-right: 10px;">Path:</span>
      <n-button strong secondary size="small" type="info" @click="select">{{ data.machoPath ? data.machoPath :
        "选择MachO" }}</n-button>
    </div>
    <div v-if="data.machos">
      <n-card v-for=" item in data.machos" :bordered="false">
        <n-descriptions label-placement="left" :title="item.cpu" :column=1 size="small" bordered>
          <n-descriptions-item label="Macho信息">
            <n-tag :bordered="false" type="default" size="small">
              {{ item.cpu }}
            </n-tag>|
            <n-tag :bordered="false" type="default" size="small">
              {{ item.magic }}
            </n-tag>|
            <n-tag :bordered="false" type="default" size="small">
              {{ item.type }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="dylinker">
            <n-tag size="small" :bordered="false" type="warning">{{ item.dylibs.dylinker }}</n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="rpaths">
            <n-ul>
              <n-li v-for="paths in item.dylibs.rpaths">
                path:
                <n-tag size="small" :bordered="false" round type="success">
                  {{ paths }}
                </n-tag>
              </n-li>
            </n-ul>
          </n-descriptions-item>
          <n-descriptions-item label="loads">
            <n-ul>
              <n-li v-for="dylib in item.dylibs.loads">
                path:
                <n-tag size="small" :bordered="false" round type="success">
                  {{ dylib.name }}
                </n-tag>
                version: <n-tag size="small" :bordered="false" round type="info">
                  {{ dylib.current_version }}
                </n-tag>
              </n-li>
            </n-ul>
          </n-descriptions-item>
          <n-descriptions-item label="weaks">
            <n-ul>
              <n-li v-for="dylib in item.dylibs.weaks">
                path:
                <n-tag size="small" :bordered="false" round type="success">
                  {{ dylib.name }}
                </n-tag>
                version: <n-tag size="small" :bordered="false" round type="info">
                  {{ dylib.current_version }}
                </n-tag>
              </n-li>
            </n-ul>
          </n-descriptions-item>
          <n-descriptions-item label="codesign">
            <n-list>
              <n-list-item>
                id: <n-tag :bordered="false" type="error" size="small">{{ item.codesign.id }}</n-tag>
              </n-list-item>
              <n-list-item>
                team-id: <n-tag :bordered="false" type="info" size="small">{{ item.codesign.team_id }}</n-tag>
              </n-list-item>
              <n-list-item>
                flags: <n-tag :bordered="false" type="info" size="small">0x{{ item.codesign.flags.toString(16)
                }} </n-tag>
                - <n-tag :bordered="false" type="error" size="small">{{ item.codesign.flags_string }}</n-tag>
              </n-list-item>
              <n-list-item>
                runtime-version: <n-tag :bordered="false" type="info" size="small">{{ item.codesign.runtime_version
                }}</n-tag>
              </n-list-item>
            </n-list>
          </n-descriptions-item>
          <n-descriptions-item label="entitlements">
            <n-ul>
              <n-li v-for="k, v in item.codesign.entitlements">
                <n-tag size="small" :bordered="false" round type="error">
                  {{ v }}
                </n-tag>:
                <n-tag size="small" :bordered="false" round type="default">
                  {{ k }}
                </n-tag>
              </n-li>
            </n-ul>
          </n-descriptions-item>
        </n-descriptions>

      </n-card>
    </div>
  </main>
</template>

<style scoped></style>
