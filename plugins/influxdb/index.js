const {InfluxDB, Point} = require('@influxdata/influxdb-client')
const _ = require("lodash");

const defaultOptions = {};

const cfg = load_config("influxdb");

const options = _.defaultsDeep({}, cfg, defaultOptions);

const client = new InfluxDB(options);

exports.writeApi = client.getWriteApi(options.org, options.bucket)
exports.queryApi = client.getQueryApi(options.org);



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
                point.floatField(v.name, v.value);
                break;
            case 'uint8':
            case 'uint16':
            case 'uint32':
            case 'uint64':
                point.uintField(v.name, v.value);
                break;
            case 'int8':
            case 'int16':
            case 'int32':
            case 'int64':
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
            + Object.keys(tags).map(t=>`|> filter(fn: (r) => r["${t}"] == "${tags[t]}")`).join('\n') +
            //+ tags.map(t => `|> filter(fn: (r) => r["${t.name}"] == "${t.value}")`).join('\n') +
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
