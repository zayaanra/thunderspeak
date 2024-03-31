function API_CreateRoom() {
    // TODO - should probably specify whether this is for chat or video room
    const username = document.getElementById("username").value;
    let createRoom = async () => {
        const response = await fetch('/api/createRoom', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                // NOTE: this should be a CSRF token
            },
            body: JSON.stringify({'username': username})
        })
        window.location.href = response.url;
    }
    createRoom();
}

function API_JoinRoom(event) {
    console.log("testr");
    if (event.key === "Enter") {
        console.log("test");
        // TODO - should probably specify whether this is for chat or video room
        const username = document.getElementById("username").value;
        const roomCode = document.getElementById("roomCode").value;
        let joinRoom = async () => {
            const response = await fetch('/api/joinRoom', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    // NOTE: this should be a CSRF token
                },
                body: JSON.stringify({'username': username, "room_code": roomCode})
            })
            window.location.href = response.url;
        }
        joinRoom();
        }
}