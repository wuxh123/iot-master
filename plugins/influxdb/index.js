const {InfluxDB, Point} = require('@influxdata/influxdb-client')

const defaultOptions = {};
exports.options = Object.assign({}, defaultOptions);

exports.config = function (options) {
    exports.options = Object.assign({}, defaultOptions, options);
    const client = new InfluxDB(exports.options)
    exports.writeApi = client.getWriteApi(exports.options.org, exports.options.bucket)
    exports.queryApi = client.getQueryApi(exports.options.org)
}


//Object.assign(module.exports, Influx);
exports.write = function (table, tags, values) {
    const point = new Point(table);
    tags.forEach(t => point.tag(t.name, t.value));
    values.forEach(v => {
        switch (v.type) {
            case 'boolean':
                point.booleanField(v.name, v.value);
                break;
            case 'word':
            case 'dword':
                point.uintField(v.name, v.value);
                break;
            case 'float':
            case 'double':
            case 'le-float':
            case 'le-double':
                point.floatField(v.name, v.value);
                break;
            case 'uint8':
            case 'uint16':
            case 'uint32':
            case 'uint64':
            case 'le-uint8':
            case 'le-uint16':
            case 'le-uint32':
            case 'le-uint64':
                point.uintField(v.name, v.value);
                break;
            case 'int8':
            case 'int16':
            case 'int32':
            case 'int64':
            case 'le-int8':
            case 'le-int16':
            case 'le-int32':
            case 'le-int64':
                point.intField(v.name, v.value);
                break;
        }
    })
    exports.writeApi.writePoint(point)
}

exports.query = function (table, tags, field, window, start, stop) {
    return new Promise((resolve, reject) => {
        const query = `
            from(bucket: "${exports.options.bucket}")
              |> range(start: ${start}, stop: ${stop})
            `
            + tags.map(t => `|> filter(fn: (r) => r["${t.name}"] == "${t.value}")`).join('\n') +
            ` |> filter(fn: (r) => r["_field"] == "${field}")
              |> aggregateWindow(every: ${window}, fn: mean, createEmpty: false)
              |> yield(name: "mean")
            `;

        const results = [];
        exports.queryApi.queryRows(query, {
            next(row, tableMeta) {
                //console.log(row)
                const obj = tableMeta.toObject(row)
                results.push({
                    [field]: obj._value,
                    time: obj._time
                })
            },
            error(error) {
                reject(error);
            },
            complete() {
                resolve(results)
            },
        });
    })
}
