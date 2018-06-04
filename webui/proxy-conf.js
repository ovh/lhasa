/**
 * gateway deployment simulation
 * @param {*} req
 * @param {*} res
 * @param {*} proxyOptions
 */
const redirect = function (req, res, proxyOptions) {
  console.log("Catch all.", req.originalUrl);
  res.set('location', 'http://localhost:4200/?redirect=' + req.originalUrl.replace('?', '&'));
  res.status(301)
};

const injectCreds = function (req, res, proxyOptions) {
  console.log("Inject.");
  req.headers["X-Remote-User"] = "fabien.meurillon";
};

const PROXY_CONFIG = {
  "/api": {
    "target": "http://localhost:8081",
    "secure": false
  },
  "/ui//assets": {
    "target": "http://localhost:4200",
    "pathRewrite": {"^/ui//assets": "/assets"},
    "secure": false
  },
  "/uservice/gateway/appcatalog/api": {
    "target": "http://localhost:8081",
    "secure": false,
    "changeOrigin": true,
    "logLevel": "info",
    "pathRewrite": {"^/uservice/gateway/appcatalog/api": "/api"},
    "bypass": injectCreds
  },
  "/all": {
    "secure": false,
    "bypass": redirect
  },
}

module.exports = PROXY_CONFIG;
