<template>
    <ul class="buckets-list">
        <li v-for="bucket in buckets" :key="bucket.name" class="bucket" :data-bucket="bucket.name">
            <RouterLink :to="`/${bucket.name}/`">{{ bucket.name }}</RouterLink>
        </li>
    </ul>
</template>


<script>
export default {
    props: {
        buckets: {
            type: Array,
            required: true
        }
    }
};
</script>

<style lang="scss" scoped>
.buckets-list {
    list-style-type: none;
    padding: 0;

    .bucket {
        margin: 0.5em 0;
        padding: 0.5em;
    }
}
</style>