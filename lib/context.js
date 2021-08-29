const Observer = require("./observer");

module.exports = class Context extends Observer{
    constructor() {
        super();
    }

    clone() {
        return Object.assign({}, this, this.values());
    }
}