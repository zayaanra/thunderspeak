function API_CreateRoom() {
    // TODO - should probably specify whether this is for chat or video room
    let createRoom = async () => {
        const response = await fetch('/api/createRoom', {
            method: 'GET',
        })
        window.location.href = response.url;
    }
    createRoom();
}