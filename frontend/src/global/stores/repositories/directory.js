
import Client from "@/global/networking/client";
export default class DirectoryRepository {
    #_client;
    #_base = '/directory';

    constructor() {
        this.#_client = Client.getClient();
    }


    async fetchDirectoryList(bucket, path, recursive = false) {
        return this.#_client.get(`${this.#_base}`, {
            params: {
                'bucket': bucket,
                'path': path,
                'recursive': recursive
            }
        });
    }

    async fetchDirectoryTree(bucket, path, recursive = false) {
        return this.#_client.get(`${this.#_base}/tree`, {
            params: {
                'bucket': bucket,
                'path': path,
                'recursive': recursive
            }
        });
    }
}