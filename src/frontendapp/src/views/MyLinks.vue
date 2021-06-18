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
                links: [],
            };
        },
        computed: {
            noResults() {
                return this.links.length === 0
            }
        },
        created() {
            this.links = LinkService.getUserLinks()
        },
        methods: {
            getShortUrt(hash) {
                return window.location.origin + '/' + hash
            },
            async copyShortURL(shortUrl) {
                const { toClipboard } = useClipboard();
                await toClipboard(shortUrl)
            },
        }
    }
</script>
