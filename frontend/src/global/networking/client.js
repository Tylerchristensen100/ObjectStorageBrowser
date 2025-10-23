import axios from "axios";

const API = {
    getClient: function () {
        let ops = {
            baseURL: import.meta.env.VITE_BASE_URL + '/api',
        };

        return axios.create(ops);
    }
};

export default API;
