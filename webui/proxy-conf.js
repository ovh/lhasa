/**
 * gateway deployment simulation
 * @param {*} req
 * @param {*} res
 * @param {*} proxyOptions
 */

const injectCreds = function (req, res, proxyOptions) {
  console.log("Inject.");
  req.headers["X-Remote-User"] = "fabien.meurillon";
};

const PROXY_CONFIG = {
  "/api": {
    "target": "http://localhost:8081",
    "secure": false
  },
  "/all": {
    "secure": false
  },
}

module.exports = PROXY_CONFIG;
