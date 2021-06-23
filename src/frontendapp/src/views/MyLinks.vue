<template>
    <header-component />
    <main>
        <div class="container">
            <h3>My links</h3>
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
            <h4 v-if="noResults">No links found.</h4>
        </div>

    </main>
</template>

<script>
    import HeaderComponent from '../components/Header';
    import LinkService from '../services/link';
    import useClipboard from 'vue-clipboard3'

    export default {
        name: 'MyLinks',
        components: {
            HeaderComponent
        },
        data() {
            return {
                limit: 10,
                hashes: [],
                links: [],
            };
        },
        computed: {
            noResults() {
                return this.links.length === 0
            },
            totalLinks() {
                return this.hashes.length
            },
            hasMore() {
                return this.totalLinks > this.links.length
            }
        },
        created() {
            this.hashes = LinkService.getUserLinksHashes();
            this.getUserLinks()
        },
        methods: {
            getShortUrt(hash) {
                return window.location.origin + '/' + hash
            },
            showMore() {
                this.getUserLinks()
            },
            async copyShortURL(shortUrl) {
                const { toClipboard } = useClipboard();
                await toClipboard(shortUrl)
            },
            async getUserLinks() {
                if(this.totalLinks === 0) {
                    return;
                }

                let ids = this.hashes.slice(this.links.length, this.links.length + this.limit);

                let result;
                try {
                    result = await LinkService.getUserLinks(ids)
                }
                catch(e) {
                    return
                }

                this.links.push(...result.data);
            }
        }
    }
</script>
