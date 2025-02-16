const express = require("express")
const cors = require("cors")
const upload = require("./middleware/multer")
const { default: axios } = require("axios")
require("dotenv").config()

const port = 3000
const app = express()

app.use(cors())

app.post("/upload_csv",upload.single("file"),async(req,res) => {
    try {
        const file = req.file.filename; 
        await axios.post(process.env.URL, { "filename" : file  },{headers: { "Content-Type": "application/json" }});
        return res.status(200).json({ message : "OK",data : file})
    } catch (error) {
        console.error(error)
        return res.status(500).json({ error })
    }
})


app.listen(port,() => {
    console.log("server berjalan....")
})