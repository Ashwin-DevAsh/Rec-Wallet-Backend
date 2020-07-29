require("dotenv").config(".env");

const app = require("express")();
const server = require("http").createServer(app);
const io = require("socket.io")(server);
const postgres = require("./Database");

io.on("connection", (client) => {
  console.log("connected ", client.socketID);

  client.on("getInformation", (data) => {
    var id = data["id"];
    var token = data["fcmToken"];
    client.join(id,(err)=>{
        if(err){
            console.log(err)
        }
    })
    console.log(id, token);
    updateOnline(id, client.socketID, token, true);
  });

  client.on("disconnect", () => {
      updateOffline(db, client.socketID);
  });
});

async function updateOnline(id, socketID, fcmToken, isOnline) {
  var insertStatement = `update info set isonline=$4 , socketid=$2 , fcmToken = $3  where id=$1`;
  try {
    await postgres.query(insertStatement, [id, socketID, fcmToken, isOnline]);
  } catch (e) {
    console.log(e);
  }
}

function updateOffline(socketID) {
    insertStatement = `update info set socketid=null , isonline=false where socketid=$1`
    try{
	  await postgres.query(insertStatement, [socketID])
    }catch(e){
        console.log(e)
    }
		
}


server.listen(7000, () => {
  console.log("Listing at 7000");
});

console.log("Last....");
