<script setup>
import Material from './components/Material.vue'
import HelloWorld from './components/HelloWorld.vue'
// import { useDark, useToggle } from '@vueuse/core'

// const isDark = useDark()
// const toggleDark = useToggle(isDark)
</script>

<template>
  <!-- <button @click="toggleDark()">
    <i inline-block align-middle i="dark:carbon-moon carbon-sun" />

    <span class="ml-2">{{ isDark ? 'Dark' : 'Light' }}</span>
  </button> -->
  <!-- <img id="logo" alt="Wails logo" src="./assets/images/logo-universal.png" /> -->
  <!-- <HelloWorld /> -->
  <Material />
</template>

<style>
/* 确保根元素占满整个视口 */
html,
body,
#app {
  margin: 0;
  padding: 0;
  height: 100vh;
  width: 100%;
}

/* #logo {
  display: block;
  width: 50%;
  height: 50%;
  margin: auto;
  padding: 10% 0 0;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
} */
</style>
