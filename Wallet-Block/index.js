const express = require("express");
const redisClient = require("./Database/Connections/redisConnection");
const transactionBlocks = require("./Routes/Transactions");

const app = express();

app.get("/", (req, res) => {
  res.send({ message: "wellcome to rec-wallet block" });
});

app.use(transactionBlocks);

app.listen(9000, () => {
  console.log("listing on 9000....");
});
