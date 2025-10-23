import { defineStore, acceptHMRUpdate } from 'pinia';
import DirectoryRepository from './repositories/directory.js';
import ObjectRepository from './repositories/object.js';

export const directoryStore = defineStore('directory', {
    state: () => ({
        directory: null,
        error: null
    }),
    getters: {
        loaded: (state) => state.directory != null,
    },
    actions: {
        async fetchDirectory(bucket, path) {
            const repo = new DirectoryRepository();
            try {
                const res = await repo.fetchDirectoryTree(bucket, path);
                if (res.data instanceof Object) {
                    this.directory = res.data?.children;
                } else if (res.data instanceof Array) {
                    this.directory = res.data;
                } else {
                    throw new Error("Unexpected data format received from server.");
                }
                this.error = null;
                console.log("Fetched directory:", this.directory);
            } catch (error) {
                console.error("Error fetching directory:", error);
                this.error = {
                    message: "Failed to fetch directory.",
                    details: error,
                }
            }
        },
        async updateDirectory(directory) {

        },
        async downloadObject(bucket, path) {
            const repo = new ObjectRepository();
            try {
                const res = await repo.fetchObject(bucket, path);
                if (res.status === 200) {
                    console.log("Downloaded object:", path);
                    const blob = new Blob([res.data]);
                    const url = window.URL.createObjectURL(blob);
                    const link = document.createElement('a');
                    link.href = url;
                    link.setAttribute('download', path.split('/').pop());
                    document.body.appendChild(link);
                    link.click();
                    link.parentNode.removeChild(link);
                    window.URL.revokeObjectURL(url);
                }
                
            } catch (error) {
                console.error("Error downloading object:", error);
                this.error = {
                    message: "Failed to download object.",
                    details: error,
                };
            }
        }
    },
});

if (
    import.meta.hot) {
    import.meta.hot.accept(acceptHMRUpdate(directoryStore,
        import.meta.hot));
}
if (
    import.meta.webpackHot) {
    import.meta.webpackHot.accept(acceptHMRUpdate(directoryStore,
        import.meta.webpackHot));
}