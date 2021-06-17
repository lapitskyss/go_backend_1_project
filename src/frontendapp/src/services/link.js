import axios from 'axios'

const LinkService = {
    create(url) {
        return axios.post('/links', { url })
    },
    findBy(page, limit, query) {
        const params = {
            page,
            limit
        };
        if(query) {
            params.query = query;
        }

        return axios.get('/links', {
            params
        })
    }
};

export default LinkService;