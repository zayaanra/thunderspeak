function API_CreateRoom() {
    // TODO - should probably specify whether this is for chat or video room
    let createRoom = async () => {
        const response = await fetch('http://localhost:8080/api/createRoom', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                roomName: 'test'
            })
        })
        const data = await response.json();
    }
    createRoom();
}