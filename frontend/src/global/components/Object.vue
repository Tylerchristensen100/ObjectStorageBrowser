<template>
    <div class="object" :data-depth="item?.depth" :data-isDir="isDir" :data-isFile="file"
        :style="`--item-depth: ${item?.depth}`">

        <div class="header">
            <i class="icon" />

            <strong @click="dialogOpen = true">{{ item?.path }}</strong>
            <span v-if="isDir" class="size">({{ item?.stats?.size ?? "?" }} bytes)</span>
            <button v-if="file" @click="download">Download</button>
        </div>

    </div>
    <!-- <dialog :open="dialogOpen">
        <pre>{{ JSON.stringify(object, null, 2) }}</pre>
    </dialog> -->

    <div v-if="isDir" :data-parent="object?.path">
        <ul :data-count="item?.children?.length ?? 0" class="children">
            <li v-for="(child, index) in item?.children" :key="index">
                <ObjectListing :item="child" />
            </li>
        </ul>
    </div>
</template>

<script>
import ObjectHelpers from '../helpers/object_helpers';

export default {
    name: "ObjectListing",
    props: {
        item: {
            type: Object,
            required: true
        },
    },
    data() {
        return {
            dialogOpen: false,
        };
    },
    computed: {
        bucket() {
            return this.$route.params.bucket;
        },
        isDir() {
            return this.item?.dir === true;
        },
        file() {
            return this.item?.file === true;
        },
    },
    methods: {
        download() {
            ObjectHelpers.download(this.bucket, this.item?.full_path);
        }
    },
    mounted() {
        console.log(this.item?.depth, this.item);
    },
    components: {
        ObjectListing: () => import('@/global/components/Object.vue'),
    },
};
</script>


<style lang="scss">
.object {
    --height: 16px;
    --icon-width: 16px;


    display: flex;
    flex-direction: row;
    align-items: center;
    margin-block: 10px;
    // max-width: 360px;




    &[data-isDir="true"] .header .icon {
        content: 'Dir Icon';
    }

    &[data-isDir="false"] .header .icon {
        content: 'File Icon';
    }

    .header {
        display: grid;

        grid-template-columns: var(--icon-width) auto 80px;
        align-items: center;
        justify-items: start;
        gap: 4px;

        .icon {
            display: inline-block;
            content: 'ICON';
            font-style: normal;
            margin-inline-end: 4px;
            width: var(--icon-width);
            height: var(--height);
            display: inline-block;
            text-align: center;
            font-size: 24px;
            color: #fafafa;
            background-color: #01019c;
        }

        strong {
            font-weight: 600;
            font-size: 1rem;
            margin-inline-end: 4px;
        }

        .size {
            color: #adabab;
            font-size: 0.625rem;
            margin-inline: 4px;
        }
    }

    &:before {
        // Indentation spacer
        content: " ";
        display: inline-block;
        width: calc(var(--item-depth) * 20px);
        height: var(--height);
        border-bottom: 1px solid #707070;
        // background-color: #707070;
        margin-inline-end: 6px;
    }
}

ul.children {
    margin-block: 12px;
    border-radius: 4px;
    list-style: none;
    padding: 8px;

}
</style>