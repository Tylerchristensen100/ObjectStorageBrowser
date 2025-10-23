<template>
    <h1>Login</h1>
</template>


<script>
import { mapStores, mapState } from "pinia";
import { userStore } from "@/global/stores/user_store.js";
export default {
    computed: {
        ...mapStores(userStore),
        ...mapState(userStore, ['user', 'loaded', 'isLoggedIn', 'error']),
    },
    created() {
        if (!this.loaded) {
            this.userStore.fetchUser();
        }
        if (!this.isLoggedIn) {
            window.location.href = '/login';
        }
    }
};
</script>