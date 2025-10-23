
import Client from "@/global/networking/client";


export default class ObjectRepository {
    #_client;
    #_base = '/object';

    constructor() {
        this.#_client = Client.getClient();
    }


    async fetchObject(bucket, path) {
        return this.#_client.get(`${this.#_base}`, {
            params: {
                'bucket': bucket,
                'path': path,
            }
        });
    }

    async uploadObject(bucket, path, file) {

        let data = new FormData();
        data.append('file', file);

        return this.#_client.post(`${this.#_base}`, data, {
            params: {
                'bucket': bucket,
                'path': path,
            }
        });
    }

    async deleteObject(bucket, path) {
        return this.#_client.delete(`${this.#_base}`, {
            params: {
                'bucket': bucket,
                'path': path,
            }
        });

    }
}