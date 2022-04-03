import axios from 'axios';

const client = axios.create({
    baseURL: 'https://localhost.hedwig100.tk:8081',
    headers: {
        'Content-Type': 'application/json',
    },
    withCredentials: true
})

export { client }