import { defineStore, acceptHMRUpdate } from 'pinia';
import BucketsRepository from './repositories/buckets.js';

export const bucketsStore = defineStore('buckets', {
    state: () => ({
        buckets: null,
        error: null,
    }),
    getters: {
        loaded: (state) => state.buckets != null,
        hasError: (state) => state.error != null,
    },
    actions: {
        async fetchBuckets() {
            const repo = new BucketsRepository();
            try {
                const res = await repo.fetchBuckets();
                this.buckets = res.data;
                console.log('Fetched buckets:', this.buckets);
            } catch (error) {
                console.error('Error fetching buckets:', error);
                this.error = {
                    message: 'Failed to fetch buckets.',
                    details: error,
                };
            }

        },
    },
});

if (
    import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(bucketsStore,
        import.meta.hot));
}
if (
    import.meta.webpackHot) {
    import.meta.webpackHot.accept(acceptHMRUpdate(bucketsStore,
        import.meta.webpackHot));
}