import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/tracks',
      name: 'tracks',
      component: () => import(/* webpackChunkName: "tracks" */ './views/Tracks.vue')
    },
    {
      path: '/about',
      name: 'about',
      component: () => import(/* webpackChunkName: "tracks" */ './views/About.vue')
    }
  ]
})
