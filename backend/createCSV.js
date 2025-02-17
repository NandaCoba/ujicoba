const fs = require("fs");
const path = require("path");

const createBatch = async () => {
    const filePath = path.join(__dirname, "data.csv");
    const writeStream = fs.createWriteStream(filePath, { flags: "a" });

    for (let index = 1; index <= 5_000; index++) {
        writeStream.write(`nanda_${index},Software Engineer,IT,7000000\n`);
    }

    writeStream.end(); 
    console.log("Data batch selesai ditulis.");
};

createBatch();