<template>
    <footer class="footer">
        <p>© 2025 Object Storage Browser</p>
        <ul>
            <li><a href="https://github.com/Tylerchristensen100/ObjectStorageBrowser" target="_blank">GitHub</a></li>
        </ul>
    </footer>
</template>

<script setup></script>

<style lang="scss" scoped>
.footer {
    background-color: #03164b;
    color: #fafafa;
    padding: 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-direction: column;

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