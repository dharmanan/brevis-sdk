const express = require('express');
const { exec } = require('child_process');
const path = require('path');
const cors = require('cors');


const app = express();
const port = 3001;

const corsOptions = {
    origin: 'https://silver-fiesta-xjrw6r57gxwcp4pv-5173.app.github.dev',
    methods: "GET,HEAD,PUT,PATCH,POST,DELETE",
    preflightContinue: false,
    optionsSuccessStatus: 204
};
app.use(cors(corsOptions));
app.use(express.json());

// 4. Define API endpoint
// Frontend will send requests to this address
app.post('/prove', (req, res) => {
    const { txHash } = req.body;

    if (!txHash) {
        return res.status(400).json({ status: 'error', message: 'txHash is required' });
    }

    console.log(`[${new Date().toISOString()}] Received request to prove tx: ${txHash}`);
    console.log(`[${new Date().toISOString()}] Starting ZK proof generation... This will take a few minutes.`);

    // Path to the ZK circuit folder.
    const circuitPath = path.join(__dirname, '..', 'examples', 'uniswap-prime');

    // Command to run: change directory and run go test.
    const command = `cd ${circuitPath} && go test -v`;

    exec(command, (error, stdout, stderr) => {
        if (error) {
            // Command failed (exit code != 0)
            console.error(`[${new Date().toISOString()}] Proof FAILED. Exit code: ${error.code}`);
            res.status(500).json({ status: 'error', message: 'Proof generation failed.', details: stderr || stdout });
            return;
        }
    // Command succeeded (exit code == 0)
        console.log(`[${new Date().toISOString()}] Proof for ${txHash} SUCCEEDED.`);
        res.json({ status: 'success', message: 'Proof generated and verified successfully!', details: stdout });
    });
});

// 5. Start the server
app.listen(port, () => {
    console.log(`Backend server running on http://localhost:${port}`);
});
