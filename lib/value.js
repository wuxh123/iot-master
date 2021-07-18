
module.exports = class Value {
    /**
     * @type Device
     */
    device;

    name;
    last;
    readonly;
    watchers = [];
}