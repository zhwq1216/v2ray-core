import axios from 'axios';

function getBaseURL() {
    let baseURL = localStorage.getItem("baseURL") || "http://localhost:8035/v2ray";
    return baseURL;
}

axios.defaults.baseURL = getBaseURL();

export default {
    async post(url, data, options = {}) {
        try {
            return await axios.post(url, data, options);
        } catch (e) {
            return {
                status: 500,
                statusText: e.toString(),
            }
        }

    },
    async get(url, options) {
        try {
            return await axios.get(url, options);
        } catch (e) {
            return {
                status: 500,
                statusText: e.toString(),
            }
        }
    },
    setBaseURL(baseURL) {
        axios.defaults.baseURL = baseURL;
        localStorage.setItem("baseURL", baseURL);
    },
    getBaseURL,
};
