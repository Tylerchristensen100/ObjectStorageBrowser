
import Client from "@/global/networking/client";
export default class BucketsRepository {
    #_client;
    #_base = '/buckets';

    constructor() {
        this.#_client = Client.getClient();
    }


    async fetchBuckets() {
        return this.#_client.get(`${this.#_base}`);
    }

    async fetchBucket(bucket) {
        return this.#_client.get(`${this.#_base}/${bucket}`);
    }

    async createBucket(name) {
        throw 'Not implemented [createBucket]';

    }

}