const express = require('express');
const app = express();

app.listen(process.env.PORT, () => {
    console.log(`${process.env.NAME} stated on PORT ${process.env.PORT}`)
})

app.get('/healthcheck', (req, res) => {
    res.status(200).json({
        server: process.env.NAME,
        status: "success"
    })
})

app.get('/', (req, res) => {
    res.status(200).json({
        server: process.env.NAME,
        status: "success",
        body: "hello world"
    })
})
