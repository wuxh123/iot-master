const Influx = require("influx");

const influx = new Influx.InfluxDB({
    host: 'localhost',
    database: 'express_response_db',
    schema: [
        {
            measurement: 'response_times',
            fields: {
                path: Influx.FieldType.STRING,
                duration: Influx.FieldType.INTEGER
            },
            tags: [
                'host'
            ]
        }
    ]
})

Object.assign(module.exports, Influx);
