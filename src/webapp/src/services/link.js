import axios from 'axios'

const LinkService = {
    create(url) {
        return axios.post('/links', { url })
    },
    findBy(page, limit, query, sort = 'created_at', order= 'desc') {
        const params = {
            page,
            limit,
            sort,
            order
        };
        if(query) {
            params.query = query;
        }

        return axios.get('/links/search', {
            params
        })
    },
    saveUserLink(linkHash) {
        let userLinks = this.getUserLinksHashes();
        if(userLinks.includes(linkHash) === false) {
            userLinks.unshift(linkHash);

            localStorage.setItem("_links", JSON.stringify(userLinks));
        }
    },
    getUserLinksHashes() {
        let links = localStorage.getItem("_links");
        if(!links) {
            return []
        }

        return JSON.parse(links);
    },
    getUserLinks(ids) {
        const params = {
            ids: ids.join()
        };

        return axios.get('/links', {
            params
        })
    }
};

export default LinkService;