module.exports = class Observer {
    $$values = {};
    $$changes = {};

    define(name, value) {
        const that = this;
        that.$$values[name] = value;
        Object.defineProperty(this, name, {
            enumerable: true,
            //writable: true,
            //value: value,
            get() {
                return that.$$values[name];
            },
            set(value) {
                //console.log('Observer set', name, value);
                that.$$values[name] = value;
                that.$$changes[name] = value;
            }
        })
    }

    set(values) {
        Object.assign(this.$$values, values);
    }

    values() {
        return this.$$values;
    }

    changes(){
        const cs = this.$$changes;
        this.$$changes = {};
        return cs;
    }

}