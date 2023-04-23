<script lang="ts" setup>
import { NTabs, NTabPane } from 'naive-ui';
import Macho from './Macho.vue';
import AppVuln from './AppVuln.vue';

</script>

<template>
    <main>
        <n-tabs type="segment">
            <n-tab-pane name="app-vuln" tab="App扫描">
                <AppVuln />
            </n-tab-pane>
            <n-tab-pane name="macho" tab="MachO分析">
                <Macho />
            </n-tab-pane>
        </n-tabs>
    </main>
</template>

<style scoped></style>
