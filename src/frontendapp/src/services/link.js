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
    },
    saveUserLinks(link) {
        let userLinks = this.getUserLinks();
        userLinks.push(link);

        localStorage.setItem("_links", JSON.stringify(userLinks));
    },
    getUserLinks() {
        let links = localStorage.getItem("_links");
        if(!links) {
            return []
        }

        return JSON.parse(links);
    }
};

export default LinkService;