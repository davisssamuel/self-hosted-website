const http = require('http');
const fs = require('fs');
const path = require('path');

const port = 3000;

const server = http.createServer(function(req, res) {
    if (req.url === '/' || req.url === '/index.html') {
        res.writeHead(200, { 'Content-Type': 'text/html' });
        fs.readFile('index.html', function(error, data) {
            if (error) {
                res.writeHead(404);
                res.write('Error: File Not Found');
            } else {
                res.write(data);
            }
            res.end();
        });
    } else if (req.url === '/style.css') {
        const cssPath = path.join(__dirname, 'style.css');
        res.writeHead(200, { 'Content-Type': 'text/css' });
        fs.readFile(cssPath, function(error, data) {
            if (error) {
                res.writeHead(404);
                res.write('Error: File Not Found');
            } else {
                res.write(data);
            }
            res.end();
        });
    } else {
        // Handle other routes or resources as needed
        res.writeHead(404);
        res.write('Error: Route Not Found');
        res.end();
    }
});

server.listen(port, function(error) {
    if (error) {
        console.log('Something went wrong', error);
    } else {
        console.log('Server is listening on port ' + port);
    }
});


// const http = require('http')
// const fs = require('fs')
// const port = 3000

// const server = http.createServer(function(req, res) {
//     res.writeHead(200, {'Content-Type': 'text/html'})
//     fs.readFile('index.html', function(error, data) {
//         if (error) {
//             res.writeHead(404)
//             res.write('Error: File Not Found')
//         } else {
//             res.write(data)
//         }
//         res.end()
//     }) 
// })

// server.listen(port, function(error) {
//     if (error) {
//         console.log('Something went wrong', error)
//     } else {
//         console.log('Server is listening on port ' + port)
//     }
// })