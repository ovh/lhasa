/**
 * fix BASE HREF problems
 */
const href = document.location.pathname.split('/');
let base = '';
let appname = 'appcatalog';
let i = 0;
for (i = 1; i < href.length; i++) {
  base += '/' + href[i];
  if (href[i].endsWith('appcatalog')) {
    appname = href[i];
    break;
  }
}

/**
 * fix API access routes
 */
export const environment = {
  production: false,
  tracing: true,
  href: base,
  appcatalog: {
    baseUrlApi: '/uservice/gateway/' + appname,
    baseUrlUi: '/uservice/gateway/' + appname
  },
  appcatalogCompanion: {
    baseUrlApi: '/uservice/gateway/' + appname + '-companion',
    baseUrlUi: '/uservice/gateway/' + appname + '-companion'
  }
};
