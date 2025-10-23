<template>
    <article class="directory">
        <div>
            <p>Currently viewing bucket: <strong>{{ bucket }}</strong></p>
            <p>Current path: <pre>{{ path }}</pre></p>
        </div>

        <div v-if="loaded && directory">
            <div v-if="directory.length === 0">No files found.</div>
            <ul v-else class="directory-list">
                <li v-for="item in directory" :key="item.name">
                    <DirectoryListing :item="item" />
                </li>
            </ul>
        </div>
        <div v-else-if="error">
            <p>Error loading directory: {{ error }}</p>
        </div>
        <div v-else>
            <p>Loading directory...</p>
        </div>
        <Error :error="error" />
    </article>
</template>

<script>
import { mapStores, mapState } from "pinia";
import { directoryStore } from "@/global/stores/directory_store.js";
import DirectoryListing from "./components/DirectoryListing.vue";
import Error from "@/global/components/Error.vue";

export default {
    components: {
        DirectoryListing,
        Error
    },
    computed: {
        ...mapStores(directoryStore),
        ...mapState(directoryStore, ['directory', 'loaded', 'error']),
        bucket() {
            return this.$route.params.bucket;
        },
        path() {
            return this.$route.params.path || '';
        }
    },
    created() {
        if (this.bucket === undefined) {
            console.error('Bucket Undefined, redirecting to /buckets');
            this.$router.push('/buckets');
        }
        console.log('Directory.vue created with bucket:', this.bucket, 'and path:', this.path);
        if (!this.loaded) {
            this.directoryStore.fetchDirectory(this.bucket, this.path);
        }
    }
};
</script>

<style lang="scss" scoped>
.directory {
    text-align: left;
    ul.directory-list {
        list-style-type: none;
        padding: 0;
    }
}

</style>