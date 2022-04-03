import axios from 'axios';
import https from 'https';

const httpsAgent = new https.Agent({
    requestCert: true,
    rejectUnauthorized: false
})
const client = axios.create({
    baseURL: 'https://localhost:8081',
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true,
    httpsAgent: httpsAgent
})

export { client }