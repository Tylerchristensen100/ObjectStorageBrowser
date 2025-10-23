import { createRouter, createWebHistory } from 'vue-router';
import Directory from '../screens/directory/Directory.vue';
import Buckets from '../screens/buckets/Buckets.vue';
import Login from '../screens/Login.vue';


const routes = [
    {
        path: '/:bucket/:path(.*)?',
        name: 'Directory',
        component: Directory,
    },
    {
        path: '/buckets',
        name: 'Buckets',
        component: Buckets,
    },
    {
        path: '/',
        redirect: '/buckets',
    },
    {
        path: '/login',
        name: 'Login',
        component: Login,
    },
];


const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
});


router.beforeEach((to, from, next) => {
    if (to.path.includes('//')) {
        const fixedPath = to.path.replaceAll('//', '/');
        next(fixedPath);
    } else {
        next();
    }
    // if (to.path != '/login' && !localStorage.getItem('auth_token')) {
    //     next('/login');
    // } else {
    //     next();
    // }
});

window.history.scrollRestoration = 'auto';

export default router;
