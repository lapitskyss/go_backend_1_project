import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home';
import Search from '../views/Search';
import NotFound from '../views/NotFound';

const routes = [
    {
        path: '/',
        name: 'home',
        component: Home
    },
    {
        path: '/search',
        name: 'search',
        component: Search,
    },
    {
        path: '/:pathMatch(.*)*',
        name: 'not-found',
        component: NotFound,
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router