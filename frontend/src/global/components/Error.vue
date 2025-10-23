<template>
    <div v-if="hasError" class="error">
        <p class="message">{{ error.message }}</p>
        <pre>{{ JSON.stringify(error.details, null, 2) }}</pre>
    </div>

</template>

<script>
export default {
    name: "Error",
    props: {
        error: {
            type: Object,
            required: true,
            struct: {
                message: String,
                details: Object,
            }
        }
    },
    computed: {
        hasError() {
            return this.error != null && this.error?.message != null;
        }
    }
};
</script>

<style lang="scss" scoped>
.error {
    background-color: #a10606;
    padding: 1rem;
    border-radius: 4px;


    .message {
        font-size: 1.25rem;
        font-weight: 600;
        color: #fafafa;
        margin-bottom: 0.5rem;
    }

    pre {
        background-color: #2d2d2d;
        color: #f8f8f2;
        padding: 1rem;
        border-radius: 4px;
        overflow-x: auto;
    }
}
</style>