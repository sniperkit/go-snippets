import App from './containers/App';
import PublicRecord from './containers/PublicRecord';
import NotFound from './components/NotFound';

function errorLoading(err) {
  console.error('Dynamic page loading failed', err);
}

function loadRoute(cb) {
  return module => cb(null, module.default);
}

export default {
  path: "/",
  component: App,
  indexRoute: { component: PublicRecord },
  // getIndexRoute(partialState, cb) {
  //  import('./containers/PublicRecord').then(loadRoute(cb)).catch(errorLoading)
  // },
  childRoutes: [
    {
      path: '/login',
      getComponent(location, cb) {
        import('./containers/Login').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/url',
      getComponent(location, cb) {
        import('./containers/Url').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: "/add-dataset",
      getComponent(location, cb) {
        import('./containers/ArchiveUrl').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/collections',
      getComponent(location, cb) {
        import('./containers/Collections').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/collections/:id',
      getComponent(location, cb) {
        import('./containers/Collection').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/content/:hash',
      getComponent(location, cb) {
        import('./containers/Content').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/communities',
      getComponent(location, cb) {
        import('./containers/Communities').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/coverage',
      getComponent(location,cb) {
        import('./containers/Coverage').then(loadRoute(cb)).catch(errorLoading);
      }
    },
    {
      path: '/settings',
      getComponent(location, cb) {
        import('./containers/UserSettings').then(loadRoute(cb)).catch(errorLoading);
      },
      childRoutes: [
        {
          path: "/keys",
          getComponent(location, cb) {
            import('./containers/Keys').then(loadRoute(cb)).catch(errorLoading);
          },
        },
      ],
    },
    {
      path: '/users/:user',
      getComponent(location, cb) {
        import('./containers/User').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/users/:user/settings',
      getComponent(location, cb) {
        import('./containers/UserSettings').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/signup',
      getComponent(location, cb) {
        import('./containers/Signup').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/primers',
      getComponent(location, cb) {
        import('./containers/Primers').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/primers/:id',
      getComponent(location, cb) {
        import('./containers/Primer').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/sources/:id',
      getComponent(location, cb) {
        import('./containers/Source').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/tasks',
      getComponent(location, cb) {
        import('./containers/Tasks').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/uncrawlables',
      getComponent(location, cb) {
        import('./containers/Uncrawlables').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/uncrawlables/:id',
      getComponent(location, cb) {
        import('./containers/Uncrawlable').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/edit',
      getComponent(location, cb) {
        import('./containers/FileEditor').then(loadRoute(cb)).catch(errorLoading);
      },
    },
    {
      path: '/*',
      component: NotFound,
    },
  ],
};
