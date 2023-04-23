<script lang="ts" setup>
import { dateZhCN, GlobalThemeOverrides, zhCN } from "naive-ui";
import Home from "./Home.vue";

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColorPressed: "#101010",
    primaryColor: "#464646",
    primaryColorHover: "#7A7B7AFF",
    primaryColorSuppl: "#7A7B7AFF",
  }
}
</script>

<template>
  <n-config-provider :theme-overrides="themeOverrides" :locale="zhCN" :date-locale="dateZhCN">
    <n-message-provider>
      <Home />
    </n-message-provider>
  </n-config-provider>
</template>

<style></style>
