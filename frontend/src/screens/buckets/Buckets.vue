<template>
    <article class="buckets">
        <h1>Welcome to the Object Storage Browser</h1>
        <p>
            Use the navigation above to browse directories and manage your object storage.
        </p>
        <div v-if="loaded">
            <BucketsList :buckets="buckets" />
        </div>
        <div v-else>
            loading...
        </div>
        <Error :error="error" />
    </article>
</template>


<script>
import { mapStores, mapState } from "pinia";
import { bucketsStore } from "@/global/stores/buckets_store.js";
import BucketsList from "./components/BucketsList.vue";
import Error from "@/global/components/Error.vue";


export default {
    components: {
        BucketsList,
        Error
    },
    computed: {
        ...mapStores(bucketsStore),
        ...mapState(bucketsStore, ['buckets', 'loaded', 'error']),
    },
    async created() {
        if (!this.loaded) {
            await this.bucketsStore.fetchBuckets();
        }
    }
};
</script>
