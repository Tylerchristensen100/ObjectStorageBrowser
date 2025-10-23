<template>
  <Nav />
  <main>
    <router-view v-slot="{ Component, route }">
      <transition :key="route.fullPath" name="fade" mode="out-in">
        <component :is="Component" :key="`Component for ${route.fullPath}`" />
      </transition>
    </router-view>
  </main>
  <Footer />
</template>


<script setup>
import Nav from '@/global/components/Nav.vue';
import Footer from '@/global/components/Footer.vue';
</script>


<style>
main {
  padding: 1rem;
  max-width: 1280px;
}
</style>