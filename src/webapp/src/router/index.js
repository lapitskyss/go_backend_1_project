import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home';
import Search from '../views/Search';
import NotFound from '../views/NotFound';
import MyLinks from "../views/MyLinks";

const routes = [
    {
        path: '/',
        name: 'home',
        component: Home
    },
    {
        path: '/my-links',
        name: 'my-links',
        component: MyLinks,
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