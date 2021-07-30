const fs = require('fs');
const mammoth = require("mammoth");

exports.convert = function (path) {
    return new Promise(function (resolve, reject) {
        const filename = path + '.html';
        fs.access(filename, fs.constants.F_OK, (err) => {
            if (err) {
                mammoth.convertToHtml({path}).then(function (result){
                    resolve(result.value);
                    fs.writeFile(filename, result.value, function (err){
                        //console.error
                    })
                }).catch(reject)
            } else {
                fs.readFile(filename, function (err, data){
                    if (err) {
                        reject(err)
                    } else {
                        resolve(data);
                    }
                })
            }
        });
    });
}
