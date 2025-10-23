<template>
    <ul class="root-listing" :data-bucket="bucket">
        <li v-for="(child, index) in item.children" :key="index">
            <ObjectListing :item="child" />
        </li>
    </ul>
</template>

<script>
import ObjectListing from '@/global/components/Object.vue';
export default {
    name: "DirectoryListing",
    props: {
        item: {
            type: Object,
            required: true
        }
    },
    computed: {
        bucket() {
            return this.$route.params.bucket;
        },
    },
    components: {
        ObjectListing
    }
};

</script>

<style lang="scss" scoped>
.root-listing {
    list-style-type: none;
    padding-left: 0;
}
</style>