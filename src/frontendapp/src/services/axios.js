import axios from 'axios';
import { API_URL } from './config';

const initAxios = () => {
    axios.defaults.baseURL = API_URL;
};

export default initAxios;