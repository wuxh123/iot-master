const mongo = require_plugin("mongodb");

exports.createEvent = function (event) {
    mongo.db.collection("event").insertOne(event).then(()=>{
    }).catch(console.error);
}