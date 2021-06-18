<template>
    <header-component/>
    <main class="nh-100">
        <div class="container nh-100">
            <div class="row align-items-center nh-100">
                <div class="col-lg-8 mx-auto">
                    <h1>Shortener</h1>
                    <form @submit.prevent="createShortURL" v-if="isShortURLCreated === false">
                        <div class="input-group">
                            <input type="text" class="form-control" :class="{'is-invalid': isError}" placeholder="Past URL here ..." v-model="url">
                            <button class="btn btn-outline-success" type="submit" v-if="isLoading === false">Create!</button>
                            <button class="btn btn-outline-success" type="button" disabled v-if="isLoading">
                                <span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
                                <span class="visually-hidden">Loading...</span>
                            </button>
                            <div class="invalid-tooltip" v-show="isError">
                                {{ error }}
                            </div>
                        </div>
                    </form>
                    <div class="input-group" v-if="isShortURLCreated">
                        <input type="text" class="form-control is-valid" placeholder="Here is yiu short URL" v-model="shortUrl">
                        <button class="btn btn-outline-success" type="button" @click="copyShortURL">
                            <i class="fa fa-copy"></i>
                        </button>
                        <button class="btn btn-outline-success" type="button" @click="createNewOne">
                            <i class="fa fa-plus"></i>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </main>
</template>

<script>
    import HeaderComponent from '../components/Header';
    import LinkService from '../services/link';
    import useClipboard from 'vue-clipboard3';

    export default {
        name: 'Home',
        components: {
            HeaderComponent
        },
        data() {
            return {
                url: '',
                shortUrl: '',
                error: '',
                isError: false,
                isLoading: false,
                isShortURLCreated: false,
            };
        },
        methods: {
            async createShortURL() {
                let newLink;
                this.isLoading = true;
                this.isError = false;
                try {
                    newLink = await LinkService.create(this.url)
                }catch (e) {
                    this.isLoading = false;
                    this.isError = true;
                    this.error = e.response.data.error;
                    return;
                }

                this.isLoading = false;
                this.isShortURLCreated = true;
                this.shortUrl = window.location.origin + '/' + newLink.data.hash;
                LinkService.saveUserLinks(newLink.data)
            },
            async copyShortURL() {
                const { toClipboard } = useClipboard();
                await toClipboard(this.shortUrl)
            },
            createNewOne() {
                this.url = '';
                this.shortUrl = '';
                this.isShortURLCreated = false;
            }
        }
    }
</script>

<style>
    .nh-100 {
        height: calc(100% - 65px);
    }
</style>
