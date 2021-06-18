<template>
    <header>
        <nav class="navbar navbar-expand-lg navbar-light bg-light">
            <div class="container-fluid">
                <router-link class="navbar-brand" to="/">Shortener</router-link>
                <ul class="navbar-nav me-auto mb-2 mb-lg-0">
                    <li class="nav-item">
                        <router-link class="nav-link active" aria-current="page" to="/my-links">My links</router-link>
                    </li>
                </ul>
                <form class="d-flex" @submit.prevent="onSearch">
                    <input class="form-control me-2" type="search" placeholder="Search" v-model="search">
                    <button class="btn btn-outline-success" type="submit">Search</button>
                </form>
            </div>
        </nav>
    </header>
</template>

<script>
    export default {
        name: 'HeaderComponent',
        data() {
            return {
                search: ''
            };
        },
        methods: {
            onSearch() {
                this.$router.push({ name: 'search', query: { q: this.search }});
                this.$emit('search-links', { search: this.search });
                this.search = ''
            }
        }
    }
</script>
