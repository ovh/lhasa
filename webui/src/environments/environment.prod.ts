/**
 * fix BASE HREF problems
 */
const href = document.location.pathname.split('/');
let base = '';
let appname = 'appcatalog';
let prefix = '';
let i = 0;
for (i = 1; i < href.length; i++) {
  base += '/' + href[i];
  if (href[i].endsWith('appcatalog')) {
    appname = href[i];
    break;
  }
}

// Check DEV env
if (document.location.pathname.indexOf('/dev/') >= 0) {
  prefix = '/dev';
}

/**
 * fix API access routes
 */
export const environment = {
  production: true,
  tracing: false,
  href: base,
  appcatalog: {
    baseUrlApi: prefix + '/uservice/gateway/' + appname,
    baseUrlUi: prefix + '/uservice/gateway/' + appname
  },
  appcatalogCompanion: {
    baseUrlApi: prefix + '/uservice/gateway/' + appname + '-companion',
    baseUrlUi: prefix + '/uservice/gateway/' + appname + '-companion'
  }
};
