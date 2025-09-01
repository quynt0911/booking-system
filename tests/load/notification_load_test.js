const autocannon = require('autocannon');
const WebSocket = require('ws');

// WebSocket Load Test
async function websocketLoadTest() {
    const connections = 1000;
    const messages = 10000;
    const clients = [];

    // Create WebSocket connections
    for (let i = 0; i < connections; i++) {
        const ws = new WebSocket('ws://localhost:8084/ws');
        clients.push(ws);
    }

    // Send messages
    for (let i = 0; i < messages; i++) {
        const client = clients[i % connections];
        client.send(JSON.stringify({
            type: 'test',
            content: `Message ${i}`
        }));
    }
}

// HTTP Load Test
async function httpLoadTest() {
    const result = await autocannon({
        url: 'http://localhost:8084/notifications',
        connections: 100,
        duration: 10,
        method: 'POST',
        headers: {
            'content-type': 'application/json'
        },
        body: JSON.stringify({
            type: 'email',
            recipient: 'test@example.com',
            subject: 'Load Test',
            content: 'Test Content'
        })
    });

    console.log(result);
}

// Run tests
async function runTests() {
    console.log('Starting WebSocket load test...');
    await websocketLoadTest();
    
    console.log('Starting HTTP load test...');
    await httpLoadTest();
}

runTests(); 