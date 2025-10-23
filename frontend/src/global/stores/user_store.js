import { defineStore, acceptHMRUpdate } from 'pinia';
import Client from "@/global/networking/client";

export const userStore = defineStore('user', {
    state: () => ({
        user: null,
        error: null
    }),
    getters: {
        loaded: (state) => state.user != null,
        isLoggedIn: (state) => state.user != null && state.user.roles.length > 0,
    },
    actions: {
        async fetchUser() {
            const client = Client.getClient();
            try {
                const res = await client.get('/user');
                this.user = res.data;
                this.error = null;
            } catch (error) {
                console.error("Error fetching user:", error);
                this.error = {
                    message: "Failed to fetch user.",
                    details: error,
                };
            }
        },
    },
});


if (
    import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(userStore,
        import.meta.hot));
}
if (
    import.meta.webpackHot) {
    import.meta.webpackHot.accept(acceptHMRUpdate(userStore,
        import.meta.webpackHot));
}