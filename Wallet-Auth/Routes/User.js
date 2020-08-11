const app = require("express").Router();
const Users = require("../Schemas/users");
const Otp = require("../Schemas/otp");
const jwt = require("jsonwebtoken");
const postgres = require("../Database/postgresql");
const axios = require("axios");

app.post("/addUser", async (req, res) => {
  var user = req.body;
  console.log(user);
  if (
    !user.name ||
    !user.email ||
    !user.number ||
    !user.password ||
    !user.fcmToken ||
    !user.qrCode
  ) {
    res.status(200).send([{ message: "error" }]);
    return;
  }
  var userID = `rpay@${user.number}`;

  try {
    var otp = (
      await postgres.query(
        "select * from otp where number = $1 and verified=true",
        [user.number]
      )
    ).rows;

    if (otp.length == 0) {
      res.json([{ message: "failed" }]);
      return;
    }

    var testUser = (
      await postgres.query("select * from users where id = $1 ", [userID])
    ).rows;

    if (testUser.length != 0) {
      res.json([{ message: "User already exist" }]);
      return;
    }
    await postgres.query("delete from info where id=$1;", [userID]);
    await postgres.query("delete from amount where id=$1;", [userID]);
    await postgres.query("insert into info values($1,$2,null,null)", [
      userID,
      user.fcmToken,
    ]);
    await postgres.query("insert into amount(id,balance) values($1,0)", [
      userID,
    ]);

    var token = jwt.sign(
      {
        name: user.name,
        id: userID,
        number: user.number,
        email: user.email,
      },
      process.env.PRIVATE_KEY
    );

    var blockResult = await axios.post(
      "http://wallet-block:9000/addUserBlock",
      {
        id: userID,
        initialAmount: 0,
      }
    );

    if ((blockResult.data["message"] = "done")) {
      await postgres.query(
        `insert into users(name,number,email,password,id,qrCode) values($1,$2,$3,$4,$5,$6)`,
        [user.name, user.number, user.email, user.password, userID, user.qrCode]
      );
      res.json([{ message: "done", token }]);
    } else {
      res.json([{ message: "failed" }]);
    }
  } catch (err) {
    console.log(err);
    res.json([{ message: "failed" }]);
  }
});

app.get("/getUsers", (req, res) => {
  Users.find({}, ["name", "number", "email", "imageURL", "id"])
    .exec()
    .then((doc) => {
      console.log(doc);
      res.status(200).send(doc);
    })
    .catch((e) => {
      console.log(e);
      res.send(e);
    });
});

app.post("/changePassword", (req, res) => {
  console.log(req.get("token"));
  jwt.verify(req.get("token"), process.env.PRIVATE_KEY, function (
    err,
    decoded
  ) {
    if (err) {
      console.log(err);
      res.status(200).send({ message: "error" });
      return;
    } else {
      console.log("Changing password...");
      var data = req.body;
      console.log(data);
      if (!data.id || !data.oldPassword || !data.newPassword) {
        res.status(200).send({ message: "error" });
        return;
      }

      Users.findOne({ id: data.id })
        .exec()
        .then((docs) => {
          console.log("data = ", docs);
          if (docs == null) {
            res.status(200).send({ message: "error" });
            return;
          }

          if (docs.password == data.oldPassword) {
            Users.findOneAndUpdate(
              { id: data.id },
              { password: data.newPassword },
              (err, doc) => {
                if (err) {
                  console.log(err);
                  res.status(200).send({ message: "error" });
                  return;
                } else {
                  res.status(200).send({ message: "done" });
                  return;
                }
              }
            );
          } else {
            console.log("data = ", null);

            res.status(200).send({ message: "error" });
            return;
          }
        });
    }
  });
});

module.exports = app;
