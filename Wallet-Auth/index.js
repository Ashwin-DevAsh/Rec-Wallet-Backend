const express = require("express");
const bodyParser = require("body-parser");
const authEndPoint = require("./Routes/Auth");
const userEndPoint = require("./Routes/User");
const merchantEndPoint = require("./Routes/Merchants");
const mongoose = require("./Database/mongodb");

const app = express();

app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());

app.get("/", (req, res) => {
  res.status(200).json({ Message: "Welcome to rec wallet" });
});

app.use(userEndPoint);
app.use(authEndPoint);
app.use(merchantEndPoint);

app.listen(8000, () => {
  console.log("connecte at port 8000");
});
