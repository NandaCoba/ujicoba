const multer = require("multer");

const upload = multer({
    storage : multer.diskStorage({
           filename : (req,file,cb) => {
               cb(null,file.originalname.split(".csv").join("_") + Date.now() + ".csv")
           },
           destination : (req,file,cb) => {
               cb(null,"worker/uploads")
           }
   }),
   fileFilter :(req,file,cb) => {
    const allowedTypes = ["text/csv", "application/vnd.ms-excel"];
       if(!allowedTypes.includes(file.mimetype)) return cb("file type is not allowed",null)
        cb(null,true)
   }
})


module.exports = upload


