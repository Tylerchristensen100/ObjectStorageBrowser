import { directoryStore } from '@/global/stores/directory_store';

export default class ObjectHelpers {

    static download(bucket, path) {
        console.log('Downloading object:', bucket, path);
        const store = directoryStore();
        store.downloadObject(bucket, path);
    }
}