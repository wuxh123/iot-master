const Observer = require("./observer");

module.exports = class Context extends Observer {

    constructor(obj) {
        super();

        obj && Object.assign(this, obj)
    }

    clone() {
        return Object.assign({}, this, this.values());
    }
}