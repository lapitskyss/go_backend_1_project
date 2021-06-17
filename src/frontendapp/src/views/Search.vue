<template>
    <header-component @search-links="onSearchLinks"/>
    <main>
        <div class="container">
            <h3>Search result for: {{ query }}</h3>
            <table class="table">
                <thead>
                <tr>
                    <th scope="col">#</th>
                    <th scope="col">Link</th>
                    <th scope="col">Short URL</th>
                    <th scope="col">Copy!</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="(link, index) in links" :key="index">
                    <th scope="row">{{ index + 1 }}</th>
                    <td>{{ link.url }}</td>
                    <td>{{ getShortUrt(link.hash) }}</td>
                    <td>
                        <button class="btn btn-sm btn-outline-success" type="button" @click="copyShortURL(getShortUrt(link.hash))">
                            <i class="fa fa-copy"></i>
                        </button>
                    </td>
                </tr>
                </tbody>
            </table>
            <div class="d-grid" v-if="hasMore">
                <button class="btn btn-outline-success" @click="showMore">Show more</button>
            </div>
            <h4 v-if="noResults">No result found</h4>
        </div>

    </main>
</template>

<script>
    import HeaderComponent from '../components/Header';
    import LinkService from '../services/link';
    import useClipboard from 'vue-clipboard3'

    export default {
        name: 'Search',
        components: {
            HeaderComponent
        },
        data() {
            return {
                query: '',
                page: 1,
                limit: 10,
                pages: 1,
                links: [],
            };
        },
        computed: {
            hasMore() {
                return this.pages > this.page
            },
            noResults() {
                return this.links.length === 0
            }
        },
        created() {
            this.query = this.$route.query.q;
            this.fetchData()
        },
        methods: {
            getShortUrt(hash) {
                return window.location.origin + '/' + hash
            },
            onSearchLinks({ search }) {
                this.query = search;
                this.fetchData();
            },
            async showMore() {
                this.page = this.page + 1;
                let data = await LinkService.findBy(this.page, this.limit, this.query);
                this.links.push(...data.data.links)
            },
            async fetchData() {
                this.page = 1;
                let data = await LinkService.findBy(this.page, this.limit, this.query);
                this.links = data.data.links;
                this.pages = data.data.pages
            },
            async copyShortURL(shortUrl) {
                const { toClipboard } = useClipboard();
                await toClipboard(shortUrl)
            },
        }
    }
</script>
