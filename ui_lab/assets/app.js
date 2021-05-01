var app = new Vue({
    el: '#app',
    data: {
        ws: null,
        serverUrl: "ws://localhost:8095/ws/",
        roomInput: null,
        rooms: [],
        user: {
            name: "alameddinc"
        },
        users: [],
    },
    mounted: function () {
        this.connectToWebsocket();
        window.addEventListener("keypress", function(e) {
            this.handleKeyboard(e.key)
        }.bind(this));
    },
    methods: {
        connectToWebsocket() {
            this.ws = new WebSocket(this.serverUrl + this.user.name);
            this.ws.addEventListener('open', (event) => {
                this.onWebsocketOpen(event)
            });
            this.ws.addEventListener('message', (event) => {
                this.handleNewMessage(event)
            });
        },
        onWebsocketOpen() {
            console.log("connected to WS!");
        },
        handleKeyboard(k){
            let actionId = 0
            switch (k){
                case "w":
                    actionId = 1;
                    break;
                case "s":
                    actionId = 2;
                    break;
                case "d":
                    actionId = 4;
                    break;
                case "a":
                    actionId = 3;
                    break;
                default:
                    return;
            }
            this.ws.send(JSON.stringify({
                type_id: actionId,
                user_id: this.user.name,
            }));
        },
        handleNewMessage(event) {
            let data = event.data;
            data = data.split(/\r?\n/);

            for (let i = 0; i < data.length; i++) {
                let msg = JSON.parse(data[i]);
                console.log(msg)

                let x = document.getElementById(msg.user_id).style.top
                let y = document.getElementById(msg.user_id).style.left
                switch (msg.type_id) {
                    case 1:
                        console.log("Ust")
                        console.log(parseInt(x.slice(0,x.length-2))+20)
                        document.getElementById(msg.user_id).style.top = (parseInt(x.slice(0,x.length-2))-10).toString() + "px"
                        break;
                    case 2:
                        document.getElementById(msg.user_id).style.top =(parseInt(x.slice(0,x.length-2))+10).toString() + "px"
                        break;
                    case 3:
                        document.getElementById(msg.user_id).style.left = (parseInt(y.slice(0,y.length-2))-10).toString() + "px"
                        break;
                    case 4:
                        document.getElementById(msg.user_id).style.left = (parseInt(y.slice(0,y.length-2))+10).toString() + "px"
                        break;
                    default:
                        break;
                }

            }
        },
        sendMessage(room) {
            if (room.newMessage !== "") {
                this.ws.send(JSON.stringify({
                    action: 'send-message',
                    message: room.newMessage,
                    target: {
                        id: room.id,
                        name: room.name
                    }
                }));
                room.newMessage = "";
            }
        }
    }
})