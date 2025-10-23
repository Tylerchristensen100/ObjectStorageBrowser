<template>
    <nav class="navbar">
        <ul>
            <li><router-link to="/">Directory</router-link></li>
            <li><router-link to="/buckets">Buckets</router-link></li>
            <li><a href="https://github.com/Tylerchristensen100/ObjectStorageBrowser" target="_blank">GitHub</a></li>

        </ul>
    </nav>
</template>

<script setup></script>

<style lang="scss" scoped>
.navbar {
    background-color: #03164b;
    color: #fafafa;
    padding: 1rem;

    ul {
        list-style: none;
        display: flex;
        gap: 1rem;
        margin: 0;
        padding: 0;

        li {
            a {
                color: #fafafa;
                text-decoration: none;

                &:hover {
                    text-decoration: underline;
                }
            }
        }
    }
}
</style>