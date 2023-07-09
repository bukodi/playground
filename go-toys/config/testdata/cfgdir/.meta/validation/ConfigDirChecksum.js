const fs = require('fs'),
    path = require('path'),
    crypto = require('crypto');

let startDir = './../..'
let checksum = new Uint8Array(32)

let processDir = function (dir) {
    for (let item of fs.readdirSync(dir)) {
        if (item.startsWith('.')) {
            continue;
        }
        // get the item path
        let itemPath = path.join(dir, item);
        let stats = fs.statSync(itemPath);
        if (stats.isDirectory()) {
            processDir(itemPath);
            continue;
        }
        
        const hash = crypto.createHash('sha256');
        const relativeToStart = path.relative(startDir, itemPath)
        hash.update((new TextEncoder('utf-8')).encode(relativeToStart));
        hash.update(fs.readFileSync(itemPath));
        const itemHash = hash.digest();

        console.log(Buffer.from(itemHash).toString('hex'), relativeToStart );

        xorItemHash(itemHash)
    }
}

let xorItemHash = function (itemHash) {
    for (let i = 0; i < 32; i++) {
        checksum[i] = checksum[i] ^ itemHash[i]
    }
}

processDir(startDir);
console.log("Checksum :", Buffer.from(checksum).toString('hex'));
