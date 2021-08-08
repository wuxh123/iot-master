const config = require("./config");
const utils = require("./utils");


module.exports = {
  config,
  utils,
  miniprogram: require("./miniprogram"),
  official: require("./official"),
};
