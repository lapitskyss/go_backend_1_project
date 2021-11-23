import axios from 'axios';

const initAxios = () => {
    axios.defaults.baseURL = window["__API_URL__"];
};

export default initAxios;