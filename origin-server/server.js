const express = require('express');
const path = require('path');
const app = express();

app.use(express.static('/app/content'));

app.get('/api/info', (req, res) => {
  res.json({
    message: 'Hello from Origin Server!',
    timestamp: new Date().toISOString(),
    server: 'origin-server'
  });
});

app.get('/sample.json', (req, res) => {
  const freshData = {
    message: "This is fresh content!",
    timestamp: new Date().toLocaleString("en-IN", {
      timeZone: "Asia/Kolkata",
      hour12: false
    }),
    cached: false,
    requestId: Math.random().toString(36),  
    server: "origin-server"
  };

  console.log(`Origin: Generated fresh data at ${freshData.timestamp}`);
  res.json(freshData);
});

app.get('/health', (req, res) => {
  res.send('Origin Server OK');
});

app.use((req, res, next) => {
  console.log(`Origin Server: ${req.method} ${req.path}`);
  next();
});

const PORT = 3000;
app.listen(PORT, () => {
  console.log(`Origin Server running on port ${PORT}`);
});

